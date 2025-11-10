.PHONY: dev migrate codegen test

DEV_CMD ?= go run ./cmd/api

export GO111MODULE=on

dev:
	$(DEV_CMD)

migrate:
	go run ./cmd/api -migrate-only

codegen:
	cd ../.. && pnpm codegen

test:
	go test ./...
	go test -tags=integration ./...
