# Project Overview

This repository is a production-shaped personal site built with Go and deployed as a single binary behind Caddy.

## MVP scope

Public routes:

- `/`
- `/projects`
- `/projects/{slug}`
- `/writing`
- `/writing/{slug}`
- `/about`
- `/contact`
- `/healthz`
- `/version`

## Core product decisions

- Server-rendered HTML only (`html/template`)
- HTMX is progressive enhancement, not a rendering/runtime dependency
- Content is code-backed (`internal/content`) for MVP (no DB)
- Deploy target is one Debian VPS with systemd + Caddy

## MVP done criteria

- Site is live on production domain with valid HTTPS
- All public routes respond with expected status codes
- Navigation works with and without JavaScript
- Deploy path is repeatable via `make deploy`
- Rollback is available via `~/site.prev`

## Explicitly deferred

- Database and persistent user state
- Markdown ingestion pipeline
- Auth and admin surface
- Contact form backend
- Analytics and RSS
- Multi-service orchestration

## Key references

- `docs/architecture.md`
- `docs/content-model.md`
- `docs/frontend.md`
- `docs/operations.md`
- `docs/deployment.md`
- `docs/deploy_checklist.md`
