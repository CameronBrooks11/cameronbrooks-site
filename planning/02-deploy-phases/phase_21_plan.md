# Phase 21 — First Deploy & Smoke Test

**Goal:** Run `make deploy` to push the real binary to the VPS, verify every route returns the correct response, confirm HTMX navigation and TLS work in-browser, test the no-JS fallback, and validate structured logs via `make logs`. This is the final phase — at the end, the site is live.

**Exit gate:** `https://cameronbrooks.net` loads with a valid TLS certificate; all 6 public pages return `200`; `/healthz` returns `ok`; `/version` returns build info; a request to a nonexistent path returns the custom 404 page; HTMX navigation works; no-JS navigation works; `make logs` shows well-formed structured JSON per request.

---

## Prerequisites

- Phase 20 complete: VPS is running, `sudo cloud-init status` → `done`, DNS A records resolve to VPS IP
- Phase 11 complete: Makefile has `deploy` target, `VPS` variable set to `deploy@<vps-ip>`
- Phase 10 complete: HTMX and progress bar working locally
- Local machine: `make build` exits 0

---

## Step 1 — Final local build check

Before deploying, confirm a clean build with current source:

```sh
make build
```

Expected: exits 0, `bin/site` updated. Check the binary is linux/amd64 (not Windows):

```sh
# PowerShell — check GOOS/GOARCH from Makefile build target output
# Or verify the binary is non-empty and recently modified
(Get-Item bin/site).Length / 1MB
```

Expected: 5–15 MB. If it shows 0 bytes or fails to build, resolve locally before deploying.

---

## Step 2 — Run `make deploy`

```sh
make deploy
```

What happens:

1. `make build` — compiles `linux/amd64` binary with `-ldflags` setting `Version` and `BuildTime`
2. `scp bin/site deploy@<vps-ip>:~/site` — copies binary to server (overwrites placeholder or previous version)
3. `ssh deploy@<vps-ip> "sudo systemctl restart site"` — restarts the service

Expected output (example):

```
GOOS=linux GOARCH=amd64 go build -ldflags "..." -o bin/site ./cmd/site
bin/site                   100%   10MB  15.3MB/s   00:00
```

No error output. The restart command produces no output on success.

Confirm the service is running:

```sh
make ssh
# Once on VPS:
sudo systemctl status site
```

Expected: `active (running)` — not restarting.

Exit the SSH session.

---

## Step 3 — Verify TLS certificate issuance

Caddy issues a Let's Encrypt certificate on the first HTTPS request. Open a browser and navigate to:

```
http://cameronbrooks.net
```

Caddy should redirect immediately to `https://`. The first load may take 3–10 seconds as Caddy completes the ACME challenge. Subsequent loads are instant.

Verify the certificate from the terminal:

```sh
curl -I https://cameronbrooks.net
```

Expected response headers:

```
HTTP/2 200
content-type: text/html; charset=utf-8
strict-transport-security: max-age=31536000; includeSubDomains
x-content-type-options: nosniff
x-frame-options: SAMEORIGIN
referrer-policy: strict-origin-when-cross-origin
```

Status must be `200`. If you see `curl: (60) SSL certificate problem` the cert may not have issued yet — wait 30 seconds and retry.

If the cert consistently fails to issue:

