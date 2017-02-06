// https://tylerchr.blog/golang-18-whats-coming/

// It is now possible to call srv.Close() to halt an
// http.Server immediately, or srv.Shutdown(ctx) to stop
// and gracefully drain the server of connections

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

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

/**
 * Handlers
 */

// Index handler handles GET /
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	time.Sleep(1 * time.Second)
	fmt.Fprint(w, "Welcome!\n")
}

// Hello handler handles GET /hello/:name
func Hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}

/**
 * Main
 */

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// SIGINT or SIGTERM handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Error handling
	errChan := make(chan error, 10)

	/**
	 * HTTP Router
	 */

	r := httprouter.New()

	// Routes
	r.GET("/", Index)
	r.GET("/hello/:name", Hello)

	/**
	 * Negroni Middleware Stack
	 */

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	/**
	 * Service connections
	 */

	// Create Server
	s := &http.Server{
		Addr:           ":8080",
		Handler:        n, // pass negroni
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("%s - Starting server...", hostname)

	// Run server
	go func() {
		errChan <- s.ListenAndServe()
	}()

	// Handle channels
	for {
		select {
		case err := <-errChan:
			if err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		case <-sigs:
			fmt.Println("")
			log.Printf("%s - Shutdown signal received, exiting...\n", hostname)
			// shut down gracefully, but wait no longer than 5 seconds before halting
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.Shutdown(ctx)
			if err := s.Shutdown(ctx); err != nil {
				log.Fatalf("Server could not shutdown: %v", err)
			}
			log.Printf("%s - Server gracefully stopped.\n", hostname)
		}
	}
}
