// Package tlv handles decoding tlv packets into more useful data
package tlv

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

// decode from ascii representation of the byte
func asciiToInt(b byte) (int, error) {
	return strconv.Atoi(string(b))
}

// decode byte array into single integer
func byteArrayToInt(b []byte) (int, error) {
	return -1, nil
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
	tag, err := asciiToInt(rawTag)
	if err != nil {
		logrus.Errorf("misreading tag %d", rawTag)
		return nil, err
	}

	rawLength := packet[1]
	length, err := asciiToInt(rawLength)
	if err != nil {
		logrus.Errorf("misreading length %d", rawLength)
		return nil, err
	}

	var value int
	switch length {
	case 1:
		// static value, doesn't matter
		value = 1
	case 4:
		// convert it to an integer
		rawValue := packet[2:]
		value, err = byteArrayToInt(rawValue)
		if err != nil {
			logrus.Errorf("unable to decode %s", rawValue)
			return nil, err
		}
	default:
		err := fmt.Errorf("unsupported value %d", value)
		return nil, err
	}
	//value := make([]byte, length)
	//for i := 0; i < length; i++ {
	//	value[i] = packet[2+i]
	//}

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
