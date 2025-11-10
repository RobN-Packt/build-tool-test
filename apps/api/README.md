# Book API

Minimal CRUD API for managing books using [Huma v2](https://huma.rocks/) and PostgreSQL.

## Prerequisites

- Go 1.24+
- PostgreSQL database
- Environment variables (see `.env.example`)

## Setup

```bash
cp .env.example .env
# adjust DB_DSN if needed
```

## Run the API

```bash
make dev
```

The server listens on `:$PORT` (defaults to `8080`). OpenAPI docs are available at `/openapi.json`. Health check: `/healthz`.

## Migrations

```bash
make migrate
```

This applies embedded SQL migrations to the database specified by `DB_DSN`.

## Testing

- Unit and integration tests:

  ```bash
  # Set TEST_DB_DSN to point at a disposable database before running.
  TEST_DB_DSN="postgres://postgres:postgres@localhost:5432/bookapi_test?sslmode=disable" make test
  ```

  The integration test truncates the `books` table after running.

## Build

```bash
make build
```

The compiled binary is written to `bin/api`.

## Docker

```bash
docker build -t book-api .
docker run --rm -p 8080:8080 -e DB_DSN="postgres://..." book-api
```

