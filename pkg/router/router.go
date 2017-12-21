package router

import (
	"expvar"
	"log"
	"net/http"
	"sync/atomic"

	handle "github.com/dstroot/simple-go-webserver/pkg/handlers"
	"github.com/dstroot/simple-go-webserver/pkg/health"
	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
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

	// Prometheus metrics
	r.Handler("GET", "/metrics", prometheus.Handler())

	// health and readiness
	r.Handler("GET", "/healthz", health.HandlerFunc())

	isReady := &atomic.Value{}
	isReady.Store(true)
	r.Handler("GET", "/readyz", health.ReadyFunc(isReady))

	// handler for serving static files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	// handle 404's gracefully
	r.NotFound = http.HandlerFunc(handle.NotFound)

	return r
}
