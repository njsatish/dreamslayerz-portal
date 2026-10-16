# Reusable Booksy Live Availability and Date/Time Preselection Runbook

**Created:** October 16, 2026
**Based on:** Dream Slayerz portal implementation
**Purpose:** Reuse this procedure for future booking portals so live availability, Booksy widget preselection, local testing, security, and deployment can be implemented without repeating the same debugging cycle.

---

## 1. What the completed implementation does

The working booking flow now performs the following sequence:

1. Loads portal configuration from the backend.
2. Lets the customer choose a provider and service.
3. Calls the local or deployed availability API with the selected provider, service, and Booksy service variant.
4. The backend validates that the provider, service, and variant belong together.
5. The backend calls Booksy using the business ID, variant ID, staffer ID, date range, and required Booksy request headers.
6. The frontend displays available dates and time slots returned by Booksy.
7. The customer selects a date and time.
8. The frontend opens the Booksy widget with the business, country, language, service variant, date, and time already selected.
9. Booksy displays the selected service, staff member, date, time, price, and duration before the customer continues.

### Verified working behavior

- Live availability returns dates and time slots.
- Date selection works.
- Time selection works.
- The Booksy widget opens inside the same-page modal.
- The selected service is retained.
- The selected date and time are retained.
- Closing and reopening the modal works after the permanent locale fix.

---

## 2. Architecture used

### Frontend

- React
- Vite
- Booking component similar to `frontend/src/booking/AdaptiveBooking.jsx`
- Same-page modal containing a Booksy iframe

### Backend

- Go HTTP server for local development
- Local entry point similar to `./backend/cmd/local`
- Availability handler similar to `backend/internal/api/availability.go`
- Booksy client similar to `backend/internal/booksy/client.go`

### External integration

- Booksy customer API for availability
- Booksy widget for the final booking flow

---

## 3. Required identifiers and configuration

Every new project must have verified values for all of the following.

### Business-level configuration

```json
{
  "booksy": {
    "businessId": 1051605,
    "widgetId": 1051605,
    "country": "us",
    "language": "en"
  },
  "booking": {
    "liveAvailability": true,
    "samePageWidget": true,
    "availabilityFailureMode": "widget",
    "availabilityWindowDays": 14
  }
}
```

Replace the example business and widget IDs for each project.

### Provider-level configuration

Each provider needs:

```json
{
  "id": "provider-slug",
  "enabled": true,
  "onlineBooking": true,
  "booksyStafferId": 1234567
}
```

### Service-variant configuration

Each selectable variant needs:

```json
{
  "id": 15427301,
  "providerIds": ["provider-slug"],
  "availability": {
    "enabled": true
  }
}
```

### Validation rule

Never assume that a business ID, widget ID, staffer ID, and variant ID belong together. Validate the exact combination in Booksy before writing portal configuration.

---

## 4. Required Booksy environment variables

The Go Booksy client requires these three runtime values:

```text
BOOKSY_API_KEY
BOOKSY_APP_VERSION
BOOKSY_FINGERPRINT
```

The client initialization should fail safely when any value is absent:

```go
func NewFromEnv() (*Client, error) {
    client := &Client{
        HTTP:        &http.Client{Timeout: 12 * time.Second},
        APIKey:      strings.TrimSpace(os.Getenv("BOOKSY_API_KEY")),
        AppVersion:  strings.TrimSpace(os.Getenv("BOOKSY_APP_VERSION")),
        Fingerprint: strings.TrimSpace(os.Getenv("BOOKSY_FINGERPRINT")),
    }

    if client.APIKey == "" || client.AppVersion == "" || client.Fingerprint == "" {
        return nil, errors.New("Booksy availability credentials are not configured")
    }

    return client, nil
}
```

### Important

Do not use invented placeholder values. Invalid values only change a local `503` configuration error into an upstream authentication or request failure.

---

## 5. Safe local credential file

Create a project-root file named:

```text
.env.local
```

Example structure:

```dotenv
BOOKSY_API_KEY=real-value
BOOKSY_APP_VERSION=real-value
BOOKSY_FINGERPRINT=real-value
```

### Protect the file

Add this line to `.gitignore`:

```gitignore
.env.local
```

Set restrictive permissions:

```bash
chmod 600 .env.local
```

Verify Git protection:

```bash
git check-ignore -v .env.local

if git status --short --untracked-files=all |
  grep -F '.env.local' >/dev/null
then
  echo "ERROR: .env.local is visible to Git"
  exit 1
else
  echo "PASS: .env.local is protected from Git"
fi
```

Never run `git add .` until this check passes.

---

