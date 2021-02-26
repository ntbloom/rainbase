package messenger

// Message defines what an individual message looks like

import (
	"encoding/json"
	"time"

	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/ntbloom/rainbase/pkg/paho"

	"github.com/ntbloom/rainbase/pkg/tlv"
	"github.com/sirupsen/logrus"
)

// values for static status messages
const (
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
	val := string(payload)
	logrus.Debug(val)
	return payload, nil
}

// SensorEvent gives static message about what's happening to the sensor
type SensorEvent struct {
	Topic     string
	Status    string
	Timestamp time.Time
}

// TemperatureEvent sends current temperature in Celsius
type TemperatureEvent struct {
	Topic     string
	TempC     int
	Timestamp time.Time
}

// RainEvent sends message about rain event
type RainEvent struct {
	Topic       string
	Millimeters string // send as a string to avoid floating point weirdness
	Timestamp   time.Time
}

// GatewayStatus sends "OK" message at regular intervals
type GatewayStatus struct {
	Topic     string
	OK        bool
	Timestamp time.Time
}

// SensorStatus sends "OK" if sensor is reachable, else "Bad"
type SensorStatus struct {
	Topic     string
	OK        bool
	Timestamp time.Time
}

// SensorEvent.Process turn static value into mqtt payload
func (s *SensorEvent) Process() ([]byte, error) {
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

// GatewayStatus.Process turn gateway status message into mqtt payload
func (gs *GatewayStatus) Process() ([]byte, error) {
	return process(gs)
}

// SensorStatus.Process turn sensor status message into mqtt payload
func (ss *SensorStatus) Process() ([]byte, error) {
	return process(ss)
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
			paho.RainTopic,
			viper.GetString(configkey.SensorRainMetric),
			now,
		}
		payload, err = rain.Process()
		topic = rain.Topic
		if err != nil {
			return nil, err
		}
	case tlv.Temperature:
		temp := TemperatureEvent{
			paho.TemperatureTopic,
			packet.Value,
			now,
		}
		payload, err = temp.Process()
		topic = temp.Topic
		if err != nil {
			return nil, err
		}
	case tlv.SoftReset:
		soft := SensorEvent{
			paho.SensorEvent,
			SensorSoftReset,
			now,
		}
		payload, err = soft.Process()
		topic = soft.Topic
		if err != nil {
			return nil, err
		}
	case tlv.HardReset:
		hard := SensorEvent{
			paho.SensorEvent,
			SensorHardReset,
			now,
		}
		payload, err = hard.Process()
		topic = hard.Topic
		if err != nil {
			return nil, err
		}
	case tlv.Pause:
		pause := SensorEvent{
			paho.SensorEvent,
			SensorPause,
			now,
		}
		payload, err = pause.Process()
		topic = pause.Topic
		if err != nil {
			return nil, err
		}
	case tlv.Unpause:
		unpause := SensorEvent{
			paho.SensorEvent,
			SensorUnpause,
			now,
		}
		payload, err = unpause.Process()
		topic = unpause.Topic
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