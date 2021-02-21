// Package mqtt sends data packets to a cloud server using the MQTT protocol
package mqtt

import (
	"fmt"

	"github.com/sirupsen/logrus"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// Connection is the entrypoint for all things MQTT
type Connection struct {
	Options *paho.ClientOptions
	Client  *paho.Client
}

// NewConnection creates a new MQTT connection or error
func NewConnection(scheme string, broker string, port int) *Connection {
	server := fmt.Sprintf("%s://%s:%d", scheme, broker, port)
	logrus.Info("opening MQTT connection at %s", server)
	options := paho.NewClientOptions()
	options.AddBroker(server)
	client := paho.NewClient(options)
	conn := Connection{
		options,
		&client,
	}

	return &conn
}
