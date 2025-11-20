# Book API & Web UI

Monorepo containing a Go/Huma API with PostgreSQL storage and a Next.js (TypeScript) frontend for a dummy book store.

## Prerequisites

- Go 1.24+
- Node.js 18+ with `corepack enable` (for pnpm)
- Docker & Docker Compose

## First Run

```bash
corepack pnpm install
(cd apps/api && go mod tidy)
pnpm run codegen
make migrate
```

`make migrate` applies SQL migrations using the `DB_DSN` in your environment (defaults to `postgres://postgres:postgres@localhost:5432/books?sslmode=disable`).

## Local Development (without Docker)

- API: `make dev`
- Web: `pnpm --filter book-web dev`

Ensure a PostgreSQL instance is running and reachable at the DSN noted above.

## Local Development (Docker)

```bash
docker compose up --build
```

This starts:
- `postgres` on port 5432 (`books` database)
- `api` on http://localhost:8080 (migrations run automatically)
- `web` on http://localhost:3000

## Testing

- API: `make test-api`
- Web: `make test-web`
- Full suite: `make test`

## Code Generation

The OpenAPI spec lives at `apps/api/openapi/openapi.yaml`. Regenerate Go models and TypeScript types after editing the spec:

```bash
pnpm run codegen
```

Outputs:
- `apps/api/openapi/gen.models.go`
- `apps/web/lib/api/types.ts`