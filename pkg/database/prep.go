package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const permissions = 0666

type DBConnector struct {
	Db        *sql.DB  // pointer to sqlite file
	File      *os.File // pointer to actual file
	BaseName  string   // name of file without extension or directory
	FullPath  string   // full POSIX path of sqlite file
	BackupDir string   // full POSIX path of backup directory
}

// NewDBConnector makes a new databaseconnector struct
func NewDBConnector(fullPath, backupDir string, clobber bool) (*DBConnector, error) {
	if clobber {
		err := os.Remove(fullPath)
		if err != nil {
			// ignore the error but log it
			logrus.Debugf("%s doesn't exist; ignoring", fullPath)
		}
	}
	// make the file
	file, err := os.Create(fullPath)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	baseName := "ADD/REGEX/HERE" // TODO: add me!

	db, err := sql.Open("sqlite", fullPath)
	if err != nil {
		logrus.Errorf("problem opening database: %s", err)
		return nil, err
	}

	// make the backup directory
	err = os.MkdirAll(backupDir, permissions)
	if err != nil {
		return nil, err
	}

	return &DBConnector{
		Db:        db,
		File:      file,
		BaseName:  baseName,
		FullPath:  fullPath,
		BackupDir: backupDir,
	}, nil
}

// Backup creates a backup file in d.BackupDir
func (d *DBConnector) Backup() error {
	//now := time.Now()
	//timestamp := fmt.Sprintf("%s-%d", d.FullPath, now.Unix())
	return nil
}
