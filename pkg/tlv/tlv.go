// Package tlv handles decoding tlv packets into more useful data
package tlv

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

// Decode from ascii representation of the byte
func ascii() *map[int]int {
	ascii := map[int]int{
		48: 0,
		49: 1,
		50: 2,
		51: 3,
		52: 4,
		53: 5,
		54: 6,
		55: 7,
		56: 8,
		57: 9,
		65: 10,
		66: 11,
		67: 12,
		68: 13,
		69: 14,
		70: 15,
	}
	return &ascii
}

// TLV: tag, length, value encoding for binary packets received over serial
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

// ValToInt: convert TLV value to 64-bit integer
func (tlv *TLV) ValToInt() int64 {
	return 1
}
