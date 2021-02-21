package paho_test

import (
	"testing"

	_ "github.com/ntbloom/rainbase/pkg/paho"
	_ "github.com/spf13/viper"
)

func TestMQTTConnection(t *testing.T) {
	//t.Fail()
	// connect to paho
	// listen for messages
	// publish hello
	// confirm message is received
}
