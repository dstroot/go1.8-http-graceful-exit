package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dstroot/go1.8-http-graceful-exit/pkg/handlers"
	"github.com/dstroot/go1.8-http-graceful-exit/pkg/tmpl"
)

func TestRoutes(t *testing.T) {

	handlers.Render = tmpl.New(
		tmpl.Options{
			TemplateDirectory: "../../templates",
		},
	)

	// test data
	var routes = []struct {
		method string
		route  string
		status int
	}{
		{"GET", "/", http.StatusOK},
		{"POST", "/", http.StatusMethodNotAllowed},
		{"GET", "/page", http.StatusOK},
		{"GET", "/hello/Dan", http.StatusOK},
		{"GET", "/debug/vars", http.StatusOK},
		{"GET", "/nonexistant", http.StatusNotFound},
	}

	// instantiate a router
	router := New()

	// test routes
	for _, r := range routes {
		//The response recorder used to record HTTP responses
		respRec := httptest.NewRecorder()

		// create request
		req, err := http.NewRequest(r.method, r.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		// serve request and capture response
		router.ServeHTTP(respRec, req)

		// check response
		if respRec.Code != r.status {
			t.Fatal("route ", r.route, " returned ", respRec.Code, " instead of ", r.status)
		}
	}
}
