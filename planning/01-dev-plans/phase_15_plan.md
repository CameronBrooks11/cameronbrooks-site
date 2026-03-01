# Phase 15 - Test Expansion And Local Smoke Harness

**Goal:** Add practical launch-focused tests and a repeatable smoke command set so regressions are caught before CI/deploy.

**Why now:** Existing package-level tests are good, but launch needs route-level behavior and middleware integration checks.

---

## Scope

- Add HTTP route integration tests against real mux wiring.
- Verify key response contracts:
  - `X-Request-ID` present
  - expected status for core routes and not-found routes
  - `/healthz` and `/version` contracts
- Add a simple smoke script or Make target for local preflight.

---

## Files to add/update

- `cmd/site/main_test.go` (or `internal/app/router_test.go` if router extraction is done)
- Optional helper:
  - `scripts/smoke_local.sh` (or PowerShell equivalent)
- Optional Make target:
  - `Makefile` (`smoke`)

---

## Test matrix (minimum)

- `GET /` -> 200
- `GET /projects` -> 200
- `GET /writing` -> 200
- `GET /healthz` -> 200 + `ok`
- `GET /version` -> 200 + JSON payload
- `GET /does-not-exist` -> 404
- response contains `X-Request-ID` on all above

---

## Verification

```sh
go test ./...
```

If smoke script is added:

```sh
./scripts/smoke_local.sh
```

Expected: all checks pass without manual editing.

---

## Exit gate

- [ ] Route integration tests exist and pass
- [ ] Middleware behavior (`X-Request-ID`) is asserted at integration level
- [ ] Local smoke command exists and is repeatable
- [ ] `go test ./...` remains green
