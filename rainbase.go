package main

import (
	"rainbase/pkg/serial"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getUartConfig() {
	viper.SetConfigName("arduino")
	viper.AddConfigPath("./config/")
	// add additional config locations for prod

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("config not loaded")
	}
}

func main() {
	getUartConfig()

	port := viper.GetString("port")
	baudrate := uint16(viper.GetInt("baudrate"))
	uart, err := serial.NewConnection(port, baudrate)
	if err != nil {
		logrus.Fatal(err)
	}
	defer uart.Close()

	for i := 0; i < 10; i++ {
		err = uart.ReadFile()
		if err != nil {
			logrus.Fatal("need to handle arduino resetting")
		}
	}
}
