# Evaluation Checklist

Use this checklist to compare build tooling and CI/CD options. Capture notes for both cold and warm runs.

## Build Tooling

- **Bazel**
  - Cold build duration (`bazel build //:build_all`)
  - Warm build duration with cache
  - Remote cache configuration tested? (Y/N)
  - Key DX notes (error output, incremental workflow)

- **Nx**
  - Cold run `pnpm nx run-many --target=build`
  - Warm run after no-op change
  - Task graph visualization helpful? (Y/N)
  - Executor extensibility observations

- **Task**
  - Cold `task build`
  - Warm `task build` with no changes
  - Concurrency usage / need for manual orchestration
  - Simplicity vs automation trade-offs

## Testing
- Go unit/integration test duration (`bazel test //apps/api:tests` or `task test:api`)
- Frontend unit test duration (`pnpm --filter web test`)
- E2E test runtime (Playwright)
- Worker unit test duration
- Flake rate or stability issues observed

## CI/CD Metrics
- GitHub Actions total runtime (branch vs main)
- Buildkite pipeline duration per step
- CodeBuild total runtime
- Terraform apply duration
- ECS deployment duration (service stable)
- Lambda update duration

## Caching & Artifacts
- Which caches provided biggest wins (Go module, pnpm, Bazel)?
- Artifact storage locations (ECR, S3) validated?
- Any cache invalidation issues encountered?

## Developer Experience Notes
- Ease of onboarding (docs, scripts)
- Local vs CI parity
- Observed pain points (tooling conflicts, dependency resolution)
- Suggested improvements or future experiments
