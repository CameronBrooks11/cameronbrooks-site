# Runbook

Operational reference: local development, first VPS bring-up, DNS, and ongoing deploys.

---

## Prerequisites (local machine)

| Tool    | Version | Purpose             |
| ------- | ------- | ------------------- |
| Go      | 1.22+   | Build and local run |
| Make    | any     | Task runner         |
| SSH key | —       | VPS access          |
| git     | any     | Source control      |

Verify Go version: `go version` — must be 1.22 or higher for method+pattern routing.

---

## Local development

### Run

```sh
make dev
# or directly:
go run ./cmd/site
```

Server starts on `http://localhost:8080`. No live-reload in v1 — stop and restart after template or Go changes. CSS changes take effect on browser refresh (no build step).

### Makefile targets

| Target        | What it does                                              |
| ------------- | --------------------------------------------------------- |
| `make dev`    | `go run ./cmd/site`                                       |
| `make build`  | `GOOS=linux GOARCH=amd64 go build -o bin/site ./cmd/site` |
| `make deploy` | build + scp binary to VPS + restart service               |
| `make ssh`    | open SSH session to the deploy user on the VPS            |
| `make logs`   | tail journalctl logs from the site service via SSH        |

### Environment variables (local)

The app reads from environment variables. For local dev, set them in a `.env` file (not committed) and load with `export $(cat .env)` before running, or put them directly in the shell session.

| Variable    | Example | Required | Notes                         |
| ----------- | ------- | -------- | ----------------------------- |
| `SITE_ADDR` | `:8080` | No       | Defaults to `:8080`           |
| `SITE_ENV`  | `dev`   | No       | `dev` enables verbose logging |

No secrets required for v1 (no email, no DB). Add variables here as features are added.

---

## Build

```sh
make build
```

Produces `bin/site` — a statically linked linux/amd64 binary. Templates (`internal/views/`) and static files (`static/`) are embedded into the binary via `//go:embed`. The binary is fully self-contained: deploying it is the only step required.

---

## First VPS bring-up

### Step 1 — Choose and create a VPS

Recommended providers (provider-agnostic, pick one):

- Hetzner CX11 or CAX11 (~€4/mo) — best value
- Vultr 1GB — solid, more US regions
- DigitalOcean Basic Droplet — most docs

OS: **Debian 12** (Bookworm). 1 vCPU, 1GB RAM is sufficient for a personal site.

### Step 2 — Prepare cloud-init

Edit `deploy/cloud-init.yaml`:

- Replace `YOUR_SSH_PUBLIC_KEY` with the contents of your `~/.ssh/id_ed25519.pub` (or equivalent)
- Replace `YOUR_DOMAIN` with your domain (e.g. `cameronbrooks.com`)

Paste the entire contents of `cloud-init.yaml` into the provider's **user data** field during VPS creation.

### Step 3 — Create the VPS

Submit. Provisioning takes 1–3 minutes. cloud-init runs on first boot — give it an additional 1–2 minutes after the VPS shows "running" before SSH-ing in.

### Step 4 — Verify cloud-init ran

```sh
ssh deploy@<vps-ip>
sudo cloud-init status
# should show: status: done
sudo systemctl status site
# should show: active (running)
sudo systemctl status caddy
# should show: active (running)
```

### Step 5 — Point DNS

At your DNS provider (not the VPS provider), create:

| Type | Name  | Value      | TTL |
| ---- | ----- | ---------- | --- |
| A    | `@`   | `<vps-ip>` | 300 |
| A    | `www` | `<vps-ip>` | 300 |

TTL 300 (5 min) is safe for initial setup. Raise to 3600 once stable.

> DNS propagation typically takes 1–30 minutes depending on your registrar and TTL.

### Step 6 — First deploy

Once DNS is resolving to your VPS:

```sh
make deploy
```

This builds the binary locally, scps it to `deploy@<vps-ip>:~/site`, and SSHes in to restart the systemd service.

### Step 7 — Verify HTTPS

Visit `https://<your-domain>` in a browser. Caddy handles certificate issuance automatically via Let's Encrypt on first request. The first load may take 5–10 seconds while the cert is issued — subsequent loads are instant.

Check cert: `curl -I https://<your-domain>` — should show `HTTP/2 200` and a valid TLS cert.

---

## Ongoing deploys

For any code change (Go, templates, CSS):

```sh
make deploy
```

Typical deploy time: 10–20 seconds (build + scp + restart).

The systemd service is configured with `Restart=on-failure`, so crashes auto-recover. A clean `make deploy` does a graceful restart, not a crash — existing connections are given time to complete. `deploy/site.service` sets `TimeoutStopSec=15`, giving the app's 10-second graceful shutdown context a 5-second buffer before systemd force-kills it.

---

## VPS maintenance

### View logs

```sh
make logs
# or directly:
ssh deploy@<vps-ip> "journalctl -u site -f"
# Caddy logs:
ssh deploy@<vps-ip> "journalctl -u caddy -f"
```

### SSH in

```sh
make ssh
# or:
ssh deploy@<vps-ip>
```

### Restart service manually

```sh
ssh deploy@<vps-ip> "sudo systemctl restart site"
```

### Update Caddy or system packages

```sh
ssh deploy@<vps-ip>
sudo apt update && sudo apt upgrade -y
sudo systemctl restart caddy  # if caddy was upgraded
```

Unattended security upgrades run automatically (configured via cloud-init). Caddy and application updates are manual.

---

## Secrets management

All runtime secrets are stored in a systemd drop-in env file on the server, never in the repo:

```txt
/etc/systemd/system/site.service.d/env.conf
```

Contents:

```ini
[Service]
Environment="SOME_SECRET=value"
```

After editing: `sudo systemctl daemon-reload && sudo systemctl restart site`

To add a new secret: SSH to the server, edit this file, reload. Do not put secrets in `cloud-init.yaml` or commit them to git.

---

## Re-provisioning a VPS from scratch

If you need to rebuild the server:

1. Note the current VPS IP
2. Create a new VPS with the same `cloud-init.yaml` user-data
3. Re-add secrets to `/etc/systemd/system/site.service.d/env.conf`
4. Update DNS A record to the new IP (or keep the same IP if your provider supports re-attaching)
5. Run `make deploy`
6. Verify HTTPS
7. Destroy old VPS

Total time: ~10 minutes.
