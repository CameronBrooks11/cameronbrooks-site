# Tech Stack & Setup

## Stack

| Layer               | Choice                     | Notes                                                          |
| ------------------- | -------------------------- | -------------------------------------------------------------- |
| Language            | Go (stdlib only)           | `net/http`, `html/template`, `log/slog`, no third-party deps   |
| Logging             | `log/slog` (stdlib)        | Structured logging; Go 1.21+, within version requirement       |
| Frontend JS         | HTMX (vendored)            | Vendored to `static/htmx.min.js`; no external requests at load |
| CSS                 | Hand-rolled, no build step | Plain CSS file(s) in `static/`                                 |
| Server              | Go binary on Debian VPS    | Runs as non-root systemd service on `127.0.0.1:8080`           |
| TLS + reverse proxy | Caddy                      | Automatic HTTPS, HTTP→HTTPS redirect, proxies to Go            |
| Bootstrap / IaC     | cloud-init                 | YAML user-data file; paste into provider UI on VPS creation    |
| Ongoing deploy      | Makefile / shell script    | Build → `scp` binary → `systemctl restart site`                |
| Source              | GitHub                     | Canonical source; GitHub Actions CI for future artifact builds |

## Go version

**Go 1.22+** is required. The key reason: Go 1.22 introduced method+pattern routing in `net/http` (`GET /projects/{slug}`), making a third-party router unnecessary. Confirm with `go version` before starting.

## Dependency policy

Go stdlib only for the initial site. Add a dep only when the stdlib genuinely cannot do the job and the benefit is clear. Current status:

- `net/http` — routing, serving, method+pattern matching (Go 1.22+)
- `html/template` — server-rendered HTML
- `log/slog` — structured logging (stdlib, Go 1.21+)
- No ORM, no router framework, no template engine

## Code → Production chain

```txt
GitHub (source)
    │
    │  go build -o site ./cmd/site   (local or CI)
    ▼
linux/amd64 binary
    │
    │  scp / rsync  (or future: GH Actions artifact + SSH deploy)
    ▼
Debian VPS  (Vultr / Hetzner / DigitalOcean — provider-agnostic)
    │
    ├── systemd  →  runs ./site on 127.0.0.1:8080
    └── Caddy    →  443/80 → proxy → 127.0.0.1:8080
```

## VPS bootstrap (cloud-init)

A single `deploy/cloud-init.yaml` file covers first-boot setup:

- Create deploy user, install SSH key
- Install Caddy from official repo
- Install systemd unit for the Go binary
- Drop Caddy config (Caddyfile)
- Set firewall rules: allow 22, 80, 443 only
- Enable unattended security upgrades

Paste the file into the VPS provider's user-data field at creation time. No Terraform, no Ansible, no extra toolchain.

## Project layout

```txt
cmd/site/           main.go — wiring, server start
internal/
  handlers/         page handlers, render() helper, error handlers
  middleware/       request ID, structured logging per request
  services/         business logic; calls content, returns view-ready data
  views/            *.gohtml templates
    views.go        //go:embed *.gohtml — embeds templates into binary
  content/          data structs + lookup helpers for projects / writing
static/             css, js, images, favicon, htmx.min.js
  static.go         //go:embed ... — embeds static assets into binary
deploy/
  cloud-init.yaml   first-boot VPS setup
  Caddyfile         reverse proxy + TLS config
  site.service      systemd unit
Makefile            build, deploy, restart targets
```

## Key constraints

1. Caddy proxies to Go on `127.0.0.1:8080` — never expose Go directly
2. All pages render inside `<main id="main">` — single layout contract for HTMX partial swaps
3. HTMX uses `hx-boost` + `hx-target="#main"` — handlers return full page or partial based on `HX-Request` header
4. Secrets live in systemd drop-in env file (`/etc/systemd/system/site.service.d/env.conf`), never in the repo

## Structured logging

Use `log/slog` (stdlib) throughout. No `log.Printf`.

Logging middleware wraps all handlers and emits one structured log line per request:

| Field         | Source                                                             |
| ------------- | ------------------------------------------------------------------ |
| `method`      | `r.Method`                                                         |
| `path`        | `r.URL.Path`                                                       |
| `status`      | response status code                                               |
| `duration_ms` | handler execution time                                             |
| `request_id`  | generated per request; also sent as `X-Request-ID` response header |
| `remote_ip`   | `X-Forwarded-For` set by Caddy, fallback to `r.RemoteAddr`         |

Application-level log entries use `level`, `msg`, and `error` (when applicable). Default level `INFO` in production, `DEBUG` when `SITE_ENV=dev`.
