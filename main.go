// https://tylerchr.blog/golang-18-whats-coming/

// It is now possible to call srv.Close() to halt an
// http.Server immediately, or srv.Shutdown(ctx) to stop
// and gracefully drain the server of connections

package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dstroot/go1.8-http-graceful-exit/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

/**
 * Handlers
 */

// index handler handles GET /
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// If you want to debug your HTTP requests, then import the
	// net/http/httputil package, and invoke DumpRequest
	// with the parameter *http.Request and a boolean to specify if you
	// want to dump the request body as well. The function returns a
	// []byte, error. You could use it like this:
	dump := func(r *http.Request) {
		output, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("Error dumping request:", err)
			return
		}
		fmt.Println(string(output))
	}

	// The function call will dump your request method, URI with
	// query parameters, headers and request body if you have one.
	dump(r)

	// page data to render page
	data := map[string]interface{}{
		"title": "The most popular HTML, CSS, and JS library in the world.",
		"Key":   "Value",
		"Slice": []string{"One", "Two", "Three"},
	}

	// render page template
	err := tmpl.RenderTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering:", err)
		return
	}
}

func page(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// page data to render page
	data := map[string]interface{}{
		"title": "Page 2",
		"Key":   "Value",
		"Slice": []string{"One", "Two", "Three"},
	}

	// render page template
	err := tmpl.RenderTemplate(w, "page.html", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// hello handler handles GET /hello/:name
func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}

/**
 * Main
 */

func main() {

	// template processing
	tmpl.InitBufferPool()
	tmpl.LoadTemplates()

	// get hostname
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
	r.GET("/", index)
	r.GET("/page", page)
	r.GET("/hello/:name", hello)

	// handler for serving expvar data
	r.Handler("GET", "/debug/vars", expvar.Handler())

	// handler for serving static files
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	/**
	 * Negroni Middleware Stack
	 */

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(r)

	/**
	 * Service http connections
	 */

	// Create Server
	s := &http.Server{
		Addr:           ":8000",
		Handler:        n, // pass in negroni
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second, // >Go 1.8
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("%s - Starting server on port %v...", hostname, s.Addr)

	// Run server
	go func() {
		errChan <- s.ListenAndServe()
	}()

	// Handle channels/graceful shutdown
	for {
		select {
		case err := <-errChan:
			if err != http.ErrServerClosed { // Go 1.8
				log.Fatalf("listen: %s\n", err)
			}
		case <-sigs:
			fmt.Println("")
			log.Printf("%s - Shutdown signal received, exiting...\n", hostname)
			// shut down gracefully, but wait no longer than 5 seconds before halting
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.Shutdown(ctx)
			if err := s.Shutdown(ctx); err != nil {
				log.Fatalf("Server could not shutdown: %v", err)
			}
			log.Printf("%s - Server gracefully stopped.\n", hostname)
			os.Exit(0)
		}
	}
}
