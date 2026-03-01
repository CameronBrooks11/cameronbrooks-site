package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerWithRequestIDPreservesStatusAndSetsHeader(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	h := Chain(base, RequestID, Logger)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusTeapot {
		t.Fatalf("unexpected status: got %d want %d", rr.Code, http.StatusTeapot)
	}
	if rr.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID response header to be set")
	}
}

func TestRealIPPrefersXForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "10.0.0.2:9999"
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.2")

	got := realIP(req)
	if got != "203.0.113.10" {
		t.Fatalf("unexpected real IP: got %q want %q", got, "203.0.113.10")
	}
}

func TestRealIPFallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "10.0.0.2:9999"

	got := realIP(req)
	if got != "10.0.0.2" {
		t.Fatalf("unexpected fallback IP: got %q want %q", got, "10.0.0.2")
	}
}
