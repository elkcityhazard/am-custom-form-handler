package database

import (
	"database/sql"
	"testing"

	"github.com/elkcityhazard/am-form/internal/config"
)

var tmpDB *sql.DB

func Test_NewDatabase(t *testing.T) {
	var mockApp = config.AppConfig{}

	NewDatabase(&mockApp)

	tmpDB = mockApp.DB

	if tmpDB == nil {
		t.Error("DB is nil")
	}

}

func Test_DatabaseConnection(t *testing.T) {

	var mockApp = config.AppConfig{}

	NewDatabase(&mockApp)

	tmpDB = DatabaseConnection()

	if tmpDB == nil {
		t.Error("DB is nil")
	}
}

func Test_DatabaseHealthCheck(t *testing.T) {

	var mockApp = config.AppConfig{}

	NewDatabase(&mockApp)
	err := DatabaseHealthCheck()

	if err == nil {
		t.Error("Error: ", err)
	}

}
