/*
Package router implements a library to manage our application's routes.
*/
package router

import (
	"net/http"
	"sync/atomic"

	handle "github.com/dstroot/simple-go-webserver/pkg/handlers"
	"github.com/dstroot/simple-go-webserver/pkg/health"
	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// New creates a new router with our routes
func New() *httprouter.Router {

	r := httprouter.New()

	// // Instrument the handlers with all the metrics, injecting the "handler"
	// // label by currying.
	// indexChain := promhttp.InstrumentHandlerDuration(m.duration.MustCurryWith(prometheus.Labels{"handler": "index"}),
	// 	promhttp.InstrumentHandlerCounter(m.counter,
	// 		promhttp.InstrumentHandlerResponseSize(m.size, handle.Index),
	// 	),
	// )

	// application routes
	r.GET("/", handle.Index)
	r.GET("/page", handle.Page)
	r.GET("/hello/:name", handle.Hello)

	// handler for serving info
	r.Handler("GET", "/info", info.HandlerFunc())

	// Prometheus metrics
	r.Handler("GET", "/metrics", promhttp.Handler())

	// healthz (for Kubernetes)
	r.Handler("GET", "/healthz", health.HandlerFunc())

	// readyz (for Kubernetes).
	// For the readiness probe we might need to wait for some event
	// (e.g. the database is ready) to be able to serve traffic. We
	// return 200 only if the variable "isReady" is true.
	// Here we simply set isReady to true.
	isReady := &atomic.Value{}
	isReady.Store(true)
	r.Handler("GET", "/readyz", health.ReadyFunc(isReady))

	// handler for serving static files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	// handle 404's gracefully
	r.NotFound = http.HandlerFunc(handle.NotFound)

	return r
}
