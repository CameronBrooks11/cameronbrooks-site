# Template Structure

Defines the Go template file layout, naming conventions, PageData struct, HTMX full-vs-partial execution model, and startup cache strategy. This doc bridges STACK.md and CONTENT.md — read both first.

---

## File layout

```txt
internal/views/
  layout.gohtml          full page shell (html, head, header, nav, footer)
  home.gohtml             /
  projects.gohtml         /projects
  project.gohtml          /projects/:slug
  writing.gohtml          /writing
  post.gohtml             /writing/:slug
  about.gohtml            /about
  contact.gohtml          /contact
```

All templates use `.gohtml` extension. Go's `html/template` auto-escapes, giving XSS protection by default.

---

## Naming contract

Every page template defines a single named block:

```html
{{define "content"}} ... page-specific HTML ... {{end}}
```

The layout template references this block:

```html
{{define "layout"}} ...
<main id="main">{{template "content" .}}</main>
... {{end}}
```

This means every full-page render is `Execute("layout", data)` and every HTMX partial render is `Execute("content", data)`. One source of HTML, two render paths.

---

## layout.gohtml

```html
{{define "layout"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}}{{if .Title}} — {{end}}Cameron Brooks</title>
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
          <li><a href="/projects">Projects</a></li>
          <li><a href="/writing">Writing</a></li>
          <li><a href="/about">About</a></li>
          <li><a href="/contact">Contact</a></li>
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

Note: `{{.Year}}` is provided by `PageData` (current year as string) so the footer copyright stays current without editing templates.

---

## PageData struct

```go
// internal/handlers/render.go

package handlers

import "html/template"

// PageData is passed as the root data object to every template execution.
type PageData struct {
    Title       string // used in <title> and <h1>; empty on home page
    Description string // used in <meta name="description">
    Year        string // current year, for footer copyright; injected by render()
    ActivePath  string // current URL path; used for nav active state
    Data        any    // page-specific payload (see per-page types below)
}
```

### Per-page Data types

| Page           | `Data` type                                                   |
| -------------- | ------------------------------------------------------------- |
| Home           | `HomeData{Featured []content.Project, Recent []content.Post}` |
| Projects list  | `[]content.Project`                                           |
| Project detail | `content.Project`                                             |
| Writing list   | `[]content.Post`                                              |
| Post detail    | `content.Post`                                                |
| About          | `nil`                                                         |
| Contact        | `nil`                                                         |

Define `HomeData` in `internal/handlers/home.go` since it is only used by that handler.

---

## Template cache

Templates are parsed **once at startup** from embedded files and stored in a map. Using `//go:embed` makes the binary self-contained — no template or static files need to be shipped alongside it on the VPS.

```go
// internal/views/views.go

package views

import "embed"

//go:embed *.gohtml
var FS embed.FS
```

```go
// static/static.go

package static

import "embed"

//go:embed css js images favicon.ico htmx.min.js
var FS embed.FS
```

```go
// internal/handlers/render.go

import (
    "html/template"
    "net/http"
    "strconv"
    "time"

    "github.com/CameronBrooks11/cameronbrooks-site/internal/views"
)

// tmplFull: page name → template parsed with layout (for full-page renders)
// tmplPart: page name → template parsed alone (for HTMX partial renders)
var (
    tmplFull map[string]*template.Template
    tmplPart map[string]*template.Template
)

var pages = []string{
    "home", "projects", "project", "writing", "post", "about", "contact",
    "notFound", "error",
}

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

        full, err := template.New("layout").Parse(string(layoutSrc))
        if err != nil {
            return err
        }
        if _, err = full.New(name).Parse(string(pageSrc)); err != nil {
            return err
        }
        tmplFull[name] = full

        part, err := template.New(name).Parse(string(pageSrc))
        if err != nil {
            return err
        }
        tmplPart[name] = part
    }
    return nil
}
```

Call `handlers.InitTemplates()` in `main.go` before starting the HTTP server. If it errors, fatal-exit immediately — a missing template is not a recoverable runtime error.

---

## render() helper

```go
// internal/handlers/render.go

func render(w http.ResponseWriter, r *http.Request, page string, data PageData) {
    data.Year = strconv.Itoa(time.Now().Year())
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if r.Header.Get("HX-Request") == "true" {
        if err := tmplPart[page].ExecuteTemplate(w, "content", data); err != nil {
            http.Error(w, "render error", http.StatusInternalServerError)
        }
        return
    }

    if err := tmplFull[page].ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, "render error", http.StatusInternalServerError)
    }
}
```

