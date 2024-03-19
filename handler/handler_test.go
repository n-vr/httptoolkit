package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n-vr/httptoolkit/handler"
)

func setupTestServer(registerRoutes func(mux *http.ServeMux)) (*httptest.Server, func()) {
	mux := http.NewServeMux()

	registerRoutes(mux)

	server := httptest.NewServer(mux)
	return server, server.Close
}

func TestHandler_withoutError(t *testing.T) {
	server, teardown := setupTestServer(func(mux *http.ServeMux) {
		mux.Handle("GET /test", handler.Handler(func(w http.ResponseWriter, r *http.Request) error {
			return nil
		}))
	})
	defer teardown()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHandler_withError(t *testing.T) {
	server, teardown := setupTestServer(func(mux *http.ServeMux) {
		mux.Handle("GET /test", handler.Handler(func(w http.ResponseWriter, r *http.Request) error {
			return handler.NewError(nil, http.StatusTeapot)
		}))
	})
	defer teardown()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("expected status %d, got %d", http.StatusTeapot, resp.StatusCode)
	}
}
