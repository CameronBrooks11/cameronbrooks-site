# Operations

## Local prerequisites

- Go 1.22+
- Make
- SSH key configured for deploy target

## Local commands

Run dev server:

```sh
make dev
```

Quality checks:

```sh
go test ./...
go vet ./...
go build ./...
```

Smoke checks (server running locally):

```sh
make smoke
```

## Build artifact

```sh
make build
```

Outputs `bin/site` (linux/amd64) with embedded templates and static assets.

## Deploy workflow

```sh
make deploy VPS=deploy@<vps-ip>
```

Deploy path:

1. validate VPS target is not placeholder
2. build binary
3. snapshot previous remote binary to `~/site.prev`
4. upload new binary
5. restart `site` systemd service

Detailed procedure: `docs/deployment.md`

## Logs

```sh
make logs VPS=deploy@<vps-ip>
```

Expect structured JSON request logs with request IDs.

## Secrets

Store runtime secrets on server only:

```txt
/etc/systemd/system/site.service.d/env.conf
```

After changes:

```sh
sudo systemctl daemon-reload
sudo systemctl restart site
```

## Rollback

```sh
ssh deploy@<vps-ip> "test -f ~/site.prev && cp ~/site.prev ~/site"
ssh deploy@<vps-ip> "sudo systemctl restart site"
```

Then recheck `/healthz` and key routes.

## First-time VPS setup

Use `deploy/cloud-init.yaml` during VPS creation, then follow `docs/deployment.md`.
