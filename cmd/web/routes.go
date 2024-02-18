package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/elkcityhazard/am-form/internal/handlers"
	"github.com/elkcityhazard/am-form/static"
)

func routes() http.Handler {

	defer func() {
		if r := recover(); r != nil {
			// r contains the value passed to panic()
			log.Printf("Recovering from panic: %v", r)
		}
	}()

	var staticDir = static.GetStaticDir()

	r := http.NewServeMux()

	// static files

	staticFiles, err := fs.Sub(staticDir, "assets")

	if err != nil {
		panic(err)
	}

	if app.IsProduction {

		fileServer := http.FileServer(http.FS(staticFiles))

		r.Handle("/static/assets/", http.StripPrefix("/static/assets/", fileServer))

	} else {

		r.Handle("/static/assets/", http.StripPrefix("/static/assets/", http.FileServer(http.Dir("./static/assets/"))))
	}

	r.HandleFunc("/am-form", handlers.HandleDisplayAMForm)

	r.HandleFunc("/success", handlers.HandleDisplaySuccess)

	return StripTrailingSlash(NoSurf(SessionLoad(r)))

}
