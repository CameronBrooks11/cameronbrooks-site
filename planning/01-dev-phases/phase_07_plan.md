# Phase 07 — Handlers

**Goal:** Implement every handler in `internal/handlers/` — one file per logical group, each calling a service function and delegating to `render()`. Also implement `Healthz`, `Version`, and the package-level `notFound` and `internalError` functions. Handlers must be thin: no business logic, no direct content imports, no template manipulation.

**Exit gate:** `go vet ./internal/handlers/...` passes; `go build ./...` exits 0; every exported handler calls `render()` with the correct `ActivePath`; `grep -r "content\." internal/handlers/` returns no matches (handlers do not import the content package).

---

## Prerequisites

- Phase 05 complete (`render()`, `PageData`, `InitTemplates()` exist)
- Phase 04 complete (service functions exist: `GetProjects`, `GetFeaturedProjects`, `GetProjectBySlug`, `GetPosts`, `GetRecentPosts`, `GetPostBySlug`)

---

## Files to create in this phase

```
internal/handlers/handler.go     — Handler struct, constructor
internal/handlers/home.go        — Home, HomeData
internal/handlers/projects.go    — Projects (list), Project (detail)
internal/handlers/writing.go     — Writing (list), Post (detail)
internal/handlers/static.go      — About, Contact
internal/handlers/system.go      — Healthz, Version
internal/handlers/errors.go      — notFound, internalError
```

`render.go` already exists from Phase 05. No changes to it in this phase.

---

## Handler design rules

These apply to every handler in this phase:

1. Parse input from `r` only (path values, query params, headers). No reading body for GET handlers.
2. Keep service calls minimal and explicit: most pages call exactly one service function (static pages call zero). Home is the intentional exception and may call two (`GetFeaturedProjects` and `GetRecentPosts`).
3. Call `render()` as the final action. Never write to `w` directly except in `Healthz` and `Version`.
4. On service "not found" return → call `notFound(w, r)` then `return`. Never 500 a missing slug.
5. `ActivePath` must match the nav link href exactly: `"/"`, `"/projects"`, `"/writing"`, `"/about"`, `"/contact"`. System handlers pass `""`.

---

## Step 1 — `internal/handlers/handler.go`

The `Handler` struct is empty for now — it exists so all handlers are methods on a concrete type, making them easy to pass to `mux.HandleFunc` in Phase 08 and easy to extend with injected dependencies (e.g. a future logger or config) without changing method signatures.

**File: `internal/handlers/handler.go`**

```go
package handlers

// Handler holds application-level dependencies shared across all page handlers.
// Currently empty — will hold config, logger refs, or service interfaces if needed later.
// All page handlers are methods on this type.
type Handler struct{}

// New returns an initialized Handler ready to register on a mux.
func New() *Handler {
	return &Handler{}
}
```

---

## Step 2 — `internal/handlers/home.go`

**File: `internal/handlers/home.go`**

```go
package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Home handles GET /.
// Renders the home page with featured projects and recent posts.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// The stdlib mux routes "/" as a catch-all for unmatched paths.
	// Any path that didn't match a more specific route ends up here.
	// Return 404 for everything except exactly "/".
	if r.URL.Path != "/" {
		notFound(w, r)
		return
	}

	render(w, r, "home", http.StatusOK, PageData{
		Description: "Cameron Brooks — software engineer. Projects and writing.",
		ActivePath:  "/",
		Data: HomeData{
			Featured: services.GetFeaturedProjects(),
			Recent:   services.GetRecentPosts(5),
		},
	})
}
```

---

## Step 3 — `internal/handlers/projects.go`

**File: `internal/handlers/projects.go`**

```go
package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Projects handles GET /projects.
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	render(w, r, "projects", http.StatusOK, PageData{
		Title:       "Projects",
		Description: "A selection of things I have built.",
		ActivePath:  "/projects",
		Data:        services.GetProjects(),
	})
}

// Project handles GET /projects/{slug}.
func (h *Handler) Project(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	project, ok := services.GetProjectBySlug(slug)
	if !ok {
		notFound(w, r)
		return
	}
	render(w, r, "project", http.StatusOK, PageData{
		Title:       project.Title,
		Description: project.Description,
		ActivePath:  "/projects",
		Data:        project,
	})
}
```

---

## Step 4 — `internal/handlers/writing.go`

**File: `internal/handlers/writing.go`**

