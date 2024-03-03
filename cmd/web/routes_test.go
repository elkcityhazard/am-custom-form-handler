package main

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elkcityhazard/am-form/static"
)

// Assuming static, app, handlers, StripTrailingSlash, NoSurf, and SessionLoad are defined elsewhere in your package

// TestRoutes checks that the routes function recovers from panics and registers handlers correctly.
func Test_Routes(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(routes())
	defer ts.Close()

	// Test Panic on go:embed files

	t.Run("PanicOnGoEmbedFiles", func(t *testing.T) {

		var staticDir = static.GetStaticDir()

		staticDir = nil

		// static files

		_, err := fs.Sub(staticDir, "assets_test")

		if err != nil {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}()
			panic(err)

		}

	})

	// Test recovery from panic
	t.Run("RecoveryFromPanic", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/panic")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected to read response body without error, got %v", err)
		}

		t.Log(resp)

		// Expecting the server to recover and thus return a 500 Internal Server Error
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %v", resp.StatusCode)
		}

	})

	// Test app.IsProduction

	t.Run("IsProduction", func(t *testing.T) {
		app.IsProduction = true

		if !app.IsProduction {
			t.Error("Expected IsProduction to be true")
		}

		if app.IsProduction {

			ts := httptest.NewServer(routes())

			defer ts.Close()

			resp, err := http.Get(ts.URL + "/static/assets/scripts/index.js")
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code 200, got %v", resp.StatusCode)
			}

		}

	})

	// Test serving "/am-form" route
	t.Run("AMFormRoute", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/am-form")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %v", resp.StatusCode)
		}
	})

	// Test serving "/success" route
	t.Run("SuccessRoute", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/success")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %v", resp.StatusCode)
		}
	})

	// Here you could add more tests for other routes and static file serving
}
