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

	// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
	"github.com/pkg/errors"
)

// Server implements our HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP Server
func NewServer(hostPort string, h http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":" + hostPort,
			Handler:        h, // pass in negroni or other mux/router
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second, // Go ver >1.8
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "hostname unavailable")
	}

	// Error handling
	errChan := make(chan error, 10)

	// Run server
	go func() {
		log.Printf("%s - Starting server on port %v", hostname, s.server.Addr)
		errChan <- s.server.ListenAndServe()
	}()

	// SIGINT/SIGTERM handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		case err := <-errChan:
			if err != http.ErrServerClosed { // Go ver >1.8
				// log.Fatalf("listen: %s\n", err)
				return errors.Wrap(err, "server not closed")
			}
		case <-sigs:
			fmt.Println("")
			log.Printf("%s - Shutdown signal received, exiting...\n", hostname)
			// shut down gracefully, but wait no longer than 5 seconds before halting
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.server.Shutdown(ctx); err != nil {
				// log.Fatalf("Server could not shutdown: %v", err)
				return errors.Wrap(err, "server could not shutdown")
			}
			log.Printf("%s - Server gracefully stopped.\n", hostname)
			os.Exit(0)
		}
	}
}
