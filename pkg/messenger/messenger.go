// Package messenger ferries data between serial port, paho, and the database
package messenger

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/ntbloom/rainbase/pkg/tlv"
	"github.com/sirupsen/logrus"
)

// temperature values
const (
	DegreesF = string("\u00B0F")
	DegreesC = string("\u00B0C")
)

// values for static status messages
const (
	OK              = "gatewayOkay"
	SensorLost      = "sensorLost"
	SensorPause     = "sensorPause"
	SensorUnpause   = "sensorUnpause"
	SensorSoftReset = "sensorSoftReset"
	SensorHardReset = "sensorHardReset"
)

// Payload generic type of message we'll send over MQTT
type Payload interface {
	Process() ([]byte, error)
}

// generic wrapper for all implementations of Payload
func process(p Payload) ([]byte, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// SensorStatus gives static message about what's happening to the sensor
type SensorStatus struct {
	status    string
	timestamp time.Time
}

// TemperatureEvent sends current temperature in Celsius
type TemperatureEvent struct {
	tempC     int
	timestamp time.Time
}

// RainEvent sends message about rain event
type RainEvent struct {
	value     float32
	timestamp time.Time
}

// SensorStatus.Process turn static value into mqtt payload
func (s *SensorStatus) Process() ([]byte, error) {
	return process(s)
}

// TemperatureEvent.Process turn temp into mqtt payload
func (t *TemperatureEvent) Process() ([]byte, error) {
	return process(t)
}

// RainEvent.Process turn rain event into mqtt payload
func (r *RainEvent) Process() ([]byte, error) {
	return process(r)
}

type Message struct {
	client   mqtt.Client
	topic    string
	retained bool
	qos      byte
	payload  []byte
}

// NewMessage makes a new message from a tlv packet mqtt topic
func NewMessage(client mqtt.Client, topic string, packet *tlv.TLV, rainAmt *float32) (*Message, error) {
	now := time.Now()
	var payload []byte
	var err error

	switch packet.Tag {
	case tlv.Rain:
		rain := RainEvent{
			value:     *rainAmt,
			timestamp: now,
		}
		payload, err = rain.Process()
		if err != nil {
			return nil, err
		}
	case tlv.Temperature:
		temp := TemperatureEvent{
			tempC:     packet.Value,
			timestamp: now,
		}
		payload, err = temp.Process()
		if err != nil {
			return nil, err
		}
	case tlv.SoftReset:
		soft := SensorStatus{
			status:    SensorSoftReset,
			timestamp: now,
		}
		payload, err = soft.Process()
		if err != nil {
			return nil, err
		}
	case tlv.HardReset:
		hard := SensorStatus{
			status:    SensorHardReset,
			timestamp: now,
		}
		payload, err = hard.Process()
		if err != nil {
			return nil, err
		}
	case tlv.Pause:
		pause := SensorStatus{
			status:    SensorPause,
			timestamp: now,
		}
		payload, err = pause.Process()
		if err != nil {
			return nil, err
		}
	case tlv.Unpause:
		unpause := SensorStatus{
			status:    SensorUnpause,
			timestamp: now,
		}
		payload, err = unpause.Process()
		if err != nil {
			return nil, err
		}
	default:
		logrus.Errorf("unsupported tag %d", packet.Tag)
		return nil, nil
	}
	msg := Message{
		client:   client,
		topic:    topic,
		retained: false,
		qos:      0,
		payload:  payload,
	}
	return &msg, nil

}

// SendMessage sends payload to the cloud over mqtt
func (m *Message) SendMessage() {
	m.client.Publish(m.topic, m.qos, m.retained, m.payload)
}
