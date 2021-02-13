package main

import (
	"rainbase/pkg/serial"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// process config files
func getConfig() {
	viper.SetConfigName("arduino")
	viper.AddConfigPath("./config/")
	// add additional config locations for prod

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("config not loaded")
	}
}

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
		viper.GetString("port"),
		uint16(viper.GetInt("baudrate")),
		viper.GetInt("maxPacketLen"),
	)
	return conn, err
}

func main() {
	getConfig()
	setLogger()

	conn, err := getSerialConnection()
	if err != nil {
		logrus.Fatal(err)
	}

	// main loop with interrupt channel; for now just closing after certain time
	duration := 10
	go conn.GetMessage()
	for i := 0; i < duration; i++ {
		time.Sleep(time.Second)
	}
	conn.State <- serial.Closed

}
