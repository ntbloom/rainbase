package paho_test

import (
	"testing"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/spf13/viper"

	"github.com/ntbloom/rainbase/pkg/paho"
)

// reusable paho function
func pahoFixture() *paho.Connection {
	config.GetConfig()
	scheme := viper.GetString(config.MQTTScheme)
	broker := viper.GetString(config.MQTTBrokerIP)
	port := viper.GetInt(config.MQTTBrokerPort)
	caCert := viper.GetString(config.MQTTCaCert)
	clientCert := viper.GetString(config.MQTTClientCert)
	clientKey := viper.GetString(config.MQTTClientKey)
	return paho.NewConnection(scheme, broker, port, caCert, clientCert, clientKey)
}

func TestMQTTConnection(t *testing.T) {
	conn := pahoFixture()
	client := *conn.Client
	client.Connect()
	if !client.IsConnected() {
		t.Fail()
	}
}