## 6. Starting the backend correctly

Environment variables are inherited only by processes started from the same shell session.

```bash
cd "$HOME/Downloads/PROJECT_REPOSITORY" || return 1

set -a
source ./.env.local
set +a
```

Validate presence without displaying values:

```bash
FAILED=0

for variable in \
  BOOKSY_API_KEY \
  BOOKSY_APP_VERSION \
  BOOKSY_FINGERPRINT
do
  eval "value=\${$variable-}"

  if [ -n "$value" ]; then
    echo "PASS: $variable is set"
  else
    echo "ERROR: $variable is missing"
    FAILED=1
  fi
done

if [ "$FAILED" -ne 0 ]; then
  return 1
fi
```

Start the local Go server:

```bash
export GOCACHE="${HOME}/Library/Caches/go-build"
mkdir -p "$GOCACHE"

go run ./backend/cmd/local
```

Do not start a capture utility or Lambda entry point for local browser testing.

---

## 7. Starting the frontend

Open a second Terminal window:

```bash
cd "$HOME/Downloads/PROJECT_REPOSITORY/frontend" || return 1

npm run dev
```

Open the Vite address, normally:

```text
http://localhost:5173/
```

Keep both backend and frontend terminals running.

---

## 8. Backend availability request

The frontend should call an endpoint similar to:

```text
GET /api/availability
    ?provider=provider-slug
    &service=service-slug
    &variant=15427301
```

The backend should resolve and validate the provider, service, and variant before calling Booksy.

### Booksy request inputs

```go
booksy.Request{
    BusinessID: config.Business.Booksy.BusinessID,
    VariantID:  selected.Variant.ID,
    StafferID:  selected.Provider.BooksyStafferID,
    StartDate:  start.Format("2006-01-02"),
    EndDate:    end.Format("2006-01-02"),
}
```

### Booksy upstream endpoint pattern

```text
POST https://us.booksy.com/core/v2/customer_api/me/businesses/{BUSINESS_ID}/appointments/time_slots
```

### Required upstream headers

```text
Accept: application/json
Content-Type: application/json
Origin: https://booksy.com
Referer: https://booksy.com/
X-Api-Key: <BOOKSY_API_KEY>
X-App-Version: <BOOKSY_APP_VERSION>
X-Fingerprint: <BOOKSY_FINGERPRINT>
```

### Request body outline

```json
{
  "subbookings": [
    {
      "service_variant_id": 15427301,
      "staffer_id": 1001956,
      "combo_children": []
    }
  ],
  "start_date": "2026-10-16",
  "end_date": "2026-10-30"
}
```

---

## 9. Availability endpoint response handling

Recommended statuses:

### `200 OK`

Live availability worked.

```json
{
  "providerId": "provider-slug",
  "serviceSlug": "service-slug",
  "variantId": 15427301,
  "live": true,
  "dates": [],
  "slots": [],
  "notice": "Booksy confirms the appointment after you continue."
}
```

Empty arrays can be valid when there are no openings.

### `400 Bad Request`

The request is malformed, such as an invalid variant value.

### `404 Not Found`

The provider, service, and variant combination is not configured.

### `503 Service Unavailable`

Typical causes:

- `liveAvailability` is disabled.
- Variant availability is disabled.
- Required Booksy environment variables are missing.

### `502 Bad Gateway`

The Booksy client was configured, but the upstream request failed or returned a non-200 response.

---

## 10. Correct Booksy widget URL

The final iframe URL must contain all of these values:

```text
id
country
language
variantId
date
```

Recommended construction:

