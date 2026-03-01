# Phase 02 — Embed Infrastructure & Skeleton Binary

**Goal:** Wire the two `//go:embed` packages (`internal/views` and `static`) so the build system knows about embedded assets from day one, create stub files to satisfy compile-time embed requirements, and get a minimal HTTP server binary that compiles and starts without crashing.

**Exit gate:** `go build ./...` exits 0; `go run ./cmd/site` starts on `:8080` without crashing; Ctrl-C stops it cleanly.

---

## Why embed first

`//go:embed *.gohtml` is a **compile-time directive**. If the pattern matches zero files, `go build` fails with an error. This means stub template files must exist before any Go file that declares the embed runs through the compiler — even if those templates are empty shells.

Establishing the embed packages now means every subsequent phase compiles against the real binary shape. It also proves the single-binary invariant from day one: after this phase, `go build` produces a binary that carries its own assets.

---

## Files to create in this phase

```
internal/views/views.go          — embed declaration for templates
internal/views/layout.gohtml     — stub layout template (satisfies //go:embed *.gohtml)
static/static.go                 — embed declaration for static assets
static/css/main.css              — stub (empty CSS file; filled in Phase 09)
static/js/progress.js            — stub (filled in Phase 10)
static/images/placeholder.txt   — satisfies //go:embed images (removed when first real image added)
static/favicon.ico               — empty placeholder (replace with real favicon before go-live)
static/htmx.min.js               — stub comment (real file vendored in Phase 10)
cmd/site/main.go                 — minimal HTTP server
```

---

## Step 1 — `internal/views/views.go`

This package owns exactly one thing: the embedded filesystem of all `.gohtml` template files. Nothing else belongs here.

**File: `internal/views/views.go`**

```go
package views

import "embed"

// FS holds all .gohtml template files embedded at compile time.
// Parsed once at startup by handlers.InitTemplates().
//
//go:embed *.gohtml
var FS embed.FS
```

The `//go:embed *.gohtml` directive matches every `.gohtml` file in the `internal/views/` directory at compile time. Go will error if the pattern matches no files — that is why Step 2 creates a stub immediately.

---

## Step 2 — `internal/views/layout.gohtml` (stub)

A minimal but syntactically valid Go template. It defines the `"layout"` block that `InitTemplates()` will execute for full-page renders (Phase 05). The stub is intentionally bare — the full shell (head, nav, footer, HTMX attributes) is written in Phase 05.

**File: `internal/views/layout.gohtml`**

```html
{{define "layout"}}<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>{{.Title}}</title>
  </head>
  <body>
    <main id="main">{{template "content" .}}</main>
  </body>
</html>
{{end}}
```

**Why even this much content in a stub?** The template cache built in Phase 05 calls `ExecuteTemplate(w, "layout", data)` and `ExecuteTemplate(w, "content", data)`. If `"layout"` is not defined in the parsed template set, those calls panic at startup. Defining the blocks here means Phase 05 can override the _content_ of those blocks without changing the calling convention.

> Do not remove or rename the `{{define "layout"}}` and `{{template "content" .}}` references — they are the contract that every phase from 05 onward depends on.

---

## Step 3 — `static/static.go`

Mirrors the views package: one file, one embed declaration.

**File: `static/static.go`**

```go
package static

import "embed"

// FS holds all static assets embedded at compile time.
// Served via http.FileServer(http.FS(static.FS)) at /static/.
//
//go:embed css js images favicon.ico htmx.min.js
var FS embed.FS
```

The pattern `css js images` embeds entire directories recursively. `favicon.ico` and `htmx.min.js` are explicit files. All six paths must exist on disk at compile time — Steps 4 and 5 satisfy this requirement.

**Important:** Go's embed excludes files whose names begin with `.` or `_`. Do not use `.gitkeep` for placeholder files in directories referenced by `//go:embed` — they will be silently excluded and the directory will appear empty to the embedded FS. Use named placeholder files instead.

---

## Step 4 — Stub static files

Each file must exist so the embed directive in `static.go` finds something to embed. Comments inside the files explain their Phase 02 status.

**`static/css/main.css`** — empty placeholder

```css
/* main.css — filled in Phase 09 */
```

**`static/js/progress.js`** — empty placeholder

```js
// progress.js — HTMX progress bar; filled in Phase 10
```

**`static/images/placeholder.txt`** — satisfies `//go:embed images` requirement

```
# placeholder — delete when the first real image file is added to this directory
```

**`static/favicon.ico`** — create as an empty file for now

```sh
# PowerShell
New-Item -ItemType File -Path "static/favicon.ico" -Force | Out-Null
```

Replace with a real `.ico` file before go-live. An empty file will not crash the server — the browser just gets a 0-byte response for `/favicon.ico` requests, which is fine during development.

**`static/htmx.min.js`** — stub comment only; real file vendored in Phase 10

```js
// htmx.min.js — placeholder; replace with vendored htmx before Phase 10
// Download from: https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js
```

> **Do not download htmx yet.** The real vendoring step (pinned version, verified file) is in Phase 10. Using a stub here avoids premature version decisions and keeps this phase narrowly scoped.

---

## Step 5 — `cmd/site/main.go` (minimal)

The entire application will be wired here in Phase 08. For now: compile, start, serve nothing, stop cleanly.

**File: `cmd/site/main.go`**

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
```

**Why `log.Printf` instead of `slog` here?** Phase 08 replaces this file entirely with the production wiring including `slog`, configured log level, `http.Server` with timeouts, graceful shutdown, and middleware chain. Using `log.Printf` in the stub avoids importing `log/slog` here and immediately needing to initialize it — keeps the noise low so Phase 02's job (compile + start) is clearly complete.

---

## Step 6 — Verify the build

```sh
# From repo root
go build ./...
```

Expected: exits 0, no output. This compiles all packages including `internal/views` and `static` — the embed directives are exercised even though `cmd/site/main.go` does not import those packages yet. This is the key verification for this phase.

If `go build ./...` errors with something like:

```
pattern *.gohtml: no matching files found
```

→ `layout.gohtml` is missing or in the wrong directory. Verify it is at `internal/views/layout.gohtml`.

If it errors with:

```
pattern images: no matching files found
```

→ `static/images/placeholder.txt` is missing or the file name starts with `.`/`_`.

---

## Step 7 — Verify the server starts

```sh
go run ./cmd/site
```

Expected output:

```
2026/02/28 12:00:00 starting server on :8080
```

Server should stay running. Open `http://localhost:8080` in a browser:

- Returns `404 page not found` from Go's default ServeMux — this is correct. No routes are registered yet.

Press Ctrl-C to stop. Server should exit immediately (no graceful shutdown yet — that is Phase 08).

---

## Step 8 — Commit

```sh
git add .
git commit -m "phase 02: embed infrastructure and skeleton binary"
```

---

## Exit gate checklist

- [ ] `go build ./...` exits 0 with no errors
- [ ] `internal/views/views.go` exists with `//go:embed *.gohtml` directive
- [ ] `internal/views/layout.gohtml` defines both `{{define "layout"}}` and calls `{{template "content" .}}`
- [ ] `static/static.go` exists with `//go:embed css js images favicon.ico htmx.min.js` directive
- [ ] All six embed targets exist on disk: `static/css/main.css`, `static/js/progress.js`, `static/images/placeholder.txt`, `static/favicon.ico`, `static/htmx.min.js`
- [ ] None of the placeholder files in `static/images/` start with `.` or `_`
- [ ] `go run ./cmd/site` starts on `:8080`, serves 404 on any path, stops on Ctrl-C
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 03.
