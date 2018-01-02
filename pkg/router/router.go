/*
Package router implements a library to manage our application's routes.
*/
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
	// "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// type metrics struct {
// 	counter  *prometheus.CounterVec
// 	duration *prometheus.HistogramVec
// 	size     *prometheus.HistogramVec
// }
//
// func initMetrics() *metrics {
// 	var m metrics
// 	m.counter = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "api_requests_total",
// 			Help: "A counter for requests to the wrapped handler.",
// 		},
// 		[]string{"code", "method"},
// 	)
// 	// duration is partitioned by the HTTP method and handler. It uses custom
// 	// buckets based on the expected request duration.
// 	m.duration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "request_duration_seconds",
// 			Help:    "A histogram of latencies for requests.",
// 			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
// 		},
// 		[]string{"handler", "method"},
// 	)
//
// 	// responseSize has no labels, making it a zero-dimensional
// 	// ObserverVec.
// 	m.size = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "response_size_bytes",
// 			Help:    "A histogram of response sizes for requests.",
// 			Buckets: []float64{200, 500, 900, 1500},
// 		},
// 		[]string{},
// 	)
//
// 	// Register all of the metrics in the standard registry.
// 	prometheus.MustRegister(m.counter, m.duration, m.size)
//
// 	return &m
// }

// New creates a new router with our routes
func New() *httprouter.Router {

	r := httprouter.New()
	// m := initMetrics()

	err := info.Init()
	if err != nil {
		log.Fatalf("info could not be initialized")
	}

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

	// handler for serving expvar data
	r.Handler("GET", "/debug", expvar.Handler())

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