```jsx
const openWidget = time => {
  if (!selected) return

  const query = new URLSearchParams({
    country: 'us',
    language: 'en',
    variantId: String(selected.variant.id),
    attribution_source: `${cfg.business.project.slug}_portal`,
    utm_medium: time ? 'timeslots' : 'booking'
  })

  if (time) {
    query.set('date', `${date}T${time}`)
  }

  setWidget({
    url:
      `https://booksy.com/widget/index.html` +
      `?id=${encodeURIComponent(cfg.business.booksy.widgetId)}` +
      `&${query.toString()}`,
    name: selected.service.name
  })
}
```

### Critical lesson

The widget can display a blank white screen when `country` is missing or invalid. The browser console reported:

```text
Value of param "country" is not set or invalid!
```

Adding these permanently resolved the blank widget:

```text
country=us
language=en
```

A browser-console modification is temporary. Closing the modal destroys the iframe, and reopening it regenerates the URL from source code. Therefore, fix the URL construction in the React source, not only in DevTools.

---

## 11. Same-page modal implementation

Example:

```jsx
{widget && (
  <div className="widget-modal" role="dialog" aria-modal="true">
    <button
      className="backdrop"
      aria-label="Close"
      onClick={() => setWidget(null)}
    />

    <div className="widget-dialog">
      <header>
        <strong>{widget.name}</strong>
        <button
          ref={closeRef}
          onClick={() => setWidget(null)}
          aria-label="Close booking"
        >
          ×
        </button>
      </header>

      <iframe
        title={`Book ${widget.name}`}
        src={widget.url}
        allow="payment"
      />
    </div>
  </div>
)}
```

### Modal behavior checklist

- Set `role="dialog"`.
- Set `aria-modal="true"`.
- Provide a labeled close button.
- Close on Escape.
- Prevent background-page scrolling while open.
- Restore scrolling when closed.
- Give the iframe a meaningful title.
- Recreate the iframe from a permanently correct URL each time.

---

## 12. Fast diagnostic sequence

Use this sequence before modifying code.

### A. Test the backend directly

```bash
curl -sS -i \
  "http://localhost:8000/api/availability?provider=PROVIDER&service=SERVICE&variant=VARIANT"
```

Interpret the status before touching the frontend.

### B. Confirm the server owns port 8000

```bash
lsof -nP -iTCP:8000 -sTCP:LISTEN
```

### C. Inspect the iframe URL

In Chrome Console:

```javascript
document.querySelector('iframe')?.src
```

### D. Inspect required widget parameters

```javascript
const frame = document.querySelector('iframe')
const url = frame ? new URL(frame.src) : null

console.log({
  id: url?.searchParams.get('id'),
  country: url?.searchParams.get('country'),
  language: url?.searchParams.get('language'),
  variantId: url?.searchParams.get('variantId'),
  date: url?.searchParams.get('date')
})
```

### E. Temporary locale proof

Use only as a diagnostic, not as the permanent fix:

```javascript
const frame = document.querySelector('iframe')

if (frame) {
  const url = new URL(frame.src)
  url.searchParams.set('country', 'us')
  url.searchParams.set('language', 'en')
  frame.src = url.href
}
```

If the widget loads after this, make the same change in source code.

### F. Locate the real URL construction

```bash
cd "$HOME/Downloads/PROJECT_REPOSITORY" || return 1

grep -RInE \
  --include='*.jsx' \
  --include='*.js' \
  --include='*.tsx' \
  --include='*.ts' \
  --exclude-dir='node_modules' \
  'booksy\.com/widget|widget/index\.html|attribution_source|utm_medium|variantId' \
  frontend/src
```

Do this before writing an automated patch. Do not assume variable names or formatting.

---

## 13. Failures encountered and what each one meant

### Error: `live availability is not configured`

Cause:

```text
BOOKSY_API_KEY
BOOKSY_APP_VERSION
BOOKSY_FINGERPRINT
```

were absent from the backend process.

Fix:

- Create a protected `.env.local`.
- Load it into the shell.
- Start the Go server from that same shell.

### Error: dates and times load, but final widget is blank

Cause:

The Booksy iframe URL was missing:

```text
country=us
```

Fix:

Add permanent country and language parameters to the source-generated widget URL.

### Symptom: DevTools has no iframe

Possible causes:

- The modal is not open.
- The widget element was not inserted.
- The iframe lives in another component or shadow root.

Action:

- Keep the modal open.
- Inspect Elements.
- Search the repository for `widget/index.html`.

### Symptom: a browser-console URL change works only once

Cause:

The change affected only the current DOM node. Closing the modal destroyed that node.

Fix:

Modify the React source that constructs the URL.

### Error: patch script cannot find code

Cause:

The source is formatted differently than the patch assumed. In this project, `openWidget` was compressed onto one line.

Fix:

- Inspect the exact source first.
- Match the real source structure.
- Require an exact single replacement.
- Validate the result and restore the backup if validation fails.

### Error: `stat ... backend/cmd/server: directory not found`

Cause:

The assumed local server package was wrong.

Fix:

Find Go main packages:

```bash
grep -RIl \
  --include='*.go' \
  '^package main$' \
  . \
  --exclude-dir='.git' \
  --exclude-dir='node_modules'
```

For this implementation, the correct local package was:

```text
./backend/cmd/local
```

---

## 14. Security requirements

- Never commit Booksy credentials.
- Never print credential values in logs or chat.
- Redact values when checking files.
- Use file permission `600` for local secret files.
- Prefer AWS Secrets Manager or Systems Manager Parameter Store for deployed applications.
- Do not store plaintext secrets in deployment scripts, backup JSON files, or Git history.
- Rotate credentials if they have been copied broadly.
- Use a dedicated IAM role or user for deployment instead of root credentials.

---

## 15. Pre-commit validation

Run from the repository root:

```bash
git status --short

