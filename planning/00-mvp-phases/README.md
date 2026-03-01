# MVP Foundation Phase Plans (01-11)

Baseline implementation plan that produced the current codebase foundation.

These phases are complete and are kept as the historical build record. New work should use:
- `planning/01-dev-plans/` for pre-deploy hardening (phases 12-19)
- `docs/deployment.md` for deployment execution

---

## Build approach

Phases 01-11 were ordered bottom-up so each layer had stable dependencies before the next layer was added.

```txt
Infrastructure  (01-02)  -> module, embed skeleton, binary boots
Application     (03-08)  -> content, services, templates, middleware, handlers, routing
Presentation    (09-10)  -> CSS + HTMX/progress
Deploy Config   (11)     -> cloud-init, Caddyfile, systemd unit
```

---

## Hard sequencing rules used in 01-11

| Rule | Reason |
| --- | --- |
| 02 before 05 | `//go:embed *.gohtml` needs template files present |
| 03 before 04 | services depend on content model/types |
| 04 before 07 | handlers call service functions |
| 05 before 07 | handlers rely on `render()` and template caches |
| 06 before 08 | `main.go` wiring uses `middleware.Chain` |
| 08 before 09 | CSS validation requires live routes/pages |
| 10 before 11 | deploy files should reflect final app/runtime contracts |

---

## Registry (01-11)

| # | Plan file | Focus |
| --- | --- | --- |
| 01 | `phase_01_plan.md` | Repository scaffold and Makefile baseline |
| 02 | `phase_02_plan.md` | Embed infrastructure + skeleton binary |
| 03 | `phase_03_plan.md` | Content model and lookup helpers |
| 04 | `phase_04_plan.md` | Services layer and trust-boundary conversion |
| 05 | `phase_05_plan.md` | Template system and render pipeline |
| 06 | `phase_06_plan.md` | Middleware (request ID, logger, chain) |
| 07 | `phase_07_plan.md` | Handlers and error/system endpoints |
| 08 | `phase_08_plan.md` | Routing and server wiring |
| 09 | `phase_09_plan.md` | Production stylesheet |
| 10 | `phase_10_plan.md` | HTMX vendoring and progress behavior |
| 11 | `phase_11_plan.md` | Deploy infrastructure files |

---

## Next planning track

Continue with:
1. `planning/01-dev-plans/README.md` (phases 12-19)
2. `docs/deployment.md`
3. `docs/deploy_checklist.md`
