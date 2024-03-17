# HTTP Toolkit

This is a toolkit to reduce boilerplate code for common HTTP related things I need to do in go.

## Subpackages

### `github.com/n-vr/httptoolkit/handler`

Package handler implements an HTTP handler that can return an error, while staying 100% compatible with the standard library's net/http package.

### `github.com/n-vr/httptoolkit/problem`

Package problem implements RFC 9457 errors that can be returned from a handler.
