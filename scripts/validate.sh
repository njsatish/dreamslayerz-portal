#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
test "$(git remote get-url origin)" = "https://github.com/njsatish/dreamslayerz-portal.git"
./scripts/validate-config.sh
./scripts/validate-hero-config.sh
./scripts/validate-demo-content.sh
gofmt -w backend
go test ./...
(cd frontend && npm run build)
echo "Validation passed."
