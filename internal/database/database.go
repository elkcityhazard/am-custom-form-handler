package database

import (
	"database/sql"
	"log"

	"github.com/elkcityhazard/am-form/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

var app *config.AppConfig

var db *sql.DB

func NewDatabase(a *config.AppConfig) {
	conn, err := sql.Open("mysql", a.DSN)

	if err != nil {
		log.Fatal(err)
	}

	a.DB = conn

	db = conn
	app = a
}

func DatabaseConnection() *sql.DB {
	return db
}

func DatabaseHealthCheck() error {
	return app.DB.Ping()
}
