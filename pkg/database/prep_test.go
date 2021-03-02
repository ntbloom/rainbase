package database_test

import (
	"os"
	"testing"

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

// connector fixture for DBConnector
func connector(t *testing.T) *database.DBConnector {
	getConfig(t)
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	backupDir := viper.GetString(configkey.DatabaseLocalDevBackupDir)
	db, _ := database.NewDBConnector(sqliteFile, backupDir, false)
	return db
}

// TestDBPrep create and destroy sqlite file 5 times
func TestDatabasePrep(t *testing.T) {
	getConfig(t)
	sqliteFile := viper.GetString(configkey.DatabaseLocalDevFile)
	backupDir := viper.GetString(configkey.DatabaseLocalDevBackupDir)

	// clean up when finished
	defer removeFile(sqliteFile)
	defer removeDir(backupDir)

	// create and destroy 5 times
	for i := 0; i < 5; i++ {
		db, err := database.NewDBConnector(sqliteFile, backupDir, true)
		if err != nil || db == nil {
			t.Fail()
		}
		_, err = os.Stat(sqliteFile)
		if err != nil {
			t.Fail()
		}
	}
}

// TestBackups can we create a backup for the file?
func TestDatabaseBackups(t *testing.T) {
	getConfig(t)

}
