# Phase 08 — Routing & main.go

**Goal:** Replace the Phase 02 skeleton `main.go` with the complete production wiring: all routes registered on a `net/http` mux, `http.Server` with timeouts, graceful shutdown via `signal.NotifyContext`, `slog` with environment-controlled log level, `InitTemplates()` fatal on error, and `middleware.Chain` wrapping the mux. Also update the Makefile `build` target to inject `Version` and `BuildTime` via `-ldflags`.

**Exit gate:** `go run ./cmd/site` with `SITE_ENV=dev` serves all routes; `/healthz` → `200 ok`; `/version` → build vars; an unmatched path → 404 template; Ctrl-C prints a shutdown log line and exits cleanly; `make build` produces `bin/site` with `-ldflags` embedded.

---

## Prerequisites

- Phase 06 complete (`middleware.Chain`, `RequestID`, `Logger` exist)
- Phase 07 complete (all handlers exist on `*handlers.Handler`)
- Phase 05 complete (`handlers.InitTemplates()` exists)

---

## Files to modify in this phase

```
cmd/site/main.go    — REPLACE Phase 02 skeleton with complete wiring
Makefile            — update build target with -ldflags (already present as scaffolding; confirm vars)
```

---

## Step 1 — Complete `cmd/site/main.go`

This file replaces the Phase 02 skeleton entirely. Every import and every decision documented in TEMPLATES.md is wired here.

**File: `cmd/site/main.go`**

```go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/handlers"
	"github.com/CameronBrooks11/cameronbrooks-site/internal/middleware"
	"github.com/CameronBrooks11/cameronbrooks-site/static"
)

// Version and BuildTime are injected at compile time via -ldflags.
// During `go run` (development) they hold their zero values — that is expected.
//
//	make build
//	→ -ldflags="-X main.Version=<git-sha> -X main.BuildTime=<iso8601>"
var (
	Version   string
	BuildTime string
)

func main() {
	// --- Logging setup -------------------------------------------------------
	// Default level: INFO. Set SITE_ENV=dev to enable DEBUG level.
	logLevel := slog.LevelInfo
	if os.Getenv("SITE_ENV") == "dev" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// --- Template cache ------------------------------------------------------
	// Fatal on any parse error — a missing or broken template is not recoverable.
	if err := handlers.InitTemplates(); err != nil {
		slog.Error("failed to initialize templates", "err", err)
		os.Exit(1)
	}

	// --- Handler setup -------------------------------------------------------
	h := handlers.New()
	h.AppVersion = Version
	h.AppBuildTime = BuildTime

	// --- Routing -------------------------------------------------------------
	mux := http.NewServeMux()

	// Page routes
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /projects", h.Projects)
	mux.HandleFunc("GET /projects/{slug}", h.Project)
	mux.HandleFunc("GET /writing", h.Writing)
	mux.HandleFunc("GET /writing/{slug}", h.Post)
	mux.HandleFunc("GET /about", h.About)
	mux.HandleFunc("GET /contact", h.Contact)

	// System routes
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /version", h.Version)

	// Static assets from embedded FS
	mux.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))),
	)

	// --- Server --------------------------------------------------------------
	addr := os.Getenv("SITE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      middleware.Chain(mux, middleware.RequestID, middleware.Logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// --- Graceful shutdown ---------------------------------------------------
	// NotifyContext cancels ctx on SIGTERM or SIGINT (Ctrl-C).
	// Systemd sends SIGTERM on `systemctl restart` / `systemctl stop`.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		slog.Info("starting server", "addr", srv.Addr, "version", Version)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Block until signal received
	<-ctx.Done()

	slog.Info("shutting down", "reason", ctx.Err())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
	slog.Info("shutdown complete")
}
```

---

## Step 2 — Add `Version` and `BuildTime` fields to `Handler`

Phase 07 noted this adjustment. Update `internal/handlers/handler.go`:

**File: `internal/handlers/handler.go`** (replace)

```go
package handlers

// Handler holds application-level dependencies shared across all page handlers.
type Handler struct {
	AppVersion   string // injected from main.Version (set via -ldflags at build time)
	AppBuildTime string // injected from main.BuildTime (set via -ldflags at build time)
}

// New returns an initialized Handler. Set AppVersion and AppBuildTime after construction
// if build metadata is needed (see cmd/site/main.go).
func New() *Handler {
	return &Handler{}
}
```

---

## Step 3 — Update `internal/handlers/system.go`

Remove the package-level `Version` and `BuildTime` vars and read from the receiver instead:

