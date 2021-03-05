package database_test

import (
	"os"
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
func getConfig(t *testing.T) {
	config.Configure()
}

// connectorFixture makes a reusable DBConnector object
func connectorFixture(t *testing.T) *database.DBConnector {
	getConfig(t)
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	db, _ := database.NewDBConnector(sqliteFile, true)
	return db
}

/* TESTS */

// create and destroy sqlite file 5 times, get DBCOnnector struct
func TestDatabasePrep(t *testing.T) {
	getConfig(t)
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
	db := connectorFixture(t)
	if foreignKeys := db.ForeignKeysAreImplemented(); !foreignKeys {
		logrus.Error("sqlite is not enforcing foreign_key constraints")
		t.Fail()
	}
}

// make sure Entry.Record() interface is implemented correcly
func TestRainEntry(t *testing.T) {
	test := func(reps uint8) bool {
		db := connectorFixture(t)
		count := int(reps)
		logrus.Debugf("count=%d", count)
		var val int

		for i := 0; i < count; i++ {
			if res, err := db.MakeRainEntry(); err != nil || res == nil {
				logrus.Error(err)
				return false
			}
		}
		if val = db.GetRainEntries(); val == -1 {
			logrus.Error("gave -1")
			return false
		}
		logrus.Debugf("val=%d", val)
		return val == count
	}
	if err := quick.Check(test, &quick.Config{
		MaxCount: 15,
	}); err != nil {
		t.Error(err)
	}
}

//func TestRainEntry(t *testing.T) {
//	db := connectorFixture(t)
//	count := 10
//	for i := 0; i < count; i++ {
//		if res, err := db.MakeRainEntry(); err != nil || res == nil {
//			t.Error(err)
//		}
//	}
//	if val := db.GetRainEntries(); val == -1 {
//		t.Error("-1")
//	}
//}
