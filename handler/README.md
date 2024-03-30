# github.com/n-vr/httptoolkit/handler

[![Go Reference](https://pkg.go.dev/badge/github.com/n-vr/httptoolkit/handler.svg)](https://pkg.go.dev/github.com/n-vr/httptoolkit/handler)

This package can be used to create HTTP handlers that can return errors.

It exposes an error handler function (`handler.ErrorHandler`) that you can replace assign another function to.

The default error handler will return the error message. A status code can also be set if you return a `handler.Error` (with the `handler.NewError` function).

## Example

```golang
package main

import (
	"errors"
	"net/http"

	"github.com/n-vr/httptoolkit/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", handler.Handler(handleHome))
	mux.Handle("GET /error", handler.Handler(handleError))

	http.ListenAndServe(":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Hello, World!"))
	return nil
}

func handleError(w http.ResponseWriter, r *http.Request) error {
	err := errors.New("an error occurred")
	return handler.NewError(err, http.StatusTeapot)
}
```