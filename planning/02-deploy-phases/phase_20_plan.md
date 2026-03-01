# Phase 20 — VPS Provisioning

**Goal:** Destroy the existing Vultr VPS (the old boilerplate site, committed to Git and safe to nuke), recreate it with `deploy/cloud-init.yaml` as user-data, SSH in to verify cloud-init completed successfully, and update the `cameronbrooks.net` DNS A records to the new VPS IP. At the end of this phase you have a fresh server ready to receive the first deploy — no application code is running yet (the placeholder binary loop-restarts until Phase 21 replaces it, which is expected and harmless).

**Exit gate:** `sudo cloud-init status` → `done`; `sudo systemctl status site` → `active (running)` (or `activating` / restarting on the placeholder binary — both OK); `sudo systemctl status caddy` → `active (running)`; `dig +short cameronbrooks.net` resolves to the VPS IP from your local machine.

---

## Prerequisites

- Phase 11 complete: `deploy/cloud-init.yaml` is written with your real SSH public key substituted in and `YOUR_DOMAIN` replaced with `cameronbrooks.net` (no placeholder strings remain)
- Vultr account with the existing VPS (Debian, currently running the old boilerplate site)
- Access to the `cameronbrooks.net` DNS settings at your registrar

---

## Step 1 — Destroy the existing Vultr VPS

The existing VPS is running the old boilerplate site. The source is committed to GitHub — the server is safe to delete.

