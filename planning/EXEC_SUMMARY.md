# Executive Summary — cameronbrooks-site

Personal site build. Six planning documents cover the full scope. This summary is the read-once overview; the referenced docs are the implementation source of truth.

---

## Overview

A minimal, production-shaped personal site built on Go, HTMX, and a single Debian VPS. The driving constraint throughout is KISS: Go stdlib only (no third-party dependencies in v1), no containers, no orchestration, no build pipeline for the frontend. The architecture is intentionally small but structured to extend — the same code and infra pattern is the foundation for future subdomains and services without rewrites.

**Six pages:** home, projects (list + detail), writing (list + detail), about, contact.

**Key decisions made and locked:**

- Go 1.22+ stdlib routing — no router framework needed, method+pattern matching is built in
- `html/template` — no template engine dep
- HTMX vendored into `static/` — no external requests at load time
- `//go:embed` for templates and static files — single self-contained binary, no file sync on deploy
- No database — content is Go struct literals, redeploy to update
- Cloud-init only for infrastructure — no Terraform, no Ansible for a single VPS
- Caddy for TLS + reverse proxy — automatic HTTPS, zero config cert management
- Hand-rolled CSS, single file, no build step

---

## ARCHITECTURE.md — Architectural Invariants

**What it decides:** the fixed constraints of the system — load-bearing decisions that cascade if violated.

**The invariants:**

| #   | Invariant                                     | One-line summary                                                                     |
| --- | --------------------------------------------- | ------------------------------------------------------------------------------------ |
| 1   | Server owns all rendering                     | No client-side rendering; HTMX swaps server-rendered fragments only                  |
| 2   | One deployable binary                         | Templates + static assets embedded via `//go:embed`; no runtime file reads           |
| 3   | Caddy is the permanent public interface       | Go port never exposed; Caddy owns TLS, headers, CSP                                  |
| 4   | Content is trusted after the service boundary | `Body string` in storage; `template.HTML` only after explicit conversion in services |
| 5   | Progressive enhancement                       | Site is fully functional without JS; HTMX is additive                                |
| 6   | Handlers are thin                             | Logic lives in services; handlers parse input, call service, call render()           |
| 7   | One layout contract                           | All page content lives inside `<main id="main">`; never broken                       |
| 8   | No secrets in the repository                  | Runtime config via systemd env drop-in only                                          |
| 9   | No runtime filesystem dependency              | Binary is self-contained; no "forgot to sync assets" failure mode                    |

**Why this doc exists:** constraints are easy to enforce when written down. Without this list, future changes accumulate small exceptions that compound into drift. The doc also maps constraint interactions and lists the common first violations to watch for.

---

## STACK.md — Tech Stack & Setup

**What it decides:** every layer of the stack, project layout, deploy chain, and non-negotiable constraints.

**Stack summary:**

| Layer       | Choice                                          |
| ----------- | ----------------------------------------------- |
| Language    | Go 1.22+ (stdlib only)                          |
| Templates   | `html/template`                                 |
| Logging     | `log/slog` (stdlib, Go 1.21+)                   |
| Frontend JS | HTMX vendored to `static/htmx.min.js`           |
| CSS         | Single hand-rolled file, no build step          |
| Runtime     | Go binary → systemd service on `127.0.0.1:8080` |
| TLS + Proxy | Caddy (port 80/443)                             |
| IaC         | cloud-init YAML user-data (no Terraform)        |
| Deploy      | `make deploy` → scp binary → systemctl restart  |
| Source      | GitHub                                          |

**Project layout:**

```txt
cmd/site/           main.go, graceful shutdown
internal/handlers/  page handlers + render helper
internal/middleware/ request ID, structured logging per request
internal/services/  business logic; content lookup + view model construction
internal/views/     *.gohtml + views.go (embed)
internal/content/   data structs + lookup helpers
static/             css, js, images + static.go (embed)
deploy/             cloud-init.yaml, Caddyfile, site.service
Makefile
```

**Hard constraints:** Caddy always proxies to Go on `127.0.0.1:8080` — Go is never exposed directly. Secrets live in a systemd drop-in env file on the server, never in the repo.

---

## ROADMAP.md — Now vs Future

**What it decides:** scope boundary for v1 and what future expansion the current design must not block.

**V1 is done when:** site is live at domain with HTTPS, all six pages work with and without JS, binary runs as a systemd service with auto-restart, deploy is `make deploy`.

**V1 explicitly excludes:** database, markdown rendering, CI pipeline, containers, auth, analytics, contact form, RSS.

**Future items the current design accommodates without rewrites:**

- Markdown content: swap struct literals for startup-parsed `.md` files; nothing else changes
- CI deploys: GitHub Actions builds the binary and SCPs it; replaces the local build step only
- Subdomains: each new service is a new Go binary + new Caddy block, no changes to the existing app
- Auth, persistence (SQLite), contact form — all have clean insertion points in the current layout

**Intent:** avoid over-building now, avoid painting into corners. The structure is already correct for the expanded future.

---

## UI_UX.md — UI/UX Plan

**What it decides:** visual language, layout, page-by-page content patterns, HTMX navigation behavior, CSS architecture, and accessibility baseline.

**Visual language:**

- Colors: 6-token CSS custom property palette, light mode only in v1. Accent is `#0f766e` (restrained teal; dark mode tokens defined now, rendering deferred to post-v1).
- Typography: system font stack only — no web fonts, no CLS. 6-step type scale in rem, base 16px.
- Spacing: 8px base unit, 7 named tokens covering all real use cases (`--space-1` through `--space-10`).
- Max-widths: `680px` for reading content, `900px` for site container.

