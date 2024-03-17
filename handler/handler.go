// Package handler implements an HTTP handler that can return an error,
// while staying 100% compatible with the standard library's net/http package.
package handler

import "net/http"

// Handler is an HTTP handler that can return an error
// and implements the http.Handler interface.
type Handler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP implements the http.Handler interface.
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		ErrorHandler(err, w)
	}
}
