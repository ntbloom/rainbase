package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

/* Wrap queries in methods so we don't expose the actual databse to the rest of the application */

// connection gets a DB and Conn struct for a sqlite file
type connection struct {
	database *sql.DB
	conn     *sql.Conn
}

func (db *DBConnector) newConnection() (*connection, error) {
	// get variables ready
	var (
		database *sql.DB
		conn     *sql.Conn
	)

	// get a DB struct
	database, err := sql.Open("sqlite", db.fullPath)
	if err != nil {
		logrus.Error("unable to open database")
		return nil, err
	}

	// make a Conn
	conn, err = database.Conn(db.ctx)
	if err != nil {
		logrus.Error("unable to get a connection struct")
		return nil, err
	}
	return &connection{database, conn}, nil
}

func (c *connection) disconnect() {
	if err := c.conn.Close(); err != nil {
		logrus.Error(err)
	}
	if err := c.database.Close(); err != nil {
		logrus.Error(err)
	}
}

// enterData enters data into the database without returning any rows
func (db *DBConnector) enterData(cmd string) (sql.Result, error) {
	var c *connection
	var err error

	if c, err = db.newConnection(); err != nil {
		return nil, err
	}
	defer c.disconnect()

	// enforce foreign keys
	safeCmd := strings.Join([]string{foreignKey, cmd}, " ")

	return c.conn.ExecContext(db.ctx, safeCmd)
}

// tally runs sql command to tally database entries for a given topic; essentially a dummy function for testing
func (db *DBConnector) tally(tag int) int {
	var rows *sql.Rows
	var err error

	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)
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
