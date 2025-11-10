PNPM ?= pnpm

.PHONY: migrate dev test test-api test-web codegen

migrate:
	@cd apps/api && make migrate

dev:
	@cd apps/api && make dev

test: test-api test-web

test-api:
	@cd apps/api && make test

test-web:
	@$(PNPM) --filter book-web test

codegen:
	@$(PNPM) run codegen
