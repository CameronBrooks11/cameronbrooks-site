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
// During go run (no -ldflags), returns version=dev and build_time=unknown.
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
	_ = json.NewEncoder(w).Encode(map[string]string{
		"version":    v,
		"build_time": bt,
	})
}