**File: `internal/handlers/system.go`** (replace)

```go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Healthz handles GET /healthz.
// Returns 200 OK with plain-text body "ok".
// Used by Caddy and external monitoring to confirm the process is alive.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

// Version handles GET /version.
// Returns the git SHA and build timestamp injected via -ldflags at compile time.
// During `go run` (no -ldflags), returns version=dev and build_time=unknown.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := h.AppVersion
	if v == "" {
		v = "dev"
	}
	bt := h.AppBuildTime
	if bt == "" {
		bt = "unknown"
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"version":    v,
		"build_time": bt,
	})
}
```

---

## Step 4 — Confirm the Makefile `build` target

The Makefile was scaffolded in Phase 01 with the ldflags already written. Verify the `build` target looks exactly like this (tabs, not spaces):

```makefile
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY) $(CMD)
```

And the variable definitions at the top:

```makefile
VERSION   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
LDFLAGS   := -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME)"
```

These are already present from Phase 01. No edits needed unless they were accidentally changed.

> **Windows note:** `date -u +%Y-%m-%dT%H:%M:%SZ` is a Unix command. On Windows, `make build` must be run from WSL or Git Bash, not PowerShell. An alternative for PowerShell:
>
> ```powershell
> $env:GOOS="linux"; $env:GOARCH="amd64"
> $version = git rev-parse --short HEAD
> $buildtime = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")
> go build "-ldflags=-X main.Version=$version -X main.BuildTime=$buildtime" -o bin/site ./cmd/site
> ```
>
> Either way works — `make build` is the canonical target; the PowerShell form is the fallback.

---

## Step 5 — Full verification

### 5a. Compile

```sh
go build ./...
```

Expected: exits 0.

### 5b. Run in development mode

```sh
$env:SITE_ENV="dev"; go run ./cmd/site
```

Expected slog output (JSON format):

```json
{
  "time": "...",
  "level": "INFO",
  "msg": "starting server",
  "addr": ":8080",
  "version": ""
}
```

### 5c. Smoke test routes

In a second terminal (server must be running):

```sh
# All should return HTTP 200
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/projects
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/projects/cameronbrooks-site
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/writing
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/writing/hello-world
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/about
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/contact
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/healthz

# Should return 200 with plain-text body
curl http://localhost:8080/healthz
# Expected: ok

curl http://localhost:8080/version
# Expected: {"build_time":"unknown","version":"dev"}

# Should return 404
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/does-not-exist
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/projects/no-such-slug
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/writing/draft-post
```

### 5d. Graceful shutdown

Stop the server with Ctrl-C. Expected slog output:

```json
{"time":"...","level":"INFO","msg":"shutting down","reason":"context canceled"}
{"time":"...","level":"INFO","msg":"shutdown complete"}
```

### 5e. Request logs

Each request to the server (from 6c) should produce one structured log line like:

```json
{
  "time": "...",
  "level": "INFO",
  "msg": "request",
  "method": "GET",
  "path": "/healthz",
  "status": 200,
  "duration_ms": 0,
  "request_id": "a1b2c3d4e5f6a7b8",
  "remote_ip": "127.0.0.1"
}
```

### 5f. Build with ldflags

```sh
# From WSL or Git Bash
make build
./bin/site &
curl http://localhost:8080/version
# Expected: {"version":"<git-sha>","build_time":"<iso8601>"}
kill %1
```

---

## Step 7 — Commit

```sh
git add cmd/site/main.go internal/handlers/handler.go internal/handlers/system.go
git commit -m "phase 08: routing and main.go wiring"
```

---

## Exit gate checklist

- [ ] `go build ./...` exits 0
- [ ] All nine page routes return HTTP 200
- [ ] `/healthz` returns `200` with body `ok`
- [ ] `/version` returns `200` with JSON `{"version":"dev","build_time":"unknown"}` during `go run`
- [ ] `/does-not-exist` returns `404` and renders `notFound.gohtml`
- [ ] `/projects/no-such-slug` returns `404`
- [ ] `/writing/draft-post` returns `404` (draft not exposed)
- [ ] Ctrl-C prints `shutting down` and `shutdown complete` log lines
- [ ] Each request produces one structured JSON log line from the Logger middleware
- [ ] `X-Request-ID` header is present on all responses
- [ ] `make build` produces `bin/site` at `bin/site` (linux/amd64)
- [ ] `bin/site` after `make build` — `curl /version` shows real git SHA in JSON (not `"dev"`)
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 09.
