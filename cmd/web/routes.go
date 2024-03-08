package main

import (
	"io/fs"
	"net/http"

	"github.com/elkcityhazard/am-form/internal/handlers"
	"github.com/elkcityhazard/am-form/static"
)

func routes() http.Handler {

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

	r.HandleFunc("/", handlers.HandleDisplayAMForm)

	r.HandleFunc("/success", handlers.HandleDisplaySuccess)

	r.HandleFunc("/panic", panicHandler)

	return PanicRecovery(StripTrailingSlash(NoSurf(SessionLoad(r))))

}

func panicHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusInternalServerError)
	defer func() {
		panic("Internal Server Error")
	}()
}
