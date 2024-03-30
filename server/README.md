# github.com/n-vr/httptoolkit/server

[![Go Reference](https://pkg.go.dev/badge/github.com/n-vr/httptoolkit/server.svg)](https://pkg.go.dev/github.com/n-vr/httptoolkit/server)

This package can be used to start an HTTP server that will gracefully shutdown when a context is canceled.

## Example

```golang
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/n-vr/httptoolkit/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleHome)

	// Create a context that will be canceled when the program receives an interrupt signal.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server.ListenAndServe(ctx, ":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
```