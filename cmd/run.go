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

	"github.com/spf13/viper"

	"github.com/StudioAquatan/hacku2020/pkg/email"
	"github.com/spf13/cobra"
)

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
	_ = viper.BindEnv("run.server", "EMAIL_SERVER")
	_ = viper.BindEnv("run.addr", "EMAIL_ADDR")
	_ = viper.BindEnv("run.password", "EMAIL_PASSWORD")
	_ = viper.BindEnv("run.box", "EMAIL_BOX")
}

func runServer() {
	server := viper.GetString("run.server")
	addr := viper.GetString("run.addr")
	pass := viper.GetString("run.password")
	box := viper.GetString("run.box")
	body := make(chan string)

	go email.WatchEmail(body, server, box, addr, pass)

	for {
		log.Printf("body: %s", <-body)
	}

}
