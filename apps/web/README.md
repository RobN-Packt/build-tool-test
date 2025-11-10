# Book Web UI

Minimal Next.js app for managing books via the Book API.

## Requirements

- Node.js 18+
- Running Book API (default: `http://localhost:8080`)

## Setup

```bash
cd apps/web
cp .env.example .env.local
# adjust NEXT_PUBLIC_API_BASE_URL if the API runs elsewhere
(cd ../.. && pnpm install)
```

## Development

```bash
# from repository root
pnpm --filter book-web dev
```

Visit http://localhost:3000 for the UI:
- `/` lists the books in a table.
- `/admin/new` provides a form to create a book.

## Testing

```bash
# from repository root
pnpm --filter book-web test
```

Runs Vitest with React Testing Library.

## Build

```bash
# from repository root
pnpm --filter book-web build
pnpm --filter book-web start
```

## Docker

```bash
docker build -t book-web .
docker run --rm -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL="http://host.docker.internal:8080" book-web
```

## Code Generation

The frontend consumes types generated from the OpenAPI spec. After updating
`apps/api/openapi/openapi.yaml`, run:

```bash
pnpm run codegen
```

This refreshes `lib/api/types.ts` and keeps the client in sync with the API.
