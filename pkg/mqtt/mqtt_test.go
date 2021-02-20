package mqtt_test

import (
	"testing"

	_ "github.com/ntbloom/rainbase/pkg/mqtt"
	_ "github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	// get config from viper
}

// can we communicate with the cloud?
func TestMQTTConnection(t *testing.T) {
	//t.Fail()
	// connect to mqtt
	// listen for messages
	// publish hello
	// confirm message is received
}
