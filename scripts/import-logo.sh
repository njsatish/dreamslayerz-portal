#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SOURCE="${1:-}"
[ -n "$SOURCE" ] || { echo "Usage: $0 /path/to/logo.svg|png|jpg|jpeg|webp" >&2; exit 2; }
[ -f "$SOURCE" ] || { echo "ERROR: Logo file not found: $SOURCE" >&2; exit 1; }
cd "$ROOT"
EXT="$(printf '%s' "${SOURCE##*.}" | tr '[:upper:]' '[:lower:]')"
case "$EXT" in svg|png|jpg|jpeg|webp) ;; *) echo "ERROR: Supported logo types: svg, png, jpg, jpeg, webp" >&2; exit 1;; esac
MIME="$(file -b --mime-type "$SOURCE")"
case "$MIME" in image/svg+xml|image/png|image/jpeg|image/webp) ;; *) echo "ERROR: File content is not a supported image: $MIME" >&2; exit 1;; esac
STAMP="$(date +%Y%m%d-%H%M%S)"
mkdir -p frontend/public/assets/brand
[ ! -f frontend/public/assets/brand/logo.svg ] || cp frontend/public/assets/brand/logo.svg "/tmp/booksy-portal-logo-$STAMP.svg"
find frontend/public/assets/brand -maxdepth 1 -type f -name 'logo.*' -delete
TARGET="frontend/public/assets/brand/logo.$EXT"
cp "$SOURCE" "$TARGET"
python3 - "$EXT" "$MIME" <<'PY'
import json,sys
from pathlib import Path
ext,mime=sys.argv[1:]
p=Path('frontend/public/config/branding.json')
data=json.loads(p.read_text())
src=f'/assets/brand/logo.{ext}'
data['logo'].update({'src':src,'mode':'business'})
data['favicon'].update({'src':src,'type':mime})
p.write_text(json.dumps(data,indent=2)+'\n')
html=Path('frontend/index.html')
s=html.read_text()
import re
s=re.sub(r'<link rel="icon"[^>]*>',f'<link rel="icon" type="{mime}" href="{src}">',s,count=1)
s=re.sub(r'<link rel="apple-touch-icon"[^>]*>',f'<link rel="apple-touch-icon" href="{src}">',s,count=1)
html.write_text(s)
PY
./scripts/validate-branding.sh
echo "Logo imported once and configured for header, footer, favicon, and Apple touch icon: $TARGET"
