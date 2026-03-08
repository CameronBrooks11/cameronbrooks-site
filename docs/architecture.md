# Architecture

## System boundary

```txt
Internet -> Caddy (80/443) -> Go app (127.0.0.1:8080)
```

Go is never exposed directly to the public internet.

## Invariants

1. Server owns rendering. Pages are rendered on the server; HTMX only swaps HTML fragments.
2. Single deployable binary. Templates and static assets are embedded at compile time.
3. Progressive enhancement. Site remains fully usable without JavaScript.
4. Thin handlers. Request parsing + service call + render only; business logic belongs in services.
5. One layout contract. Page content is swapped in `<main id="main">`.
6. Trust boundary in services. Raw content `Body` is converted to trusted HTML only in `internal/services`.
7. No repository secrets. Runtime secrets live in server-side systemd drop-ins.

- No runtime filesystem dependency. Content and templates are embedded at compile time; Markdown is parsed once at startup, not per-request.

## Runtime composition

- `cmd/site`: server wiring, startup, graceful shutdown
- `internal/handlers`: route handlers + render orchestration
- `internal/services`: view-model preparation and trust conversion boundary
- `internal/content`: Post type, Markdown loader (parses `writing/*.md` at startup via embed)
- `internal/middleware`: request ID + structured request logging
- `internal/views`: embedded templates
- `static`: embedded CSS/JS/assets
- `deploy`: cloud-init, Caddyfile, systemd unit

## Logging contract

Per request log entries include:

- `method`
- `path`
- `status`
- `duration_ms`
- `request_id`
- `remote_ip`

`X-Request-ID` is returned on responses.

## Guardrails

Avoid these drift patterns:

- Rendering HTML directly in handlers
- Moving trust conversion into storage structs
- Adding JS-only critical paths
- Reintroducing runtime file reads for templates/assets
- Exposing the Go listener directly
