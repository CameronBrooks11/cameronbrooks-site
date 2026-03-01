# cameronbrooks-site

Personal site. Go 1.22+, stdlib only, HTMX, single Debian VPS.

## Run locally

```sh
make dev
```

Requires Go 1.22+. Server starts at http://localhost:8080.

## Smoke check

With the server running locally:

```sh
make smoke
```

## Build

```sh
make build
```

Produces `bin/site` — a self-contained linux/amd64 binary with embedded templates and static assets.

## Deploy

```sh
make deploy
```

Builds locally, scps binary to VPS, restarts systemd service. See `planning/RUNBOOK.md` for first-time VPS setup.

## Planning

See `planning/` for full architecture, stack, content model, templates, UI/UX, roadmap, and operational runbook.
