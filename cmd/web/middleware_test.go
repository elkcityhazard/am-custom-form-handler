package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_StripTrailingSlash(t *testing.T) {

	// arrange

	tests := []struct {
		name        string
		handler     http.Handler
		method      string
		url         string
		expectedURL string
		wantedCode  int
	}{
		{
			name: "test home page",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.Write([]byte("/"))

			}),
			method:      http.MethodGet,
			url:         "/",
			expectedURL: "/",
			wantedCode:  http.StatusOK,
		},
		{
			name: "test blog page",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.Write([]byte("/blog/"))

			}),
			method:      http.MethodGet,
			url:         "/blog/",
			expectedURL: "/blog/",
			wantedCode:  http.StatusOK,
		},
		{
			name: "test am-form",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			}),
			method:      http.MethodGet,
			url:         "/am-form/",
			expectedURL: "/am-form",
			wantedCode:  http.StatusMovedPermanently,
		},
		{
			name: "test multiple trailing slashes",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.Write([]byte("/blog//"))
			}),
			method:      http.MethodGet,
			url:         "/blog//",
			expectedURL: "/blog/",
			wantedCode:  http.StatusMovedPermanently,
		},
		{
			name: "test blank space",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			}),
			method:      http.MethodGet,
			url:         "/blog/some-fancy-test-with-trailing-slash/  ",
			expectedURL: "/blog/some-fancy-test-with-trailing-slash",
			wantedCode:  http.StatusMovedPermanently,
		},
		{
			name: "test url with blank spaces",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				w.Write([]byte(" "))
			}),
			method:      http.MethodGet,
			url:         "/blog/some-fancy-test-with-trailing slash/  ",
			expectedURL: "/blog/some-fancy-test-with-trailing-slash",
			wantedCode:  http.StatusMovedPermanently,
		},
		{
			name: "serve http",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wasCalled := true

				if wasCalled {
					w.Write([]byte("serve http"))
				}
			}),
			method:      http.MethodGet,
			url:         "/blog/some-standard-url",
			expectedURL: "/blog/some-standard-url",
			wantedCode:  http.StatusOK,
		},
	}
	// act

	// assert

	for _, tt := range tests {

		// set up the test server

		srv := httptest.NewServer(StripTrailingSlash(tt.handler))

		// use appropriate middleware

		defer srv.Close()

		// run the test

		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(tt.method, srv.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := StripTrailingSlash(tt.handler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantedCode {
				t.Errorf("got %d want %d", rr.Code, tt.wantedCode)
			}

			if tt.name == "serve http" {
				if rr.Body.String() != "serve http" {
					t.Errorf("handler was not called")
				}
			}

			if rr.Code != 200 {

				location, err := rr.Result().Location()

				if err != nil {
					t.Errorf("got %s want %s", rr.Header().Get("Location"), tt.expectedURL)
				}

				if location.String() != tt.expectedURL {
					t.Errorf("got %s want %s", location.String(), tt.expectedURL)
				}

			}

		})
	}
}

func Test_NoSurf(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Test setting base cookie with HttpOnly, Path, and Secure attributes
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := NoSurf(mockHandler)
	handler.ServeHTTP(rr, req)
	// Add assertions for cookie attributes

	// Test using middleware in routes()
	// Add assertions for middleware usage in routes
}

func Test_SessionLoad(t *testing.T) {

	// mock app config

	// mock hander

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.SessionManager.Put(r.Context(), "foo", "bar")

		if !app.SessionManager.Exists(r.Context(), "foo") {
			t.Error("foo does not exist in session")
		}

		t.Log(app)

		w.WriteHeader(http.StatusOK)
	})

	// wrap mock handler in the SessionLoadAndSave Func

	handler := SessionLoad(mockHandler)

	ts := httptest.NewServer(handler)

	defer ts.Close()

	rr := httptest.NewRecorder()

	req := httptest.NewRequest("GET", ts.URL, nil)

	handler.ServeHTTP(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("got %d want %d", res.StatusCode, http.StatusOK)
	}

}

func Test_PanicRecovery(t *testing.T) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	handler := PanicRecovery(mockHandler)

	ts := httptest.NewServer(handler)

	defer ts.Close()

	rr := httptest.NewRecorder()

	req := httptest.NewRequest("GET", ts.URL, nil)

	handler.ServeHTTP(rr, req)

	res := rr.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("got %d want %d", res.StatusCode, http.StatusInternalServerError)
	}

}
