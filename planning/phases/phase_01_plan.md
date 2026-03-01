# Phase 01 — Repo Scaffold

**Goal:** Empty repository → fully laid-out module with every directory, a working Makefile, `.gitignore`, and a clean initial commit. No application code yet — just the skeleton everything else will be built inside.

**Exit gate:** `go mod tidy` exits 0; every directory from the STACK.md layout exists; `git status` is clean after the initial commit.

---

## Prerequisites

- Go 1.22+ installed (`go version` confirms)
- Git initialized in the repo root (GitHub already creates this — confirm with `git log`)
- Working directory: repo root (`d:\repos\cameronbrooks-site` or equivalent)

---

## Step 1 — Initialize the Go module

```sh
go mod init github.com/CameronBrooks11/cameronbrooks-site
```

This creates `go.mod`. No dependencies yet — stdlib only, so there is nothing to add.

```sh
go mod tidy
```

Creates an empty `go.sum`. Exits 0. Stop here if it errors — likely a Go version or network issue.

---

## Step 2 — Create the full directory tree

Create every directory from the STACK.md project layout. Directories with no files are tracked with a `.gitkeep` placeholder so they survive `git add`.

```sh
# From repo root
mkdir -p cmd/site
mkdir -p internal/handlers
mkdir -p internal/middleware
mkdir -p internal/services
mkdir -p internal/views
mkdir -p internal/content
mkdir -p static/css
mkdir -p static/js
mkdir -p static/images
mkdir -p deploy
mkdir -p planning/phases   # already exists from planning work; ensure it's present
mkdir -p bin               # build output — gitignored
```

Add `.gitkeep` to directories that will be empty at commit time (Go source will be added in later phases):

```sh
# PowerShell
@("internal/handlers", "internal/middleware", "internal/services", "internal/content", "static/images", "deploy") | ForEach-Object {
    New-Item -ItemType File -Path "$_/.gitkeep" -Force | Out-Null
}
```

> `static/css`, `static/js`, `internal/views`, and `cmd/site` will get real files in Phase 02 — no `.gitkeep` needed there.

---

## Step 3 — Create `.gitignore`

**File: `.gitignore`**

```gitignore
# Build output — never commit compiled binaries
bin/

# Go binaries and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Code coverage profiles and other test artifacts
*.out
coverage.*
*.coverprofile
profile.cov

# Go workspace file
go.work
go.work.sum

# Environment / secrets — never commit
.env
*.env

# Editor / IDE
.idea/
.vscode/

# Vim swap files
*.swp
*.swo
*~

# Logs
*.log
logs/

# OS artifacts
.DS_Store
Thumbs.db
```

---

## Step 4 — Create root `README.md`

**File: `README.md`**

````markdown
# cameronbrooks-site

Personal site. Go 1.22+, stdlib only, HTMX, single Debian VPS.

## Run locally

```sh
make dev
```
````

Requires Go 1.22+. Server starts at http://localhost:8080.

## Build

```sh
make build
```

Produces `bin/site` — a self-contained linux/amd64 binary with embedded templates and static assets.

## Deploy

```sh
make deploy
```

Builds locally, scps binary to VPS, restarts systemd service. See `planning/RUNBOOK.md` for first-time VPS setup.

## Planning

See `planning/` for full architecture, stack, content model, templates, UI/UX, roadmap, and operational runbook.

````

---

## Step 5 — Create Makefile

The Makefile defines every task runner command used throughout development and deployment. Commands reference paths that will exist after Phase 02+ — they are correct as written, just not runnable against working code until those phases complete.

**File: `Makefile`**

> Use **tabs** for recipe indentation — Make requires tabs, not spaces. Copy carefully.

```makefile
# cameronbrooks-site Makefile
# Requires Go 1.22+, ssh access to VPS.

# --- Configuration -----------------------------------------------------------
# Override on the command line: make deploy VPS=deploy@1.2.3.4
VPS ?= deploy@YOUR_VPS_IP
BINARY    = bin/site
CMD       = ./cmd/site

# Build-time version injection
VERSION   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
LDFLAGS   := -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME)"

# --- Local development -------------------------------------------------------
.PHONY: dev
dev:
	go run $(CMD)

# --- Build -------------------------------------------------------------------
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY) $(CMD)

# --- Deploy ------------------------------------------------------------------
.PHONY: deploy
deploy: build
	scp $(BINARY) $(VPS):~/site
	ssh $(VPS) "sudo systemctl restart site"

# --- VPS access --------------------------------------------------------------
.PHONY: ssh
ssh:
	ssh $(VPS)

.PHONY: logs
logs:
	ssh $(VPS) "journalctl -u site -f"

# --- Cleanup -----------------------------------------------------------------
.PHONY: clean
clean:
	rm -f $(BINARY)
````

**Notes on the Makefile:**

- `VPS` is a variable intentionally left as a placeholder here — it will be set to the real VPS IP in Phase 11.
- `LDFLAGS` injects `Version` and `BuildTime` into `main` package vars; this is wired in Phase 08.
- The `date` command uses UTC format compatible with Linux/macOS. **On Windows, `make build` must be run from WSL or Git Bash**, not PowerShell — `date -u` is a Unix command. See Phase 08 Step 4 for a PowerShell alternative.
- `build` cross-compiles to linux/amd64 regardless of the local machine's OS. Go's cross-compile is zero-config.

---

## Step 6 — Verify the scaffold

```sh
go mod tidy
```

Should exit 0. `go.sum` will be minimal or empty (no third-party deps).

```sh
# Confirm all directories exist
Get-ChildItem -Directory -Recurse | Select-Object FullName
```

Expected: `cmd/site`, `internal/handlers`, `internal/middleware`, `internal/services`, `internal/views`, `internal/content`, `static/css`, `static/js`, `static/images`, `deploy`, `planning`, `planning/phases`, `bin` (gitignored).

There is no Go code to build yet — `go build ./...` will produce nothing but also error on nothing (empty packages are skipped).

---

## Step 7 — Initial commit

```sh
git add .
git commit -m "phase 01: repo scaffold"
```

Run `git status` after — should show `nothing to commit, working tree clean`.

**Verify `bin/` is gitignored:**

```sh
git status bin/
# Should show: nothing to commit (ignored)
```

---

## Exit gate checklist

- [ ] `go mod tidy` exits 0
- [ ] `go.mod` contains `module github.com/CameronBrooks11/cameronbrooks-site` and `go 1.22`
- [ ] All directories from STACK.md layout exist (run `ls` check above)
- [ ] `Makefile` uses tab indentation (not spaces) — verify with a hex editor or `cat -A Makefile | head` looking for `^I`)
- [ ] `.gitignore` excludes `bin/`, `.env`
- [ ] `git status` is clean after commit
- [ ] `bin/` is not tracked by git

All boxes checked → proceed to Phase 02.
