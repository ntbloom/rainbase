package serial

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Uart struct {
	fileDescriptor string
	baudrate       uint16
	data           []byte
	file           *os.File
}

// NewConnection: create a new serial connection with a unix-style file descriptor
func NewConnection(filename string, baudrate uint16) (*Uart, error) {
	var data []byte

	_, err := os.Stat(filename)
	if err != nil {
		logrus.Errorf("file descriptor %s does not exist", filename)
		return nil, err
	}

	open, err := os.Open(filename)
	if err != nil {
		logrus.Errorf("problem opening file at %s: %s", filename, err)
		return nil, err
	}

	uart := &Uart{filename, baudrate, data, open}

	return uart, nil
}

// Close: close the serial connection
func (uart *Uart) Close() {
	err := uart.file.Close()
	if err != nil {
		logrus.Errorf("problem closing %s: %s", uart.fileDescriptor, err)
	}
}

// Print the file descriptor
func (uart *Uart) GetFileDescriptor() string {
	return uart.fileDescriptor
}

// open a file descriptor and poll the file

// process rain event

// process temperature event

// process soft reset

// process hard reset

// process pause

// process unpause
