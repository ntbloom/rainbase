package serial

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Uart struct {
	port     string
	baudrate uint16
	data     []byte
	file     *os.File
}

// NewConnection: create a new serial connection with a unix-style file descriptor
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

// Print the file descriptor
func (uart *Uart) GetFileDescriptor() string {
	return uart.port
}

// open a file descriptor and poll the file

// process rain event

// process temperature event

// process soft reset

// process hard reset

// process pause

// process unpause
