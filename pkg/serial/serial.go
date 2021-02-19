// Package serial controls serial communication over USB
package serial

import (
	"os"
	"time"

	"github.com/ntbloom/rainbase/pkg/exitcodes"
	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"
)

// uint8 for state of serial port
const Closed = 1

// temperature values
const (
	DegreesF = string("\u00B0F")
	DegreesC = string("\u00B0C")
)

// Serial communicates with a serial port
type Serial struct {
	port         string
	maxPacketLen int
	timeout      time.Duration
	data         []byte
	file         *os.File
	State        chan uint8
}

// NewConnection: create a new serial connection with a unix filename
func NewConnection(port string, maxPacketLen int, timeout time.Duration) (*Serial, error) {
	checkPortStatus(port, timeout)
	logrus.Infof("opening connection on `%s`", port)
	var data []byte

	// attempt to connect until timeout is exhausted

	file, err := os.Open(port)
	if err != nil {
		logrus.Errorf("problem opening port `%s`: %s", port, err)
		return nil, err
	}

	state := make(chan uint8)

	uart := &Serial{
		port,
		maxPacketLen,
		timeout,
		data,
		file,
		state,
	}

	return uart, nil
}

// Close: close the serial connection
func (serial *Serial) Close() {
	logrus.Infof("closing serial port `%s`", serial.port)
	err := serial.file.Close()
	if err != nil {
		logrus.Errorf("problem closing `%s`: %s", serial.port, err)
	}
}

// GetMessage: read the file contents
func (serial *Serial) GetMessage() {
	checkPortStatus(serial.port, serial.timeout)

	logrus.Tracef("reading contents of `%s`", serial.port)
	for {
		packet := make([]byte, serial.maxPacketLen)
		_, err := serial.file.Read(packet)
		if err != nil {
			// connection to file was lost, attempt reconnection
			logrus.Infof("connection lost, attempting reconnection")
			checkPortStatus(serial.port, serial.timeout)
			_ = serial.reopenConnection()
			continue
		}

		tlvPacket, err := tlv.NewTLV(packet)
		if err != nil {
			logrus.Errorf("unexpected TLV packet: %s", err)
		}
		tag := tlvPacket.Tag
		switch tag {
		case tlv.Rain:
			go serial.HandleRain()
		case tlv.Temperature:
			go serial.HandleTemp(tlvPacket.Value)
		case tlv.SoftReset:
			go serial.HandleSoftReset()
		case tlv.HardReset:
			go serial.HandleHardReset()
		case tlv.Pause:
			go serial.HandlePause()
		case tlv.Unpause:
			go serial.HandleUnpause()
		default:
			logrus.Errorf("unsupported tag `%d`", tag)
		}

		// run forever until uninterrupted by close signal
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
	logrus.Debugf("checking if `%s` exists", port)
	start := time.Now()
	for {
		_, err := os.Stat(port)
		if err == nil {
			logrus.Debugf("found port `%s`", port)
			return
		}
		logrus.Tracef("file `%s` doesn't exist on first look, re-checking for %s", port, timeout)
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
		logrus.Debugf("port `%s` temporarily down: %s", serial.port, err)
		return err
	}
	serial.file = file
	return nil
}

// HandlePortFailure: what to do when sensor is unresponsive?
func HandlePortFailure(port string) {
	logrus.Fatalf("unable to locate sensor at `%s`", port)

	// for now...
	os.Exit(exitcodes.SerialPortNotFound)
}

// HandleRain: process rain event
func (serial *Serial) HandleRain() {
	logrus.Debug("calling HandleRain")
}

// HandleTemp: process temperature measurement
func (serial *Serial) HandleTemp(value int) {
	logrus.Debug("calling HandleTemp")
	f := ((9 * value) / 5) + 32
	logrus.Debugf("temp is %d%s/%d%s", value, DegreesC, f, DegreesF)
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
