// Package tlv handles decoding tlv packets into more useful data
package tlv

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// tags for TLV packets
const (
	Rain        = 0
	Temperature = 1
	SoftReset   = 2
	HardReset   = 3
	Pause       = 4
	Unpause     = 5
)

// packet length for determining how to process value
const (
	constant = 1
	variable = 4
)

const maxInt = 65535

// decode from ascii representation of the byte
func asciiToInt(b byte) int {
	// use a hash table for faster lookups than other conversion
	dict := map[byte]int{
		48: 0,
		49: 1,
		50: 2,
		51: 3,
		52: 4,
		53: 5,
		54: 6,
		55: 7,
		58: 8,
		59: 9,
		65: 10,
		66: 11,
		67: 12,
		68: 13,
		69: 14,
		70: 15,
	}
	return dict[b]
}

//  concatenate a 4-byte array into its integer equivalent
func concatenateBytesToInt(b []byte) int {
	asNums := make([]int, 4)
	for idx, val := range b {
		asNums[idx] = int(asciiToInt(val))
	}
	value := asNums[0] << 12
	value = value | (asNums[1] << 8)
	value = value | (asNums[2] << 4)
	value = value | (asNums[3])

	// account for negative numbers
	if asNums[0] > 0 {
		value -= maxInt
	}
	return value

}

// TLV: tag, length, value encoding for binary packets received over serial
type TLV struct {
	Tag    int
	Length int
	Value  int
}

// NewTLV: make a new TLV packet
func NewTLV(packet []byte) (*TLV, error) {
	rawTag := packet[0]
	tag := asciiToInt(rawTag)

	rawLength := packet[1]
	length := asciiToInt(rawLength)

	var value int
	switch length {
	case constant:
		// static value, doesn't matter
		value = 1
	case variable:
		// convert it to an integer
		rawValue := packet[2:5] // packet[6] is newline
		value = concatenateBytesToInt(rawValue)
	default:
		err := fmt.Errorf("unsupported value %d", value)
		return nil, err
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
