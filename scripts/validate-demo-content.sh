#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
reviews=json.loads(Path('frontend/public/config/reviews.json').read_text())
gallery=json.loads(Path('frontend/public/config/gallery.json').read_text())
assert reviews['mode'] in {'demo','verified'}
assert gallery['mode'] in {'demo','business'}
ids=set()
for r in reviews['reviews']:
    assert r['id'] not in ids and 1 <= r['rating'] <= 5
    ids.add(r['id'])
    if reviews['mode']=='demo': assert r['verified'] is False and r['source']=='demo'
ids=set()
for item in gallery['items']:
    assert item['id'] not in ids and item['src']
    ids.add(item['id'])
    if gallery['mode']=='demo': assert item['generated'] is True
print('Demo content validation passed.')
PY
