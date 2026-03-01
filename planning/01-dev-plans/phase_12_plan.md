# Phase 12 - Planning Alignment

**Goal:** Make planning docs internally consistent after the directory split and deploy-doc relocation. Ensure there is one canonical pre-deploy sequence and no stale references to old paths/numbers.

**Why now:** If planning references are inconsistent, execution drifts and phase gates become unreliable.

---

## Scope

- Align links from `planning/` docs to:
  - `planning/01-dev-plans/`
  - `docs/deployment.md`
- Remove stale references to deleted deploy phase files.
- Keep historical docs intact where they are explicitly archival; annotate if needed.

---

## Files to update

- `planning/README.md`
- `planning/ROADMAP.md`
- `planning/RUNBOOK.md` (only if it references old phase IDs/paths)
- `planning/01-dev-phases/README.md` (mark archival or completed baseline)
- `docs/deployment.md`
- `docs/deploy_checklist.md`

---

## Steps

1. Establish canonical index location for active plans (`planning/01-dev-plans/README.md`).
2. Replace broken/stale links (for example old `planning/phases/...`).
3. Move deployment execution references to docs (`docs/deployment.md`) while preserving deployment requirements.
4. Add one explicit note describing that phases 01-11 are complete baseline implementation.

---

## Verification

Run link/reference checks:

```sh
rg -n "planning/phases|phase_12_plan\.md|phase_13_plan\.md" planning
```

Expected:
- No stale path references to `planning/phases`.
- No accidental references that treat deploy docs as deleted phase files.

Manual check:
- Open `planning/README.md` and confirm it points to active plan indexes.

---

## Exit gate

- [ ] Planning docs reference current directories and numbering
- [ ] Deploy docs are referenced via `docs/deployment.md` (not deleted phase files)
- [ ] No stale `planning/phases` references remain
- [ ] Team can navigate from `planning/README.md` to the full plan sequence without dead links
