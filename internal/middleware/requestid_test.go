package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRequestIDUsesIncomingHeader(t *testing.T) {
	const incoming = "req-123"

	var ctxID string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID = RequestIDFrom(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", incoming)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if got := rr.Header().Get("X-Request-ID"); got != incoming {
		t.Fatalf("unexpected response X-Request-ID: got %q want %q", got, incoming)
	}
	if ctxID != incoming {
		t.Fatalf("unexpected context request id: got %q want %q", ctxID, incoming)
	}
}

func TestRequestIDGeneratesHexIDWhenMissing(t *testing.T) {
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	id := rr.Header().Get("X-Request-ID")
	if id == "" {
		t.Fatal("expected generated request id header")
	}
	matched, err := regexp.MatchString("^[0-9a-f]{16}$", id)
	if err != nil {
		t.Fatalf("failed to compile regexp: %v", err)
	}
	if !matched {
		t.Fatalf("expected 16-char lowercase hex request id, got %q", id)
	}
}
