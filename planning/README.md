# Planning

Working planning documents for the site. Intended for both human reference and AI agent context.

## Documents

| File                                 | Contents                                                                                                  |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| [EXEC_SUMMARY.md](EXEC_SUMMARY.md)   | One-page overview + per-doc summary of all decisions — read this first                                    |
| [ARCHITECTURE.md](ARCHITECTURE.md)   | Explicit architectural invariants — the load-bearing constraints that must not be violated                |
| [STACK.md](STACK.md)                 | Tech stack, Go version requirement, project layout, deploy chain, key constraints                         |
| [ROADMAP.md](ROADMAP.md)             | What is being built now vs. what is deferred and why                                                      |
| [UI_UX.md](UI_UX.md)                 | Visual language, design tokens, layout, page patterns, HTMX navigation, CSS structure, accessibility      |
| [CONTENT.md](CONTENT.md)             | Go struct definitions for Project and Post, storage strategy, lookup helpers, migration path to markdown  |
| [TEMPLATES.md](TEMPLATES.md)         | Template file naming, layout structure, PageData, `//go:embed` cache strategy, `render()` helper, routing |
| [RUNBOOK.md](RUNBOOK.md)             | Local dev setup, VPS first bring-up, DNS, ongoing deploys, secrets management                             |
| [00-mvp-phases/README.md](00-mvp-phases/README.md) | Baseline implementation plan (phases 01-11) that produced the current app foundation |
| [01-dev-plans/README.md](01-dev-plans/README.md) | Remaining pre-deploy development/hardening plan (phases 12-19) |
| [02-deploy-phases/phase_20_plan.md](02-deploy-phases/phase_20_plan.md) | Deploy phase 20: VPS provisioning and DNS cutover |
| [02-deploy-phases/phase_21_plan.md](02-deploy-phases/phase_21_plan.md) | Deploy phase 21: first deploy and live smoke test |
