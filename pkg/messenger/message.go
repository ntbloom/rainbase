package messenger

// Message defines what an individual message looks like

import (
	"encoding/json"
	"time"

	"github.com/ntbloom/rainbase/pkg/paho"

	"github.com/ntbloom/rainbase/pkg/config"

	"github.com/ntbloom/rainbase/pkg/tlv"
	"github.com/sirupsen/logrus"
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
	topic     string
	status    string
	timestamp time.Time
}

// TemperatureEvent sends current temperature in Celsius
type TemperatureEvent struct {
	topic     string
	tempC     int
	timestamp time.Time
}

// RainEvent sends message about rain event
type RainEvent struct {
	topic     string
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
	topic    string
	retained bool
	qos      byte
	payload  []byte
}

// NewMessage makes a new message from a tlv packet mqtt topic
func NewMessage(packet *tlv.TLV) (*Message, error) {
	now := time.Now()
	var payload []byte
	var topic string
	var err error

	switch packet.Tag {
	case tlv.Rain:
		rain := RainEvent{
			topic:     paho.RainTopic,
			value:     float32(config.RainAmt),
			timestamp: now,
		}
		payload, err = rain.Process()
		topic = rain.topic
		if err != nil {
			return nil, err
		}
	case tlv.Temperature:
		temp := TemperatureEvent{
			topic:     paho.TemperatureTopic,
			tempC:     packet.Value,
			timestamp: now,
		}
		payload, err = temp.Process()
		topic = temp.topic
		if err != nil {
			return nil, err
		}
	case tlv.SoftReset:
		soft := SensorStatus{
			topic:     paho.SoftResetTopic,
			status:    SensorSoftReset,
			timestamp: now,
		}
		payload, err = soft.Process()
		topic = soft.topic
		if err != nil {
			return nil, err
		}
	case tlv.HardReset:
		hard := SensorStatus{
			topic:     paho.HardResetTopic,
			status:    SensorHardReset,
			timestamp: now,
		}
		payload, err = hard.Process()
		topic = hard.topic
		if err != nil {
			return nil, err
		}
	case tlv.Pause:
		pause := SensorStatus{
			topic:     paho.PauseTopic,
			status:    SensorPause,
			timestamp: now,
		}
		payload, err = pause.Process()
		topic = pause.topic
		if err != nil {
			return nil, err
		}
	case tlv.Unpause:
		unpause := SensorStatus{
			topic:     paho.UnpauseTopic,
			status:    SensorUnpause,
			timestamp: now,
		}
		payload, err = unpause.Process()
		topic = unpause.topic
		if err != nil {
			return nil, err
		}
	default:
		logrus.Errorf("unsupported tag %d", packet.Tag)
		return nil, nil
	}
	msg := Message{
		topic:    topic,
		retained: false,
		qos:      0,
		payload:  payload,
	}
	return &msg, nil

}
