// Package messenger ferries data between serial port, paho, and the database
package messenger

import (
	"github.com/ntbloom/rainbase/pkg/tlv"
	"github.com/sirupsen/logrus"
)

// temperature values
const (
	DegreesF = string("\u00B0F")
	DegreesC = string("\u00B0C")
)

func Handle(packet *tlv.TLV) {
	switch packet.Tag {
	case tlv.Rain:
		handleRain()
	case tlv.Temperature:
		handleTemp(packet.Value)
	case tlv.SoftReset:
		handleSoftReset()
	case tlv.HardReset:
		handleHardReset()
	case tlv.Pause:
		handlePause()
	case tlv.Unpause:
		handleUnpause()
	default:
		logrus.Errorf("unsupported tag %d", packet.Tag)
	}
}

func handleRain() {
	logrus.Debug("calling handleRain")
}

func handleTemp(value int) {
	logrus.Debug("calling handleTemp")
	f := ((9 * value) / 5) + 32
	logrus.Debugf("temp is %d%s/%d%s", value, DegreesC, f, DegreesF)
}

func handleSoftReset() {
	logrus.Debug("calling handleSoftReset")
}

func handleHardReset() {
	logrus.Debug("calling handleHardReset")
}

func handlePause() {
	logrus.Debug("calling handlePause")
}

func handleUnpause() {
	logrus.Debug("calling HandleUnpause")
}
