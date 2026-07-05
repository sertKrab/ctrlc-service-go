# __PROJECT_NAME__ Backend

Go + Gin + GORM backend. Module: `git.trovefin.com/poc/ctrlc-service-go`

## Quick start

```bash
cp .env.example .env
# Edit .env — set DB_NAME, DB_PASSWORD, JWT_SECRET

make docker-up     # start Postgres
make migrate-up    # run migrations
make seed          # seed roles, permissions, admin user
make run           # start API on :__API_PORT__
```

## Default admin
- Email: `admin@localhost`
- Password: `Admin1234!`

## API response envelope
```json
{ "success": true, "data": {}, "error": null, "timestamp": 1700000000000 }
```

## Backend test response capture
Generated backend HTTP tests should capture the primary asserted response with
the shared helper:

```go
import "git.trovefin.com/poc/ctrlc-service-go/internal/testutil"

func TestAuthHandler_AUTH_TC_001_LoginSuccess(t *testing.T) {
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    testutil.RecordBEAPIResponse(t, "AUTH_TC_001", w)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }
}
```

The helper writes `.ctrlc/e2e/api-responses/be_AUTH_TC_001.json` using the same
`{"status": <int>, "body": <json-or-text>}` shape as frontend captures, and
redacts sensitive fields such as tokens, cookies, OTPs, PINs, passwords, and
secrets.

## Auth endpoints
| Method | Path | Auth |
|--------|------|------|
| POST | /api/v1/auth/login | — |
| POST | /api/v1/auth/refresh | cookie |
| POST | /api/v1/auth/logout | Bearer |
| GET  | /api/v1/auth/me | Bearer |

## Redis (optional)
```bash
# .env: ENABLE_REDIS=true
make docker-up-redis
```

## Scaffold placeholders
Run after scaffold: `make check-placeholders`
