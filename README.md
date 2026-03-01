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

Builds locally, snapshots the previous VPS binary for rollback (`~/site.prev`), uploads the new binary, then restarts the systemd service. See `docs/deployment.md` for VPS setup and first deploy steps.

## Documentation

See `docs/` for architecture, content model, frontend contracts, operations, deployment, and deploy checklist.
