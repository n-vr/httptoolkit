// Package servertest provides utilities for testing the server package.
//
// Heavily inspired by the httptest package from the standard library.
package servertest

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/n-vr/httptoolkit/server"
)

type Server struct {
	URL      string
	Config   *server.Server
	Started  chan struct{}
	err      chan error
	listener net.Listener
}

// NewServer creates a new server with the given handler
// and starts listening on a random port.
//
// The URL of the server is available in the URL field.
//
// The server is closed by calling the Close method.
func NewServer(ctx context.Context, handler http.Handler) *Server {
	srv := &Server{
		listener: newListener(),
		Config: &server.Server{
			Server: &http.Server{
				Handler: handler,
			},
		},
		Started: make(chan struct{}, 1),
		err:     make(chan error, 1),
	}

	go func() {
		srv.err <- srv.start(ctx)
	}()

	return srv
}

func (s *Server) Error() error {
	if s.err == nil {
		return nil
	}

	if len(s.err) == 1 {
		return <-s.err
	}
	return nil
}

func (s *Server) start(ctx context.Context) error {
	if s.URL != "" {
		panic("servertest: server already started")
	}
	s.URL = "http://" + s.listener.Addr().String()

	// Close the started channel to signal that the server has started.
	close(s.Started)

	return s.Config.Serve(ctx, s.listener)
}

// Close closes the server.
func (s *Server) Close() {
	if s.URL == "" {
		panic("servertest: server already closed")
	}
	s.listener.Close()
	s.URL = ""
	close(s.err)
}

func newListener() net.Listener {
	var (
		l   net.Listener
		err error
	)

	// Try to listen on a random port 5 times.
	// This is to avoid flaky tests due to port collisions.
	for i := 0; i < 5; i++ {
		l, err = net.Listen("tcp", "localhost:0")
		if err == nil {
			break
		}
	}
	if err != nil {
		panic(fmt.Sprintf("servertest: failed to listen: %v", err))
	}

	return l
}
