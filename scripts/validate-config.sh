#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
go test ./backend/internal/config -run TestLoadDreamSlayerzConfig -v
echo "Configuration validation passed."
