# Now vs Future

## What we're building now

A minimal but production-shaped personal site. The goal is a live, well-structured codebase — not a prototype.

### Pages

- `/` — home
- `/projects` and `/projects/:slug`
- `/writing` and `/writing/:slug`
- `/about`
- `/contact` — mailto link and social links (no form yet)

### Features

- Server-rendered HTML with Go `html/template`
- HTMX progressive enhancement (`hx-boost`, partial `#main` swaps)
- Hand-rolled CSS, no build step
- Content stored as Go structs in `internal/content/` (no DB, no markdown parsing yet)
- Single Debian VPS, single Go binary, Caddy in front
- cloud-init for first-boot provisioning
- Makefile for build and deploy

### What "done" looks like

- Site is live at the domain with HTTPS
- All pages load and navigate correctly with and without JS
- Binary is deployed as a systemd service, restarts on failure
- Caddy handles TLS automatically
- Deploy is a one-command process (`make deploy`)

---

## Future expansion — keep in mind, do not build yet

These are the things to avoid painting yourself into a corner on. The current design already accommodates them.

### Markdown content

When writing volume increases, add a markdown render step. Options:

- Render `.md` files at startup, cache in memory (no build step, minimal dep)
- Add a small `goldmark` dep when the time is right
- Keep content loading behind `internal/content/` so the rest of the app is unaffected

### CI / artifact-based deploys

Replace manual `scp` with:

- GitHub Actions builds a `linux/amd64` binary on push to `main`
- Uploads as a release artifact or SCP directly to the VPS
- This eliminates Go toolchain requirement on the server

### Subdomains and additional services

Caddy already supports host-based routing. Each new service is just:

- A new Go binary running on a different port as its own systemd service
- A new Caddy reverse_proxy block for that subdomain

No containers needed unless you outgrow single-host simplicity.

### Contact form

Add a handler that accepts POST, validates server-side, and sends email.
Go stdlib `net/smtp` works; no external dep needed for basic sending.

### Authentication

If private subdomains or admin pages are ever needed:

- Simple session-based auth in the Go app is sufficient for a single-user site
- Caddy can also enforce basic auth at the edge for quick internal tools

### Persistence

If state is needed (comments, analytics, etc.):

- SQLite first — single file, no separate service, well-supported by `database/sql`
- Keep storage behind a clean `internal/storage/` boundary so it can be swapped later

### More IaC if needed

cloud-init handles single-VPS setup fine. If you ever manage multiple servers:

- Ansible is the natural next step (idempotent, no heavy toolchain)
- Terraform only makes sense if you're creating and destroying infra regularly

---

## Decisions deferred on purpose

| Topic                                  | Deferred reason                                                     |
| -------------------------------------- | ------------------------------------------------------------------- |
| Markdown rendering                     | Not needed until writing volume justifies it                        |
| Database                               | No persistent state in v1                                           |
| Containers                             | Single binary + systemd is simpler for one server                   |
| CI deploy pipeline                     | Manual deploy is fine to start; add when friction is felt           |
| Analytics                              | Add only when curiosity justifies the complexity                    |
| `UpdatedAt time.Time` on content types | Not needed in v1; add when SEO or feed generation requires it       |
| robots.txt + security.txt              | Simple addition; do during initial deployment, not before           |
| CSP header via Caddy                   | Straightforward since there is no inline JS; add after site is live |
