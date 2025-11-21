# Book Web UI

Minimal Next.js app for managing books via the Book API.

## Requirements

- Node.js 18+
- Running Book API (default upstream: `http://localhost:8080`, configurable via `API_SERVER_BASE_URL`)

## Setup

```bash
cd apps/web
cp .env.example .env.local
# set API_SERVER_BASE_URL to the reachable Book API URL
# optionally set NEXT_PUBLIC_API_BASE_URL if you want the browser to skip the proxy
(cd ../.. && pnpm install)
```

## API Proxy

The web app exposes `app/api/books` route handlers that forward requests to the Go API.

- Browser requests automatically target `https://<ui>/api/books/...`, eliminating mixed-content issues.
- Server components and route handlers still call `API_SERVER_BASE_URL` directly.
- If you prefer to skip the proxy (for example, when pointing to a public HTTPS load balancer), set
  `NEXT_PUBLIC_API_BASE_URL` to that absolute URL.

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
docker run --rm -p 3000:3000 \
  -e API_SERVER_BASE_URL="http://host.docker.internal:8080" \
  book-web
```

## Code Generation

The frontend consumes types generated from the OpenAPI spec. After updating
`apps/api/openapi/openapi.yaml`, run:

```bash
pnpm run codegen
```

This refreshes `lib/api/types.ts` and keeps the client in sync with the API.
