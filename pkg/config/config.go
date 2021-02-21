// Package config sets and retrieves configuration values
package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Map all yaml params to a constant
const (
	Loglevel = "logger.level"

	USBPacketLengthMax   = "usb.packet.length.max"
	USBConnectionPort    = "usb.connection.port"
	USBConnectionTimeout = "usb.connection.timeout"

	MQTTScheme     = "mqtt.scheme"
	MQTTBrokerIP   = "mqtt.broker.ip"
	MQTTBrokerPort = "mqtt.broker.port"
	MQTTCaCert     = "mqtt.certs.ca"
	MQTTClientCert = "mqtt.certs.client"
	MQTTClientKey  = "mqtt.certs.key"
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
	fmt.Println(viper.GetString(MQTTBrokerIP))
}
