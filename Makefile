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
