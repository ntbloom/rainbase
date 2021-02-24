// Package messenger ferries data between serial port, paho, and the database
package messenger

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/sirupsen/logrus"
)

// Messenger receives Message from serial port, publishes to mqtt and stores locally
type Messenger struct {
	client mqtt.Client
	State  chan int
	Data   chan *Message
}

// NewMessenger get a new messenger
func NewMessenger(client mqtt.Client) *Messenger {
	state := make(chan int)
	data := make(chan *Message)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("unable to connect to MQTT: %s", token.Error())
	}
	return &Messenger{client, state, data}
}

// Wait for packet to publish or to receive signal interrupt
func (m *Messenger) Listen() {
	defer m.client.Disconnect(1000)

	// loop until signal
	for {
		select {
		//case m.State <- configkey.SerialClosed:
		case closed := <-m.State:
			if closed == configkey.SerialClosed {
				logrus.Debug("received `Closed` signal, closing mqtt connection")
				return
			}
		case msg := <-m.Data:
			logrus.Debugf("received Message from serial port: %s", msg.payload)
			m.Publish(msg)
		}
	}
}

func (m *Messenger) Publish(msg *Message) {
	m.client.Publish(msg.topic, msg.qos, msg.retained, msg.payload)
}
