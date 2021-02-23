// Package config sets and retrieves configuration values
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config files
const (
	configDir     = "/etc/rainbase/"
	mainConfig    = "rainbase"
	secretsConfig = "secrets"
)

// GetConfig process config files
func GetConfig() {
	// get the base config
	viper.SetConfigName(mainConfig)
	viper.AddConfigPath(configDir)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("config not loaded: %s", err)
	}

	// bring in secrets
	viper.SetConfigName(secretsConfig)
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatal("secrets not loaded")
	}
}
