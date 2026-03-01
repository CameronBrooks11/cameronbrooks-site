# Phase 12 - Planning Alignment

**Goal:** Make planning docs internally consistent after the directory split and deploy-phase renumbering. Ensure there is one canonical pre-deploy sequence and no stale references to old paths/numbers.

**Why now:** If planning references are inconsistent, execution drifts and phase gates become unreliable.

---

## Scope

- Align links from `planning/` docs to:
  - `planning/01-dev-plans/`
  - `planning/02-deploy-phases/`
- Update references from old `phase_12/13` deploy terminology to `phase_20/21` where appropriate.
- Keep historical docs intact where they are explicitly archival; annotate if needed.

---

## Files to update

- `planning/README.md`
- `planning/ROADMAP.md`
- `planning/RUNBOOK.md` (only if it references old phase IDs/paths)
- `planning/01-dev-phases/README.md` (mark archival or completed baseline)
- `planning/02-deploy-phases/phase_20_plan.md`
- `planning/02-deploy-phases/phase_21_plan.md`

---

## Steps

1. Establish canonical index location for active plans (`planning/01-dev-plans/README.md`).
2. Replace broken/stale links (for example old `planning/phases/...`).
3. Update deploy phase doc titles to match file numbers (`Phase 20`, `Phase 21`) while preserving content.
4. Add one explicit note describing that phases 01-11 are complete baseline implementation.

---

## Verification

Run link/reference checks:

```sh
rg -n "planning/phases|phase_12_plan\.md|phase_13_plan\.md" planning
```

Expected:
- No stale path references to `planning/phases`.
- No accidental references that treat deploy docs as phases 12/13.

Manual check:
- Open `planning/README.md` and confirm it points to active plan indexes.

---

## Exit gate

- [ ] Planning docs reference current directories and numbering
- [ ] `phase_20_plan.md` and `phase_21_plan.md` headings match numbering
- [ ] No stale `planning/phases` references remain
- [ ] Team can navigate from `planning/README.md` to the full plan sequence without dead links
