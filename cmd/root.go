// Copyright © 2017 Abdisamad Hashi <shaybix@tuta.io>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	dotenvFile string
	envConfig  map[string]string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "loadenv",
	Short: "Loadenv loads environment for a laravel project using Docker",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := load(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.loadenv.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().StringVar(&dotenvFile, "dotenv", "", "dotenv file with environment variables")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".loadenv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".loadenv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// load reads from a .env file by default unless given --file
// flag has been set.
func load() error {

	var fname string

	if dotenvFile != "" {
		// load file
		fname = dotenvFile
	} else {
		fname = ".env"
	}

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return fmt.Errorf("can not find %s file in the local directory", fname)
	}

	if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
		return fmt.Errorf("can not find Dockerfile file in the local directory")
	}

	if err := loadEnvVars(fname); err != nil {
		return err
	}

	if err := startDocker(); err != nil {
		return err
	}

	return nil
}

//loadEnvVars will load environment variables from file
func loadEnvVars(fname string) error {

	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		} else {
			envVar := strings.Split(line, "=")

			if err := os.Setenv(envVar[0], envVar[1]); err != nil {
				return err
			}
		}
	}

	return nil
}

// startDocker will orchestrate the docker containers by executing the docker-compose
// command in the shell.
func startDocker() error {

	dockerComposeBuildCmd := exec.Command("docker-compose", "build", ".")
	dockerComposeBuildCmd.Stdout = os.Stdout
	dockerComposeBuildCmd.Stderr = os.Stderr
	if err := dockerComposeBuildCmd.Run(); err != nil {
		return err
	}

	dockerComposeUpCmd := exec.Command("docker-compose", "up")
	dockerComposeUpCmd.Stdout = os.Stdout
	dockerComposeUpCmd.Stderr = os.Stderr

	if err := dockerComposeUpCmd.Run(); err != nil {
		return err
	}

	return nil
}

// stopDocker stops docker environment for the project in the
// current working directory
func stopDocker() error {

	dockerComposeDownCmd := exec.Command("docker-compose", "down")
	dockerComposeDownCmd.Stdout = os.Stdout
	dockerComposeDownCmd.Stderr = os.Stderr

	if err := dockerComposeDownCmd.Run(); err != nil {
		return err
	}

	if err := cleanup(); err != nil {
		return err
	}

	return nil
}

// cleanup will clean up files/directory created
func cleanup() error {

	// remove tmp dir in the project folder

	return nil
}
