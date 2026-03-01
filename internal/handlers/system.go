package handlers

import (
	"fmt"
	"net/http"
)

// Version and BuildTime are injected at compile time via:
//
//	-ldflags="-X github.com/CameronBrooks11/cameronbrooks-site/internal/handlers.Version=$(git rev-parse --short HEAD)"
//	-ldflags="-X github.com/CameronBrooks11/cameronbrooks-site/internal/handlers.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
//
// During go run (development), both are empty strings.
var (
	Version   string
	BuildTime string
)

// Healthz handles GET /healthz.
// Returns 200 OK with body "ok".
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

// Version handles GET /version.
// Returns the git SHA and build timestamp injected at compile time.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	version := Version
	if version == "" {
		version = "dev"
	}

	buildTime := BuildTime
	if buildTime == "" {
		buildTime = "unknown"
	}

	fmt.Fprintf(w, "version=%s build_time=%s\n", version, buildTime)
}
