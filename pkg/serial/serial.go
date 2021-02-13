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

const (
	Open   = 1
	Closed = 2
)

type Serial struct {
	port         string
	maxPacketLen int
	data         []byte
	file         *os.File
	State        chan uint8
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, maxPacketLen int) (*Serial, error) {
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

	state := make(chan uint8)
	uart := &Serial{port, maxPacketLen, data, file, state}

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
func (serial *Serial) GetMessage() {
	logrus.Tracef("reading contents of file at %s", serial.port)

	for {
		packet := make([]byte, serial.maxPacketLen)
		_, err := serial.file.Read(packet)
		if err != nil {
			logrus.Errorf("unable to open %s: %s", serial.port, err)
			return
		}

		tag := packet[0]
		length, err := strconv.Atoi(string(packet[1]))
		if err != nil {
			logrus.Errorf("bad atoi conversion: %d", length)
			return
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
		select {
		case state := <-serial.State:
			if state == Closed {
				logrus.Debug("received `Closed` signal")
				serial.Close()
				break
			}
		default:
			continue
		}
		return
	}
}

// HandleRain: process rain event
func (serial *Serial) HandleRain() {
	logrus.Debug("calling HandleRain")
}

// HandleTemp: process temperature measurement
func (serial *Serial) HandleTemp(value []byte) {
	logrus.Debug("calling HandleTemp")
	logrus.Debugf("write code to process %d", value)
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