1. Log in to the [Vultr console](https://my.vultr.com)
2. Navigate to **Products → Compute** and select the current server
3. Go to **Settings → Danger Zone** (or use the server actions menu) → **Destroy Server**
4. Confirm the destruction

Note the current VPS IP before destroying — DNS may already point there, but the IP will change with the new server (Vultr does not re-assign IPs automatically after destroy/recreate). You will update DNS in Step 7 with the new IP.

> The old server's SSH host key will change. If your local `~/.ssh/known_hosts` has an entry for the old IP, you may get a host key warning when you SSH to the new server at the same IP range. Remove the old entry with `ssh-keygen -R <old-ip>` if needed.

---

## Step 2 — Copy your SSH public key

On your local machine, read your public key:

```sh
# PowerShell
Get-Content ~/.ssh/id_ed25519.pub
```

Copy the entire line (starts with `ssh-ed25519 AAAA...` and ends with your email or a comment). This should already be in `deploy/cloud-init.yaml` under `ssh_authorized_keys` from Phase 11.

If you do not have an SSH key:

```sh
ssh-keygen -t ed25519 -C "your@email.com"
# Accept default location (~/.ssh/id_ed25519)
# Set a passphrase if desired
```

---

## Step 3 — Create the new Vultr VPS with cloud-init user-data

1. In the Vultr console → **Deploy New Server** (or **+ Deploy**)
2. **Server type:** Cloud Compute (shared CPU) is sufficient for a personal site
3. **Location:** same region as before (or whichever suits you)
4. **Image:** Debian 12 x64
5. **Plan:** 1 vCPU, 1GB RAM ($6/mo) is more than enough
6. **Additional Features → User Data** — enable and paste the entire contents of `deploy/cloud-init.yaml`
7. **SSH Keys** — optionally add your key here too (cloud-init already installs it, adding it here is harmless)
8. **Server Hostname:** `cameronbrooks-site`
9. Click **Deploy Now**

**Record the new VPS IP address** — it will be different from the old one. You will need it for DNS (Step 7) and for the Makefile `VPS` variable (Step 6).

---

## Step 4 — Wait for provisioning

After the VPS shows **Running** in the provider dashboard:

1. Wait an additional **2–3 minutes** for cloud-init to complete. cloud-init runs on first boot and installs Caddy, configures UFW, and starts services. SSH-ing in too early may interrupt it.

2. Check cloud-init status via console (most providers offer a web console/VNC) or just wait the full 3 minutes before SSHing.

---

## Step 5 — SSH in and verify

```sh
ssh deploy@<vps-ip>
```

If you get a connection refused, wait another minute and retry — the firewall or sshd may still be initialising.

Once in, run each of these in order:

### Verify cloud-init completed

```sh
sudo cloud-init status
```

Expected: `status: done`

If it shows `status: running`, wait 1–2 minutes and re-run.

If it shows `status: error`, check the log:

```sh
sudo cat /var/log/cloud-init-output.log | tail -50
```

Common causes of failure:

- Package download error (retry by running the failed command manually)
- Syntax error in cloud-init.yaml (rare if you validated YAML in Phase 11)

### Verify Caddy is running

```sh
sudo systemctl status caddy
```

Expected: `active (running)`

If it shows `failed`, check:

```sh
sudo journalctl -u caddy --no-pager | tail -30
```

Common cause: syntax error in Caddyfile (e.g. wrong domain format). Fix the Caddyfile at `/etc/caddy/Caddyfile` and restart: `sudo systemctl restart caddy`.

### Verify site service started (may be restarting — that is expected)

```sh
sudo systemctl status site
```

Expected at this stage: **`active (running)`** or **`activating (auto-restart)`** — both are OK. The service is running the placeholder binary (`/home/deploy/site` created by cloud-init as an empty file), which exits immediately and triggers a restart loop. `Restart=on-failure` and `RestartSec=5` will keep it retrying. This is harmless — Phase 21 replaces the binary.

If the service shows `failed` with many restart attempts, check:

```sh
sudo journalctl -u site --no-pager | tail -20
```

Expected output: repeated `started → exited → restarting` cycle. That is correct.

### Verify UFW firewall

```sh
sudo ufw status
```

Expected:

```
Status: active

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere
80/tcp                     ALLOW       Anywhere
443/tcp                    ALLOW       Anywhere
```

### Verify Caddy is listening on 443

```sh
sudo ss -tlnp | grep -E "80|443"
```

Expected: `caddy` bound to `0.0.0.0:80` and `0.0.0.0:443`.

If DNS is not pointing to this server yet, Caddy will show a TLS certificate error or challenge-in-progress state — that is normal at this stage.

---

## Step 6 — Update Makefile with VPS IP

Back on your local machine, open `Makefile` and update the `VPS` variable with the actual IP:

```makefile
VPS ?= deploy@<your-actual-vps-ip>
```

Verify the SSH shortcut works:

```sh
make ssh
# Should open a shell on the VPS — exit when done
```

---

## Step 7 — Update DNS A records

At your registrar, update the existing A records for `cameronbrooks.net` to point to the **new** VPS IP (the old IP is now gone):

| Type | Name  | Value      | TTL |
| ---- | ----- | ---------- | --- |
| A    | `@`   | `<vps-ip>` | 300 |
| A    | `www` | `<vps-ip>` | 300 |

Use TTL 300 (5 minutes) — easy to adjust if needed. Raise to 3600 once the site is stable.

> If your DNS is managed through Cloudflare, set the proxy status to **DNS only (grey cloud)** for both records. This lets Caddy complete the ACME/Let's Encrypt HTTP-01 challenge directly on port 80. You can enable the Cloudflare proxy later if desired, but it conflicts with Caddy's automatic TLS issuance.

---

## Step 8 — Wait for DNS propagation and verify

DNS propagation typically takes 1–30 minutes. Check from your local machine:

```sh
# PowerShell
Resolve-DnsName cameronbrooks.net -Type A
```

Or:

```sh
nslookup cameronbrooks.net
```

Expected: the new VPS IP address in the answer.

For a second opinion (bypasses local DNS cache):

```sh
nslookup cameronbrooks.net 8.8.8.8
```

Once both `cameronbrooks.net` and `www.cameronbrooks.net` resolve to the new VPS IP, DNS is ready.

---

## Exit gate checklist

Run these checks from inside the VPS via SSH:

- [ ] `sudo cloud-init status` → `status: done`
- [ ] `sudo systemctl status caddy` → `active (running)`
- [ ] `sudo systemctl status site` → `active (running)` or restart loop (both OK at this stage)
- [ ] `sudo ufw status` → 22/80/443 ALLOW, all others DENY
- [ ] `sudo ss -tlnp` → Caddy bound on 80 and 443

Run these checks from your local machine:

- [ ] `make ssh` connects to VPS without password prompt
- [ ] `Resolve-DnsName cameronbrooks.net` resolves to new VPS IP
- [ ] `Resolve-DnsName www.cameronbrooks.net` resolves to new VPS IP

All boxes checked → proceed to Phase 21.
