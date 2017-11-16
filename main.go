// This is an example of a simple webserver with a graceful
// shutdown. It is now possible to call srv.Close() to halt an
// http.Server immediately, or srv.Shutdown(ctx) to stop
// and gracefully drain the server of connections. We are using
// httprouter and negroni also.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dstroot/simple-go-webserver/pkg/router"
	"github.com/urfave/negroni"
)

func main() {

	// Get hostname
	hostname, _ := os.Hostname()

	// SIGINT/SIGTERM handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Error handling
	errChan := make(chan error, 10)

	// HTTP router
	r := router.New()

	// Negroni Middleware Stack
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(r) // pass router to negroni

	// Create Server
	s := &http.Server{
		Addr:           ":8000",
		Handler:        n, // pass in negroni
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second, // Go ver >1.8
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("%s - Starting server on port %v", hostname, s.Addr)

	// Run server
	go func() {
		errChan <- s.ListenAndServe()
	}()

	// Handle channels/graceful shutdown
	for {
		select {
		case err := <-errChan:
			if err != http.ErrServerClosed { // Go ver >1.8
				log.Fatalf("listen: %s\n", err)
			}
		case <-sigs:
			fmt.Println("")
			log.Printf("%s - Shutdown signal received, exiting...\n", hostname)
			// shut down gracefully, but wait no longer than 5 seconds before halting
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil {
				log.Fatalf("Server could not shutdown: %v", err)
			}
			log.Printf("%s - Server gracefully stopped.\n", hostname)
			os.Exit(0)
		}
	}
}
