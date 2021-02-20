// Package serial controls serial communication over USB
package serial

import (
	"os"
	"time"

	"github.com/ntbloom/rainbase/pkg/messenger"

	"github.com/ntbloom/rainbase/pkg/exitcodes"
	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"
)

// uint8 for state of serial port
const Closed = 1


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

// GetTLV: read the file contents
func (serial *Serial) GetTLV() {
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
		go messenger.Handle(tlvPacket)

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

// HandlePortFailure: what to do when sensor is unresponsive?
func HandlePortFailure(port string) {
	logrus.Fatalf("unable to locate sensor at `%s`", port)

	// for now...
	os.Exit(exitcodes.SerialPortNotFound)
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
