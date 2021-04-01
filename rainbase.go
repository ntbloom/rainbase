package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/ntbloom/rainbase/pkg/database"
	"github.com/ntbloom/rainbase/pkg/messenger"
	"github.com/ntbloom/rainbase/pkg/paho"
	"github.com/ntbloom/rainbase/pkg/serial"
	"github.com/ntbloom/rainbase/pkg/timer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// connect to paho
func connectToMQTT() mqtt.Client {
	client, err := paho.NewConnection(paho.GetConfigFromViper())
	if err != nil {
		panic(err)
	}
	return client
}

// connect to the sqlite database
func connectToDatabase() *database.DBConnector {
	db, err := database.NewSqliteDBConnector(viper.GetString(configkey.DatabaseLocalDevFile), true)
	if err != nil {
		panic(err)
	}
	return db
}

// get a serial connection
func connectSerialPort(msgr *messenger.Messenger) *serial.Serial {
	conn, err := serial.NewConnection(
		viper.GetString(configkey.USBConnectionPort),
		viper.GetInt(configkey.USBPacketLengthMax),
		viper.GetDuration(configkey.USBConnectionTimeout),
		msgr,
	)
	if err != nil {
		panic(err)
	}
	return conn
}

type kill struct{ channels []chan uint8 }

func stopLoop(channels []chan uint8) {
	for _, channel := range channels {
		channel <- configkey.SerialClosed
	}
}

func (k *kill) DoAction() {
	stopLoop(k.channels)
	logrus.Info("main loop killed by timer, exiting program")
	os.Exit(0)
}

func startKillTimer(duration time.Duration, killChannels []chan uint8) *timer.Timer {
	k := kill{killChannels}
	t := timer.NewTimer(duration, &k)
	return t
}

// run main loop for number of seconds or indefinitely
// for debugging/testing purposes; not to be used for production
func listen(duration time.Duration) {
	client := connectToMQTT()
	db := connectToDatabase()
	msgr := messenger.NewMessenger(client, db)
	conn := connectSerialPort(msgr)

	// start the listening threads
	go msgr.Listen()
	go conn.GetTLV()

	killChannels := []chan uint8{conn.State, msgr.State}

	// wait for timer to interrupt
	if duration > 0 {
		logrus.Infof("program running for %s duration", duration)
		t := startKillTimer(duration, killChannels)
		t.Loop()
	}

	// or else run indefinitely until the shell is interrupted
	sigs := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		sig := <-sigs
		logrus.Infof("program received %s signal, exiting", sig)
		stopLoop(killChannels)
		done <- true
	}()
	logrus.Info("program running indefinitely")
	<-done
}

func main() {
	// read config from the config file
	config.Configure()

	// run the main listening loop
	duration := time.Second * -1
	listen(duration)
}
