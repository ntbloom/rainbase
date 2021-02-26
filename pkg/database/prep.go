package database

// Prep a database.  This is essentially a test fixture but also to be called at the start of a new deployment
import (
	"database/sql"

	"github.com/sirupsen/logrus"

	_ "modernc.org/sqlite"
)

const (
	SqliteProd = "/etc/rainbase/gateway-prod.db"
	SqliteDev  = "/etc/rainbase/gateway-dev.db"
)

// NewSqliteFile initializes new file
func NewSqliteFile(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", name)
	if err != nil {
		logrus.Errorf("problem opening database: %s", err)
		return nil, err
	}
	return db, nil
}
