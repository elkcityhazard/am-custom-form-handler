package config

import (
	"context"
	"database/sql"
	"html/template"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/elkcityhazard/am-form/internal/mailer"
)

type AppConfig struct {
	Port           string
	DSN            string
	IsProduction   bool
	Context        context.Context
	DB             *sql.DB
	WG             *sync.WaitGroup
	Mutex          *sync.Mutex
	TemplateCache  map[string]*template.Template
	SessionManager *scs.SessionManager
	Mailer         *mailer.Mailer
}
