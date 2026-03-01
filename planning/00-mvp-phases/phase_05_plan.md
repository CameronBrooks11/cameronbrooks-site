# Phase 05 — Template System

**Goal:** Write the full `layout.gohtml` (replacing the Phase 02 stub), all nine page templates, and `internal/handlers/render.go` containing `PageData`, `InitTemplates()`, and `render()`. By the end of this phase the template cache initializes without error and every page key resolves correctly.

**Exit gate:** `InitTemplates()` returns nil; all nine keys (`home`, `projects`, `project`, `writing`, `post`, `about`, `contact`, `notFound`, `error`) are present in both `tmplFull` and `tmplPart`; `go vet ./internal/handlers/...` passes; `go build ./...` exits 0.

---

## Prerequisites

- Phase 04 complete (service view models exist — `ProjectView`, `PostView`)
- `internal/views/layout.gohtml` is the Phase 02 stub — it will be replaced entirely in this phase
- `internal/handlers/` directory exists; `.gitkeep` can be deleted

---

## Files to create / replace in this phase

```
internal/views/layout.gohtml         — REPLACE Phase 02 stub with full shell
internal/views/home.gohtml           — new
internal/views/projects.gohtml       — new
internal/views/project.gohtml        — new
internal/views/writing.gohtml        — new
internal/views/post.gohtml           — new
internal/views/about.gohtml          — new
internal/views/contact.gohtml        — new
internal/views/notFound.gohtml       — new
internal/views/error.gohtml          — new
internal/handlers/render.go          — new: PageData, per-page types, InitTemplates, render()
```

---

## Step 1 — Delete `.gitkeep` from `internal/handlers/`

```sh
Remove-Item internal/handlers/.gitkeep
```

---

## Step 2 — Replace `internal/views/layout.gohtml`

This replaces the minimal stub from Phase 02 with the full page shell. Every element documented in TEMPLATES.md and UI_UX.md is present. CSS and HTMX are not functional yet (Phase 09 and 10) but the references are correct.

**File: `internal/views/layout.gohtml`**

```html
{{define "layout"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{if .Title}}{{.Title}} — {{end}}Cameron Brooks</title>
    <meta name="description" content="{{.Description}}" />
    <link rel="stylesheet" href="/static/css/main.css" />
  </head>
  <body hx-boost="true" hx-target="#main" hx-select="#main" hx-push-url="true">
    <a href="#main" class="sr-only skip-link">Skip to content</a>

    <div id="progress-bar" aria-hidden="true"></div>

    <header>
      <nav class="container">
        <a href="/" class="nav-home">Cameron Brooks</a>
        <ul class="nav-links">
          <li><a href="/projects"{{if eq .ActivePath "/projects"}} class="active"{{end}}>Projects</a></li>
          <li><a href="/writing"{{if eq .ActivePath "/writing"}} class="active"{{end}}>Writing</a></li>
          <li><a href="/about"{{if eq .ActivePath "/about"}} class="active"{{end}}>About</a></li>
          <li><a href="/contact"{{if eq .ActivePath "/contact"}} class="active"{{end}}>Contact</a></li>
        </ul>
      </nav>
    </header>

    <main id="main">{{template "content" .}}</main>

    <footer class="container">
      <p>Cameron Brooks &middot; {{.Year}}</p>
    </footer>

    <script src="/static/htmx.min.js"></script>
    <script src="/static/js/progress.js"></script>
  </body>
</html>
{{end}}
```

**Key details:**

- `{{if .Title}}{{.Title}} — {{end}}` — home page has no `.Title` so the title renders as just `"Cameron Brooks"` without a stray `" — "` prefix
- `ActivePath` is compared server-side so the active nav state is correct even without JavaScript
- HTMX attributes (`hx-boost`, `hx-target`, `hx-select`, `hx-push-url`) are on `<body>` — they do nothing until `htmx.min.js` is present (Phase 10); the page degrades gracefully
- `#progress-bar` div is always in the DOM so the CSS transition in Phase 09 can attach immediately

---

## Step 3 — Page templates

Each template defines exactly one named block: `{{define "content"}}...{{end}}`. Nothing else. The layout shell is never duplicated.

Templates reference `{{.Data}}` via the typed value passed from the handler. The `any` type in `PageData.Data` means templates must use the concrete type's fields directly — Go's `html/template` resolves fields dynamically at execution time.

---

### `internal/views/home.gohtml`

```html
{{define "content"}}
<div class="container">
  <section class="home-intro">
    <p class="text-lg">
      Hi, I'm Cameron — a software engineer. I build tools, write about what I'm
      learning, and share the work here.
    </p>
  </section>

  {{if .Data.Featured}}
  <section class="home-section">
    <h2>Selected Projects</h2>
    <ul class="post-list">
      {{range .Data.Featured}}
      <li class="card">
        <a href="/projects/{{.Slug}}">{{.Title}}</a>
        <p class="text-muted">{{.Description}}</p>
        {{if .Tags}}
        <p class="tag-list">
          {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
        </p>
        {{end}}
      </li>
      {{end}}
    </ul>
  </section>
  {{end}} {{if .Data.Recent}}
  <section class="home-section">
    <h2>Recent Writing</h2>
    <ul class="post-list">
      {{range .Data.Recent}}
      <li>
        <a href="/writing/{{.Slug}}">{{.Title}}</a>
        <span class="text-muted"> — {{.Date}}</span>
      </li>
      {{end}}
    </ul>
  </section>
  {{end}}
</div>
{{end}}
```

