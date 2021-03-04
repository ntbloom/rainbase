package database_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/ntbloom/rainbase/pkg/config/configkey"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/ntbloom/rainbase/pkg/database"
)

// reusable configs
func getConfig(t *testing.T) {
	config.Configure()
}

// helper function removes file
func removeFile(file string) {
	_ = os.Remove(file)
}

// helper function removes backup directory
func removeDir(dir string) {
	_ = os.RemoveAll(dir)
}

// connectorFixture makes a reusable DBConnector object
func connectorFixture(t *testing.T) *database.DBConnector {
	getConfig(t)
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	db, _ := database.NewDBConnector(sqliteFile, true)
	return db
}

// TestDBPrep create and destroy sqlite file 5 times
func TestDatabasePrep(t *testing.T) {
	getConfig(t)
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)

	// clean up when finished
	defer removeFile(sqliteFile)

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
