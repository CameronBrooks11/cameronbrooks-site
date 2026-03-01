# Phase 19 - Final Go/No-Go Gate

**Goal:** Create a formal pre-deploy signoff gate so deployment is a deliberate decision, not a momentum step.

**Why now:** This is the final quality and operations checkpoint before starting deploy phases 20/21.

---

## Scope

- Consolidate all phase 12-18 outcomes into a single checklist.
- Record unresolved risks explicitly.
- Decide go/no-go with clear criteria.

---

## Files to add/update

- `docs/release_checklist.md` (recommended)
- Optional runbook link to checklist

---

## Required signoff categories

1. **Product readiness**
   - Real content/copy present
   - Public boilerplate assets complete

2. **Quality readiness**
   - Tests passing
   - CI green
   - Local release dry run passed

3. **Operational readiness**
   - Deploy docs current
   - Rollback steps tested/documented
   - Deploy files validated (no placeholders, valid YAML)

4. **Risk register**
   - list known non-blocking risks
   - owner + follow-up phase/date

---

## Verification

Checklist completion command examples:

```sh
go test ./...
go vet ./...
go build ./...
rg -n "YOUR_SSH_PUBLIC_KEY|YOUR_DOMAIN" deploy
```

Manual:
- Confirm CI status in GitHub
- Confirm release checklist has explicit go/no-go decision

---

## Exit gate

- [ ] Release checklist exists and is complete
- [ ] All blocking items are closed
- [ ] Known risks are documented and accepted
- [ ] Explicit `GO` decision recorded

After this gate, begin deploy phase 20.
