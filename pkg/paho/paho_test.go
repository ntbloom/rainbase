package paho_test

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/spf13/viper"

	"github.com/ntbloom/rainbase/pkg/paho"
)

// reusable paho function
func pahoFixture(t *testing.T) *paho.Connection {
	config.GetConfig()
	scheme := viper.GetString(config.MQTTScheme)
	broker := viper.GetString(config.MQTTBrokerIP)
	port := viper.GetInt(config.MQTTBrokerPort)
	caCert := viper.GetString(config.MQTTCaCert)
	clientCert := viper.GetString(config.MQTTClientCert)
	clientKey := viper.GetString(config.MQTTClientKey)
	conn, err := paho.NewConnection(scheme, broker, port, caCert, clientCert, clientKey)
	if err != nil {
		t.Fail()
	}
	return conn
}

// Can
func TestMQTTConnection(t *testing.T) {
	conn := pahoFixture(t)
	client := *conn.Client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Fail()

	}
	defer client.Disconnect(1000)
	if !client.IsConnected() {
		logrus.Error("failed to connect")
		t.Fail()
	}
	client.Publish("hello", 0, false, "world")
}
