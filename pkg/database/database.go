package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const (
	//nolint
	schema = `CREATE TABLE packet 
id INT PRIMARY KEY SERIAL,
tag INT NOT NULL,
value INT NOT NULL,
timestamp TEXT /* created by go */
`
)
const permissions = 0666

type DBConnector struct {
	db       *sql.DB  // pointer to sqlite file
	file     *os.File // pointer to actual file
	fullPath string   // full POSIX path of sqlite file
}

// NewDBConnector makes a new databaseconnector struct
func NewDBConnector(fullPath string, clobber bool) (*DBConnector, error) {
	logrus.Debug("making new DBConnector")
	if clobber {
		_ = os.Remove(fullPath)
	}
	// make the file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", fullPath)
	if err != nil {
		return nil, err
	}

	return &DBConnector{
		db:       db,
		file:     file,
		fullPath: fullPath,
	}, nil
}


