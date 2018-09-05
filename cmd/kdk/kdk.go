// Copyright © 2018 Cisco Systems, Inc.
//
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
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cisco-sso/kdk/pkg/kdk"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var CurrentKdkEnvConfig = kdk.KdkEnvConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kdk",
	Short: "Kubernetes Development Kit",
	Long: `

 _  __ ____  _  __
/ |/ //  _ \/ |/ /
|   / | | \||   / 
|   \ | |_/||   \ 
\_|\_\\____/\_|\_\
                  

A full kubernetes development environment in a container`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to execute RootCmd.")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	CurrentKdkEnvConfig.Init()

	rootCmd.PersistentFlags().StringVar(&CurrentKdkEnvConfig.ConfigFile.AppConfig.Name, "name", "kdk", "KDK name")
	rootCmd.PersistentFlags().BoolVarP(&CurrentKdkEnvConfig.ConfigFile.AppConfig.Debug, "debug", "d", false, "Debug Mode")
}

func initConfig() {

	if _, err := os.Stat(CurrentKdkEnvConfig.ConfigRootDir()); os.IsNotExist(err) {
		err = os.Mkdir(CurrentKdkEnvConfig.ConfigRootDir(), 0700)
		if err != nil {
			logrus.WithField("err", err).Fatal("Unable to create Config Directory")
		}
	}

	viper.SetConfigFile(CurrentKdkEnvConfig.ConfigPath())

	viper.SetEnvPrefix("kdk")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithFields(logrus.Fields{"configFileUsed": viper.ConfigFileUsed(), "err": err}).Warnln("Failed to load KDK config.")
	}

	if viper.GetBool("json") {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	if _, err := os.Stat(CurrentKdkEnvConfig.ConfigPath()); err == nil {
		// read the config.yaml file
		data, err := ioutil.ReadFile(CurrentKdkEnvConfig.ConfigPath())
		if err != nil {
			logrus.WithField("err", err).Fatalf("Failed to read configFile %v", CurrentKdkEnvConfig.ConfigPath())
		}

		err = yaml.Unmarshal(data, &CurrentKdkEnvConfig.ConfigFile)
		if err != nil {
			logrus.WithField("err", err).Error("Corrupted or deprecated kdk config file format")
			logrus.Fatal("Please rebuild config file with `kdk init`")
		}
	}
}