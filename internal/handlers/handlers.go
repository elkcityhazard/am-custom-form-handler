package handlers

import "github.com/elkcityhazard/am-form/internal/config"

var app *config.AppConfig

func NewHandlers(a *config.AppConfig) {
	app = a
}
