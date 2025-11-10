# Book API

Minimal CRUD API for managing books using [Huma v2](https://huma.rocks/) and PostgreSQL.

## Prerequisites

- Go 1.24+
- PostgreSQL database
- Environment variables (see `.env.example`)

## Setup

```bash
cp .env.example .env
# adjust DB_DSN if needed (defaults to books database)
```

## Running Locally

- Apply migrations: `make migrate`
- Start the API server: `make dev`

The server listens on `:$PORT` (defaults to `8080`). OpenAPI docs are available at `/openapi.json`. Health check: `/healthz`.

## Testing

```bash
# Set TEST_DB_DSN to point at a disposable database before running.
TEST_DB_DSN="postgres://postgres:postgres@localhost:5432/books_test?sslmode=disable" make test
```

The integration test truncates the `books` table after running.

## Build

```bash
make build
```

The compiled binary is written to `bin/api`.

## Docker

```bash
docker build -t book-api apps/api
docker run --rm -p 8080:8080 -e DB_DSN="postgres://..." book-api
```

## Code Generation

OpenAPI definitions live in `apps/api/openapi/openapi.yaml`. After editing the
spec, regenerate Go and TypeScript models from the repository root:

```bash
pnpm run codegen
```

This updates:
- `apps/api/openapi/gen.models.go`
- `apps/web/lib/api/types.ts`
