GO ?= go
DB_DSN ?= postgres://postgres:postgres@localhost:5432/books?sslmode=disable
BIN_DIR ?= bin
COMPOSE ?= docker compose

.PHONY: dev migrate test build db-up db-down

dev:
	@DB_DSN=$(DB_DSN) PORT=$${PORT:-8080} $(GO) run ./cmd/api

migrate:
	@DB_DSN=$(DB_DSN) $(GO) run ./cmd/api -migrate

test:
	@$(GO) test ./...

build:
	@mkdir -p $(BIN_DIR)
	@DB_DSN=$(DB_DSN) $(GO) build -o $(BIN_DIR)/api ./cmd/api

db-up:
	@$(COMPOSE) up -d postgres

db-down:
	@$(COMPOSE) down

