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
| [phases/README.md](phases/README.md) | Ordered phase map: 13 phases from empty repo to live deployed MVP index of all phase_XY_plan.md docs      |
