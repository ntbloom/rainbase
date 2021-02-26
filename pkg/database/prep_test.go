package database_test

import (
	"testing"

	"github.com/ntbloom/rainbase/pkg/database"
	"github.com/sirupsen/logrus"
)

func TestDBPrep(t *testing.T) {
	_, err := database.NewSqliteFile(database.SqliteDev)
	if err != nil {
		logrus.Errorf("failing on NewSqliteFile")
		t.Fail()
	}
}
