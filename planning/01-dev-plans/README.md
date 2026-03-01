# Dev Plans (Pre-Deploy)

This track defines the remaining development work that should be completed before infrastructure deployment.

Deployment starts only after this track is complete:
- [../../docs/deployment.md](../../docs/deployment.md)
- [../../docs/deploy_checklist.md](../../docs/deploy_checklist.md)

---

## Sequence

| #   | Plan file           | Focus |
| --- | ------------------- | ----- |
| 12  | `phase_12_plan.md`  | Planning alignment and phase-map consistency |
| 13  | `phase_13_plan.md`  | Real content + copy replacement (remove placeholders) |
| 14  | `phase_14_plan.md`  | Public-site boilerplate assets (favicon, robots/security, metadata hygiene) |
| 15  | `phase_15_plan.md`  | Practical test expansion + local smoke harness |
| 16  | `phase_16_plan.md`  | GitHub Actions CI baseline |
| 17  | `phase_17_plan.md`  | Deploy guardrails and operator docs hardening |
| 18  | `phase_18_plan.md`  | Release-candidate dry run (production-like local validation) |
| 19  | `phase_19_plan.md`  | Final go/no-go checklist and release freeze |

---

## Why this ordering

1. Fix planning/numbering first so the team executes from one canonical map.
2. Replace public placeholders before adding polish/metadata.
3. Add CI only after tests and smoke checks are defined.
4. Validate deploy mechanics locally before touching the VPS.
5. Use a formal go/no-go gate to prevent rushed deployment.

---

## Exit rule for this track

Do not start production deployment until all phase exit gates in 12-19 are complete and documented.
