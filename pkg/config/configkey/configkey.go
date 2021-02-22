// Package configkey maps yaml configs to a constant
package configkey

const (
	Loglevel = "logger.level"

	USBPacketLengthMax   = "usb.packet.length.max"
	USBConnectionPort    = "usb.connection.port"
	USBConnectionTimeout = "usb.connection.timeout"

	MQTTScheme     = "mqtt.scheme"
	MQTTBrokerIP   = "mqtt.broker.ip"
	MQTTBrokerPort = "mqtt.broker.port"
	MQTTCaCert     = "mqtt.certs.ca"
	MQTTClientCert = "mqtt.certs.client"
	MQTTClientKey  = "mqtt.certs.key"
)
