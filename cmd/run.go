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
	"log"
	"path/filepath"
	"time"

	"github.com/StudioAquatan/hacku2020/pkg/email"

	"github.com/StudioAquatan/hacku2020/pkg/character"
	"github.com/StudioAquatan/hacku2020/pkg/slack"

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

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.StringP("character-config-path", "c", "", "A path to the character config file")
	flags.StringP("message-num", "n", "1", "a number of slack message")

	_ = viper.BindPFlag("run.config", flags.Lookup("character-config-path"))
	_ = viper.BindPFlag("run.num", flags.Lookup("message-num"))
	_ = viper.BindEnv("run.server", "EMAIL_SERVER")
	_ = viper.BindEnv("run.addr", "EMAIL_ADDR")
	_ = viper.BindEnv("run.password", "EMAIL_PASSWORD")
	_ = viper.BindEnv("run.box", "EMAIL_BOX")
	_ = viper.BindEnv("run.token", "SLACK_TOKEN")
	_ = viper.BindEnv("run.channel", "SLACK_CHANNEL")

	_ = runCmd.MarkFlagRequired("character-config-path")
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
		ec := <-ecChan
		if !email.ClassifyMailBySubj(ec.Subject) {
			log.Printf("[INFO] Ignored email subject: %s", ec.Subject)
			continue
		}
		if !email.ClassifyMailByBody(ec.Body) {
			log.Printf("[INFO] Ignored email Body: %s", ec.Body)
			continue
		}
		mis := character.CreateMessageInfoByRandom(cis, messageNum)
		for _, mi := range *mis {
			i := slack.NewSlackMessageInfo(token, channelID, mi.Name, mi.Icon, mi.Message)
			err := i.PostMessage()
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}
