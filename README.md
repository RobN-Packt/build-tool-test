# Bookshop CRUD Platform

## Overview

This repository contains a full-stack bookshop management platform featuring:

- **Go backend (Go 1.22)** built with Chi & Huma, automatically generating an OpenAPI 3.1 contract for the book catalogue.
- **Next.js 16 frontend (TypeScript, App Router)** offering an inventory dashboard to create, read, update, and delete books.
- **PostgreSQL** for persistent storage, provisioned and orchestrated via Docker Compose.
- **Comprehensive test coverage** including Go unit/integration tests (with optional Testcontainers) and frontend Jest tests.

## Project Structure

```
backend/   # GoFr REST API, migrations, and Go tests
frontend/  # Next.js UI with React Testing Library suites
docker-compose.yml
```

## Prerequisites

- Go 1.22 (for local backend work)
- Node.js 20+ and npm (for local frontend work)
- Docker & Docker Compose (for containerised workflows)

## Running the Stack with Docker Compose

```bash
docker-compose up --build
```

Services will be available at:

- Backend API: http://localhost:8080
- Frontend UI: http://localhost:3000
- Database: exposed internally (`db` service)

Initial seed data is automatically applied on backend startup via embedded migrations.

## Backend Development

```bash
cd backend

# Run all Go tests (integration tests auto-skip when Docker is unavailable)
GOTOOLCHAIN=local go test ./...

# Lint/format
gofmt -w $(find cmd internal tests -name '*.go')
```

### Key Environment Variables

| Variable       | Description                  | Default (docker-compose) |
| -------------- | ---------------------------- | ------------------------ |
| `DB_HOST`      | Database hostname             | `db`                     |
| `DB_PORT`      | Database port                 | `5432`                   |
| `DB_USER`      | Database user                 | `bookshop`               |
| `DB_PASSWORD`  | Database password             | `bookshop`               |
| `DB_NAME`      | Database name                 | `bookshop`               |
| `DB_SSL_MODE`  | SSL mode for Postgres         | `disable`                |
| `HTTP_PORT`    | Exposed API port              | `8080`                   |
| `PUBLIC_BASE_URL` | Base URL advertised in OpenAPI docs | `http://localhost:8080` |

### REST Endpoints

`/books`

- `GET /books` → `200` `{ "books": Book[] }`
- `GET /books/{id}` → `200` `Book`
- `POST /books` → `201` `Book`
- `PUT /books/{id}` → `200` `Book`
- `DELETE /books/{id}` → `204` *(no content)*

An OpenAPI 3.1 specification is served at `GET /openapi.json` (and YAML/HTML at `/openapi.yaml` / `/docs`), generated directly from the Huma route definitions.

## Frontend Development

```bash
cd frontend
npm install

# Launch dev server
npm run dev

# Run Jest test suites
npm test

# Lint the project
npm run lint
```

Set `NEXT_PUBLIC_API_URL` to point at the backend (defaults to `http://localhost:8080` in development, overridden via Docker Compose).

## Testing Strategy

- **Backend**
  - Unit tests cover validation, service logic, and repository behaviours.
  - Integration tests leverage Testcontainers (auto-skipped when Docker is not reachable) to run the API against a real PostgreSQL instance.

- **Frontend**
  - Jest with React Testing Library validates data loading, form submission, and optimistic UI updates.

All test commands are wired into project scripts and the CI-ready `go test`/`npm test` workflows.

## Branch & Deployment

After implementing changes:

```bash
git checkout -b feature/bookshop-app
git add .
git commit -m "feat: implement bookshop platform"
git push origin feature/bookshop-app
```

Replace the branch name as required before pushing.
