package main

import (
	"os"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/elkcityhazard/am-form/internal/config"
)

//	TestMain will always be executed before the tests run

func TestMain(m *testing.M) {

	app.SessionManager = scs.New()

	os.Exit(m.Run())

}

func ReturnAppConfig() config.AppConfig {
	return app
}
