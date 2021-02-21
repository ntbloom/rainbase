package paho_test

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"

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

// Can
func TestMQTTConnection(t *testing.T) {
	conn := pahoFixture()
	client := *conn.Client
	token := client.Connect()
	fmt.Println(token)
	defer client.Disconnect(1000)
	if !client.IsConnected() {
		logrus.Error("failed to connect")
		t.Fail()
	}
}
