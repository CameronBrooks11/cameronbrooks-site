# Phase 13 - Content And Copy Finalization

**Goal:** Remove placeholder public content and ship real copy/data for launch quality.

**Why now:** Deploying placeholder content creates immediate credibility issues and rework after go-live.

---

## Scope

- Replace bracket placeholders in templates (`[location]`, etc.).
- Replace placeholder contact email.
- Replace placeholder project/post entries with real entries (or intentionally remove extras).
- Keep one explicit draft post only if intentionally used for route gating tests.

---

## Files to update

- `internal/views/about.gohtml`
- `internal/views/contact.gohtml`
- `internal/content/data.go`
- Optional matching text in handlers if descriptions changed:
  - `internal/handlers/home.go`
  - `internal/handlers/projects.go`
  - `internal/handlers/writing.go`

---

## Content requirements

- About page includes real location and concise work summary.
- Contact page includes real email and active social links.
- Projects list contains only intentional entries with meaningful descriptions/tags.
- Writing list contains real summaries; no placeholder language.

---

## Verification

String scan:

```sh
rg -n "\[location\]|\[brief description\]|placeholder|example\.com" internal/views internal/content
```

Expected:
- No matches in production-facing content.
- If a draft test entry remains, the word `draft` is acceptable but not `placeholder`.

Run tests/build:

```sh
go test ./internal/content ./internal/services ./internal/handlers
go build ./...
```

Manual checks:
- `/about`, `/contact`, `/projects`, `/writing` all read as final copy.

---

## Exit gate

- [ ] No placeholder copy in `internal/views` and `internal/content`
- [ ] Real contact details are in place
- [ ] Content pages render with final text and valid links
- [ ] Tests/build pass after copy/data updates
