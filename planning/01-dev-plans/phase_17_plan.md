# Phase 17 - Deploy Guardrails And Runbook Hardening

**Goal:** Add operator guardrails so deployment is predictable and low-risk.

**Why now:** A technically-correct deploy can still fail operationally without clear checks and rollback guidance.

---

## Scope

- Add/update runbook sections for:
  - pre-deploy checklist
  - rollback procedure
  - post-deploy verification
- Add a validation command block for deploy artifacts (placeholders, YAML, binary build).
- Ensure docs clearly separate dev phases (12-19) from deploy phases (20-21).

---

## Files to update

- `planning/RUNBOOK.md`
- `docs/README.md` (if needed for discoverability)
- Optional helper doc:
  - `docs/deploy_checklist.md`

---

## Mandatory operator procedures

1. **Pre-deploy checks**
   - clean git state (or known intentional changes)
   - all CI checks green
   - deploy config validated

2. **Deploy**
   - single command path (`make deploy`)

3. **Rollback**
   - restore previous binary on VPS and restart service
   - verify health endpoints and key pages

4. **Post-deploy**
   - smoke routes
   - logs inspection (`make logs`)

---

## Verification

Manual runbook walkthrough should be executable line-by-line without ambiguous steps.

---

## Exit gate

- [ ] Runbook contains explicit pre-deploy, deploy, rollback, and post-deploy sections
- [ ] Deploy artifact validation commands are documented
- [ ] Dev/deploy phase boundaries are unambiguous in docs
