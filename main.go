// This is an example of a simple webserver with a graceful
// shutdown. We are using httprouter for routing and negroni
// for middleware.
package main

import (
	"log"

	"github.com/dstroot/simple-go-webserver/pkg/info"
	// "github.com/dstroot/simple-go-webserver/pkg/metrics"
	"github.com/dstroot/simple-go-webserver/pkg/router"
	"github.com/urfave/negroni"
)

func main() {
	// HTTP router (a mux)
	r := router.New()

	// Negroni Middleware Stack
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	// n.Use(metrics.NewMetrics(info.Report.HostName, info.Report.Program))
	n.Use(negroni.NewLogger())
	n.UseHandler(r) // pass mux to negroni

	s := NewServer(info.Report.Port, n) // pass port and mux
	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
