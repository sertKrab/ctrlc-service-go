# API response contract

Every endpoint returns the same envelope:

```json
{ "success": true, "data": {}, "error": null, "timestamp": 1700000000000 }
```

| Field | Type | Description |
|-------|------|-------------|
| success | bool | true = 2xx, false = 4xx/5xx |
| data | any | payload on success, null on error |
| error | { code, message } or null | error detail on failure |
| timestamp | int64 | Unix milliseconds |

## Frontend axios unwrap

```typescript
const res = await axios.post('/api/v1/auth/login', body)
// res.data = Response envelope
const { success, data, error } = res.data
```

## Scaffold placeholders

| Placeholder | Replaced in |
|-------------|-------------|
| `__PROJECT_NAME__` | README, Makefile, Dockerfile, log strings |
| `__PROJECT_LOWER__` | .env.example, docker-compose.yml, Makefile, Dockerfile |
| `__MODULE_PATH__` | go.mod + all *.go imports (find/replace by scaffold_setup) |
| `__DB_NAME__` | .env.example, docker-compose.yml |
| `__DB_HOST__` | .env.example |
| `__DB_PORT__` | .env.example, docker-compose.yml |
| `__API_PORT__` | .env.example, docker-compose.yml, Dockerfile |
| `__JWT_SECRET__` | .env.example |
| `__ENABLE_REDIS__` | .env.example |
