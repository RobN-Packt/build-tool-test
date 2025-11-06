#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

if ! command -v pnpm >/dev/null 2>&1; then
  echo "pnpm is required for this target" >&2
  exit 1
fi

cd "$ROOT_DIR"
pnpm "$@"
