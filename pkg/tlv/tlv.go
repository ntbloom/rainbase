// Package tlv handles decoding tlv packets into more useful data
package tlv

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type TLV struct {
	Tag    byte
	Length int
	Value  []byte
}

// NewTLV: make a new TLV packet
func NewTLV(packet []byte) (*TLV, error) {

	tag := packet[0]
	length, err := strconv.Atoi(string(packet[1]))
	if err != nil {
		logrus.Errorf("bad atoi conversion for Length=`%d`", length)
		return nil, err
	}
	value := make([]byte, length)
	for i := 0; i < length; i++ {
		value[i] = packet[2+i]
	}

	logrus.Tracef("packet=%s", string(packet))
	logrus.Tracef("Tag=%d", tag)
	logrus.Tracef("Length=%d", length)
	logrus.Tracef("Value=%d", value)
	tlv := &TLV{
		tag,
		length,
		value,
	}
	return tlv, nil
}
