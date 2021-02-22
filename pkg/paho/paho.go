// Package paho wraps the Eclipse Paho code for handling mqtt messaging
package paho

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

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
	scheme     string
	broker     string
	port       int
	caCert     string
	clientCert string
	clientKey  string
}

// GetConfigFromViper get paho configuration details from viper directly
func GetConfigFromViper() *ConnectionConfig {
	return &ConnectionConfig{
		scheme:     viper.GetString(configkey.MQTTScheme),
		broker:     viper.GetString(configkey.MQTTBrokerIP),
		port:       viper.GetInt(configkey.MQTTBrokerPort),
		caCert:     viper.GetString(configkey.MQTTCaCert),
		clientCert: viper.GetString(configkey.MQTTClientCert),
		clientKey:  viper.GetString(configkey.MQTTClientKey),
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

	client := mqtt.NewClient(options)
	return &Connection{
		options,
		&client,
	}, nil
}

// get a new config for ssl
func configureTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	// import CA from file
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		logrus.Errorf("problem reading CA file at %s: %s", caCertFile, err)
		return nil, err
	}
	certpool.AppendCertsFromPEM(ca)

	// match client cert and key
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		logrus.Errorf("problem with cert/key pair: %s", err)
		return nil, err
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true, //nolint:gosect
		Certificates:       []tls.Certificate{cert},
	}, nil
}
