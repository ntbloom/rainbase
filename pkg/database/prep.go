package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const permissions = 0666

type DBConnector struct {
	db        *sql.DB
	Name      string
	BackupDir string
}

// NewDBConnector makes a new databaseconnector struct
func NewDBConnector(name, backupDir string, clobber bool) (*DBConnector, error) {
	// delete and start fresh if necessary
	if clobber {
		// remove the file and make it clean again
		err := os.Remove(name)
		if err != nil {
			// ignore the error but log it
			logrus.Debugf("%s doesn't exist; ignoring", name)
		}
		err = ioutil.WriteFile(name, nil, permissions)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	// make the backup directory
	err := os.MkdirAll(backupDir, permissions)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", name)
	if err != nil {
		logrus.Errorf("problem opening database: %s", err)
		return nil, err
	}

	return &DBConnector{
		db:        db,
		Name:      name,
		BackupDir: backupDir,
	}, nil
}

// Backup creates a backup file in d.BackupDir
func (d *DBConnector) Backup() error {
	return nil
}
