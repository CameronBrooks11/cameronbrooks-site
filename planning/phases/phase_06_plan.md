# Phase 06 — Middleware

**Goal:** Implement `internal/middleware/` — the `Middleware` type alias, the `Chain` helper, `RequestID` (crypto/rand, 8-byte hex, context propagation), and `Logger` (per-request structured log line using `slog`, with a `statusRecorder` to capture the response code). These wrap the mux in Phase 08.

**Exit gate:** `go vet ./internal/middleware/...` passes; `go build ./...` exits 0; `Chain(handler, RequestID, Logger)` compiles and wraps a test handler without panicking.

---

## Prerequisites

- Phase 05 complete
- `internal/middleware/` directory exists; `.gitkeep` can be deleted

---

## Files to create in this phase

```
internal/middleware/middleware.go   — Middleware type, Chain helper
internal/middleware/requestid.go    — RequestID middleware
internal/middleware/logger.go       — Logger middleware, statusRecorder
```

---

## Step 1 — Delete `.gitkeep`

```sh
Remove-Item internal/middleware/.gitkeep
```

---

## Step 2 — `internal/middleware/middleware.go`

The `Chain` function applies middlewares right-to-left so that the first middleware listed is the outermost (first to receive a request, last to finish). This is the standard composition order — `Chain(mux, RequestID, Logger)` means `RequestID` wraps `Logger` wraps `mux`, so the request ID is set before the logger runs and is available to it.

**File: `internal/middleware/middleware.go`**

```go
package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares around h in declaration order:
// the first middleware is the outermost (first in, last out).
//
//	Chain(mux, RequestID, Logger)
//	→ RequestID(Logger(mux))
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
```

---

## Step 3 — `internal/middleware/requestid.go`

Reads the incoming `X-Request-ID` header (set by Caddy or a calling service) and falls back to generating a crypto-random 16-character hex string. The ID is:

1. Written as `X-Request-ID` on the response
2. Stored in the request context for the logger and any handler that needs it

**File: `internal/middleware/requestid.go`**

```go
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

// RequestID is a middleware that ensures every request has a unique ID.
// It reads X-Request-ID from the incoming request (forwarded by Caddy) or
// generates a new crypto-random 8-byte (16 hex char) ID if none is present.
// The ID is set on the response header and stored in the request context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			if _, err := rand.Read(b); err != nil {
				// crypto/rand failure is extremely unlikely; fall back to a fixed marker
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

// RequestIDFrom retrieves the request ID from a context.
// Returns an empty string if no ID is present (e.g. in tests that do not use RequestID middleware).
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
```

---

## Step 4 — `internal/middleware/logger.go`

Logs one structured line per request after the handler returns. Uses `log/slog` (stdlib). The `statusRecorder` wrapper captures the status code written by the handler — without it, the logger would always record 200 even for 404s.

`realIP` extracts the client IP from `X-Forwarded-For` (set by Caddy) and falls back to `r.RemoteAddr`. Since Go always has Caddy in front, the `X-Forwarded-For` value is the actual client IP.

**File: `internal/middleware/logger.go`**

```go
package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

// Logger is a middleware that emits one structured slog line per request.
// It runs after RequestID so the request ID is available in context.
// Log fields: method, path, status, duration_ms, request_id, remote_ip.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", RequestIDFrom(r.Context()),
			"remote_ip", realIP(r),
		)
	})
}

// statusRecorder wraps http.ResponseWriter to capture the written status code.
// The default status is 200 — it is overwritten when WriteHeader is called.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// realIP returns the client IP address.
// Prefers X-Forwarded-For set by Caddy; falls back to r.RemoteAddr.
// Strips the port component if present.
func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may be a comma-separated list; take the first entry
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
```

---

## Step 5 — Verify

```sh
go vet ./internal/middleware/...
```

Expected: exits 0.

```sh
go build ./...
```

Expected: exits 0.

**Verify Chain composition compiles:**

```sh
go run - <<'EOF'
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"

    "github.com/CameronBrooks11/cameronbrooks-site/internal/middleware"
)

func main() {
    base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "ok")
    })

    handler := middleware.Chain(base, middleware.RequestID, middleware.Logger)

    req := httptest.NewRequest("GET", "/healthz", nil)
    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    id := rr.Header().Get("X-Request-ID")
    fmt.Printf("status: %d\n", rr.Code)           // expect 200
    fmt.Printf("request_id set: %v\n", id != "")  // expect true
}
EOF
```

Expected output:

```
status: 200
request_id set: true
```

The `slog.Info("request", ...)` line will also print to stdout during this test — that is expected.

---

## Step 6 — Commit

```sh
git add internal/middleware/
git commit -m "phase 06: middleware"
```

---

## Exit gate checklist

- [ ] `go vet ./internal/middleware/...` exits 0
- [ ] `go build ./...` exits 0
- [ ] `Chain(base, RequestID, Logger)` compiles and runs (verification step above)
- [ ] `statusRecorder.WriteHeader` captures non-200 codes correctly (covered by the status assertion)
- [ ] `X-Request-ID` response header is set on every request (covered by the request_id assertion)
- [ ] `RequestIDFrom` returns the ID from context (used by Logger — confirmed by non-empty log output)
- [ ] `.gitkeep` removed from `internal/middleware/`
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 07.
