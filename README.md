# Monorepo PoC: Book Purchasing Platform

This repository contains a production-like proof of concept for a book purchasing platform. It is designed to evaluate three build systems (Bazel, Nx, Task) and three CI/CD providers (GitHub Actions, Buildkite, AWS CodeBuild/CodePipeline) while remaining runnable locally.

## Contents
- `apps/api`: GoFr-based CRUD API for managing books
- `apps/web`: Next.js App Router frontend
- `apps/worker`: Node.js AWS Lambda processor for purchase events
- `packages/contracts`: OpenAPI contract and code generators
- `tooling`: Bazel, Nx, and Task orchestrations
- `infra`: Terraform infrastructure for AWS (ECS, Lambda, API Gateway, CloudFront)
- CI/CD definitions for GitHub Actions, Buildkite, and AWS CodeBuild/CodePipeline

## Prerequisites
- Go 1.22+
- Node.js 20+, pnpm 9+
- Docker + Docker Compose
- Terraform 1.6+
- AWS CLI v2 (for deploys)
- Bazelisk, Nx CLI, Task (see tooling sections)

## First Run Quickstart
```bash
pnpm install
go mod tidy ./...
pnpm --filter contracts gen
pnpm --recursive build
```

The generated clients (TypeScript + Go) are committed, but running `pnpm --filter contracts gen` refreshes them after contract changes.

### Local Development

#### Task (recommended)
```bash
task gen
task build
task test
task up
```

Visit `http://localhost:3000` for the web app and `http://localhost:8080/healthz` for the API. The worker listens for SQS messages via the ElasticMQ emulator.

#### Bazel
```bash
bazel build //:build_all
bazel test //...
```

#### Nx
```bash
pnpm nx graph
pnpm nx run-many --target=build
pnpm nx run-many --target=test
```

## Application Overview

- **Domain**: `Book { id, title, author, price, currency, stock, created_at, updated_at }`
- **API**: CRUD endpoints plus `POST /books/{id}/purchase` to enqueue purchase processing.
- **Worker**: Processes purchase messages (SQS), decrements stock via API, emits log events.
- **Frontend**: Uses OpenAPI-generated client to list, create, edit, and delete books. Includes simple admin form.

## Testing
- Go: `go test ./apps/api/...`
- Frontend unit tests: `pnpm --filter web test`
- Frontend e2e tests: `pnpm --filter web test:e2e`
- Worker tests: `pnpm --filter worker test`

Task orchestrates these via `task test`. Bazel and Nx targets invoke the same commands under their respective orchestration models.

## Build Tool Comparison

- **Bazel**: Rules for Go (`rules_go`), shell wrappers for JS builds/tests. Root `//:build_all` aggregates the main deliverables. Gains reproducibility and remote cache options.
- **Nx**: Manages TypeScript workspaces natively; custom executor bridges Go builds/tests. Provides dependency graph, caching, and affected commands.
- **Task**: Lightweight task runner for local workflows. Builds on existing project scripts for quick iteration.

See `EVALUATION_CHECKLIST.md` for benchmarking guidance and evaluation criteria.

## Docker Compose

Run the full stack locally:
```bash
docker-compose up --build
```

Services:
- `api`: GoFr API on port 8080
- `web`: Next.js app on port 3000 (depends on API)
- `worker`: Purchase processor
- `queue`: ElasticMQ providing an SQS-compatible endpoint

Environment variables are defined in `.env.example` files per app. Copy them to `.env` (or export in shell) before starting.

## Deployment

### GitHub Actions (default path)
1. Open a PR to run build/test matrix (`.github/workflows/ci.yml`).
2. Merge to `main` to trigger `deploy.yml`:
   - Authenticates via OIDC
   - Runs Terraform (`infra/terraform`)
   - Builds and pushes Docker images to ECR
   - Updates ECS service (`poc-api`) and Lambda (`poc-purchase-worker`)
   - Invalidates CloudFront for the web frontend

### Buildkite & CodePipeline
Alternative pipelines mirror the same stages. See `.buildkite/pipeline.yml`, `buildspec.yml`, and `codepipeline-template.yaml` for details.

## AWS Configuration
- `AWS_REGION`: `eu-west-1`
- `AWS_ACCOUNT_ID`: `123456789012`
- ECR repositories: `poc/api`, `poc/web`, `poc/worker`
- ECS cluster/service: `poc-cluster` / `poc-api`
- Lambda function: `poc-purchase-worker`
- Frontend uses `NEXT_PUBLIC_API_BASE_URL` to target the API

## Repository Conventions
- TypeScript projects managed via pnpm workspaces (`pnpm-workspace.yaml`).
- Go modules kept within each app to keep dependencies isolated.
- Generated code committed (`apps/web/lib/api`, `apps/api/openapi`). Regenerate via `pnpm --filter contracts gen`.
- Terraform state not committed; configure remote backend before production usage.

## Contributing
1. Create feature branch
2. Run `task gen` if contracts change
3. Ensure `task test` passes locally
4. Commit with descriptive messages
5. Open PR; include build tool results if evaluating performance

## Troubleshooting
- Run `task clean` to remove build artifacts across tools.
- Ensure Docker resources are pruned if containers fail to rebuild.
- For Bazel caching, configure `.bazelrc` (stub provided under `tooling/bazel`).

## License
MIT (see `LICENSE` if added)