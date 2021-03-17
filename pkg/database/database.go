package database

import (
	"context"
	"os"

	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"
)

const sqlite = "sqlite"
const postgres = "postgresql"

type DBConnector struct {
	dbName string          // full POSIX path of sqlite file or name of the database in Postgresql
	driver string          // change the type of database connection
	ctx    context.Context // background context
}

// NewSqliteDBConnector makes a new connector struct for sqlite
func NewSqliteDBConnector(fullPath string, clobber bool) (*DBConnector, error) {
	logrus.Debug("making new DBConnector struct for Sqlite")
	if clobber {
		_ = os.Remove(fullPath)
	}

	// connect to the file and open it
	_, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	// make a DBConnector object and make the schema if necessary
	db := DBConnector{
		dbName: fullPath,
		driver: sqlite,
		ctx:    context.Background(),
	}
	if clobber {
		_, err = db.makeSchema()
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}

/* SQL LOG ENTRIES */

// MakeRainEntry addRecord a rain event
func (db *DBConnector) MakeRainEntry() {
	_, err := db.addRecord(tlv.Rain, tlv.RainValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeSoftResetEntry addRecord a soft reset event
func (db *DBConnector) MakeSoftResetEntry() {
	_, err := db.addRecord(tlv.SoftReset, tlv.SoftResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeHardResetEntry addRecord a hard reset event
func (db *DBConnector) MakeHardResetEntry() {
	_, err := db.addRecord(tlv.HardReset, tlv.HardResetValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakePauseEntry addRecord a pause event
func (db *DBConnector) MakePauseEntry() {
	_, err := db.addRecord(tlv.Pause, tlv.Unpause)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeUnpauseEntry addRecord an unpause event
func (db *DBConnector) MakeUnpauseEntry() {
	_, err := db.addRecord(tlv.Unpause, tlv.UnpauseValue)
	if err != nil {
		logrus.Error(err)
	}
}

// MakeTemperatureEntry addRecord a temperature measurement
func (db *DBConnector) MakeTemperatureEntry(tempC int) {
	_, err := db.addRecord(tlv.Temperature, tempC)
	if err != nil {
		logrus.Error(err)
	}
}

/* GETTERS, MOSTLY FOR TESTING */
func (db *DBConnector) GetRainEntries() int {
	return db.tally(tlv.Rain)
}

func (db *DBConnector) GetSoftResetEntries() int {
	return db.tally(tlv.SoftReset)
}

func (db *DBConnector) GetHardResetEntries() int {
	return db.tally(tlv.HardReset)
}

func (db *DBConnector) GetPauseEntries() int {
	return db.tally(tlv.Pause)
}

func (db *DBConnector) GetUnpauseEntries() int {
	return db.tally(tlv.Unpause)
}

// GetLastTemperatureEntry returns last temp reading, sorted by primary key
func (db *DBConnector) GetLastTemperatureEntry() int {
	return db.getLastRecord(tlv.Temperature)
}

/* SELECTED METHODS EXPORTED FOR TEST/VERIFICATION */

// ForeignKeysAreImplemented, test function to ensure foreign key implementation
func (db *DBConnector) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.enterData(illegal)
	return res == nil && err != nil
}
