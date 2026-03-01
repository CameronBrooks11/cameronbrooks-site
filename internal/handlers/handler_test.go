package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func initHandlerTemplates(t *testing.T) {
	t.Helper()
	if err := InitTemplates(); err != nil {
		t.Fatalf("InitTemplates returned error: %v", err)
	}
}

func TestHomeReturnsNotFoundForNonRootPath(t *testing.T) {
	initHandlerTemplates(t)

	h := New()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	h.Home(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestProjectReturnsNotFoundForUnknownSlug(t *testing.T) {
	initHandlerTemplates(t)

	h := New()
	req := httptest.NewRequest(http.MethodGet, "/projects/does-not-exist", nil)
	req.SetPathValue("slug", "does-not-exist")
	rec := httptest.NewRecorder()

	h.Project(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestPostReturnsNotFoundForUnknownSlug(t *testing.T) {
	initHandlerTemplates(t)

	h := New()
	req := httptest.NewRequest(http.MethodGet, "/writing/does-not-exist", nil)
	req.SetPathValue("slug", "does-not-exist")
	rec := httptest.NewRecorder()

	h.Post(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHealthzReturnsOK(t *testing.T) {
	h := New()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	h.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), "ok")
	}
}

func TestVersionReturnsDevDefaultsWhenUnset(t *testing.T) {
	h := New()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	h.Version(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: got %q want %q", got, "application/json")
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode version response: %v", err)
	}

	if got, want := payload["version"], "dev"; got != want {
		t.Fatalf("unexpected version: got %q want %q", got, want)
	}
	if got, want := payload["build_time"], "unknown"; got != want {
		t.Fatalf("unexpected build_time: got %q want %q", got, want)
	}
}

func TestVersionReturnsInjectedBuildMetadata(t *testing.T) {
	h := New()
	h.AppVersion = "abc1234"
	h.AppBuildTime = "2026-03-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	h.Version(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode version response: %v", err)
	}

	if got, want := payload["version"], "abc1234"; got != want {
		t.Fatalf("unexpected version: got %q want %q", got, want)
	}
	if got, want := payload["build_time"], "2026-03-01T00:00:00Z"; got != want {
		t.Fatalf("unexpected build_time: got %q want %q", got, want)
	}
}

func TestNotFoundWrites404(t *testing.T) {
	initHandlerTemplates(t)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	notFound(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestInternalErrorWrites500(t *testing.T) {
	initHandlerTemplates(t)

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()

	internalError(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusInternalServerError)
	}
}
