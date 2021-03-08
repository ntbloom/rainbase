package database_test

import (
	"database/sql"
	"os"
	"sync"
	"testing"
	"testing/quick"

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
	maxCount := 5
	if testing.Short() {
		logrus.Info("skipping property tests")
		return
	} else {
		logrus.Info("doing full property test")
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

// Tests all the various entries work (except temperature). Also tests concurrent use of database
func TestStaticSQLEntries(t *testing.T) {
	db := connectorFixture()
	count := 5

	// asynchronously make an entry for each type
	var wg sync.WaitGroup
	wg.Add(5 * count)
	type addFunction func() (sql.Result, error)
	checkAdd := func(callable addFunction) {
		defer wg.Done()
		_, err := callable()
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < count; i++ {
		go checkAdd(db.MakeRainEntry)
		go checkAdd(db.MakeSoftResetEntry)
		go checkAdd(db.MakeHardResetEntry)
		go checkAdd(db.MakePauseEntry)
		go checkAdd(db.MakeUnpauseEntry)
	}
	// wait for entries to finish
	wg.Wait()

	// verify counts
	wg.Add(5)
	type getFunction func() int
	checkGet := func(callable getFunction) {
		defer wg.Done()
		tally := callable()
		if tally != count {
			t.Fail()
		}
	}
	go checkGet(db.GetRainEntries)
	go checkGet(db.GetSoftResetEntries)
	go checkGet(db.GetHardResetEntries)
	go checkGet(db.GetPauseEntries)
	go checkGet(db.GetUnpauseEntries)
	wg.Wait()
}

// tests that we can enter temperature
func TestTemperatureEntries(t *testing.T) {
	db := connectorFixture()
	vals := []int{-100, -25, -15, -1, 0, 1, 2, 20, 24, 100}
	for _, expected := range vals {
		_, err := db.MakeTemperatureEntry(expected)
		if err != nil {
			logrus.Error(err)
			t.Error(err)
		}
		if actual := db.GetLastTemperatureEntry(); expected != actual {
			logrus.Errorf("expected=%d, actual=%d", expected, actual)
			t.Fail()
		}
	}
}
