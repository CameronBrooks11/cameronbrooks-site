package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestChainAppliesMiddlewaresInDeclarationOrder(t *testing.T) {
	var calls []string

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "A:before")
			next.ServeHTTP(w, r)
			calls = append(calls, "A:after")
		})
	}
	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "B:before")
			next.ServeHTTP(w, r)
			calls = append(calls, "B:after")
		})
	}

	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
		w.WriteHeader(http.StatusNoContent)
	})

	h := Chain(base, mwA, mwB)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	want := []string{"A:before", "B:before", "handler", "B:after", "A:after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected call order: got %v want %v", calls, want)
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rr.Code, http.StatusNoContent)
	}
}
