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
