GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint
GOCACHE ?= $(CURDIR)/.gocache
GOMODCACHE ?= $(CURDIR)/.gomodcache
GO_TEST ?= go test
PKG ?=
FILE ?=
RUN ?=
TEST_FLAGS ?= -count=1

OAPI_CODEGEN ?= oapi-codegen
OPENAPI_SPEC ?= ../api-schema/generated/openapi.yaml
OAPI_CONFIG ?= oapi-codegen.yaml
PROTOC ?= protoc
PROTO_DIR ?= ../proto
PROTO_OUT_DIR ?= internal/adapter/grpc/generated
MIGRATE ?= migrate
MIGRATIONS_DIR ?= migrations
NAME ?=
DB_HOST ?= db
ENV_FILE ?= .env

ifneq (,$(wildcard .env))
include .env
export $(shell grep -E '^[A-Za-z0-9_]+=' .env | cut -d= -f1)
endif

# Build DATABASE_URL after .env is loaded, allowing override from environment.
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: lint
lint:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) ./...

# Run all tests (override flags with TEST_FLAGS, e.g. TEST_FLAGS="-count=1 -v")
.PHONY: test-all
test-all:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) $(TEST_FLAGS) ./...

# Run package-level tests. Example:
#   make test-pkg PKG=./internal/domain/note
.PHONY: test-pkg
test-pkg:
	@test -n "$(PKG)" || (echo "PKG is required, e.g. make test-pkg PKG=./internal/domain/note"; exit 1)
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) $(TEST_FLAGS) $(PKG)

# Run a single test file. Example:
#   make test-file FILE=internal/domain/note/logic_test.go
.PHONY: test-file
test-file:
	@test -n "$(FILE)" || (echo "FILE is required, e.g. make test-file FILE=internal/domain/note/logic_test.go"; exit 1)
	@dir=$$(dirname $(FILE)); \
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) $(TEST_FLAGS) ./$$dir

# Run a specific test (regexp). Example:
#   make test-run PKG=./internal/domain/note RUN=TestValidateNoteForCreate
.PHONY: test-run
test-run:
	@test -n "$(PKG)" || (echo "PKG is required, e.g. make test-run PKG=./internal/domain/note RUN=TestValidateNoteForCreate"; exit 1)
	@test -n "$(RUN)" || (echo "RUN is required, e.g. make test-run PKG=./internal/domain/note RUN=TestValidateNoteForCreate"; exit 1)
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) $(TEST_FLAGS) -run $(RUN) $(PKG)

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: oapi
oapi:
	@test -f $(OPENAPI_SPEC) || (echo "OpenAPI spec not found at $(OPENAPI_SPEC). Run api-schema/scripts/generate-openapi.sh first."; exit 1)
	$(OAPI_CODEGEN) --config $(OAPI_CONFIG) $(OPENAPI_SPEC)

.PHONY: proto
proto:
	@mkdir -p $(PROTO_OUT_DIR)/accountpb
	@export PATH="$$HOME/.local/bin:$$(go env GOPATH)/bin:$$PATH"; \
	$(PROTOC) \
		--go_out=. \
		--go_opt=module=immortal-architecture-cqrs/backend \
		--go-grpc_out=. \
		--go-grpc_opt=module=immortal-architecture-cqrs/backend \
		-I .. \
		$(PROTO_DIR)/account.proto

.PHONY: build
build:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build ./...

.PHONY: build-api
build-api:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build -o bin/api ./cmd/api

.PHONY: build-grpc
build-grpc:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build -o bin/grpc ./cmd/grpc

.PHONY: run-api
run-api:
	@go run ./cmd/api/main.go

.PHONY: run-grpc
run-grpc:
	@go run ./cmd/grpc/main.go

.PHONY: migrate-up
migrate-up:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" up'

.PHONY: migrate-down
migrate-down:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" down 1'

.PHONY: migrate-new
migrate-new:
	@test -n "$(NAME)" || (echo "NAME is required, e.g. make migrate-new NAME=add_users"; exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)
