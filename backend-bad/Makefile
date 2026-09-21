GOBIN ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOCACHE ?= $(CURDIR)/.gocache
GOMODCACHE ?= $(CURDIR)/.gomodcache
MIGRATE ?= $(GOBIN)/migrate
SQLC ?= sqlc
MIGRATIONS_DIR ?= $(CURDIR)/migrations
ENV_FILE ?= $(CURDIR)/.env

.PHONY: lint tidy migrate-create migrate-up migrate-down migrate-force sqlc build dev oapi

lint:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GOLANGCI_LINT) run ./...

tidy:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go mod tidy

sqlc:
	@$(SQLC) generate

.PHONY: build

build:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build ./...

dev:
	@$(GOBIN)/air

oapi:
	@./scripts/generate-oapi.sh

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=add_notes_table"; exit 1; fi
	@$(MIGRATE) create -seq -ext sql -dir $(MIGRATIONS_DIR) $(name)

migrate-up:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$DATABASE_URL" ]; then echo "Set DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DATABASE_URL" up'

migrate-down:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$DATABASE_URL" ]; then echo "Set DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DATABASE_URL" down 1'

migrate-force:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=20240101010101"; exit 1; fi
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$DATABASE_URL" ]; then echo "Set DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$DATABASE_URL" force $(version)'
test:
	go test ./...
