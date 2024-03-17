package handler

import (
	"net/http"

	"github.com/n-vr/httptoolkit/problem"
)

// Error is an error that also contains a status code.
// It implements the error interface and can be returned from a Handler.
// In conjuction with the default error handler,
// it will respond as plain text to the client with the status code set.
type Error struct {
	Err        error
	StatusCode int
}

// Create a new Error using err and statusCode.
func NewError(err error, statusCode int) *Error {
	return &Error{
		Err:        err,
		StatusCode: statusCode,
	}
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.StatusCode)
}

// ErrorHandlerFunc is a function that will handle an error returned by a Handler.
type ErrorHandlerFunc = func(err error, w http.ResponseWriter)

// ErrorHandler is called with an error when a Handler returns one.
var ErrorHandler ErrorHandlerFunc = defaultErrorHandler

// Default error handler.
// If the error is an *Error, it will respond as plain text with the status code set.
// Otherwise, it will use problem.HTTPErrorHandler.
func defaultErrorHandler(err error, w http.ResponseWriter) {
	if e, ok := err.(*Error); ok {
		http.Error(w, http.StatusText(e.StatusCode), e.StatusCode)
		return
	}

	problem.HTTPErrorHandler(err, w)
}
