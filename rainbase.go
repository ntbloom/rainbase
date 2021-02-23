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

// set the logger level
func setLogger() {
	level, err := logrus.ParseLevel(viper.GetString(configkey.Loglevel))
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}
	logrus.Infof("logger set to %s level", logrus.GetLevel())
}

// get a serial connection
func getSerialConnection(messenger *messenger.Messenger) (*serial.Serial, error) {
	conn, err := serial.NewConnection(
		viper.GetString(configkey.USBConnectionPort),
		viper.GetInt(configkey.USBPacketLengthMax),
		viper.GetDuration(configkey.USBConnectionTimeout),
		messenger,
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
	dataMessenger := messenger.NewMessenger(client)
	conn, err := getSerialConnection(dataMessenger)
	if err != nil {
		logrus.Fatal(err)
	}
	go dataMessenger.Listen()

	go conn.GetTLV()
	if duration > 0 {
		for i := 0; i < duration; i++ {
			time.Sleep(time.Second)
			logrus.Tracef("sleep #%d", i+1)
		}
		conn.State <- configkey.SerialClosed
		dataMessenger.State <- configkey.SerialClosed
	}
	logrus.Info("reaching end of listener, program about to end")
}

func main() {
	config.GetConfig()
	setLogger()

	// run the main listening loop
	duration := 10
	listen(duration)
}
