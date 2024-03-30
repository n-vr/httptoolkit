package server_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/n-vr/httptoolkit/server"
	"github.com/n-vr/httptoolkit/server/servertest"
)

func TestServer_shutdownWithoutOpenConnections(t *testing.T) {
	t.Parallel()

	srv, cancel := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Implicitly write status code 200.
		w.Write([]byte("OK!"))
	})
	defer srv.Close()

	// Wait for the server to start.
	<-srv.Started

	// Make a request while the server is running.
	// The request should succeed.
	makeRequest(t, srv.URL, nil, http.StatusOK)

	if err := srv.Error(); err != nil {
		t.Fatalf("Server.Err = %v, want nil", err)
	}

	// Cancel the context to shutdown the server.
	cancel()

	// Make a request while the server is shut down.
	// The request should fail with an error.
	makeRequest(t, srv.URL, &os.SyscallError{}, 0)

	if err := srv.Error(); err != nil {
		t.Fatalf("Server.Err = %v, want nil", err)
	}
}

// Shutdown deadline exceeded

func TestServer_shutdownWithOpenConnections(t *testing.T) {
	t.Parallel()

	srv, cancel := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		// Implicitly write status code 200.
		w.Write([]byte("OK!"))
	})

	// Wait for the server to start.
	<-srv.Started

	wg := sync.WaitGroup{}
	started := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		close(started)
		// Make a request while the server is running.
		// The request should succeed.
		makeRequest(t, srv.URL, nil, http.StatusOK)
	}()

	// Wait for the request to start.
	<-started
	time.Sleep(100 * time.Millisecond)

	if err := srv.Error(); err != nil {
		t.Fatalf("Server.Err = %v, want nil", err)
	}

	// Cancel the context to shutdown the server.
	// Set a timeout that is longer than
	// the time it takes to finish the in-flight request.
	srv.Config.ShutdownTimeout = 1 * time.Millisecond
	cancel()

	wg.Wait()

	// The server should have an error.
	if err := srv.Error(); err != server.ErrDeadlineExceeded {
		t.Fatalf("Server.Err = %v, want %v", err, server.ErrDeadlineExceeded)
	}
}

func TestServer_shutdownDeadlineExceeded(t *testing.T) {
	t.Parallel()

	srv, cancel := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		err := longExpensiveOperation(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Implicitly write status code 200.
		w.Write([]byte("OK!"))
	})

	// Wait for the server to start.
	<-srv.Started

	wg := sync.WaitGroup{}
	started := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		close(started)
		// Make a request while the server is running.
		// The request should succeed.
		makeRequest(t, srv.URL, nil, http.StatusInternalServerError)
	}()

	// Wait for the request to start.
	<-started
	time.Sleep(100 * time.Millisecond)

	if err := srv.Error(); err != nil {
		t.Fatalf("Server.Err = %v, want nil", err)
	}

	// Cancel the context to shutdown the server.
	// Set a timeout that is shorter than
	// the time it takes to finish the in-flight request.
	srv.Config.ShutdownTimeout = 1 * time.Second
	cancel()

	wg.Wait()

	// The server should not have an error.
	if err := srv.Error(); err != nil {
		t.Fatalf("Server.Err = %v, want nil", err)
	}
}

func setupTestServer(t *testing.T, fn http.HandlerFunc) (srv *servertest.Server, cancelCtx func()) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/test", fn)

	ctx, cancel := context.WithCancel(context.Background())
	return servertest.NewServer(ctx, mux),
		func() {
			if ctx.Err() == nil {
				t.Log("Cancelling the context")
				cancel()
			}
		}
}

func makeRequest(t *testing.T, addr string, excepctedErrType any, expectedStatus int) {
	t.Helper()

	resp, err := http.Get(addr + "/test")
	if err != nil {
		if excepctedErrType == nil {
			t.Errorf("http.Get() err = %v, want nil", err)
		}
		if !errors.As(err, &excepctedErrType) {
			t.Errorf("http.Get() err = %t, want %t", err, excepctedErrType)
		}
		// Return since the response is nil.
		return
	}

	if resp.StatusCode != expectedStatus {
		t.Errorf("http.Get() status = %v, want %v", resp.StatusCode, expectedStatus)
	}

	resp.Body.Close()
}

func longExpensiveOperation(ctx context.Context) error {
	select {
	case <-time.After(10 * time.Second):
		return nil
	case <-ctx.Done():
		return errors.New("context deadline exceeded")
	}
}
