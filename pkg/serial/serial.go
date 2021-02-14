// Package serial controls serial communication over USB
package serial

import (
	"os"
	"rainbase/pkg/exitcodes"
	"strconv"
	"time"

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

// uint8 for state of serial port
const (
	Open   = 1
	Closed = 2
)

type Serial struct {
	port         string
	maxPacketLen int
	timeout      time.Duration
	data         []byte
	file         *os.File
	State        chan uint8
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, packetLenMax int, timeout time.Duration) (*Serial, error) {
	checkPortStatus(port, timeout)
	logrus.Infof("opening connection on %s", port)
	var data []byte

	// attempt to connect until timeout is exhausted

	file, err := os.Open(port)
	if err != nil {
		logrus.Errorf("problem opening port %s: %s", port, err)
		return nil, err
	}

	state := make(chan uint8)

	uart := &Serial{
		port,
		packetLenMax,
		timeout,
		data,
		file,
		state,
	}

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
	checkPortStatus(serial.port, serial.timeout)

	logrus.Tracef("reading contents of file at %s", serial.port)
	for {
		packet := make([]byte, serial.maxPacketLen)
		_, err := serial.file.Read(packet)
		if err != nil {
			// connection to file was lost, attempt reconnection
			logrus.Errorf("error in main read loop, attempting reconnection")
			checkPortStatus(serial.port, serial.timeout)
			_ = serial.reopenConnection()
			continue
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

// checkPortStatus: keep trying to open a file until timeout is up
func checkPortStatus(port string, timeout time.Duration) {
	logrus.Debugf("checking if %s exists", port)
	start := time.Now()
	for {
		_, err := os.Stat(port)
		if err == nil {
			logrus.Debugf("found port at %s", port)
			return
		}
		logrus.Tracef("file %s doesn't exist on first look, re-checking for %s", port, timeout)
		if time.Since(start).Milliseconds() > timeout.Milliseconds() {
			HandlePortFailure(port)
			return
		}
	}
}

// reopenConnection: get another file descriptor for the port
func (serial *Serial) reopenConnection() error {
	file, err := os.Open(serial.port)
	if err != nil {
		logrus.Errorf("problem opening %s: %s", serial.port, err)
		return err
	}
	serial.file = file
	return nil
}

// HandlePortFailure: what to do when sensor is unresponsive?
func HandlePortFailure(port string) {
	logrus.Fatalf("unable to locate sensor at %s", port)

	// for now...
	os.Exit(exitcodes.SerialPortNotFound)
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
