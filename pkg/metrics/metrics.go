/*
Package metrics implements a library to expose Prometheus metrics.
It is best practice to provide an endpoint for instrumentation tools
(like prometheus). It is implemented as Negroni middleware.  Use this
middleware like the negroni.Loggermiddleware (after negroni.Recovery,
before every other middleware).

You also need to implement a corresponding route to expose the metrics:

	import (
		"github.com/julienschmidt/httprouter"
		"github.com/prometheus/client_golang/prometheus/promhttp"
	)

	// New creates a new router with our routes included
	func New() *httprouter.Router {

		r := httprouter.New()

		// Prometheus metrics
		r.Handler("GET", "/metrics", promhttp.Handler())

		return r
	}

*/
package metrics

// have a look at this package as well:
// https://github.com/thoas/stats/blob/master/stats.go
// It seems to take ideas from negroni as well as adapt to
// other middleware.

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

var (
	dflBuckets = []float64{.25, .5, 1, 2.5, 5, 10}
)

const (
	reqsName    = "requests_total"
	reqsHelp    = "HTTP requests processed, partitioned by status code, method and HTTP path."
	latencyName = "request_duration_milliseconds"
	latencyHelp = "How long it took to process the request, partitioned by status code, method and HTTP path."
)

// Metrics holds our prometheus metrics buckets
type Metrics struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
	size    *prometheus.HistogramVec
}

// NewMetrics returns a new instance of prometheus middleware for Negroni.
func NewMetrics(host string, service string, buckets ...float64) *Metrics {
	var m Metrics

	// requests
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
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
	m.latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        latencyName,
			Help:        latencyHelp,
			ConstLabels: prometheus.Labels{"host": host, "service": service},
			Buckets:     buckets,
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)

	// responseSize has no labels, making it a zero-dimensional ObserverVec.
	m.size = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "response_size_bytes",
			Help:        "A histogram of response sizes for requests.",
			ConstLabels: prometheus.Labels{"host": host, "service": service},
			Buckets:     []float64{1000, 5000, 10000, 500000},
		},
		[]string{},
	)
	prometheus.MustRegister(m.size)

	return &m
}

// Negroni middleware
func (m *Metrics) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	status := http.StatusText(res.Status())

	// captures metrics
	m.reqs.WithLabelValues(status, r.Method, r.URL.Path).Inc()
	m.latency.WithLabelValues(status, r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	m.size.WithLabelValues().Observe(float64(res.Size()))

	fmt.Println(res.Size())
}
