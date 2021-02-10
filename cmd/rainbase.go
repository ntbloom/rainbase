package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"rainbase/pkg/uart"
)

func main() {
	uart, err := uart.NewConnection("/dev/ttyACM98", 9600)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(uart.GetFileDescriptor())
}
