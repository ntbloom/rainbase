package uart

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Uart struct {
	filedescriptor string
	baudrate       uint16
	data           []byte
}

func NewConnection(file string, baudrate uint16) (*Uart, error) {
	var data []byte
	uart := &Uart{file, baudrate, data}
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return uart, nil
}

// Print the file descriptor
func (uart *Uart) GetFileDescriptor() (string, error) {
	fileinfo, err := os.Stat(uart.filedescriptor)
	if err != nil {
		logrus.Fatal(err)
	}
	return fileinfo.Name(), nil

}

// open a file descriptor and poll the file

// process rain event

// process temperature event

// process soft reset

// process hard reset

// process pause

// process unpause