---

### `internal/views/projects.gohtml`

```html
{{define "content"}}
<div class="container">
  <h1>Projects</h1>
  <p class="text-muted">A selection of things I have built.</p>

  <ul class="post-list">
    {{range .Data}}
    <li class="card">
      <a href="/projects/{{.Slug}}">{{.Title}}</a>
      <span class="text-muted"> &mdash; {{.Date}}</span>
      <p>{{.Description}}</p>
      {{if .Tags}}
      <p class="tag-list">
        {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
      </p>
      {{end}}
    </li>
    {{else}}
    <li><p class="text-muted">No projects yet.</p></li>
    {{end}}
  </ul>
</div>
{{end}}
```

---

### `internal/views/project.gohtml`

```html
{{define "content"}}
<div class="container">
  <a href="/projects" class="back-link">&larr; Projects</a>

  <article>
    <header class="article-header">
      <h1>{{.Data.Title}}</h1>
      <p class="text-muted">
        {{.Data.Date}} {{if .Data.Tags}}&mdash; {{range .Data.Tags}}<span
          class="tag"
          >{{.}}</span
        >
        {{end}}{{end}}
      </p>
      {{if .Data.Links}}
      <p class="article-links">
        {{range .Data.Links}}<a
          href="{{.URL}}"
          target="_blank"
          rel="noopener noreferrer"
          >{{.Label}}</a
        >
        {{end}}
      </p>
      {{end}}
    </header>

    <div class="article-body">{{.Data.Body}}</div>
  </article>

  <a href="/projects" class="back-link">&larr; Projects</a>
</div>
{{end}}
```

---

### `internal/views/writing.gohtml`

```html
{{define "content"}}
<div class="container">
  <h1>Writing</h1>
  <p class="text-muted">Notes and longer pieces.</p>

  <ul class="post-list">
    {{range .Data}}
    <li>
      <a href="/writing/{{.Slug}}">{{.Title}}</a>
      <span class="text-muted"> &mdash; {{.Date}}</span>
      <p>{{.Summary}}</p>
    </li>
    {{else}}
    <li><p class="text-muted">Nothing published yet.</p></li>
    {{end}}
  </ul>
</div>
{{end}}
```

---

### `internal/views/post.gohtml`

```html
{{define "content"}}
<div class="container">
  <a href="/writing" class="back-link">&larr; Writing</a>

  <article>
    <header class="article-header">
      <h1>{{.Data.Title}}</h1>
      <p class="text-muted">
        {{.Data.Date}} {{if .Data.Tags}}&mdash; {{range .Data.Tags}}<span
          class="tag"
          >{{.}}</span
        >
        {{end}}{{end}}
      </p>
    </header>

    <div class="article-body">{{.Data.Body}}</div>
  </article>

  <a href="/writing" class="back-link">&larr; Writing</a>
</div>
{{end}}
```

---

### `internal/views/about.gohtml`

```html
{{define "content"}}
<div class="container">
  <article>
    <h1>About</h1>
    <p>
      I'm Cameron Brooks — a software engineer based in [location]. I work on
      [brief description of work]. This site is where I share projects and write
      about what I'm learning.
    </p>
    <p>The best way to reach me is by <a href="/contact">email</a>.</p>
  </article>
</div>
{{end}}
```

> Replace the bracketed placeholders with real copy before go-live.

---

### `internal/views/contact.gohtml`

```html
{{define "content"}}
<div class="container">
  <article>
    <h1>Contact</h1>
    <p>The best way to reach me:</p>
    <ul>
      <li><a href="mailto:cameron@example.com">cameron@example.com</a></li>
      <li>
        <a
          href="https://github.com/CameronBrooks11"
          target="_blank"
          rel="noopener noreferrer"
          >GitHub</a
        >
      </li>
    </ul>
  </article>
</div>
{{end}}
```

> Replace the email address before go-live.

---

### `internal/views/notFound.gohtml`

```html
{{define "content"}}
<div class="container">
  <article>
    <h1>404 — Not Found</h1>
    <p class="text-muted">The page you're looking for doesn't exist.</p>
    <p><a href="/">Go home</a></p>
  </article>
</div>
{{end}}
```

---

### `internal/views/error.gohtml`

```html
{{define "content"}}
<div class="container">
  <article>
    <h1>Something went wrong</h1>
    <p class="text-muted">An internal error occurred. Try again in a moment.</p>
    <p><a href="/">Go home</a></p>
  </article>
</div>
{{end}}
```

---

## Step 4 — `internal/handlers/render.go`

This file owns three things: `PageData` (the universal template data struct), `InitTemplates()` (startup cache builder), and `render()` (HTMX-aware execution helper). Nothing else belongs here.

