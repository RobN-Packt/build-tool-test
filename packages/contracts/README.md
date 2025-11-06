# Contracts

`packages/contracts` owns the OpenAPI specification for the book purchasing API and the code generation pipeline for server and client stubs.

## Files
- `openapi.yaml`: Single source of truth API definition
- `codegen.ts`: Generates Go server types/handlers and TypeScript client bindings

## Usage

From the repository root:
```bash
pnpm --filter contracts gen
```

This script performs three steps:
1. Validates the OpenAPI document using `@apidevtools/swagger-parser`
2. Runs `oapi-codegen` to generate Go server interfaces in `apps/api/openapi` (using Go's module cache)
3. Runs `openapi-typescript` to generate types and a thin client in `apps/web/lib/api`

Generated code is committed so other build tools can consume it without requiring runtime codegen.
