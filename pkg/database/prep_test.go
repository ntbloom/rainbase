package database_test

import (
	"os"
	"testing"

	"github.com/ntbloom/rainbase/pkg/database"
	"github.com/sirupsen/logrus"
)

func TestDBCreation(t *testing.T) {
	file := database.SqliteDev
	err := os.Remove(file)
	if err != nil {
		logrus.Errorf("can't remove %s: %s", file, err)
	}
	_, err = os.Stat(database.SqliteDev)
	if err != nil {
		t.Fail()
	}
}