**File: `internal/handlers/render.go`**

```go
package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
	"github.com/CameronBrooks11/cameronbrooks-site/internal/views"
)

// PageData is passed as the root data object to every template execution.
type PageData struct {
	Title       string // used in <title> and <h1>; leave empty on home page
	Description string // used in <meta name="description">
	Year        string // current year for footer copyright; injected automatically by render()
	ActivePath  string // current URL path, e.g. "/projects"; used for nav active state
	Data        any    // page-specific payload; see per-page types defined below
}

// HomeData is the payload for the home page handler.
// Defined here because it is only used by handlers.Home.
type HomeData struct {
	Featured []services.ProjectView
	Recent   []services.PostView
}

// tmplFull maps page name → template parsed with layout (used for full-page renders).
// tmplPart maps page name → template parsed alone (used for HTMX partial renders).
var (
	tmplFull map[string]*template.Template
	tmplPart map[string]*template.Template
)

// pages lists every template name that must be present in the cache.
// Each name corresponds to a <name>.gohtml file in internal/views/.
var pages = []string{
	"home", "projects", "project", "writing", "post",
	"about", "contact", "notFound", "error",
}

// InitTemplates parses all page templates at startup and populates tmplFull and tmplPart.
// Must be called once from main() before the HTTP server starts.
// Returns an error if any template file is missing or unparseable — treat as fatal.
func InitTemplates() error {
	tmplFull = make(map[string]*template.Template)
	tmplPart = make(map[string]*template.Template)

	layoutSrc, err := views.FS.ReadFile("layout.gohtml")
	if err != nil {
		return err
	}

	for _, name := range pages {
		pageSrc, err := views.FS.ReadFile(name + ".gohtml")
		if err != nil {
			return err
		}

		// Full: parse layout first, then associate the page template with it.
		full, err := template.New("layout").Parse(string(layoutSrc))
		if err != nil {
			return err
		}
		if _, err = full.New(name).Parse(string(pageSrc)); err != nil {
			return err
		}
		tmplFull[name] = full

		// Partial: parse the page template alone (no layout wrapper).
		part, err := template.New(name).Parse(string(pageSrc))
		if err != nil {
			return err
		}
		tmplPart[name] = part
	}
	return nil
}

// render executes the named template and writes the response.
// status is the HTTP status code for full-page renders (e.g. http.StatusOK, http.StatusNotFound).
// If the request carries an HX-Request header (HTMX navigation), only the "content" block is
// executed and status is always implicitly 200 — HTMX only swaps on 2xx, so error pages are
// intentionally shown inline without changing the HTTP status.
// Year is injected automatically; callers must set Title, Description, ActivePath, and Data.
func render(w http.ResponseWriter, r *http.Request, page string, status int, data PageData) {
	data.Year = strconv.Itoa(time.Now().Year())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		// HTMX partial: always implicit 200 so the fragment is swapped into #main.
		if err := tmplPart[page].ExecuteTemplate(w, "content", data); err != nil {
			slog.Error("partial render failed", "page", page, "err", err)
		}
		return
	}

	w.WriteHeader(status)
	if err := tmplFull[page].ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("full render failed", "page", page, "err", err)
	}
}
```

---

## Step 5 — Verify

```sh
go vet ./internal/handlers/...
```

Expected: exits 0.

```sh
go build ./...
```

Expected: exits 0. `internal/handlers` is now compiled but not yet imported by `cmd/site/main.go` — that happens in Phase 08.

**Verify `InitTemplates()` succeeds** with a quick standalone check:

```sh
go run - <<'EOF'
package main

import (
    "fmt"
    "github.com/CameronBrooks11/cameronbrooks-site/internal/handlers"
)

func main() {
    if err := handlers.InitTemplates(); err != nil {
        fmt.Printf("FAIL: %v\n", err)
        return
    }
    fmt.Println("OK: InitTemplates succeeded")
}
EOF
```

Expected output: `OK: InitTemplates succeeded`

If it fails with `open notFound.gohtml: file does not exist` — a template file is missing or misspelled. Check that all nine `.gohtml` files exist in `internal/views/` and that the names match the `pages` slice exactly (case-sensitive).

---

## Step 6 — Commit

```sh
git add internal/views/ internal/handlers/render.go
git commit -m "phase 05: template system"
```

---

## Exit gate checklist

- [ ] `go vet ./internal/handlers/...` exits 0
- [ ] `go build ./...` exits 0
- [ ] `internal/views/` contains exactly 10 `.gohtml` files: `layout` + all 9 page names
- [ ] Every page template defines `{{define "content"}}` and `{{end}}`
- [ ] `layout.gohtml` defines `{{define "layout"}}`, references `{{template "content" .}}`, and includes the HTMX attributes on `<body>`
- [ ] `InitTemplates()` returns nil (standalone check above passes)
- [ ] Phase 02 stub `layout.gohtml` is fully replaced (no `<title>{{.Title}}</title>` stub line remains)
- [ ] `.gitkeep` removed from `internal/handlers/`
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 06.
