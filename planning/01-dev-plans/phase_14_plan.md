# Phase 14 - Public Boilerplate Assets

**Goal:** Complete public-facing web boilerplate needed before launch (favicon, crawler/security text endpoints, metadata hygiene).

**Why now:** These are low-effort, high-signal launch essentials and should be done before infrastructure rollout.

---

## Scope

- Add a non-empty favicon.
- Add `robots.txt` and `security.txt` (and optionally `sitemap.xml`) served by the app.
- Ensure default metadata in layout is production-ready and consistent.

---

## Files to add/update

- `static/favicon.ico` (replace empty file)
- `internal/content/` or `internal/handlers/` additions for text endpoints:
  - `internal/handlers/static_text.go` (recommended)
- `cmd/site/main.go` (route wiring for `/robots.txt`, `/.well-known/security.txt`, optional `/sitemap.xml`)
- Optional metadata touch-ups:
  - `internal/views/layout.gohtml`

---

## Recommended endpoint payloads

`/robots.txt` minimal:

```txt
User-agent: *
Allow: /
Sitemap: https://cameronbrooks.net/sitemap.xml
```

`/.well-known/security.txt` minimal:

```txt
Contact: mailto:<your-email>
Preferred-Languages: en
Canonical: https://cameronbrooks.net/.well-known/security.txt
```

---

## Verification

Binary asset check:

```sh
ls -l static/favicon.ico
```

Expected: non-zero bytes.

Endpoint checks:

```sh
curl -i http://localhost:8080/robots.txt
curl -i http://localhost:8080/.well-known/security.txt
```

Expected: `200 OK`, `text/plain`, correct body.

Build/test:

```sh
go test ./...
go build ./...
```

---

## Exit gate

- [ ] Favicon is non-empty and served correctly
- [ ] Robots and security endpoints return 200 with expected text
- [ ] Metadata defaults are launch-ready
- [ ] Full test/build passes
