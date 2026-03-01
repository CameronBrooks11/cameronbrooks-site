# Phase 11 — Deploy Infrastructure

**Goal:** Write all three deploy files (`deploy/cloud-init.yaml`, `deploy/Caddyfile`, `deploy/site.service`) with real, production-ready content (no placeholder stubs), and confirm the `deploy`, `ssh`, and `logs` Makefile targets from Phase 01 are in place. No server exists yet — this phase is about having the files ready so Phase 12 can consume them immediately.

**Exit gate:** `make build` produces `bin/site` (linux/amd64); all three deploy files exist on disk with no placeholder strings (`YOUR_SSH_PUBLIC_KEY`, `YOUR_DOMAIN`, etc.); `cloud-init.yaml` passes YAML validation; Makefile `deploy`, `ssh`, and `logs` targets are present.

---

## Prerequisites

- Phase 10 complete (CSS, HTMX, progress bar working locally)
- Phase 01/08 complete (Makefile with `build` target including `-ldflags`)
- Domain name: `cameronbrooks.net`
- You have an SSH public key at `~/.ssh/id_ed25519.pub` (or equivalent)

---

## Files to create in this phase

```
deploy/cloud-init.yaml
deploy/Caddyfile
deploy/site.service
```

> The `deploy`, `ssh`, and `logs` Makefile targets already exist from Phase 01. Step 4 of this phase confirms them — no Makefile edits expected.

---

## Step 1 — Write `deploy/site.service`

This is the simplest of the three files. Write it first because cloud-init will embed its content.

**File: `deploy/site.service`**

```ini
[Unit]
Description=cameronbrooks.net Go site
Documentation=https://github.com/CameronBrooks11/cameronbrooks-site
After=network.target

[Service]
Type=simple
User=deploy
Group=deploy
WorkingDirectory=/home/deploy
ExecStart=/home/deploy/site
Restart=on-failure
RestartSec=5
TimeoutStopSec=15

# Runtime environment
Environment="SITE_ADDR=127.0.0.1:8080"
Environment="SITE_ENV=production"

# Optional: secrets drop-in (created manually on server after first deploy)
# File: /etc/systemd/system/site.service.d/env.conf
# Contents: [Service]\nEnvironment="SOME_SECRET=value"
EnvironmentFile=-/etc/systemd/system/site.service.d/env.conf

# Hardening
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

**Notes:**

- `User=deploy` — process runs as the non-root `deploy` user
- `ExecStart=/home/deploy/site` — `make deploy` scps the binary here
- `Restart=on-failure` — auto-recovers from crashes but not from clean exits
- `RestartSec=5` — 5-second backoff before restart (prevents restart storm)
- `TimeoutStopSec=15` — gives the app's 10-second graceful shutdown a 5-second buffer before systemd force-kills
- `EnvironmentFile=-/...` — the `-` prefix makes this optional; if the file doesn't exist the service still starts
- `NoNewPrivileges=true` / `PrivateTmp=true` — baseline systemd hardening, no cost

---

## Step 2 — Write `deploy/Caddyfile`

**File: `deploy/Caddyfile`**

```caddy
cameronbrooks.net, www.cameronbrooks.net {
    reverse_proxy 127.0.0.1:8080

    # Access log — structured JSON, one line per request
    log {
        output stderr
        format json
    }

    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }
}
```

**Notes:**

- Caddy automatically issues a Let's Encrypt TLS certificate for all hostnames in the site block. No certificate configuration needed.
- `reverse_proxy 127.0.0.1:8080` — proxies all traffic to the Go binary. The Go binary must only listen on `127.0.0.1:8080` (loopback), never on `0.0.0.0`.
- `www.cameronbrooks.net` is listed alongside the apex — Caddy handles both.
- HTTP → HTTPS redirect is automatic in Caddy when HTTPS is configured — no explicit redirect block needed.
- `-Server` header removes Caddy's `Server: Caddy` response header.

---

## Step 3 — Write `deploy/cloud-init.yaml`

cloud-init performs the entire server bootstrap on first boot: user creation, firewall, Caddy install, systemd service setup. After cloud-init runs, `make deploy` is all that's needed to go live.

**Fill in before use:**

1. Replace `YOUR_SSH_PUBLIC_KEY` with the full content of your `~/.ssh/id_ed25519.pub` (the line starting with `ssh-ed25519 ...`)
2. Replace `YOUR_DOMAIN` with `cameronbrooks.net` in the Caddyfile section

**File: `deploy/cloud-init.yaml`**

```yaml
#cloud-config

