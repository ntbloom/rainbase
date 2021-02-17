package tlv

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// Test static packets like rain event, etc.
func TestStaticVals(t *testing.T) {
	exp := 1
	rain := []byte{48, 49, 49}
	if !verifyValToInt(rain, exp) {
		t.Fail()
	}
	softReset := []byte{50, 49, 49}
	if !verifyValToInt(softReset, exp) {
		t.Fail()
	}
	hardReset := []byte{51, 49, 49}
	if !verifyValToInt(hardReset, exp) {
		t.Fail()
	}
	pause := []byte{52, 49, 49}
	if !verifyValToInt(pause, exp) {
		t.Fail()
	}
}

func TestValToInt(t *testing.T) {
	// positives
	temp1 := []byte{49, 52, 48, 48, 49, 50}
	exp1 := 18
	if !verifyValToInt(temp1, exp1) {
		t.Fail()
	}
	temp2 := []byte{49, 52, 48, 48, 49, 57}
	exp2 := 25
	if !verifyValToInt(temp2, exp2) {
		t.Fail()
	}
	temp3 := []byte{49, 52, 48, 48, 49, 65}
	exp3 := 26
	if !verifyValToInt(temp3, exp3) {
		t.Fail()
	}

	// zero
	temp4 := []byte{49, 52, 48, 48, 48, 48}
	exp4 := 0
	if !verifyValToInt(temp4, exp4) {
		t.Fail()
	}

	// negatives
	temp5 := []byte{49, 52, 70, 70, 69, 55}
	exp5 := -24
	if !verifyValToInt(temp5, exp5) {
		t.Fail()
	}
}

func verifyValToInt(raw []byte, expected int) bool {
	tlv, err := NewTLV(raw)
	if tlv == nil {
		logrus.Errorf("error making tlv: %s", err)
		return false
	}
	actual := tlv.Value
	if actual != expected {
		logrus.Errorf("expected=%d, actual=%d", expected, actual)
		return false
	}
	return true
}
