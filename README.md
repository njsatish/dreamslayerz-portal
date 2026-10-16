# Dream Slayerz Portal

Fresh React + Go implementation modeled on the proven High Quality project structure, without copying business-specific identifiers or AWS resources.

## Local development

```bash
cd frontend && npm install
cd ..
./scripts/dev.sh
```

Frontend: http://localhost:5173
API: http://localhost:8000/api/health

## Validate

```bash
./scripts/validate.sh
```

## Capture a new Booksy business catalog

The capture command creates reviewable candidate configuration. It does not overwrite `config/`.

```bash
./scripts/capture-booksy-business.sh \
  1051605 \
  "Dream Slayerz 540" \
  dreamslayerz \
  "https://booksy.com/en-us/1051605_dream-slayerz-540_professional-services_134579_roanoke"
```

Required server-side environment variables for live capture:

```text
BOOKSY_API_KEY
BOOKSY_APP_VERSION
BOOKSY_FINGERPRINT
```

Output:

```text
generated/booksy-<business-id>/
  raw-business.json
  business.candidate.json
  providers.candidate.json
  services.candidate.json
  capture-report.json
```

For offline parser testing, use:

```bash
go run ./backend/cmd/booksy-capture \
  --business-id 12345 \
  --input /path/to/saved-business-response.json
```