```go
package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Writing handles GET /writing.
func (h *Handler) Writing(w http.ResponseWriter, r *http.Request) {
	render(w, r, "writing", http.StatusOK, PageData{
		Title:       "Writing",
		Description: "Notes and longer pieces.",
		ActivePath:  "/writing",
		Data:        services.GetPosts(),
	})
}

// Post handles GET /writing/{slug}.
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, ok := services.GetPostBySlug(slug)
	if !ok {
		notFound(w, r)
		return
	}
	render(w, r, "post", http.StatusOK, PageData{
		Title:       post.Title,
		Description: post.Summary,
		ActivePath:  "/writing",
		Data:        post,
	})
}
```

---

## Step 5 — `internal/handlers/static.go`

"Static" here means pages with no dynamic data — `Data` is nil.

**File: `internal/handlers/static.go`**

```go
package handlers

import "net/http"

// About handles GET /about.
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	render(w, r, "about", http.StatusOK, PageData{
		Title:       "About",
		Description: "A little about Cameron Brooks.",
		ActivePath:  "/about",
	})
}

// Contact handles GET /contact.
func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	render(w, r, "contact", http.StatusOK, PageData{
		Title:       "Contact",
		Description: "How to reach Cameron Brooks.",
		ActivePath:  "/contact",
	})
}
```

---

## Step 6 — `internal/handlers/system.go`

`Version` reads package-level vars (`Version`, `BuildTime`) injected by `-ldflags` at build time (wired in Phase 08). During `go run`, they hold their zero values (`""`); that is fine for development.

**File: `internal/handlers/system.go`**

```go
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
// During `go run` (development), both are empty strings — that is expected.
var (
	Version   string
	BuildTime string
)

// Healthz handles GET /healthz.
// Returns 200 OK with body "ok". Used by Caddy health checks and monitoring.
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
```

> **Note on ldflags path:** Phase 08 replaces this file entirely. Package-level `Version` and `BuildTime` vars here are temporary scaffolding. In Phase 08, they move to `Handler` struct fields named `AppVersion` and `AppBuildTime`, keeping the ldflags target at `main.Version` and `main.BuildTime`.

> **Simplification for Phase 08:** When wiring in Phase 08, add `AppVersion string` and `AppBuildTime string` fields to the `Handler` struct, set them from `main.go` vars, and read them in the `Version` handler method. The method name `Version` stays unchanged; the field names use the `App` prefix to avoid a field/method name collision.

---

## Step 7 — `internal/handlers/errors.go`

Package-level functions (not methods) because they are called from within other handlers and do not need a receiver.

**File: `internal/handlers/errors.go`**

```go
package handlers

import "net/http"

// notFound writes a 404 response using the notFound template.
// Call this from any handler when a requested resource does not exist.
func notFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "notFound", http.StatusNotFound, PageData{
		Title:       "Not Found",
		Description: "The page you're looking for doesn't exist.",
	})
}

// internalError writes a 500 response using the error template.
// Use sparingly — prefer returning specific errors or 404s where possible.
func internalError(w http.ResponseWriter, r *http.Request) {
	render(w, r, "error", http.StatusInternalServerError, PageData{
		Title:       "Error",
		Description: "An internal error occurred.",
	})
}
```

---

## Step 8 — Verify

```sh
go vet ./internal/handlers/...
```

Expected: exits 0.

```sh
go build ./...
```

Expected: exits 0.

**Trust boundary audit:**

```sh
# Handlers must not import the content package directly
grep -r '"github.com/CameronBrooks11/cameronbrooks-site/internal/content"' internal/handlers/
```

Expected: no output. All content access goes through `internal/services/`.

**ActivePath audit:**

```sh
grep -n "ActivePath" internal/handlers/*.go
```

Every handler that renders a nav-linked page must set `ActivePath`. Confirm that `Home`, `Projects`, `Project`, `Writing`, `Post`, `About`, and `Contact` all appear in the output with a non-empty value.

---

## Step 9 — Commit

```sh
git add internal/handlers/
git commit -m "phase 07: handlers"
```

---

## Exit gate checklist

- [ ] `go vet ./internal/handlers/...` exits 0
- [ ] `go build ./...` exits 0
- [ ] `grep ... internal/content ... internal/handlers/` returns no matches
- [ ] `ActivePath` set on all seven page handlers
- [ ] `Home` returns 404 for paths other than `"/"`
- [ ] `Project` and `Post` detail handlers return 404 for unknown slugs
- [ ] `Healthz` writes `ok` with status 200
- [ ] `Version` writes plaintext `version=dev build_time=unknown` during `go run` (no ldflags yet; JSON format is wired in Phase 08)
- [ ] `notFound` and `internalError` pass the correct HTTP status codes to `render()` (404 and 500 respectively)
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 08.
