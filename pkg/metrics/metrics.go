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

// have a look at this package as well:
// https://github.com/thoas/stats/blob/master/stats.go
// It seems to take ideas from negroni as well as adapt to
// other middleware.

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
	reqsHelp    = "HTTP requests processed, partitioned by status code, method and HTTP path."
	latencyName = "request_duration_milliseconds"
	latencyHelp = "How long it took to process the request, partitioned by status code, method and HTTP path."
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewMetrics returns a new instance of prometheus middleware for Negroni.
func NewMetrics(host string, service string, buckets ...float64) *Middleware {
	var m Middleware

	// requests
	m.reqs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        reqsName,
		Help:        reqsHelp,
		ConstLabels: prometheus.Labels{"host": host, "service": service},
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	// latency
	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        latencyHelp,
		ConstLabels: prometheus.Labels{"host": host, "service": service},
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
	res := negroni.NewResponseWriter(rw)

	// TODO I cant seem to get the actual status - it seems the response hasn't
	// been actually written yet when we gte the code. Somehow we have to wait
	// until the respnse is done. How does logger do it?
	code := res.Status()
	if code == 0 {
		code = 200
	}
	status := http.StatusText(code)

	// capture metrics
	m.reqs.WithLabelValues(status, r.Method, r.URL.Path).Inc()
	m.latency.WithLabelValues(status, r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)

	// defer func(start time.Time) { // Make sure we record a status.
	// 	duration := time.Since(start)
	// 	status := http.StatusText(res.Status())
	//
	// 	// capture metrics
	// 	m.reqs.WithLabelValues(status, r.Method, r.URL.Path).Inc()
	// 	m.latency.WithLabelValues(status, r.Method, r.URL.Path).Observe(duration.Seconds())
	// }(time.Now())
}
