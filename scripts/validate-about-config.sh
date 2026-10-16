#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
about=json.loads(Path('frontend/public/config/about.json').read_text())
providers=json.loads(Path('config/providers.json').read_text())['providers']
provider_ids={p['id'] for p in providers}
assert about['schemaVersion']==1 and about['enabled'] is True
spot=about['spotlight']
assert spot['providerId'] in provider_ids
assert spot['paragraphs'] and spot['image']
assert spot['imageMode'] in {'demo','business'}
asset=Path('frontend/public')/spot['image'].lstrip('/')
assert asset.is_file(),f'missing about image: {asset}'
trust=about['trustCard']
assert trust['mode'] in {'demo','verified'}
if trust['mode']=='demo':
    combined=' '.join(str(v) for v in trust.values()).lower()
    assert '450' not in combined and 'review count' not in combined
print('About configuration validation passed.')
PY
