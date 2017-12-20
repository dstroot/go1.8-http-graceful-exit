package router

import (
	"expvar"
	"log"
	"net/http"

	handle "github.com/dstroot/simple-go-webserver/pkg/handlers"
	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/julienschmidt/httprouter"
)

// New creates a new router with our routes
func New() *httprouter.Router {

	r := httprouter.New()

	err := info.Init()
	if err != nil {
		log.Fatalf("info could not be initialized")
	}

	// Routes
	r.GET("/", handle.Index)
	r.GET("/page", handle.Page)
	r.GET("/hello/:name", handle.Hello)

	// handler for serving expvar data
	r.Handler("GET", "/debug", expvar.Handler())

	// handler for serving info
	r.Handler("GET", "/info", info.HandlerFunc())

	// handler for serving static files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	// handle 404's gracefully
	r.NotFound = http.HandlerFunc(handle.NotFound)

	return r
}
