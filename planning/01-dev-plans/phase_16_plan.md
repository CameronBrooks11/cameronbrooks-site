# Phase 16 - GitHub CI Baseline

**Goal:** Add simple, reliable CI workflows that enforce build/test quality on every push and pull request.

**Why now:** Manual validation is error-prone; CI provides consistent quality gates before deployment.

---

## Scope

- Create GitHub Actions workflow for:
  - `go test ./...`
  - `go vet ./...`
  - `go build ./...`
  - formatting check (`gofmt -l` must be empty)
- Trigger on `push` and `pull_request`.
- Keep workflow minimal and fast.

---

## Files to add

- `.github/workflows/ci.yml`

Optional:
- `.github/workflows/security.yml` for `govulncheck` (if desired after baseline)

---

## Recommended workflow shape

- `actions/checkout@v4`
- `actions/setup-go@v5` with `go-version-file: go.mod`
- Cache enabled via setup-go
- Steps:
  1. `go test ./...`
  2. `go vet ./...`
  3. `gofmt -l .` (fail if output non-empty)
  4. `go build ./...`

---

## Verification

Local precheck before push:

```sh
go test ./...
go vet ./...
go build ./...
```

Then push branch and confirm workflow passes in GitHub UI.

---

## Exit gate

- [ ] CI workflow exists in `.github/workflows/ci.yml`
- [ ] Workflow runs on push + PR
- [ ] Workflow is green on current branch
- [ ] CI enforces tests, vet, build, and formatting
