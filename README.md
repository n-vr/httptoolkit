# HTTP Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/n-vr/httptoolkit)](https://goreportcard.com/report/github.com/n-vr/httptoolkit)

This is a toolkit to reduce boilerplate code for common HTTP related things I need to do in go.

## Subpackages

### `github.com/n-vr/httptoolkit/handler`

[![Go Reference](https://pkg.go.dev/badge/github.com/n-vr/httptoolkit/handler.svg)](https://pkg.go.dev/github.com/n-vr/httptoolkit/handler)

Package handler implements an HTTP handler that can return an error, while staying 100% compatible with the standard library's net/http package.

### `github.com/n-vr/httptoolkit/problem`

[![Go Reference](https://pkg.go.dev/badge/github.com/n-vr/httptoolkit/problem.svg)](https://pkg.go.dev/github.com/n-vr/httptoolkit/problem)

Package problem implements RFC 9457 errors that can be returned from a handler.

### `github.com/n-vr/httptoolkit/server`

[![Go Reference](https://pkg.go.dev/badge/github.com/n-vr/httptoolkit/server.svg)](https://pkg.go.dev/github.com/n-vr/httptoolkit/server)

Package server implements a simple HTTP server that features graceful shutdown.
