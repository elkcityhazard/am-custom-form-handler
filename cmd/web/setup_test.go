package main

import (
	"os"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/elkcityhazard/am-form/internal/config"
	"github.com/elkcityhazard/am-form/internal/database"
)

//	TestMain will always be executed before the tests run

func TestMain(m *testing.M) {

	var mockApp = config.AppConfig{}

	database.NewDatabase(&mockApp)

	app.SessionManager = scs.New()

	os.Exit(m.Run())

}

func ReturnAppConfig() config.AppConfig {
	return app
}
