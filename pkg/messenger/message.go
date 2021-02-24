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
	val := string(payload)
	logrus.Debug(val)
	return payload, nil
}

// SensorStatus gives static message about what's happening to the sensor
type SensorStatus struct {
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
		soft := SensorStatus{
			paho.SoftResetTopic,
			SensorSoftReset,
			now,
		}
		payload, err = soft.Process()
		topic = soft.Topic
		if err != nil {
			return nil, err
		}
	case tlv.HardReset:
		hard := SensorStatus{
			paho.HardResetTopic,
			SensorHardReset,
			now,
		}
		payload, err = hard.Process()
		topic = hard.Topic
		if err != nil {
			return nil, err
		}
	case tlv.Pause:
		pause := SensorStatus{
			paho.PauseTopic,
			SensorPause,
			now,
		}
		payload, err = pause.Process()
		topic = pause.Topic
		if err != nil {
			return nil, err
		}
	case tlv.Unpause:
		unpause := SensorStatus{
			paho.UnpauseTopic,
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
