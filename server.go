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

const (
	// NOTE: If you use docker, docker stop has a default timeout of 10 seconds,
	// so the graceful timeout should be set to expire before then.
	timeout = 5 * time.Second
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

// Run starts the HTTP server and performs a graceful shutdown
func (s *Server) Run() error {

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "hostname unavailable")
	}

	// Error handling
	listenErr := make(chan error, 1)

	// Run server
	go func() {
		log.Printf("%s - Web server available on port %v", hostname, s.server.Addr)
		log.Printf("%s - Press Ctrl+C to stop", hostname)
		listenErr <- s.server.ListenAndServe()
	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		// If server.ListenAndServe() cannot start due to errors such as "port
		// in use" it will return an error.
		case err := <-listenErr:
			return err
		// handle termination signal
		case <-osSignals:
			fmt.Printf("\n")
			log.Printf("%s - Shutdown signal received.\n", hostname)

			// Servers in the process of shutting down should disable KeepAlives.
			s.server.SetKeepAlivesEnabled(false)

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := s.server.Shutdown(ctx); err != nil {
				return err
			}

			// return any errors from this channel other than "ServerClosed"
			if err := <-listenErr; err != http.ErrServerClosed {
				return err
			}

			log.Printf("%s - Server gracefully stopped.\n", hostname)
			return nil
		}
	}
}
