// Package configkey maps yaml configs to a constant
package configkey

const (
	Loglevel = "logger.level"

	USBPacketLengthMax   = "usb.packet.length.max"
	USBConnectionPort    = "usb.connection.port"
	USBConnectionTimeout = "usb.connection.timeout"

	MQTTScheme            = "mqtt.scheme"
	MQTTBrokerIP          = "mqtt.broker.ip"
	MQTTBrokerPort        = "mqtt.broker.port"
	MQTTCaCert            = "mqtt.certs.ca"
	MQTTClientCert        = "mqtt.certs.client"
	MQTTClientKey         = "mqtt.certs.key"
	MQTTConnectionTimeout = "mqtt.connection.timeout"
	MQTTQuiescence        = "mqtt.connection.quiescence"
	MQTTQos               = "mqtt.qos"

	SensorRainCustomary = "sensor.rain.measurement.inches"
	SensorRainMetric    = "sensor.rain.measurement.mm"

	DatabaseLocalProdFile = "database.local.prod.file"
	DatabaseLocalDevFile  = "database.local.dev.file"

	MessengerStatusInterval      = "messenger.status.interval"
	MessengerTemperatureInterval = "messenger.temperature.interval"
)

// random constants, not tied to config
const (
	DegreesF          = string("\u00B0F")
	DegreesC          = string("\u00B0C")
	SerialClosed      = uint8(1)
	SendStatusMessage = uint8(2)
)
