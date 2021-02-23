// Package paho wraps the Eclipse Paho code for handling mqtt messaging
package paho

import (
	"fmt"
	"time"

	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Connection is the entrypoint for all things MQTT
type Connection struct {
	Options *mqtt.ClientOptions
	Client  *mqtt.Client
}

// ConnectionConfig configures the paho connection
type ConnectionConfig struct {
	scheme            string
	broker            string
	port              int
	caCert            string
	clientCert        string
	clientKey         string
	connectionTimeout time.Duration
}

// GetConfigFromViper get paho configuration details from viper directly
func GetConfigFromViper() *ConnectionConfig {
	return &ConnectionConfig{
		scheme:            viper.GetString(configkey.MQTTScheme),
		broker:            viper.GetString(configkey.MQTTBrokerIP),
		port:              viper.GetInt(configkey.MQTTBrokerPort),
		caCert:            viper.GetString(configkey.MQTTCaCert),
		clientCert:        viper.GetString(configkey.MQTTClientCert),
		clientKey:         viper.GetString(configkey.MQTTClientKey),
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
	}
}

// NewConnection creates a new MQTT connection or error
func NewConnection(config *ConnectionConfig) (*Connection, error) {
	options := mqtt.NewClientOptions()

	// add broker
	server := fmt.Sprintf("%s://%s:%d", config.scheme, config.broker, config.port)
	logrus.Infof("opening MQTT connection at %s", server)
	options.AddBroker(server)

	// configure tls
	tlsConfig, err := configureTLSConfig(config.caCert, config.clientCert, config.clientKey)
	if err != nil {
		return nil, err
	}
	options.SetTLSConfig(tlsConfig)

	// miscellaneous options
	options.SetConnectTimeout(config.connectionTimeout)

	client := mqtt.NewClient(options)
	return &Connection{
		options,
		&client,
	}, nil
}
