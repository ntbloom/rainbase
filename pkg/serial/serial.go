// Package serial controls serial communication over USB
package serial

import (
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

type Serial struct {
	port         string
	baudrate     uint16
	maxPacketLen int
	data         []byte
	file         *os.File
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, baudrate uint16, maxPacketLen int) (*Serial, error) {
	logrus.Infof("opening connection on %s", port)
	var data []byte

	_, err := os.Stat(port)
	if err != nil {
		logrus.Errorf("file descriptor %s does not exist", port)
		return nil, err
	}

	file, err := os.Open(port)
	if err != nil {
		logrus.Errorf("problem opening port %s: %s", port, err)
		return nil, err
	}

	uart := &Serial{port, baudrate, maxPacketLen, data, file}

	return uart, nil
}

// Close: close the serial connection
func (serial *Serial) Close() {
	logrus.Infof("closing serial port %s", serial.port)
	err := serial.file.Close()
	if err != nil {
		logrus.Errorf("problem closing %s: %s", serial.port, err)
	}
}

// GetMessage: read the file contents
func (serial *Serial) GetMessage() error {
	logrus.Tracef("reading contents of file at %s", serial.port)

	packet := make([]byte, serial.maxPacketLen)
	_, err := serial.file.Read(packet)
	if err != nil {
		logrus.Errorf("unable to open %s: %s", serial.port, err)
		return err
	}

	tag := packet[0]
	length, err := strconv.Atoi(string(packet[1]))
	if err != nil {
		logrus.Errorf("bad atoi conversion: %d", length)
		return err
	}
	value := make([]byte, length)
	for i := 0; i < length; i++ {
		value[i] = packet[2+i]
	}
	logrus.Debugf("packet=%s", string(packet))
	logrus.Debugf("tag=%d", tag)
	logrus.Debugf("length=%d", length)
	logrus.Debugf("value=%d", value)

	switch tag {
	case rain:
		go serial.HandleRain()
	case temperature:
		go serial.HandleTemp(value)
	case softReset:
		go serial.HandleSoftReset()
	case hardReset:
		go serial.HandleHardReset()
	case pause:
		go serial.HandlePause()
	case unpause:
		go serial.HandleUnpause()
	default:
		logrus.Error("unsupported tag")
	}

	return nil
}

func printBuf(buf []byte, length int) {
}

// HandleRain: process rain event
func (serial *Serial) HandleRain() {
	logrus.Debug("calling HandleRain")
}

// HandleTemp: process temperature measurement
func (serial *Serial) HandleTemp(value []byte) {
	logrus.Debug("calling HandleTemp")
	logrus.Errorf("write code to process %d", value)
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
