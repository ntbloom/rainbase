package uart

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Uart struct {
	filename string
	baudrate uint16
	data     []byte
}

func NewConnection(filename string, baudrate uint16) (*Uart, error) {
	var data []byte
	uart := &Uart{filename, baudrate, data}
	_, err := os.Stat(filename)
	if err != nil {
		logrus.Fatal("file descriptor does not exist")
		return nil, err
	}
	return uart, nil
}

// Print the file descriptor
func (uart *Uart) GetFileDescriptor() string {
	return uart.filename
}

// open a file descriptor and poll the file

// process rain event

// process temperature event

// process soft reset

// process hard reset

// process pause

// process unpause
