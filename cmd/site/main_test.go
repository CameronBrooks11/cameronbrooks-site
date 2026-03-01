package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/handlers"
)

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	if err := handlers.InitTemplates(); err != nil {
		t.Fatalf("InitTemplates returned error: %v", err)
	}

	h := handlers.New()
	return newHTTPHandler(h)
}

func TestRouteContracts(t *testing.T) {
	app := newTestHandler(t)

	type routeTest struct {
		name   string
		path   string
		status int
		check  func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := []routeTest{
		{name: "home", path: "/", status: http.StatusOK},
		{name: "projects list", path: "/projects", status: http.StatusOK},
		{name: "writing list", path: "/writing", status: http.StatusOK},
		{
			name:   "healthz",
			path:   "/healthz",
			status: http.StatusOK,
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				if got := strings.TrimSpace(rec.Body.String()); got != "ok" {
					t.Fatalf("unexpected healthz body: got %q want %q", got, "ok")
				}
			},
		},
		{
			name:   "version",
			path:   "/version",
			status: http.StatusOK,
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				var payload map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
					t.Fatalf("failed to decode version JSON: %v", err)
				}
				if payload["version"] == "" {
					t.Fatalf("version field should not be empty")
				}
				if payload["build_time"] == "" {
					t.Fatalf("build_time field should not be empty")
				}
			},
		},
		{
			name:   "robots",
			path:   "/robots.txt",
			status: http.StatusOK,
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				body := rec.Body.String()
				if !strings.Contains(body, "User-agent: *") || !strings.Contains(body, "Allow: /") {
					t.Fatalf("unexpected robots body: %q", body)
				}
			},
		},
		{
			name:   "security",
			path:   "/.well-known/security.txt",
			status: http.StatusOK,
			check: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				body := rec.Body.String()
				if !strings.Contains(body, "Contact: mailto:") {
					t.Fatalf("unexpected security.txt body: %q", body)
				}
			},
		},
		{name: "not found", path: "/does-not-exist", status: http.StatusNotFound},
		{name: "project missing", path: "/projects/no-such-project", status: http.StatusNotFound},
		{name: "post missing", path: "/writing/no-such-post", status: http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			if rec.Code != tc.status {
				t.Fatalf("unexpected status for %s: got %d want %d", tc.path, rec.Code, tc.status)
			}
			if got := rec.Header().Get("X-Request-ID"); got == "" {
				t.Fatalf("missing X-Request-ID header for %s", tc.path)
			}
			if tc.check != nil {
				tc.check(t, rec)
			}
		})
	}
}