# === User setup ===
users:
  - name: deploy
    groups: [sudo]
    shell: /bin/bash
    sudo: "ALL=(ALL) NOPASSWD:ALL"
    ssh_authorized_keys:
      - YOUR_SSH_PUBLIC_KEY

# === System packages ===
package_update: true
package_upgrade: true
packages:
  - ufw
  - curl
  - debian-keyring
  - debian-archive-keyring
  - apt-transport-https

# === First-boot commands ===
runcmd:
  # --- Install Caddy from official apt repo ---
  - curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
  - curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
  - apt-get update
  - apt-get install -y caddy

  # --- Write Caddyfile ---
  - |
    cat > /etc/caddy/Caddyfile << 'EOF'
    YOUR_DOMAIN, www.YOUR_DOMAIN {
        reverse_proxy 127.0.0.1:8080

        log {
            output stderr
            format json
        }

        header {
            Strict-Transport-Security "max-age=31536000; includeSubDomains"
            X-Content-Type-Options "nosniff"
            X-Frame-Options "SAMEORIGIN"
            Referrer-Policy "strict-origin-when-cross-origin"
            -Server
        }
    }
    EOF

  # --- Create placeholder site binary (service won't crash on first start) ---
  - touch /home/deploy/site
  - chmod +x /home/deploy/site
  - chown deploy:deploy /home/deploy/site

  # --- Install systemd service ---
  - |
    cat > /etc/systemd/system/site.service << 'EOF'
    [Unit]
    Description=cameronbrooks.net Go site
    After=network.target

    [Service]
    Type=simple
    User=deploy
    Group=deploy
    WorkingDirectory=/home/deploy
    ExecStart=/home/deploy/site
    Restart=on-failure
    RestartSec=5
    TimeoutStopSec=15
    Environment="SITE_ADDR=127.0.0.1:8080"
    Environment="SITE_ENV=production"
    EnvironmentFile=-/etc/systemd/system/site.service.d/env.conf
    NoNewPrivileges=true
    PrivateTmp=true

    [Install]
    WantedBy=multi-user.target
    EOF

  # --- Create secrets drop-in directory ---
  - mkdir -p /etc/systemd/system/site.service.d

  # --- Enable and start services ---
  - systemctl daemon-reload
  - systemctl enable site
  - systemctl start site
  - systemctl enable caddy
  - systemctl restart caddy

  # --- Firewall: allow SSH, HTTP, HTTPS only ---
  - ufw default deny incoming
  - ufw default allow outgoing
  - ufw allow 22/tcp
  - ufw allow 80/tcp
  - ufw allow 443/tcp
  - ufw --force enable

  # --- Unattended security upgrades ---
  - apt-get install -y unattended-upgrades
  - dpkg-reconfigure -f noninteractive unattended-upgrades
