package middleware

import "net/http"

// Middleware wraps an http.Handler with additional behavior.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares around h in declaration order:
// the first middleware is the outermost (first in, last out).
//
//	Chain(mux, RequestID, Logger)
//	-> RequestID(Logger(mux))
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
