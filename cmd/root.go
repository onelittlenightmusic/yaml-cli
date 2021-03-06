/*
Copyright © 2021 Hiroyuki Osaki <hiroyuki.osaki@gmail.com>

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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yaml-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

    if fileName == "" {
        fmt.Println("Please provide yaml file by using -f option")
        return
    }

    yamlFile, _ := ioutil.ReadFile(fileName)
		serializedCommands := SerializedCommandsFile{}
    err := yaml.Unmarshal(yamlFile, &serializedCommands)
    if err != nil {
        fmt.Printf("Error parsing YAML file: %s\n", err)
    }
		str, _ := json.Marshal(serializedCommands)
		fmt.Print(string(str))

		for _, v := range serializedCommands.Spec.Cmds {
			out, err := runCommand(v)
			if err != nil {
				fmt.Errorf("%v", err)
			}	else {
				fmt.Println(out)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yaml-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".yaml-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".yaml-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func runCommand(cmd SerializedCommand) (out string, err error) {
	fmt.Println(cmd)

	plannedCommand := exec.Command(cmd.Cmd, stringifyOpts(cmd)...)
	var outBuf bytes.Buffer
	plannedCommand.Stdout = &outBuf
	err = plannedCommand.Run()

	if err == nil {
		out = outBuf.String()
	}
	return out, err
}

func stringifyOpts(cmd SerializedCommand) ([]string) {
	rtn := []string{}
	for _, v := range cmd.Opts {
		rtn = append(rtn, fmt.Sprintf("-%s", v.Opt))
	}
	return rtn
}