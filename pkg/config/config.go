// Package config sets and retrieves configuration values
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configFile = "rainbase"
)

// process config files
func GetConfig() {
	viper.SetConfigName(configFile)
	// do we need additional production paths?
	viper.AddConfigPath("pkg/config/")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("config not loaded")
	}
}
