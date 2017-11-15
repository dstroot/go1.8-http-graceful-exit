package router

import (
	"expvar"
	"net/http"

	handle "github.com/dstroot/go1.8-http-graceful-exit/pkg/handlers"
	"github.com/julienschmidt/httprouter"
)

// New creates a new router with our routes
func New() *httprouter.Router {

	r := httprouter.New()

	// Routes
	r.GET("/", handle.Index)
	r.GET("/page", handle.Page)
	r.GET("/hello/:name", handle.Hello)

	// handler for serving expvar data
	r.Handler("GET", "/debug/vars", expvar.Handler())

	// handler for serving static files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	return r
}