### Usage in a handler

```go
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
    render(w, r, "projects", PageData{
        Title:       "Projects",
        Description: "A selection of things I have built.",
        ActivePath:  "/projects",
        Data:        services.GetProjects(),
    })
}
```

---

## Routing (Go 1.22+ stdlib)

```go
// cmd/site/main.go

mux := http.NewServeMux()
// page routes
mux.HandleFunc("GET /", h.Home)
mux.HandleFunc("GET /projects", h.Projects)
mux.HandleFunc("GET /projects/{slug}", h.Project)
mux.HandleFunc("GET /writing", h.Writing)
mux.HandleFunc("GET /writing/{slug}", h.Post)
mux.HandleFunc("GET /about", h.About)
mux.HandleFunc("GET /contact", h.Contact)
// system routes
mux.HandleFunc("GET /healthz", h.Healthz)
mux.HandleFunc("GET /version", h.Version)
// static assets (embedded FS)
mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))
```

`/healthz` returns `200 OK` with body `ok`. `/version` returns the git SHA and build time injected at compile time via:

```sh
-ldflags="-X main.Version=$(git rev-parse --short HEAD) -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

Slug is extracted in the handler with `r.PathValue("slug")`.

---

## Active nav link

The current page should highlight the matching nav link. Pass an `ActivePath` field in `PageData` (e.g. `"/projects"`) and check it in the template:

```html
<a href="/projects" {{if eq .ActivePath "/projects"}}class="active"{{end}}>Projects</a>
```

Add `ActivePath string` to `PageData` and set it in each handler. Alternatively, use a JS-based approach on the client (`document.location.pathname`) to avoid the server-side field — either is fine, server-side is more correct for the no-JS requirement.

---

## 404 and error pages

Define a `notFound.gohtml` and `error.gohtml` in `internal/views/`. These use the same layout contract.

```go
func notFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    render(w, r, "notFound", PageData{Title: "Not Found"})
}
```

Do not register a separate `GET /` catch-all for `notFound` if `h.Home` is already registered on `GET /`. Keep `h.Home` on `GET /` and handle the catch-all behavior inside that handler by checking `r.URL.Path == "/"`; for all other paths, call `notFound(w, r)`.

---

## main.go wiring

```go
// cmd/site/main.go

func main() {
    if err := handlers.InitTemplates(); err != nil {
        log.Fatal("templates:", err)
    }

    // mux is constructed and routes registered as shown in the Routing section above.
    h := &http.Server{
        Addr:         os.Getenv("SITE_ADDR"), // defaults to ":8080"
        Handler:      middleware.Chain(mux, middleware.RequestID, middleware.Logger),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    if h.Addr == "" {
        h.Addr = ":8080"
    }

    // Graceful shutdown on SIGTERM / SIGINT
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    go func() {
        slog.Info("starting server", "addr", h.Addr)
        if err := h.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            slog.Error("server error", "err", err)
            os.Exit(1)
        }
    }()

    <-ctx.Done()
    slog.Info("shutting down")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _ = h.Shutdown(shutdownCtx)
}
```

Systemd sends `SIGTERM` on `systemctl restart`. The 10-second shutdown context lets in-flight requests complete. Set `TimeoutStopSec=15` in `site.service` to match.

---

## Middleware

Defined in `internal/middleware/`. Applied as a chain wrapping the mux in `main.go`.

### Chain helper

```go
// internal/middleware/middleware.go

type Middleware func(http.Handler) http.Handler

// Chain applies middlewares right-to-left so the first argument is the outermost.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
```

### Request ID

```go
// internal/middleware/requestid.go

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            b := make([]byte, 8)
            _, _ = rand.Read(b) // crypto/rand — no dep
            id = hex.EncodeToString(b)
        }
        w.Header().Set("X-Request-ID", id)
        r = r.WithContext(context.WithValue(r.Context(), requestIDKey, id))
        next.ServeHTTP(w, r)
    })
}
```

### Logger

```go
// internal/middleware/logger.go

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rec, r)
        slog.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", rec.status,
            "duration_ms", time.Since(start).Milliseconds(),
            "request_id", RequestIDFrom(r.Context()),
            "remote_ip", realIP(r), // checks X-Forwarded-For (set by Caddy), fallback r.RemoteAddr
        )
    })
}
```

`statusRecorder` is a thin `http.ResponseWriter` wrapper that captures the written status code.