- Confirm DNS has propagated (`nslookup cameronbrooks.net` resolves to VPS IP)
- Confirm port 80 is open in UFW (`sudo ufw status`)
- Check Caddy logs: `ssh deploy@<vps-ip> "sudo journalctl -u caddy --no-pager | tail -40"`
- Common cause: DNS not propagated yet or port 80 blocked (Let's Encrypt HTTP-01 challenge requires port 80)

---

## Step 4 — Full route smoke test

Test every route from the terminal using `curl`.

### Public pages (200)

```sh
# Home
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/

# Projects list
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/projects

# Project detail (use the slug from your data.go)
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/projects/cameronbrooks-site

# Writing list
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/writing

# Post detail (use the slug from your data.go)
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/writing/hello-world

# About
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/about

# Contact
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/contact
```

All must return `200`.

### System routes

```sh
# Healthz
curl -s https://cameronbrooks.net/healthz
# Expected body: ok

# Version
curl -s https://cameronbrooks.net/version
# Expected body: JSON with "version" and "build_time" fields
```

### 404 handling

```sh
# Unknown page
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/does-not-exist

# Unknown project slug
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/projects/no-such-project

# Unknown post slug
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/writing/no-such-post
```

All must return `404` — not `200` or `500`.

Verify the 404 response is the styled custom page (not a plain text error):

```sh
curl -s https://cameronbrooks.net/does-not-exist | Select-String -Pattern "<html|nav|#progress-bar"
```

Expected: matches showing this is a full HTML page with nav.

### Static assets

```sh
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/static/css/main.css
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/static/htmx.min.js
curl -o /dev/null -s -w "%{http_code}  %{url_effective}\n" https://cameronbrooks.net/static/js/progress.js
```

All must return `200`.

---

## Step 5 — Browser smoke test

Open `https://cameronbrooks.net` in Chrome or Firefox. Open DevTools (F12).

### Visual checks

- [ ] Home page renders with correct layout (nav, hero, sections, footer)
- [ ] Correct colours (teal accent `#0f766e`, not blue or black)
- [ ] Nav logo links to `/`
- [ ] Nav links: Projects, Writing, About, Contact all present
- [ ] Footer present

### Navigation checks

Open DevTools → Network tab.

1. Click "Projects" → URL changes to `/projects`, page updates, no full reload (no Document request in Network tab — only a Fetch)
2. Click a project card → URL changes to `/projects/<slug>`, partial swap
3. Click "← Projects" back-link → URL returns, partial swap
4. Click "Writing" → URL changes to `/writing`
5. Click "About" → URL changes to `/about`
6. Click browser Back button → previous page restored

On each navigation: the `<header>` and `<footer>` do not re-render (they do not flicker).

### Active nav state

- On `/projects`, the "Projects" link should have the teal active colour and `font-weight: 500`
- On `/writing`, the "Writing" link should be active
- On `/about`, "About" should be active

### Progress bar

Throttle the network (DevTools → Network → Slow 3G) and navigate between pages. The teal progress bar should:

1. Appear at the top (2px line, ~80% width)
2. Complete to full width
3. Fade out

---

## Step 6 — No-JavaScript fallback test

In Chrome DevTools → Settings (⚙ gear icon) → Preferences → Debugger → **Disable JavaScript**.

Or add `?nojs` and use the network tab to intercept and cancel JS requests — or just use the DevTools setting.

With JS disabled:

- [ ] All pages load via normal `<a>` navigation (full page loads)
- [ ] Nav links work
- [ ] Back button works
- [ ] No broken UI, no "JavaScript required" overlay
- [ ] Progress bar is invisible (correct — no JS to trigger it)

Re-enable JavaScript after testing.

---

## Step 7 — Check structured logs

```sh
make logs
```

This runs `ssh deploy@<vps-ip> "journalctl -u site -f"` — tails the service log in real time.

While `make logs` is running, open a browser tab and visit a few pages.

Each page request should produce a log line like:

```json
{
  "time": "2026-02-28T18:00:00Z",
  "level": "INFO",
  "msg": "request",
  "method": "GET",
  "path": "/projects",
  "status": 200,
  "duration_ms": 1,
  "request_id": "a3f7c2b1d4e9f021",
  "remote_ip": "1.2.3.4"
}
```

Verify:

- [ ] Output is one JSON object per line (structured, not free text)
- [ ] `"method"`, `"path"`, `"status"`, `"duration_ms"`, `"request_id"`, `"remote_ip"` fields all present
- [ ] HTMX navigation requests appear with `"path"` matching the navigated page
- [ ] Static asset requests appear (or are filtered — either is acceptable, depends on Phase 08 configuration)
- [ ] No `"level":"ERROR"` entries

Press Ctrl-C to exit `make logs`.

---

## Step 8 — `/version` output check

```sh
curl -s https://cameronbrooks.net/version | python -m json.tool
```

Expected (example — `version` is the short git SHA injected by `-ldflags`):

```json
{
  "build_time": "2026-02-28T17:45:00Z",
  "version": "abc1234"
}
```

> `python -m json.tool` sorts keys alphabetically, so `build_time` appears before `version`. The raw response is a single-line JSON object with `Content-Type: application/json`.

The `version` field comes from `ldflags -X main.Version` and the `build_time` from `ldflags -X main.BuildTime`. If both show `dev` / empty string, the binary was built without `-ldflags` — confirm `make build` was used (not `go run`).

---

## Step 9 — Commit

The site is live. No new files were created in this phase — it is purely operational. Mark the phase complete with a tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

---

## Exit gate checklist

**TLS and server:**

- [ ] `https://cameronbrooks.net` loads in browser with valid TLS cert (padlock in address bar)
- [ ] `curl -I https://cameronbrooks.net` → `HTTP/2 200`
- [ ] Security headers present in response: `strict-transport-security`, `x-content-type-options`, `x-frame-options`

**Routes:**

- [ ] `/` → `200`
- [ ] `/projects` → `200`
- [ ] `/projects/<valid-slug>` → `200`
- [ ] `/writing` → `200`
- [ ] `/writing/<valid-slug>` → `200`
- [ ] `/about` → `200`
- [ ] `/contact` → `200`
- [ ] `/healthz` → `200`, body `ok`
- [ ] `/version` → `200`, JSON with `version` and `build_time`
- [ ] `/does-not-exist` → `404`, custom HTML page (not plain text)
- [ ] `/projects/no-such-slug` → `404`
- [ ] `/static/css/main.css` → `200`
- [ ] `/static/htmx.min.js` → `200`

**Browser:**

- [ ] Nav links change URL without full page reload (HTMX partial swap)
- [ ] Active nav link highlighted correctly on each route
- [ ] Progress bar animates on navigation (visible with network throttling)
- [ ] JS disabled: all pages load and navigate as plain HTML

**Logs:**

- [ ] `make logs` shows structured JSON per request
- [ ] All expected fields present: `method`, `path`, `status`, `duration_ms`, `request_id`, `remote_ip`
- [ ] No ERROR-level entries

**Version:**

- [ ] `/version` JSON shows non-empty `version` and `build_time` (not `dev`)

All boxes checked → **the site is live. Phase 21 complete. Deploy track complete.**

---

**Ongoing operations reference**

For any future change (Go code, templates, CSS, content):

```sh
make deploy
```

That's it. Build, copy, restart — ~15 seconds.

To view logs at any time:

```sh
make logs
```

To SSH into the server:

```sh
make ssh
```

For troubleshooting reference, see `planning/RUNBOOK.md`.
