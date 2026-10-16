#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
python3 - <<'PY'
import json
from pathlib import Path
providers=json.loads(Path('config/providers.json').read_text())['providers']
services=json.loads(Path('config/services.json').read_text())['services']
ids=[provider['id'] for provider in providers]
assert len(ids)==len(set(ids)),'duplicate provider IDs'
for provider in providers:
    assert provider.get('id') and provider.get('name')
    photo=provider.get('presentation',{}).get('photo')
    if provider.get('enabled') and photo:
        assert (Path('frontend/public')/photo.lstrip('/')).is_file(),f'missing provider photo: {photo}'
used={provider_id for service in services for variant in service.get('variants',[]) for provider_id in variant.get('providerIds',[])}
unknown=used-set(ids)
assert not unknown,f'unknown provider IDs in variants: {sorted(unknown)}'
app=Path('frontend/src/App.jsx').read_text()
assert "import ProviderPage from './pages/ProviderPage'" in app
assert "value.startsWith('providers/')" in app
assert 'provider:<ProviderPage' in app
assert 'providerId={providerId}' in app
about=Path('frontend/src/pages/AboutPage.jsx').read_text()
assert 'providers.length===1' in about and 'providers.length>1' in about
profile=Path('frontend/src/pages/ProviderPage.jsx').read_text()
assert 'routeProviderId' in profile
assert 'variant.providerIds' in profile
print('Provider page validation passed.')
PY
