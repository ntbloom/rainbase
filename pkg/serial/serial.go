// Package serial controls serial communication over USB
package serial

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// tags for TLV packets, uint8 numbers coded in ASCII
const (
	rain        = 48
	temperature = 49
	softReset   = 50
	hardReset   = 51
	pause       = 52
	unpause     = 53
)

type Uart struct {
	port     string
	baudrate uint16
	data     []byte
	file     *os.File
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, baudrate uint16) (*Uart, error) {
	logrus.Debugf("opening connection on %s", port)
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
	logrus.Debugf("closing file %s", uart.port)
	err := uart.file.Close()
	if err != nil {
		logrus.Errorf("problem closing %s: %s", uart.port, err)
	}
}

// ReadFile: read the file contents
func (uart *Uart) ReadFile() error {
	logrus.Tracef("reading from file %s", uart.port)

	bufLength := 7
	buf := make([]byte, bufLength)
	_, err := uart.file.Read(buf)
	if err != nil {
		logrus.Errorf("unable to open %s: %s", uart.port, err)
		return err
	}

	tag := buf[0]
	logrus.Debugf("tag=%d", tag)
	var name string
	switch tag {
	case rain:
		name = "rain"
		go uart.HandleRain()
	case temperature:
		name = "temperature"
		go uart.HandleTemp()
	case softReset:
		name = "soft reset"
		go uart.HandleSoftReset()
	case hardReset:
		name = "hard reset"
		go uart.HandleHardReset()
	case pause:
		name = "pause"
		go uart.HandlePause()
	case unpause:
		name = "unpause"
		go uart.HandleUnpause()
	default:
		logrus.Error("unsupported tag")
	}

	length, err := strconv.Atoi(string(buf[1]))
	if err != nil {
		logrus.Errorf("bad atoi conversion: %d", length)
	}
	val := make([]byte, length)
	for i := 0; i < length; i++ {
		val[i] = buf[2+i]
	}

	// print the tag and value
	fmt.Printf("\ntag=%s\nvalue=", name)
	for _, v := range val {
		fmt.Printf("%s", string(v))
	}
	fmt.Print("\n")
	return nil
}

// HandleRain: process rain event
func (uart *Uart) HandleRain() {
	logrus.Debug("calling HandleRain")
}

// HandleTemp: process temperature measurement
func (uart *Uart) HandleTemp() {
	logrus.Debug("calling HandleTemp")
}

// HandleSoftReset: process soft reset
func (uart *Uart) HandleSoftReset() {
	logrus.Debug("calling HandleSoftReset")
}

// HandleHardReset: process hard reset
func (uart *Uart) HandleHardReset() {
	logrus.Debug("calling HandleHardReset")
}

// HandlePause: process pause
func (uart *Uart) HandlePause() {
	logrus.Debug("calling HandlePause")
}

// HandleUnpause: process unpause
func (uart *Uart) HandleUnpause() {
	logrus.Debug("calling HandleUnpause")
}
