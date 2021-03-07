package database_test

import (
	"os"
	"sync"
	"testing"
	"testing/quick"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/ntbloom/rainbase/pkg/config/configkey"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/ntbloom/rainbase/pkg/database"
)

/* FIXTURES */

// reusable configs
func getConfig() {
	config.Configure()
}

// connectorFixture makes a reusable DBConnector object
func connectorFixture() *database.DBConnector {
	getConfig()
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	db, _ := database.NewDBConnector(sqliteFile, true)
	return db
}

/* TESTS */

// create and destroy sqlite file 5 times, get DBCOnnector struct
func TestDatabasePrep(t *testing.T) {
	getConfig()
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)

	// clean up when finished
	defer func() { _ = os.Remove(sqliteFile) }()

	// create and destroy 5 times
	for i := 0; i < 5; i++ {
		db, err := database.NewDBConnector(sqliteFile, true)
		if err != nil || db == nil {
			logrus.Error("database not created")
			t.Error(err)
		}
		_, err = os.Stat(sqliteFile)
		if err != nil {
			logrus.Error("sqlite file doesn't exist")
			t.Error(err)
		}
	}
}

// make sure foreign key contraints are enforced
func TestForeignKeysEnforced(t *testing.T) {
	db := connectorFixture()
	if foreignKeys := db.ForeignKeysAreImplemented(); !foreignKeys {
		logrus.Error("sqlite is not enforcing foreign_key constraints")
		t.Fail()
	}
}

// Property-based test for creating a bunch of rows and making sure the data get put in
func TestRainEntry(t *testing.T) {
	var maxCount int
	if testing.Short() {
		maxCount = 1
	} else {
		maxCount = 5
	}

	db := connectorFixture()
	var total int
	test := func(reps uint8) bool {
		count := int(reps)
		for i := 0; i < count; i++ {
			if res, err := db.MakeRainEntry(); err != nil || res == nil {
				logrus.Error(err)
				return false
			}
		}
		var val int
		if val = db.GetRainEntries(); val == -1 {
			logrus.Error("gave -1")
			return false
		}
		logrus.Debugf("val=%d, count=%d, total=%d", val, count, total)
		total += count
		return val == total
	}
	if err := quick.Check(test, &quick.Config{
		MaxCount: maxCount,
	}); err != nil {
		t.Error(err)
	}
}

// test concurrency
func TestConcurrentEntries(t *testing.T) {
	db := connectorFixture()
	expected := 5
	timeout := 5
	total := make(chan int)
	var mu sync.Mutex
	tally := 0

	// loop <count> times
	for i := 0; i < expected; i++ {
		go func() {
			_, err := db.MakeRainEntry()
			if err != nil {
				t.Error(err)
			}
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
	}

	// wait for them to finish
	var collected bool
	for i := timeout; i != 0; i-- {
		finished := <-total
		logrus.Info(finished)
		if finished == expected {
			collected = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if !collected {
		logrus.Error("all loops not finished")
		t.Fail()
	}
	actual := db.GetRainEntries()
	if actual != expected {
		logrus.Errorf("actual=%d, expected=%d", actual, expected)
		t.Fail()
	}
}

// Tests all the various entries work (except temperature)
func TestStaticSQLEntries(t *testing.T) {
	db := connectorFixture()

	count := 5
	timeout := 30

	var rain, soft, hard, pause, unpause int
	var err error
	check := func(e error) {
		if e != nil {
			t.Error(err)
		}
	}
	total := make(chan int)
	tally := 0
	var mu sync.Mutex

	// TODO: Refactor with cleaner function and waitgroups
	for i := 0; i < count; i++ {
		go func() {
			_, err = db.MakeRainEntry()
			check(err)
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
		go func() {
			_, err = db.MakeSoftResetEntry()
			check(err)
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
		go func() {
			_, err = db.MakeHardResetEntry()
			check(err)
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
		go func() {
			_, err = db.MakePauseEntry()
			check(err)
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
		go func() {
			_, err = db.MakeUnpauseEntry()
			check(err)
			mu.Lock()
			tally++
			total <- tally
			mu.Unlock()
		}()
	}
	// wait for entries to finish
	done := false
	for i := timeout; i != 0; i-- {
		finished := <-total
		logrus.Info(finished)
		if finished == count*5 {
			done = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if !done {
		t.Fail()
	}

	rain = db.GetRainEntries()
	soft = db.GetSoftResetEntries()
	hard = db.GetHardResetEntries()
	pause = db.GetPauseEntries()
	unpause = db.GetUnpauseEntires()

	if rain != count || soft != count || hard != count || pause != count || unpause != count {
		t.Fail()
	}
}
