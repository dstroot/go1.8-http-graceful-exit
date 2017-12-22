/*
Package metrics implements a library to expose Prometheus metrics.
It is best practice to provide an endpoint for instrumentation tools
(like prometheus). It is implemented as Negroni middleware.  Use this
middleware like the negroni.Loggermiddleware (after negroni.Recovery,
before every other middleware).

You also need to implement a corresponding route to expose the metrics:

	import (
		"github.com/julienschmidt/httprouter"
		"github.com/prometheus/client_golang/prometheus"
	)

	// New creates a new router with our routes included
	func New() *httprouter.Router {

		r := httprouter.New()

		// Prometheus metrics
		r.Handler("GET", "/metrics", prometheus.Handler())

		return r
	}

*/
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
)

const (
	reqsName    = "requests_total"
	latencyName = "request_duration_milliseconds"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewMetrics returns a new instance of prometheus middleware for Negroni.
func NewMetrics(name string, buckets ...float64) *Middleware {
	var m Middleware

	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)

	return &m
}

// ServeHTTP method captures the metrics are interested in...
func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)
	res := rw.(negroni.ResponseWriter)

	code := http.StatusText(res.Status())
	m.reqs.WithLabelValues(code, r.Method, r.URL.Path).Inc()
	m.latency.WithLabelValues(code, r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}