```

**Key design decisions:**

| Decision                              | Reason                                                                                                                                                                                                                                                                                                   |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `NOPASSWD:ALL` sudo for `deploy`      | Allows `make deploy` to restart the service without a password prompt over SSH. Acceptable for a single-operator personal site. Remove after setup if you prefer.                                                                                                                                        |
| `touch /home/deploy/site` placeholder | The service unit starts immediately. Without a binary, `ExecStart` fails. The placeholder lets the service start (and immediately exit/restart) before the real binary is deployed via `make deploy`. Alternatively, set `Restart=on-failure` — it will keep retrying until the real binary is deployed. |
| Caddy from official apt repo          | Ensures Caddy is updated via `apt upgrade`; the Debian default Caddy package in some releases is outdated.                                                                                                                                                                                               |
| Firewall before services start        | The firewall `ufw enable` comes last in `runcmd` — services are up first so cloud-init doesn't lock itself out. UFW drop rules take effect immediately.                                                                                                                                                  |

> **Replace `YOUR_SSH_PUBLIC_KEY` and `YOUR_DOMAIN` before committing.** If these placeholders remain, the exit gate check (see below) will catch it.

---

## Step 4 — Confirm Makefile `deploy`, `ssh`, `logs` targets

The Makefile written in Phase 01 already contains the `VPS` variable and `deploy`, `ssh`, and `logs` targets. **No edits are needed here** — this is a confirmation step.

Open `Makefile` and verify these targets are present:

```makefile
# Override on the command line: make deploy VPS=deploy@1.2.3.4
VPS ?= deploy@YOUR_VPS_IP

.PHONY: deploy
deploy: build
	scp $(BINARY) $(VPS):~/site
	ssh $(VPS) "sudo systemctl restart site"

.PHONY: ssh
ssh:
	ssh $(VPS)

.PHONY: logs
logs:
	ssh $(VPS) "journalctl -u site -f"
```

If they are not there (e.g. the Makefile was reset), add them — but this should not be necessary if Phase 01 was followed.

The `VPS` placeholder (`deploy@YOUR_VPS_IP`) will be replaced with the real IP in Phase 12 Step 6, once the new server exists. Until then, `make deploy` will fail to connect — that is correct and expected.

**How `make deploy` works:**

1. `make build` compiles a `linux/amd64` binary at `bin/site` with `-ldflags` injecting `Version` and `BuildTime` (from Phase 08)
2. `scp bin/site $(VPS):~/site` copies the binary to the deploy user's home directory, overwriting the previous binary
3. `ssh $(VPS) "sudo systemctl restart site"` restarts the service, which picks up the new binary

The old binary is overwritten atomically — there is no "stop, replace, start" window because `scp` writes to the same path and systemd's restart is triggered afterward.

---

## Step 5 — Replace placeholder strings

Before committing, verify no placeholder strings remain:

```sh
# Check for remaining placeholders
Select-String -Path "deploy\*" -Pattern "YOUR_SSH_PUBLIC_KEY|YOUR_DOMAIN|YOUR_VPS_IP"
```

Expected: zero matches (if you have filled them all in).

If running `make deploy` in Phase 13 from a machine where the VPS IP isn't yet known, set the IP then. The SSH key and domain must be filled in before Phase 12 (cloud-init is pasted into the provider UI at VPS creation time).

---

## Step 6 — Validate YAML syntax

```sh
# Install yq if not present (optional — any YAML validator works)
# Or use Python's yaml module:
python -c "import yaml, sys; yaml.safe_load(open('deploy/cloud-init.yaml'))" && echo "YAML valid"
```

Or use an online YAML validator (paste the file contents). The file must parse without errors.

---

## Step 7 — Confirm `make build` works

```sh
make build
ls -lh bin/site
```

Expected: `bin/site` exists, size 5–15 MB (Go binary with embedded assets). This confirms the binary is ready to be deployed in Phase 13.

---

## Step 8 — Commit

```sh
git add deploy/cloud-init.yaml deploy/Caddyfile deploy/site.service Makefile
git commit -m "phase 11: deploy infrastructure"
```

---

## Exit gate checklist

- [ ] `deploy/cloud-init.yaml` exists and is valid YAML (no parse errors)
- [ ] `deploy/Caddyfile` exists with correct domain name
- [ ] `deploy/site.service` exists with `User=deploy`, `Restart=on-failure`, `TimeoutStopSec=15`
- [ ] No placeholder strings remain in any deploy file (`YOUR_SSH_PUBLIC_KEY`, `YOUR_DOMAIN`, `YOUR_VPS_IP`)
- [ ] Makefile has `deploy`, `ssh`, `logs` targets
- [ ] `make build` exits 0 and produces `bin/site` (linux/amd64)
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 12.
