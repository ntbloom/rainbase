package tlv

import (
	"fmt"
	"testing"
)

// Test static packets like rain event, etc.
func TestStaticVals(t *testing.T) {
	var exp int32 = 1
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

// Test temperature packets
func Test18(t *testing.T) {
	var exp int32 = 1
	// positives
	temp1 := []byte{49, 52, 48, 48, 49, 50}
	exp = 18
	if !verifyValToInt(temp1, exp) {
		t.Fail()
	}
}
func Test25(t *testing.T) {
	temp2 := []byte{49, 52, 48, 48, 49, 57}
	var exp int32 = 25
	if !verifyValToInt(temp2, exp) {
		t.Fail()
	}
}
func Test26(t *testing.T) {
	temp3 := []byte{49, 52, 48, 48, 49, 65}
	var exp int32 = 26
	if !verifyValToInt(temp3, exp) {
		t.Fail()
	}
}

func Test0(t *testing.T) {
	// zero
	temp4 := []byte{49, 52, 48, 48, 48, 48}
	var exp int32 = 0
	if !verifyValToInt(temp4, exp) {
		t.Fail()
	}
}

func TestMinus24(t *testing.T) {
	// negatives
	temp5 := []byte{49, 52, 70, 70, 69, 55}
	var exp int32 = -24
	if !verifyValToInt(temp5, exp) {
		t.Fail()
	}
}

func verifyValToInt(raw []byte, expected int32) bool {
	tlv, err := NewTLV(raw)
	if tlv == nil {
		fmt.Printf("error making tlv: %s\n", err)
		return false
	}
	actual := tlv.Value
	if actual != expected {
		fmt.Printf("expected=%d, actual=%d\n", expected, actual)
		return false
	}
	return true
}
