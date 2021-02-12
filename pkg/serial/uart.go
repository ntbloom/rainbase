// Package serial controls serial communication over UART
package serial

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Uart struct {
	port     string
	baudrate uint16
	data     []byte
	file     *os.File
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, baudrate uint16) (*Uart, error) {
	var data []byte

	_, err := os.Stat(port)
	if err != nil {
		logrus.Errorf("file descriptor %s does not exist", port)
		return nil, err
	}

	open, err := os.Open(port)
	if err != nil {
		logrus.Errorf("problem opening file at %s: %s", port, err)
		return nil, err
	}

	uart := &Uart{port, baudrate, data, open}

	return uart, nil
}

// Close: close the serial connection
func (uart *Uart) Close() {
	err := uart.file.Close()
	if err != nil {
		logrus.Errorf("problem closing %s: %s", uart.port, err)
	}
}

// GetFileDescriptor: print the file descriptor
func (uart *Uart) GetFileDescriptor() string {
	return uart.port
}

// ReadFile: read the file contents
func (uart *Uart) ReadFile() error {
	bufLength := 10
	buf := make([]byte, bufLength)
	_, err := uart.file.Read(buf)
	if err != nil {
		logrus.Errorf("problem reading %s: %s", uart.port, err)
		return err
	}
	tag := buf[0]
	length, err := strconv.Atoi(string(buf[1]))
	if err != nil {
		logrus.Errorf("bad conversion: %d", length)
	}
	val := make([]byte, length)
	for i := 0; i < length; i++ {
		val[i] = buf[2+i]
	}

	// print the tag and value
	fmt.Printf("tag=%d\nvalue=", tag)
	for _, v := range val {
		fmt.Printf("%s", string(v))
	}
	fmt.Print("\n")
	return nil
}

// process rain event

// process temperature event

// process soft reset

// process hard reset

// process pause

// process unpause
