# github.com/n-vr/httptoolkit/problem

This package can be used in addition to the `github.com/n-vr/httptoolkit/handler` package to return a response according to the poblem details format (RFC 9457).

## Example

```golang
package main

import (
	"errors"
	"net/http"

	"github.com/n-vr/httptoolkit/handler"
	"github.com/n-vr/httptoolkit/problem"
)

// userNotFoundProblemType is a problem type for user not found.
// This can be reused when creating a Problem.
//
// This is an example of how to create a problem type.
// The URI should should uniquely identify the problem type
// and dereferencing it should provide human-readable documentation.
var userNotFoundProblemType = problem.NewType("http://localhost/user-not-found", "User not found")

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", handler.Handler(handleHome))
	mux.Handle("GET /simple-error", handler.Handler(handleSimpleError))
	mux.Handle("GET /complex-error", handler.Handler(handleComplexError))

	// Set the error handler to the problem.ErrorHandler.
	// This is needed to handle problem.Problem errors.
	handler.ErrorHandler = problem.ErrorHandler

	http.ListenAndServe(":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Hello, World!"))
	return nil
}

func handleSimpleError(w http.ResponseWriter, r *http.Request) error {
	err := errors.New("an error occurred")
	return problem.New(err, http.StatusTeapot)
}

func handleComplexError(w http.ResponseWriter, r *http.Request) error {
	err := errors.New("an error occurred")
	return problem.New(err, http.StatusNotFound,
		problem.WithType(userNotFoundProblemType),
		problem.WithInstance("/users/123"),
		problem.WithExtension("user", "123"))
}
```