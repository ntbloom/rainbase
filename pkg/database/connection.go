package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4" // driver for postgresql
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite" // driver for sqlite
)

// Connector interface that abstracts the type of database implemented, in this case postgresql vs. sqlite
type Connector interface {
	connect() (*connection, error)
}

// Sqlite struct, simple connector for sqlite files
type SqliteConnector struct {
	dbFullPath string // full POSIX path of database file
}

// connection gets a DB and Conn struct for a sqlite file
type connection struct {
	database *sql.DB
	conn     *sql.Conn
}

// connect makes a connection with sqlite file
func (s *SqliteConnector) connect() (*connection, error) {
	var (
		database *sql.DB
		conn     *sql.Conn
		err      error
	)
	database, err = sql.Open("sqlite", s.dbFullPath)
	if err != nil {
		logrus.Error("unable to open database")
		return nil, err
	}

	// make a Conn
	conn, err = database.Conn(context.Background())
	if err != nil {
		logrus.Error("unable to get a connection struct")
		return nil, err
	}
	return &connection{database, conn}, nil
}

// disconnect breaks connection with the database
func (db *DBConnector) disconnect(c *connection) {
	if err := c.conn.Close(); err != nil {
		logrus.Error(err)
	}
	if err := c.database.Close(); err != nil {
		logrus.Error(err)
	}
}
