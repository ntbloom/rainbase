package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"context"
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const permissions = 0666

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

// exec runs sql command on the database without returning rows
func (d *DBConnector) exec(cmd string) (sql.Result, error) {
	// get variables ready
	var (
		db   *sql.DB
		conn *sql.Conn
	)
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
		if err := db.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	// connect to the database
	db, err := sql.Open("sqlite", d.fullPath)
	if err != nil {
		logrus.Error("unable to open database")
		return nil, err
	}
	conn, err = db.Conn(d.ctx)
	if err != nil {
		logrus.Error("unable to get a connection struct")
		return nil, err
	}
	return conn.ExecContext(d.ctx, cmd)
}

// makeSchema puts the schema in the sqlite file
func (d *DBConnector) makeSchema() (sql.Result, error) {
	return d.exec(sqlschema)
}
