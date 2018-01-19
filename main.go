// This is an example of a simple webserver with a graceful
// shutdown. We are using httprouter for routing and negroni
// for middleware.
package main

import (
	"log"
	"net/http"

	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/dstroot/simple-go-webserver/pkg/metrics"
	"github.com/dstroot/simple-go-webserver/pkg/router"
	"github.com/dstroot/simple-go-webserver/pkg/tracing"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	stats "github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/go-kit"
	"github.com/uber/jaeger-lib/metrics/go-kit/expvar"
	"github.com/urfave/negroni"
)

var (
	metricsFactory stats.Factory
)

func main() {
	// create an HTTP router (a mux)
	r := router.New()

	// negroni middleware stack
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(metrics.NewMetrics(info.Report.HostName, info.Report.Program))
	n.Use(negroni.NewLogger())
	n.UseHandler(r) // pass mux to negroni

	// create a tracer
	metricsFactory = xkit.Wrap("", expvar.NewFactory(10)) // 10 buckets for histograms
	tracer, closer := tracing.Init(
		info.Report.Program,
		metricsFactory.Namespace(info.Report.Program, nil),
	)
	defer closer.Close()

	// instrument the router for tracing
	mw := nethttp.Middleware(
		tracer,
		n, // pass in negroni
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		}),
	)

	// run our server
	s := NewServer(info.Report.Port, mw) // pass port and mux
	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
