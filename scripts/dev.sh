#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
trap 'kill 0 2>/dev/null || true' EXIT
(cd "$ROOT" && go run ./backend/cmd/local) &
(cd "$ROOT/frontend" && npm run dev) &
wait
