package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/n-vr/httptoolkit/handler"
	"github.com/n-vr/httptoolkit/problem"
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

func TestHandler_withHandlerError(t *testing.T) {
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

func TestHandler_withProblemError(t *testing.T) {
	server, teardown := setupTestServer(func(mux *http.ServeMux) {
		mux.Handle("GET /test", handler.Handler(func(w http.ResponseWriter, r *http.Request) error {
			return problem.New(errors.New("problem error"), http.StatusTeapot)
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

	if resp.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("expected Content-Type %q, got %q", "application/problem+json", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"detail":"problem error","status":418,"title":"I'm a teapot","type":"https://httpstatuses.com/418"}`
	if strings.TrimSpace(string(body)) != want {
		t.Errorf("expected body %q, got %q", want, body)
	}
}
