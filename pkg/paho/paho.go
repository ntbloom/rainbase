// Package paho wraps the Eclipse Paho code for handling mqtt messaging
package paho

import (
	"fmt"

	"github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Connection is the entrypoint for all things MQTT
type Connection struct {
	Options    *mqtt.ClientOptions
	Client     *mqtt.Client
	caCert     string
	clientCert string
	clientKey  string
}

// NewConnection creates a new MQTT connection or error
func NewConnection(scheme string, broker string, port int, caCert string, clientCert string, clientKey string) *Connection {
	server := fmt.Sprintf("%s://%s:%d", scheme, broker, port)
	logrus.Infof("opening MQTT connection at %s", server)
	options := mqtt.NewClientOptions()
	options.AddBroker(server)
	client := mqtt.NewClient(options)
	return &Connection{
		options,
		&client,
		caCert,
		clientCert,
		clientKey,
	}
}
