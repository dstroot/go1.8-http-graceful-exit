// This is an example of a simple webserver with a graceful
// shutdown. It is now possible to call srv.Close() to halt an
// http.Server immediately, or srv.Shutdown(ctx) to stop
// and gracefully drain the server of connections. We are using
// httprouter and negroni also.
package main

import (
	"github.com/dstroot/simple-go-webserver/pkg/info"
	"github.com/dstroot/simple-go-webserver/pkg/metrics"
	"github.com/dstroot/simple-go-webserver/pkg/router"
	"github.com/urfave/negroni"
)

func main() {
	// HTTP router
	r := router.New()

	// Negroni Middleware Stack
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(metrics.NewMiddleware(info.Report.Program))
	n.Use(negroni.NewLogger())
	n.UseHandler(r) // pass router to negroni

	// Create and run server
	s := NewServer(":8000", n) // pass port and negroni
	s.Run()
}
