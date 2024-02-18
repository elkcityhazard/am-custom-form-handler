package main

import (
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
)

func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/":
			next.ServeHTTP(w, r)
			return
		case "/blog/":
			next.ServeHTTP(w, r)
			return
		default:
			if strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
				return
			}
		}

		next.ServeHTTP(w, r)

	})
}

func NoSurf(next http.Handler) http.Handler {
	// new csrf handler
	csrfHandler := nosurf.New(next)
	// set the base cookie

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	// don't forget to use middleware in routes()
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return app.SessionManager.LoadAndSave(next)
}
