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
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	slog := logger.Sugar()
	// logger.Infow("Failed to fetch URL.",
	// 	// Structured context as loosely-typed key-value pairs.
	// 	"url", url,
	// 	"attempt", retryNum,
	// 	"backoff", time.Second,
	// )
	slog.Infof("Failed to fetch URL")

	// logger := zap.New(
	// 	zap.NewJSONEncoder(
	// 		zap.RFC3339Formatter("time"), // human-readable timestamps
	// 		zap.MessageKey("msg"),        // customize the message key
	// 		zap.LevelString("level"),     // stringify the log level
	// 	),
	// )

	// logger = logger.With(
	// 	zap.Int("pid", os.Getpid()),
	// 	zap.String("exe", path.Base(os.Args[0])),
	// )
	//
	// textLogger := zap.New(zap.NewTextEncoder(
	// 	zap.TextTimeFormat(time.RFC822),
	// ))

	// textLogger.Debug("This is debug data.", zap.Int("foo", 42))
	// textLogger.Info("This is a text log.", zap.Int("bar", 12))
	//
	// logger.Warn("Log without structured data...")
	// logger.Warn(
	// 	"Or use strongly-typed wrappers to add structured context.",
	// 	zap.String("library", "zap"),
	// 	zap.Duration("latency", time.Nanosecond),
	// )
	//
	// // Avoid re-serializing the same data repeatedly by creating a child logger
	// // with some attached context. That context is added to all the child's
	// // log output, but doesn't affect the parent.
	// child := logger.With(
	// 	zap.String("user", "jane@test.com"),
	// 	zap.Int("visits", 42),
	// )
	// child.Error("Oh no!")

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
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second, // Go 1.8
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
