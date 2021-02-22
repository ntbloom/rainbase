// Package paho wraps the Eclipse Paho code for handling mqtt messaging
package paho

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Connection is the entrypoint for all things MQTT
type Connection struct {
	Options *mqtt.ClientOptions
	Client  *mqtt.Client
}

// NewConnection creates a new MQTT connection or error
func NewConnection(scheme, broker string, port int, caCert, clientCert, clientKey string) (*Connection, error) {
	options := mqtt.NewClientOptions()

	// add broker
	server := fmt.Sprintf("%s://%s:%d", scheme, broker, port)
	logrus.Infof("opening MQTT connection at %s", server)
	options.AddBroker(server)

	// configure tls
	tlsConfig, err := configureTlsConfig(caCert, clientCert, clientKey)
	if err != nil {
		return nil, err
	}
	options.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(options)
	return &Connection{
		options,
		&client,
	}, nil
}

// get a new config for ssl
func configureTlsConfig(caCertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	// import CA from file
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		logrus.Errorf("problem reading CA file at %s: %s", caCertFile, err)
	}
	certpool.AppendCertsFromPEM(ca)

	// match client cert and key
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		logrus.Errorf("problem with cert/key pair: %s", err)
		return nil, err
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}, nil

}
