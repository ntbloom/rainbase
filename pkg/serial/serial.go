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

// signals for serial channel
const (
	stop  = 1
	start = 2
)

type Serial struct {
	port     string
	baudrate uint16
	data     []byte
	file     *os.File
	signal   chan uint8
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, baudrate uint16) (*Serial, error) {
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

	signal := make(chan uint8)
	uart := &Serial{port, baudrate, data, open, signal}

	return uart, nil
}

// Close: close the serial connection
func (serial *Serial) Close() {
	logrus.Debugf("closing file %s", serial.port)
	err := serial.file.Close()
	if err != nil {
		logrus.Errorf("problem closing %s: %s", serial.port, err)
	}
}

// Loop: read the file contents
func (serial *Serial) Loop() error {
	logrus.Tracef("reading from file %s", serial.port)

	bufLength := 7
	buf := make([]byte, bufLength)
	_, err := serial.file.Read(buf)
	if err != nil {
		logrus.Errorf("unable to open %s: %s", serial.port, err)
		return err
	}

	tag := buf[0]
	logrus.Debugf("tag=%d", tag)
	var name string
	switch tag {
	case rain:
		name = "rain"
		go serial.HandleRain()
	case temperature:
		name = "temperature"
		go serial.HandleTemp()
	case softReset:
		name = "soft reset"
		go serial.HandleSoftReset()
	case hardReset:
		name = "hard reset"
		go serial.HandleHardReset()
	case pause:
		name = "pause"
		go serial.HandlePause()
	case unpause:
		name = "unpause"
		go serial.HandleUnpause()
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
func (serial *Serial) HandleRain() {
	logrus.Debug("calling HandleRain")
}

// HandleTemp: process temperature measurement
func (serial *Serial) HandleTemp() {
	logrus.Debug("calling HandleTemp")
}

// HandleSoftReset: process soft reset
func (serial *Serial) HandleSoftReset() {
	logrus.Debug("calling HandleSoftReset")
}

// HandleHardReset: process hard reset
func (serial *Serial) HandleHardReset() {
	logrus.Debug("calling HandleHardReset")
}

// HandlePause: process pause
func (serial *Serial) HandlePause() {
	logrus.Debug("calling HandlePause")
}

// HandleUnpause: process unpause
func (serial *Serial) HandleUnpause() {
	logrus.Debug("calling HandleUnpause")
}