**Layout:** single column throughout. Sticky nav (not fixed — no body padding hack, no anchor offset bugs). One `<h1>` per page. No sidebar, no grid.

**Page patterns:** all six pages specified — what elements appear, what's the content hierarchy, what's deferred. Home shows 2–3 featured projects and 3–5 recent posts. Detail pages are `<article>` elements at reading width.

**HTMX navigation:** `hx-boost` on `<body>`, swaps `<main id="main">` only, pushes URL history. One template per page — handlers detect `HX-Request` header and render full page or `content` block only. Minimal CSS progress bar for loading state, no library.

**CSS architecture:** single `main.css`, six ordered layers (tokens → reset → base → layout → components → utilities). Target under 300 lines for v1.

**Accessibility:** semantic HTML, one `<h1>` per page, skip-to-content link, visible focus styles matching accent color, WCAG AA contrast minimum.

---

## CONTENT.md — Content Model

**What it decides:** the exact Go data structures for all site content, storage strategy, lookup helpers, and the migration path when markdown is needed.

**Two types:**

`Project` — `Slug`, `Title`, `Description` (one sentence, for list view), `Body` (`string`, raw HTML stored in content struct), `Tags []string`, `Date` (`time.Time`, enables sorting and feeds), `Links []Link` (source/demo/etc.), `Featured bool` (home page inclusion).

`Post` — `Slug`, `Title`, `Summary` (one sentence, for list view), `Body` (`string`, raw HTML), `Tags []string`, `Date` (`time.Time`), `Published bool` (false = draft, never exposed via any route).

**Trust boundary:** `Body` is stored as `string`. The explicit conversion to `template.HTML` — which disables Go's auto-escaping — happens in `internal/services/` at the service boundary, not in the storage struct. This keeps XSS guarantees intact and the trust grant auditable.

**Storage in v1:** Go slice literals in `internal/content/data.go`. Changing content means changing code and redeploying. Accepted trade-off for a low-volume personal site with no editorial workflow.

**Lookup helpers defined:** `ProjectBySlug`, `FeaturedProjects`, `PublishedPosts`, `PostBySlug` — these are the only surfaces handlers call.

**Migration path:** when markdown is wanted, populate the same `[]Project` and `[]Post` slices from `//go:embed`-loaded `.md` files parsed at startup. Handler and template signatures are unchanged.

---

## TEMPLATES.md — Template Structure

**What it decides:** template file naming and location, the layout/content block contract, the `PageData` struct, `//go:embed`-backed template cache, the `render()` helper that implements full-vs-partial HTMX logic, Go 1.22 routing table, active nav approach, and error page handling.

**Key contract:** every page template defines `{{define "content"}}...{{end}}`. The layout template calls `{{template "content" .}}` inside `<main id="main">`. Full-page render executes `"layout"`; HTMX partial render executes `"content"` only. One template file per page, two code paths, no duplication.

**`PageData`** passed to every execution:

```go
type PageData struct {
    Title       string
    Description string
    Year        string  // injected automatically by render()
    ActivePath  string  // for nav highlight
    Data        any     // page-specific payload
}
```

**Template cache:** `InitTemplates()` called once at startup, reads `views.FS` (embed), parses `layout.gohtml` + each page file into two maps (`tmplFull`, `tmplPart`). Fatal exit on any parse error — missing template is not a runtime condition.

**`render()` helper:** sets `Year`, sets `Content-Type`, checks `HX-Request` header, executes correct template map. Three lines of logic; all handlers call this one function.

**Routing:** Go 1.22 `net/http` pattern matching — nine routes (seven page routes, `/healthz`, `/version`) plus static file serving from `http.FS(static.FS)`. No third-party router.

---

## RUNBOOK.md — Operational Runbook

**What it decides:** how to run locally, how to build, step-by-step first VPS bring-up (including DNS), ongoing deploy process, log access, secrets management, and full re-provisioning procedure.

**Local dev:** `go run ./cmd/site` or `make dev`. No live reload. Two optional env vars (`SITE_ADDR`, `SITE_ENV`). No secrets needed in v1.

**Build:** `make build` → `GOOS=linux GOARCH=amd64 go build` → `bin/site`. Binary is self-contained via `//go:embed`; nothing else needs to be transferred.

**First VPS bring-up (ordered):**

1. Create Debian 12 VPS at chosen provider (Hetzner/Vultr/DO), paste `deploy/cloud-init.yaml` as user-data
2. Wait ~3 minutes for cloud-init to complete; verify with `cloud-init status` and `systemctl status site`
3. Point DNS A records (`@` and `www`) to VPS IP at registrar; TTL 300 initially
4. Run `make deploy` once DNS resolves
5. Verify HTTPS — Caddy issues the Let's Encrypt cert on first request (~5–10s delay on first load)

**Ongoing deploy:** `make deploy` — build locally, scp binary, restart service. ~15 seconds. Existing connections drain gracefully.

**Secrets:** stored in `/etc/systemd/system/site.service.d/env.conf` on the server. Edit on server, `daemon-reload`, restart. Never in repo.

**Re-provisioning from scratch:** create new VPS with same cloud-init → add secrets → `make deploy` → update DNS. ~10 minutes total.

---

## Decisions still open

One item not yet resolved, noted for implementation:

**Nav link label style:** the nav is `--text-sm` but whether labels are uppercase or sentence-case with semibold weight is not locked. This is a feel decision to make when the CSS is first rendered. It does not affect any other doc or any code.
