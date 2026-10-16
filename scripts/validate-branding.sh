#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
config=json.loads(Path('frontend/public/config/branding.json').read_text())
assert config['schemaVersion']==1
assert config['logo']['mode'] in {'demo','business'}
assert config['logo']['src']==config['favicon']['src'], 'one logo file must drive logo and favicon'
asset=Path('frontend/public')/config['logo']['src'].lstrip('/')
assert asset.is_file(),f'missing brand asset: {asset}'
html=Path('frontend/index.html').read_text()
assert config['favicon']['src'] in html
print('Branding validation passed: one asset drives header, footer, and favicon.')
PY
