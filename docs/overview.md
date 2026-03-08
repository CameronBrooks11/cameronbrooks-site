# Project Overview

This repository is a production-shaped personal site built with Go and deployed as a single binary behind Caddy.

## Public routes

- `/`
- `/writing`
- `/writing/{slug}`
- `/about`
- `/contact`
- `/healthz`
- `/version`

## Core product decisions

- Server-rendered HTML only (`html/template`)
- HTMX is progressive enhancement, not a rendering/runtime dependency
- Writing posts are Markdown files in `internal/content/writing/`, embedded and parsed at startup
- Deploy target is one Debian VPS with systemd + Caddy

## Explicitly deferred

- Database and persistent user state
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
- `docs/deploy-checklist.md`
