package mymiddleware

import (
	"net/http"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain - chains all middleware functions right to left
// https://husobee.github.io/golang/http/middleware/2015/12/22/simple-middleware.html
func Chain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise run recursively over nested handlers
	return m[0](Chain(f, m[1:cap(m)]...))
}
