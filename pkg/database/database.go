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
	res, err := db.enter(illegal)
	return res == nil && err != nil
}

// MakeRainEntry record a rain event
func (db *DBConnector) MakeRainEntry() (sql.Result, error) {
	return db.record(tlv.Rain, tlv.RainValue)
}

// GetRainEntries get total rain events
func (db *DBConnector) GetRainEntries() int {
	// TODO: implement me!
	return -1
}

// enter runs sql command on the database without returning rows
func (db *DBConnector) enter(cmd string) (sql.Result, error) {
	// get variables ready
	var (
		database *sql.DB
		conn     *sql.Conn
	)
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
		if err := database.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	// connect to the database
	database, err := sql.Open("sqlite", db.fullPath)
	if err != nil {
		logrus.Error("unable to open database")
		return nil, err
	}
	conn, err = database.Conn(db.ctx)
	if err != nil {
		logrus.Error("unable to get a connection struct")
		return nil, err
	}

	// enforce foreign keys
	safeCmd := strings.Join([]string{foreignKey, cmd}, " ")

	return conn.ExecContext(db.ctx, safeCmd)
}

// tally runs sql command to tally database entries for a given topic; essentially a dummy function for testing
func (db *DBConnector) tally(tag int) int {
	return -1
}

// makeSchema puts the schema in the sqlite file
func (db *DBConnector) makeSchema() (sql.Result, error) {
	return db.enter(sqlschema)
}

// record makes an entry into the databse
func (db *DBConnector) record(tag int, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%db, %db, %s);", tag, value, timestamp)
	return db.enter(cmd)
}
