# Deployment Guide

Canonical guide for provisioning the VPS, cutting DNS, deploying the site, and validating production behavior.

## Prerequisites

- `deploy/cloud-init.yaml` contains your real SSH public key and real domain (no `YOUR_SSH_PUBLIC_KEY` or `YOUR_DOMAIN` placeholders)
- You have access to your VPS provider account and DNS registrar
- Local machine has working `go`, `make`, `ssh`, and `scp`
- Local checks pass:

```sh
go test ./...
go vet ./...
go build ./...
```

## 1) Provision the VPS

Use Debian 12, 1 vCPU / 1 GB RAM, and pass the full contents of `deploy/cloud-init.yaml` as user-data at server creation.

After the server shows running, wait 2-3 minutes, then SSH in:

```sh
ssh deploy@<vps-ip>
```

Validate bootstrapping:

```sh
sudo cloud-init status
sudo systemctl status caddy
sudo systemctl status site
sudo ufw status
sudo ss -tlnp | grep -E "80|443"
```

Expected:

- `cloud-init` is `done`
- `caddy` is `active (running)`
- `site` is either `active` or restart-looping on placeholder binary before first real deploy
- firewall allows 22/80/443
- caddy listens on 80 and 443

## 2) Point DNS

Set A records to your VPS IP:

- `@ -> <vps-ip>`
- `www -> <vps-ip>`

Use TTL `300` initially. Verify propagation:

```sh
nslookup cameronbrooks.net
nslookup www.cameronbrooks.net
```

## 3) Set deploy target locally

Either set `VPS` in `Makefile` or pass it at command time:

```sh
make deploy VPS=deploy@<vps-ip>
```

`make deploy` now does the full safe path:

1. validates `VPS` is not placeholder
2. builds linux/amd64 binary
3. snapshots current remote binary to `~/site.prev` (if present)
4. uploads new binary to `~/site`
5. restarts systemd service

## 4) First deploy

```sh
make deploy VPS=deploy@<vps-ip>
```

Confirm service health:

```sh
ssh deploy@<vps-ip> "sudo systemctl status site --no-pager"
curl -i https://cameronbrooks.net/healthz
curl -i https://cameronbrooks.net/version
```

## 5) Production smoke checks

Route checks:

```sh
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/projects
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/writing
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/about
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/contact
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/does-not-exist
```

Expected: `200` for valid routes, `404` for `does-not-exist`.

Static assets:

```sh
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/static/css/main.css
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/static/htmx.min.js
curl -o /dev/null -s -w "%{http_code}\n" https://cameronbrooks.net/static/js/progress.js
```

TLS and headers:

```sh
curl -I https://cameronbrooks.net
```

Expect `HTTP/2 200` and security headers from Caddy (`strict-transport-security`, `x-content-type-options`, `x-frame-options`, `referrer-policy`).

## 6) Log verification

```sh
make logs VPS=deploy@<vps-ip>
```

Confirm request logs are structured JSON and include expected fields (`method`, `path`, `status`, `duration_ms`, `request_id`, `remote_ip`).

## 7) Rollback (if needed)

```sh
ssh deploy@<vps-ip> "test -f ~/site.prev && cp ~/site.prev ~/site"
ssh deploy@<vps-ip> "sudo systemctl restart site"
curl -i https://cameronbrooks.net/healthz
```

If `~/site.prev` is missing, deploy the previous known-good commit from git.

## Related docs

- `docs/deploy_checklist.md` for operator signoff checklist
- `docs/operations.md` for local-dev and maintenance reference
