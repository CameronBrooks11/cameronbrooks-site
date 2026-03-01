# Architecture Invariants

These are the fixed constraints of the system. They are not implementation preferences — they are load-bearing decisions that everything else is built on top of. Violating one cascades into multiple other things breaking or needing rewrites.

When a future change is "just a quick exception," check it against this list first.

---

## Invariants

### 1. The server owns all rendering

There is no client-side rendering. HTML is assembled on the server and sent to the browser. HTMX replaces target DOM fragments with server-rendered HTML — it does not manage state or drive rendering logic.

The browser holds no authoritative application state. A page refresh must always produce a correct, fully-rendered result.

### 2. One deployable binary

The application ships as a single compiled Go binary. Templates (`internal/views/`) and static assets (`static/`) are embedded at compile time via `//go:embed`. No files are read from disk at runtime. No file sync step is required on deploy.

Deploy = copy one binary + restart service.

### 3. Caddy is the permanent public interface

```txt
Internet → Caddy → Go app (127.0.0.1:8080)
```

Go's HTTP port is never exposed directly, in any environment. Caddy handles TLS termination, HTTP→HTTPS redirect, and reverse proxying. This boundary is permanent — future services are added as additional Caddy upstreams, not by exposing Go ports.

Caddy also owns: HTTP hardening headers, cache-control policies, CSP (future).

### 4. Content is trusted after the service boundary

Storage structs (`internal/content/`) hold `Body` as `string`. The explicit conversion to `template.HTML` — which disables Go's auto-escaping — happens only in `internal/services/`, after any sanitization decision. This trust grant is never made in the storage layer.

`html/template` auto-escaping is the default for all other fields throughout the application.

### 5. Progressive enhancement — JS is additive

Every page must be fully readable and navigable without JavaScript. HTMX adds app-like navigation feel; it does not gate access to any content. Disabling JS produces a functional plain HTML site.

This is a correctness invariant, not a preference. It prevents dependency on JS availability for content delivery and ensures SEO correctness.

### 6. Handlers are thin

Handlers (`internal/handlers/`) do three things: parse request input, call a service function, call `render()`. Business logic — content lookup, filtering, view-model construction — lives in `internal/services/`. Data structures live in `internal/content/`.

A handler that grows beyond ~20 lines of logic is a signal that service layer work has leaked upward.

### 7. One layout contract

Every page's dynamic content lives inside `<main id="main">`. The HTMX partial swap replaces this element. The `{{define "content"}}` block in each page template is the only thing that changes between pages.

Nothing inside `<main>` touches or references the layout shell (header, nav, footer).

### 8. No secrets in the repository

Runtime secrets are environment variables injected via systemd drop-in at `/etc/systemd/system/site.service.d/env.conf`. No secret, credential, key, or token is ever committed to the repository — not even in example files.

### 9. No runtime filesystem dependency

The binary loads nothing from disk after startup. Templates parsed at startup from embedded FS. Static files served from embedded FS. Content loaded from compiled-in Go data.

This makes the binary portable, container-ready if ever needed, and immune to deploy ordering issues (no "forgot to sync assets" failures).

---

## Constraint interaction map

Understanding how these connect helps evaluate future changes:

```txt
[3] Caddy boundary
    └─ enforces [5] progressive enhancement at the edge (future CSP)

[2] Single binary (embed)
    └─ requires [9] no runtime FS dependency
    └─ enables simple deploy from [3]

[4] Trust boundary at services
    └─ depends on [6] thin handlers (handlers don't make trust decisions)

[1] Server-owned rendering
    └─ depends on [7] layout contract (consistent swap target)
    └─ depends on [5] JS-optional (rendering doesn't require client execution)
```

---

## Common drift patterns to watch for

These are the typical first violations — the "just this once" exceptions that compound:

| Temptation                                            | Violates | Why it matters                       |
| ----------------------------------------------------- | -------- | ------------------------------------ |
| Render HTML in a handler directly (skip render())     | #1, #7   | Breaks HTMX swap contract            |
| Read a template file at request time for "hot reload" | #9       | Runtime FS dependency                |
| Store `template.HTML` in a content struct             | #4       | Disables escaping at storage layer   |
| Add a feature that only works with JS                 | #5       | Breaks no-JS correctness requirement |
| Expose Go port in Caddy config "temporarily"          | #3       | Permanent boundary erosion           |
| Commit a `.env` file with real values                 | #8       | Credential exposure                  |
| Put DB query logic in a handler                       | #6       | Handler thickness violation          |
