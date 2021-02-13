package main

import (
	"rainbase/pkg/serial"

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

func main() {
	getConfig()
	setLogger()

	// open a new serial connection
	conn, err := serial.NewConnection(viper.GetString("port"), uint16(viper.GetInt("baudrate")))
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()

	// main loop, will need a refactor
	for i := 0; i < 10; i++ {
		err = conn.Loop()
		if err != nil {
			logrus.Fatal("need to handle arduino resetting")
		}
	}
}
