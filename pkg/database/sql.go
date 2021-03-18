package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

/* All raw SQL lives here to keep an abstraction layer between DBConnector objects and the actual database access */

const sqliteForeignKeyPragma = `PRAGMA foreign_keys = ON;`

// enterData enters data into the database without returning any rows
func (db *DBConnector) enterData(cmd string) (sql.Result, error) {
	var c *connection
	var err error

	// manually add foreign key pragma for sqlite
	if db.driver == sqlite {
		cmd = strings.Join([]string{sqliteForeignKeyPragma, cmd}, " ")
	}
	if c, err = db.connector.connect(); err != nil {
		return nil, err
	}
	defer db.disconnect(c)

	return c.conn.ExecContext(context.Background(), cmd)
}

// makeSchema puts the schema in the sqlite file
func (db *DBConnector) makeSchema() (sql.Result, error) {
	return db.enterData(sqlschema)
}

// addRecord makes an entry into the databse
// base command for all logging
func (db *DBConnector) addRecord(tag, value int) (sql.Result, error) {
	timestamp := time.Now().Format(time.RFC3339)
	cmd := fmt.Sprintf("INSERT INTO log (tag, value, timestamp) VALUES (%d, %d, \"%s\");", tag, value, timestamp)
	return db.enterData(cmd)
}

/* QUERYING METHODS, MOSTLY FOR TESTING */

// tally runs sql command to count database entries for a given topic
func (db *DBConnector) tally(tag int) int {
	query := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE tag = %d;", tag)
	return db.getSingleInt(query)
}

// getLastRecord gets the last record for a given tag
func (db *DBConnector) getLastRecord(tag int) int {
	cmd := fmt.Sprintf(`SELECT value FROM log WHERE tag = %d ORDER BY id DESC LIMIT 1;`, tag)
	return db.getSingleInt(cmd)
}

// getSingleInt returns the first result of any SQL query that gives at least one integer result
// simple function for confirming correct value was entered for, say, temperature
func (db *DBConnector) getSingleInt(query string) int {
	var rows *sql.Rows
	var err error

	c, _ := db.connector.connect() // don't handle the error, just return -1
	defer db.disconnect(c)

	if rows, err = c.conn.QueryContext(context.Background(), query); err != nil {
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
