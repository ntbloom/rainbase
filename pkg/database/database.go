package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"context"
	"database/sql"
	"os"

	"github.com/ntbloom/rainbase/pkg/tlv"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const permissions = 0666
const foreignKey = `PRAGMA foreign_keys = ON;`

type DBConnector struct {
	file     *os.File        // pointer to actual file
	fullPath string          // full POSIX path of sqlite file
	ctx      context.Context // background context
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

// ForeignKeysAreImplemented, test function to ensure foreign key implementation
func (db *DBConnector) ForeignKeysAreImplemented() bool {
	illegal := `INSERT INTO log (tag, value, timestamp) VALUES (99999,1,"timestamp");`
	res, err := db.enterData(illegal)
	return res == nil && err != nil
}

// MakeRainEntry addRecord a rain event
func (db *DBConnector) MakeRainEntry() (sql.Result, error) {
	return db.addRecord(tlv.Rain, tlv.RainValue)
}

// GetRainEntries get total rain events
func (db *DBConnector) GetRainEntries() int {
	return db.tally(tlv.Rain)
}

// makeSchema puts the schema in the sqlite file
func (db *DBConnector) makeSchema() (sql.Result, error) {
	return db.enterData(sqlschema)
}
