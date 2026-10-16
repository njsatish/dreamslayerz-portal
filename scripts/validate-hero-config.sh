#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
hero=json.loads(Path('frontend/public/config/hero.json').read_text())
assert hero['mode'] in {'auto','brand','service','provider','gallery','mixed'}
assert hero['slides'] and len({s['id'] for s in hero['slides']})==len(hero['slides'])
for slide in hero['slides']:
    assert slide['type'] in {'service','provider','gallery','booking'}
    assert slide['image'] and slide['title'] and slide['action']['target']
    asset=Path('frontend/public')/slide['image'].lstrip('/')
    assert asset.is_file(), f'missing hero asset: {asset}'
print('Hero configuration validation passed.')
PY
