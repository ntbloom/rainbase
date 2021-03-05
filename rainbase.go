package main

import (
	"time"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/ntbloom/rainbase/pkg/messenger"
	"github.com/ntbloom/rainbase/pkg/paho"
	"github.com/ntbloom/rainbase/pkg/serial"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// get a serial connection
func getSerialConnection(msgr *messenger.Messenger) (*serial.Serial, error) {
	conn, err := serial.NewConnection(
		viper.GetString(configkey.USBConnectionPort),
		viper.GetInt(configkey.USBPacketLengthMax),
		viper.GetDuration(configkey.USBConnectionTimeout),
		msgr,
	)
	return conn, err
}

// run main loop for number of seconds or indefinitely
// for debugging/testing purposes; not to be used for production
func listen(duration int) {
	client, err := paho.NewConnection(paho.GetConfigFromViper())
	if err != nil {
		panic(err)
	}
	msgr := messenger.NewMessenger(client)
	conn, err := getSerialConnection(msgr)
	if err != nil {
		panic(err)
	}
	go msgr.Listen()

	go conn.GetTLV()
	if duration > 0 {
		for i := 0; i < duration; i++ {
			time.Sleep(time.Second)
			logrus.Tracef("sleep #%d", i+1)
		}
		conn.State <- configkey.SerialClosed
		msgr.State <- configkey.SerialClosed
	}
	logrus.Info("reaching end of listener, program about to end")
}

func main() {
	config.Configure()

	// run the main listening loop
	duration := 5
	listen(duration)
}
