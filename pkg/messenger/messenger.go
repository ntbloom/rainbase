// Package messenger ferries data between serial port, paho, and the database
package messenger

import (
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/ntbloom/rainbase/pkg/paho"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	defer m.client.Disconnect(viper.GetUint(configkey.MQTTQuiescence))

	// loop until signal
	for {
		select {

		case state := <-m.State:
			switch state {
			case configkey.SerialClosed:
				logrus.Debug("received `Closed` signal, closing mqtt connection")
				return
			case configkey.SendStatusMessage:
				logrus.Debug("requesting status message")
				m.SendStatus()
			default:
				continue
			}
		case msg := <-m.Data:
			logrus.Tracef("received Message from serial port: %s", msg.payload)
			m.Publish(msg)
		}
	}
}

// Publish sends a Message over MQTT
func (m *Messenger) Publish(msg *Message) {
	logrus.Tracef("sending Message over MQTT: %s", msg.payload)
	m.client.Publish(msg.topic, msg.qos, msg.retained, msg.payload)
}

// SendStatus sends a status message about the gateway and sensor at regular interval
func (m *Messenger) SendStatus() {
	// assume if this code is running that the gateway is up
	gwStatus, _ := gatewayStatusMessage()
	m.Publish(gwStatus)

	sensorStatus, _ := sensorStatusMessage()
	m.Publish(sensorStatus)
}

// get a status message about how the gateway is doing
func gatewayStatusMessage() (*Message, error) {
	gs := GatewayStatus{
		Topic:     paho.GatewayStatus,
		OK:        true,
		Timestamp: time.Time{},
	}
	msg, err := gs.Process()
	if err != nil {
		return nil, err
	}

	return &Message{
		topic:    gs.Topic,
		retained: false,
		qos:      0,
		payload:  msg,
	}, nil
}

// get a status message about how the sensor is doing
func sensorStatusMessage() (*Message, error) {
	var up bool
	port := viper.GetString(configkey.USBConnectionPort)
	_, err := os.Stat(port)
	if err != nil {
		up = false
	} else {
		up = true
	}
	ss := SensorStatus{
		Topic:     paho.SensorStatus,
		OK:        up,
		Timestamp: time.Time{},
	}
	msg, err := ss.Process()
	if err != nil {
		return nil, err
	}
	return &Message{
		topic:    ss.Topic,
		retained: false,
		qos:      0,
		payload:  msg,
	}, nil
}
