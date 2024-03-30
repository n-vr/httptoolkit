// Package server implements a simple HTTP server
// that features graceful shutdown.
package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

var (
	// ErrDeadlineExceeded is returned when the server shutdown deadline is exceeded.
	ErrDeadlineExceeded = context.DeadlineExceeded
)

// Server is a simple HTTP server that features graceful shutdown.
// It wraps around http.Server and provides additional functionality.
//
// The server can be started using ListenAndServe or Serve and passing in a context.
// The server will be gracefully shutdown when the context is canceled.
type Server struct {
	*http.Server

	// ShutdownTimeout is the maximum duration the server is allowed to shutdown gracefully.
	ShutdownTimeout time.Duration
}

// ListenAndServe listens on the TCP network address s.Addr and handles incoming HTTP connections.
// It uses http.Server.ListenAndServe under the hood.
// The server will be gracefully shutdown when the context is canceled.
//
// The function will block until the server is shutdown.
//
// An error is returned if the server fails to start or shutdown.
// If the server is gracefully shutdown, nil is returned.
func (s *Server) ListenAndServe(ctx context.Context) error {
	return s.runUntilCancel(ctx, s.Server.ListenAndServe)
}

// Serve listens on the given listener and handles incoming HTTP connections.
// It uses http.Server.Serve under the hood.
// The server will be gracefully shutdown when the context is canceled.
//
// The function will block until the server is shutdown.
//
// An error is returned if the server fails to start or shutdown.
// If the server is gracefully shutdown, nil is returned.
func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	return s.runUntilCancel(ctx, func() error {
		return s.Server.Serve(l)
	})
}

// runUntilCancel runs the given function fn until the context is canceled.
// The server is gracefully shutdown when the context is canceled.
// This is a convenience function to avoid duplicating the same logic in ListenAndServe and Serve.
func (s *Server) runUntilCancel(ctx context.Context, fn func() error) error {
	s.setBaseContext(ctx)

	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		errChan <- fn()
	}()

	// Wait for the context to be canceled.
	// At this point, the server needs to be shutdown.
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	err := s.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	err = <-errChan
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) setBaseContext(ctx context.Context) {
	s.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
}

// ListenAndServe listens on the TCP network address addr and handles incoming HTTP connections.
// The server will be gracefully shutdown when the context is canceled.
//
// The function will block until the server is shutdown.
//
// An error is returned if the server fails to start or shutdown.
// If the server is gracefully shutdown, nil is returned.
func ListenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srv := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
	return srv.ListenAndServe(ctx)
}
