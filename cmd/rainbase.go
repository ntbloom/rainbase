package main

import (
	"fmt"
	"rainbase/pkg/serial"

	"github.com/sirupsen/logrus"
)

func main() {
	uart, err := serial.NewConnection("/dev/ttyACM99", 9600)
	if err != nil {
		logrus.Fatal(err)
	}
	defer uart.Close()

	fmt.Println(uart.GetFileDescriptor())
}