git diff --check

git diff -- \
  frontend/src/booking/AdaptiveBooking.jsx
```

Confirm `.env.local` is ignored:

```bash
git check-ignore -v .env.local
```

Run frontend validation:

```bash
cd frontend
npm run build
```

Run Go tests if present:

```bash
cd ..
go test ./...
```

Stage only intended files. Do not stage timestamped backups or `.env.local`.

---

## 16. New-project implementation checklist

### Configuration

- [ ] Business ID verified
- [ ] Widget ID verified
- [ ] Country set to `us`
- [ ] Language set to `en`
- [ ] Provider slug verified
- [ ] Booksy staffer ID verified
- [ ] Service slug verified
- [ ] Variant ID verified
- [ ] Provider-to-variant relationship verified
- [ ] `liveAvailability` enabled
- [ ] Variant availability enabled
- [ ] Availability window set between 1 and 31 days

### Credentials

- [ ] `BOOKSY_API_KEY` available
- [ ] `BOOKSY_APP_VERSION` available
- [ ] `BOOKSY_FINGERPRINT` available
- [ ] `.env.local` protected from Git
- [ ] No secrets printed in logs
- [ ] Backend started from the credential-loaded shell

### Backend

- [ ] Correct local Go package identified
- [ ] Availability endpoint returns `200`
- [ ] Provider/service/variant validation works
- [ ] Upstream errors return `502`
- [ ] Missing configuration returns `503`
- [ ] No-store caching header applied to live availability

### Frontend

- [ ] Dates render
- [ ] Times filter by selected date
- [ ] Time selection works
- [ ] Widget URL includes business ID
- [ ] Widget URL includes `country=us`
- [ ] Widget URL includes `language=en`
- [ ] Widget URL includes variant ID
- [ ] Widget URL includes `YYYY-MM-DDTHH:MM`
- [ ] Modal closes and reopens correctly
- [ ] Booksy continues to customer details

### Final validation

- [ ] Frontend production build passes
- [ ] Go tests pass
- [ ] Git diff check passes
- [ ] Only intended files are staged
- [ ] Local credentials and backups are not staged
- [ ] Development deployment passes before production deployment

---

## 17. Recommended reusable workflow

For every future Booksy portal:

1. Capture and verify business, provider, staffer, service, and variant identifiers.
2. Configure the portal without credentials in source control.
3. Store local credentials in an ignored `.env.local` file.
4. Start the backend with the environment loaded.
5. Test the availability endpoint directly before starting frontend debugging.
6. Confirm live dates and slots.
7. Build the Booksy widget URL with `id`, `country`, `language`, `variantId`, and `date`.
8. Inspect the real iframe URL in DevTools.
9. Test closing and reopening the modal.
10. Run builds and tests.
11. Review the exact Git diff.
12. Commit only intended source and configuration changes.
13. Deploy to a development environment first.
14. Validate the deployed availability API and widget.
15. Promote to production only after the full booking flow passes.

---

## 18. Final reference URL pattern

```text
https://booksy.com/widget/index.html
?id={WIDGET_ID}
&country=us
&language=en
&variantId={VARIANT_ID}
&attribution_source={PROJECT_SLUG}_portal
&utm_medium=timeslots
&date={YYYY-MM-DD}T{HH:MM}
```

Example shape:

```text
https://booksy.com/widget/index.html?id=1051605&country=us&language=en&variantId=15427301&attribution_source=dreamslayerz_portal&utm_medium=timeslots&date=2026-10-20T09%3A15
```

The order of query parameters does not matter. Their presence and valid values do.

---

## 19. Project-specific change summary for Dream Slayerz

The Dream Slayerz implementation required these changes and runtime steps:

1. Enabled live availability in the portal configuration.
2. Configured Booksy business and widget ID `1051605`.
3. Configured the provider Booksy staffer ID.
4. Configured service variant IDs and provider relationships.
5. Added and loaded the required Booksy environment variables.
6. Started the correct local backend package:

```bash
go run ./backend/cmd/local
```

7. Confirmed the local availability endpoint returned dates and time slots.
8. Added date/time selection to the React booking flow.
9. Passed the chosen date and time through the Booksy widget URL.
10. Added permanent widget locale query parameters:

```text
country=us
language=en
```

11. Verified service, provider, date, time, price, and duration in the Booksy widget.
12. Confirmed that closing and reopening the modal regenerates a working widget.

---

**End of runbook**
