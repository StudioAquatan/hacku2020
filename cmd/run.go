/*
Copyright © 2020 StudioAquatan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/StudioAquatan/hacku2020/pkg/email"
	"github.com/StudioAquatan/hacku2020/pkg/slack"

	"github.com/StudioAquatan/hacku2020/pkg/character"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type YamlConfig struct {
	Characters []character.Info
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run server",
	Long:  `run server`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

const LightEndpointPath = "/light"
const M5stackEndpointPath = "/status"

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.StringP("character-config-path", "c", "", "A path to the character config file")
	flags.StringP("message-num", "n", "1", "a number of slack message")
	flags.StringP("light-host", "l", "", "ip or hostname for lighting")
	flags.StringP("m5stack-host", "m", "", "ip or hostname to m5stack")

	_ = viper.BindPFlag("run.config", flags.Lookup("character-config-path"))
	_ = viper.BindPFlag("run.num", flags.Lookup("message-num"))
	_ = viper.BindPFlag("run.light", flags.Lookup("light-host"))
	_ = viper.BindPFlag("run.m5stack", flags.Lookup("m5stack-host"))
	_ = viper.BindEnv("run.server", "EMAIL_SERVER")
	_ = viper.BindEnv("run.addr", "EMAIL_ADDR")
	_ = viper.BindEnv("run.password", "EMAIL_PASSWORD")
	_ = viper.BindEnv("run.box", "EMAIL_BOX")
	_ = viper.BindEnv("run.token", "SLACK_TOKEN")
	_ = viper.BindEnv("run.channel", "SLACK_CHANNEL")

	_ = runCmd.MarkFlagRequired("character-config-path")
	_ = runCmd.MarkFlagRequired("light-host")
	_ = runCmd.MarkFlagRequired("m5stack-host")
}

func runServer() {
	configPath := viper.GetString("run.config")
	messageNum := viper.GetInt("run.num")
	server := viper.GetString("run.server")
	addr := viper.GetString("run.addr")
	pass := viper.GetString("run.password")
	box := viper.GetString("run.box")
	token := viper.GetString("run.token")
	channelID := viper.GetString("run.channel")
	ecChan := make(chan email.Content)

	dirPath, fileName := filepath.Split(configPath)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dirPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	var yc YamlConfig
	err = viper.Unmarshal(&yc)
	if err != nil {
		log.Fatalf("Fatal error unmarshal config file: %s \n", err)
	}
	cis := yc.Characters

	go email.WatchEmail(ecChan, server, box, addr, pass)

	for {
		var oinori bool
		ec := <-ecChan
		if !email.ClassifyScreeningMailBySubj(ec.Subject) {
			log.Printf("[INFO] Ignored email subject: %s", ec.Subject)
			continue
		}
		if email.ClassifyAcceptanceMailByBody(ec.Body) {
			oinori = true
		}
		if !email.ClassifyOinoriMailByBody(ec.Body) && !oinori {
			log.Printf("[INFO] Ignored email Body: %s", ec.Body)
			continue
		}
		if !email.ClassifyOinoriMailBySentiment(ec.Body) {
			log.Printf("[INFO] Ignored email by sentiment score: %s", ec.Body)
			continue
		}

		wg := &sync.WaitGroup{} // WaitGroupの値を作る
		go func() {
			wg.Add(1)
			notify(oinori)
			wg.Done()
		}()

		mis := character.CreateMessageInfoByRandom(cis, messageNum, oinori)
		for _, mi := range *mis {
			i := slack.NewSlackMessageInfo(token, channelID, mi.Name, mi.Icon, mi.Message)
			err := i.PostMessage()
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}
			time.Sleep(1 * time.Second)
		}
		wg.Wait()
	}
}

func notify(oinori bool) {
	yeelightAddr := viper.GetString("run.light")
	m5stackAddr := viper.GetString("run.m5stack")
	respM5stack := make(chan string)

	go notifyLight(yeelightAddr, oinori)
	go notifyM5stack(m5stackAddr, oinori, respM5stack)

	log.Printf("[INFO] m5stack api response: %s", <-respM5stack)

	return
}

func notifyLight(addr string, oinori bool) {
	src := "positive.py"
	if oinori {
		src = "negative.py"
	}

	cmd := exec.Command("python", "./light_control/%s", src, addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[ERROR] yeelight script failed: %s", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("[ERROR] yeelight script failed: %s", err)
	}

	err = cmd.Run()
	streamReader := func(scanner *bufio.Scanner, outputChan chan string, doneChan chan bool) {
		defer close(outputChan)
		defer close(doneChan)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
		doneChan <- true
	}

	// stdout, stderrをひろうgoroutineを起動
	stdoutScanner := bufio.NewScanner(stdout)
	stdoutOutputChan := make(chan string)
	stdoutDoneChan := make(chan bool)
	stderrScanner := bufio.NewScanner(stderr)
	stderrOutputChan := make(chan string)
	stderrDoneChan := make(chan bool)
	go streamReader(stdoutScanner, stdoutOutputChan, stdoutDoneChan)
	go streamReader(stderrScanner, stderrOutputChan, stderrDoneChan)

	// channel経由でデータを引っこ抜く
	stillGoing := true
	for stillGoing {
		select {
		case <-stdoutDoneChan:
			stillGoing = false
		case line := <-stdoutOutputChan:
			log.Println("[INFO] yeelight response: %s", line)
		case line := <-stderrOutputChan:
			log.Println("[INFO] yeelight response: %s", line)
		}
	}

	//一応Waitでプロセスの終了をまつ
	ret := cmd.Wait()
	if ret != nil {
		log.Printf("[ERROR] yeelight script failed: %s", err)
	}
	if err != nil {
		log.Printf("[ERROR] yeelight script failed: %s", err)
		return
	}
}

func notifyM5stack(addr string, oinori bool, respStr chan string) {
	urlVal := url.Values{}
	if oinori {
		urlVal.Add("status", "negative")
	} else {
		urlVal.Add("status", "positive")
	}
	urlStr := "http://" + addr + M5stackEndpointPath + urlVal.Encode()
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Printf("[ERROR] POST to M5stack API failed: %s", err)
		respStr <- ""
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] ioutil.ReadAll failed: %s", err)
		respStr <- ""
		return
	}
	respStr <- string(b)
}
