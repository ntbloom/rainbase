package main

import (
	"time"

	"github.com/ntbloom/rainbase/pkg/config"

	"github.com/ntbloom/rainbase/pkg/serial"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// set the logger level
func setLogger() {
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}
	logrus.Infof("logger set to %s level", logrus.GetLevel())
}

// get a serial connection
func getSerialConnection() (*serial.Serial, error) {
	conn, err := serial.NewConnection(
		viper.GetString("connection.port"),
		viper.GetInt("packet.length.max"),
		viper.GetDuration("connection.timeout"),
	)
	return conn, err
}

// run main loop for number of seconds or indefinitely
func listen(duration int) {
	conn, err := getSerialConnection()
	if err != nil {
		logrus.Fatal(err)
	}

	go conn.GetTLV()
	if duration > 0 {
		for i := 0; i < duration; i++ {
			time.Sleep(time.Second)
			logrus.Tracef("sleep #%d", i+1)
		}
		conn.State <- serial.Closed
	}
	logrus.Error("reaching end of listener, program about to end")
}

func main() {
	config.GetConfig()
	setLogger()

	// run the main listening loop
	duration := 10
	listen(duration)
}
