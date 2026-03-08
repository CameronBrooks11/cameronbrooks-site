# Deploy Checklist

Operator checklist for production deployments. Use this with `docs/deployment.md`.

## Deploy gate

- [ ] Pre-deploy hardening is complete (tests, CI, docs, and smoke checks)
- [ ] VPS is provisioned and reachable via SSH
- [ ] DNS points to the target VPS IP

## Pre-deploy checks

- [ ] `git status --short` reviewed (only intentional changes)
- [ ] `go test ./...` passed
- [ ] `go vet ./...` passed
- [ ] `go build ./...` passed
- [ ] CI is green for the commit to deploy

## Deploy artifact checks

- [ ] `rg -n "YOUR_SSH_PUBLIC_KEY|YOUR_DOMAIN" deploy` returns no matches
- [ ] `deploy/cloud-init.yaml` parses successfully
- [ ] `make build` produced a non-empty `bin/site`
- [ ] `make smoke` passed locally

## Deploy steps

- [ ] `make deploy` executed from repo root
- [ ] `~/site.prev` exists on VPS after deploy (rollback snapshot)

## Post-deploy verification

- [ ] `curl -i https://<your-domain>/healthz` is `200` with body `ok`
- [ ] `curl -i https://<your-domain>/version` returns valid JSON metadata
- [ ] `curl -i https://<your-domain>/` is `200`
- [ ] `curl -i https://<your-domain>/writing` is `200`
- [ ] `curl -i https://<your-domain>/does-not-exist` is `404`
- [ ] `make logs` shows healthy requests and no restart loop

## Rollback (if needed)

- [ ] `ssh deploy@<vps-ip> "test -f ~/site.prev && cp ~/site.prev ~/site"`
- [ ] `ssh deploy@<vps-ip> "sudo systemctl restart site"`
- [ ] `curl -i https://<your-domain>/healthz` returns healthy result after rollback
