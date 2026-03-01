package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// contextKey is an unexported type for context keys in this package.
// Prevents collisions with keys from other packages.
type contextKey string

const requestIDKey contextKey = "request_id"

// RequestID ensures every request has a unique ID.
// It reads X-Request-ID from the incoming request (if present) or
// generates a crypto-random 8-byte (16 hex char) ID.
// The ID is set on the response header and stored in request context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			if _, err := rand.Read(b); err != nil {
				// crypto/rand failure is extremely unlikely; use fixed fallback
				id = "00000000deadbeef"
			} else {
				id = hex.EncodeToString(b)
			}
		}

		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom retrieves the request ID from context.
// Returns empty string if no ID is present.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
