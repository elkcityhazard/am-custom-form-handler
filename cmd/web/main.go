package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/elkcityhazard/am-form/internal/config"
	"github.com/elkcityhazard/am-form/internal/database"
	"github.com/elkcityhazard/am-form/internal/handlers"
	"github.com/elkcityhazard/am-form/internal/mailer"
	"github.com/elkcityhazard/am-form/internal/render"
	"github.com/elkcityhazard/am-form/internal/repository/dbrepo"
	"github.com/go-mail/mail/v2"
)

var Port string

var app config.AppConfig

func main() {

	setAppConfig(&app)

	_ = dbrepo.NewDBRepo(&app)

	//	start mail configs

	go app.Mailer.ListenForMail(app.Mailer.ErrorChan, app.Mailer.MailerChan, app.Mailer.MailerDoneChan)

	//	end mail configs

	render.NewRenderer(&app)

	templateCache, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = templateCache

	database.NewDatabase(&app)

	err = database.DatabaseHealthCheck()

	if err != nil {
		log.Println("Error: ", err)
		log.Fatalln(err)
	}

	app.SessionManager.Store = mysqlstore.New(database.DatabaseConnection())

	handlers.NewHandlers(&app)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: routes(),
	}

	fmt.Println("Starting server on port 8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func parseFlags() {

	flag.StringVar(&app.DSN, "dsn", "", "Database connection string")
	flag.StringVar(&app.Port, "port", "8080", "Port to listen on")
	flag.StringVar(&app.Mailer.Host, "smtp_host", "", "SMTP Host")
	flag.IntVar(&app.Mailer.Port, "smtp_port", 0, "SMTP Port")
	flag.StringVar(&app.Mailer.Username, "smtp_user", "", "SMTP Username")
	flag.StringVar(&app.Mailer.Password, "smtp_user_pass", "", "SMTP Password")
	flag.StringVar(&app.Mailer.FromAddress, "smtp_from_address", "", "From Address")
	flag.StringVar(&app.Mailer.ToAddress, "smtp_to_address", "", "To Address")
	flag.Parse()

}

func setAppConfig(app *config.AppConfig) {
	gob.Register(mailer.EmailMessage{})
	app.IsProduction = false
	app.Context = context.Background()
	app.WG = new(sync.WaitGroup)
	app.Mutex = new(sync.Mutex)
	app.SessionManager = scs.New()

	newMailer := mailer.New()
	app.Mailer = newMailer
	parseFlags()
	app.Mailer.Dialer = mail.NewDialer(app.Mailer.Host, app.Mailer.Port, app.Mailer.Username, app.Mailer.Password)
	app.Mailer.Dialer.Timeout = 5 * time.Second

	app.SessionManager.Lifetime = 24 * time.Hour

	app.SessionManager.Cookie.Secure = true

	app.SessionManager.Cookie.SameSite = http.SameSiteLaxMode

	app.SessionManager.Cookie.Persist = true

	app.SessionManager.Cookie.SameSite = http.SameSiteLaxMode

	app.Mailer.MailerChan = make(chan *mailer.EmailMessage)
	app.Mailer.MailerDoneChan = make(chan bool)
	app.Mailer.ErrorChan = make(chan error)
}
