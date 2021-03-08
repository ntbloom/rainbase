package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite" // driver for sqlite
)

const foreignKey = `PRAGMA foreign_keys = ON;`

type DBConnector struct {
	file       *os.File        // pointer to actual file
	fullPath   string          // full POSIX path of sqlite file
	ctx        context.Context // background context
	//sync.Mutex                 // access the database serially
}

// NewDBConnector makes a new databaseconnector struct
func NewDBConnector(fullPath string, clobber bool) (*DBConnector, error) {
	logrus.Debug("making new DBConnector")
	if clobber {
		_ = os.Remove(fullPath)
	}

	// connect to the file and open it
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}

	// make a DBConnector object and make the schema if necessary
	db := DBConnector{
		file:     file,
		fullPath: fullPath,
		ctx:      context.Background(),
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
func (db *DBConnector) MakeRainEntry() (sql.Result, error) {
	return db.addRecord(tlv.Rain, tlv.RainValue)
}

func (db *DBConnector) GetRainEntries() int {
	return db.tally(tlv.Rain)
}

// MakeSoftResetEntry addRecord a soft reset event
func (db *DBConnector) MakeSoftResetEntry() (sql.Result, error) {
	return db.addRecord(tlv.SoftReset, tlv.SoftResetValue)
}

func (db *DBConnector) GetSoftResetEntries() int {
	return db.tally(tlv.SoftReset)
}

// MakeHardResetEntry addRecord a hard reset event
func (db *DBConnector) MakeHardResetEntry() (sql.Result, error) {
	return db.addRecord(tlv.HardReset, tlv.HardResetValue)
}

func (db *DBConnector) GetHardResetEntries() int {
	return db.tally(tlv.HardReset)
}

// MakePauseEntry addRecord a pause event
func (db *DBConnector) MakePauseEntry() (sql.Result, error) {
	return db.addRecord(tlv.Pause, tlv.Unpause)
}

func (db *DBConnector) GetPauseEntries() int {
	return db.tally(tlv.Pause)
}

// MakeUnpauseEntry addRecord an unpause event
func (db *DBConnector) MakeUnpauseEntry() (sql.Result, error) {
	return db.addRecord(tlv.Unpause, tlv.UnpauseValue)
}

func (db *DBConnector) GetUnpauseEntires() int {
	return db.tally(tlv.Unpause)
}

// MakeTemperatureEntry addRecord a temperature measurement
func (db *DBConnector) MakeTemperatureEntry(tempC int) (sql.Result, error) {
	return db.addRecord(tlv.Temperature, tempC)
}

/* SELECTED METHODS EXPORTED FOR TEST/VERIFICATION */

// ForeignKeysAreImplemented, test function to ensure foreign key implementation
func (db *DBConnector) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.enterData(illegal)
	return res == nil && err != nil
}

/* HELPER METHODS */

// makeSchema puts the schema in the sqlite file
func (db *DBConnector) makeSchema() (sql.Result, error) {
	return db.enterData(sqlschema)
}

// enterData enters data into the database without returning any rows
func (db *DBConnector) enterData(cmd string) (sql.Result, error) {
	var c *connection
	var err error

	// enforce foreign keys
	safeCmd := strings.Join([]string{foreignKey, cmd}, " ")

	// Disabled for now but determine if mutex is necessary or if SQLITE handles concurrent writes properly
	// db.Lock()
	// defer db.Unlock()
	if c, err = db.newConnection(); err != nil {
		return nil, err
	}
	defer c.disconnect()

	return c.conn.ExecContext(db.ctx, safeCmd)
}

// tally runs sql command to tally database entries for a given topic; essentially a dummy function for testing
func (db *DBConnector) tally(tag int) int {
	var rows *sql.Rows
	var err error

	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)

	// Disabled for now but determine if mutex is necessary or if SQLITE handles concurrent writes properly
	// db.Lock()
	// defer db.Unlock()
	c, _ := db.newConnection() // don't handle the error, just return -1
	defer c.disconnect()

	if rows, err = c.conn.QueryContext(db.ctx, query); err != nil {
		return -1
	}
	closed := func() {
		if err = rows.Close(); err != nil {
			logrus.Error(err)
		}
	}
	defer closed()
	results := make([]int, 0)
	for rows.Next() {
		var val int
		if err = rows.Scan(&val); err != nil {
			logrus.Error(err)
			return -1
		}
		results = append(results, val)
	}

	return results[0]
}

// addRecord makes an entry into the databse
func (db *DBConnector) addRecord(tag, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%d, %d, \"%s\");", tag, value, timestamp)
	return db.enterData(cmd)
}
