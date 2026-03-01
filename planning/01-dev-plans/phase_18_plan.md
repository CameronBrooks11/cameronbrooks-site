# Phase 18 - Release Candidate Dry Run

**Goal:** Rehearse a production-like release locally before touching infrastructure.

**Why now:** This catches packaging/runtime issues while rollback is trivial.

---

## Scope

- Build with release-style ldflags.
- Run binary with production-like env on a local port.
- Execute smoke tests against the built binary (not `go run`).
- Verify version/build metadata in `/version`.

---

## Files to add/update

- Optional helper script:
  - `scripts/release_dry_run.sh`
- Optional note in runbook with exact dry-run commands.

---

## Dry-run checklist

1. Build release artifact:

```sh
make build
```

2. Run artifact locally:

```sh
SITE_ADDR=:18080 SITE_ENV=production ./bin/site
```

3. Validate endpoints:

```sh
curl -i http://localhost:18080/healthz
curl -i http://localhost:18080/version
curl -i http://localhost:18080/
curl -i http://localhost:18080/does-not-exist
```

4. Validate static assets:

```sh
curl -I http://localhost:18080/static/htmx.min.js
curl -I http://localhost:18080/static/css/main.css
```

---

## Verification

- `/healthz` returns 200 + `ok`
- `/version` returns non-empty JSON fields when built with ldflags
- HTML routes and 404 route behave as expected
- Structured logs are emitted in production env

---

## Exit gate

- [ ] Release binary runs locally without `go run`
- [ ] Smoke checks pass against built artifact
- [ ] Version/build metadata behavior is confirmed
- [ ] No runtime surprises found during rehearsal
