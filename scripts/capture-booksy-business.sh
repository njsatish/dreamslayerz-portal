#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
if [ $# -lt 1 ]; then
  echo "Usage: $0 BUSINESS_ID [BUSINESS_NAME] [PROJECT_SLUG] [PROFILE_URL]" >&2
  exit 2
fi
BUSINESS_ID="$1"; BUSINESS_NAME="${2:-}"; PROJECT_SLUG="${3:-}"; PROFILE_URL="${4:-}"
ARGS=(--business-id "$BUSINESS_ID")
[ -z "$BUSINESS_NAME" ] || ARGS+=(--business-name "$BUSINESS_NAME")
[ -z "$PROJECT_SLUG" ] || ARGS+=(--project-slug "$PROJECT_SLUG")
[ -z "$PROFILE_URL" ] || ARGS+=(--profile-url "$PROFILE_URL")
go run ./backend/cmd/booksy-capture "${ARGS[@]}"
