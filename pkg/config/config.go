// Package config sets and retrieves configuration values
package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configFile = "rainbase"
	secrets    = ".secrets"
)

// process config files
func GetConfig() {
	// get the base config
	viper.SetConfigName(configFile)
	viper.AddConfigPath("pkg/config/")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("config not loaded")
	}

	// bring in secrets
	viper.SetConfigName(secrets)
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatal("secrets not loaded")
	}
	fmt.Println(viper.GetString("cloud.broker.ip"))
}
