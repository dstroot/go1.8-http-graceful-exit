// This is an example of a simple webserver with a graceful
// shutdown. We are using httprouter for routing and negroni
// for middleware.
package main

import (
	"log"
	"net/http"

	// /debug/vars and /debug/pprof
	_ "expvar"
	_ "net/http/pprof"

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
	// Let's put the expvar and pprof http server on a separate port on
	// localhost, separate from the application http server. Both register
	// handlers on the default mux automatically:
	//  - http://localhost:6060/debug/vars
	//  - http://localhost:6060/debug/pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// initialize program info
	err := info.Init()
	if err != nil {
		log.Fatalf("info could not be initialized")
	}

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
	tracer, closer, err := tracing.Init(
		info.Report.Program,
		metricsFactory.Namespace(info.Report.Program, nil),
	)
	if err != nil {
		log.Fatal(err)
	}
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
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
