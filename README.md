# BN Mobile Go Backend

Backend service for BN Mobile using:
- Gin (HTTP framework)
- GORM (ORM)
- PostgreSQL
- JWT access token + refresh token rotation
- Google OAuth login via Google ID Token verification
- Versioned SQL migrations with `golang-migrate`

## Architecture

Structure follows the architecture style from `belajar-go`:

```text
server/
├── cmd/
├── configs/
├── internal/
│   ├── delivery/
│   │   ├── data/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── router/
│   ├── domain/
│   │   ├── models/
│   │   ├── repositories/
│   │   └── services/
│   └── infrastructure/
│       └── database/
└── pkg/
```

## Endpoints

- `GET /health`
- `POST /api/v1/auth/google`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me` (Bearer token required)
- `GET /api/v1/user` (Bearer token required)
- `PATCH /api/v1/user` (Bearer token required)
- `GET /api/v1/progress` (Bearer token required)
- `POST /api/v1/progress` (Bearer token required)
- `GET /api/v1/progress/:module/:itemId` (Bearer token required)
- `GET /api/v1/bookmarks` (Bearer token required)
- `POST /api/v1/bookmarks` (Bearer token required)
- `DELETE /api/v1/bookmarks/:id` (Bearer token required)
- `GET /api/v1/dhikr/counters` (Bearer token required)
- `POST /api/v1/dhikr/counters` (Bearer token required)
- `GET /api/v1/quiz/attempts` (Bearer token required)
- `POST /api/v1/quiz/attempts` (Bearer token required)
- `GET /api/v1/quiz/stats` (Bearer token required)
- `GET /api/v1/audio-proxy`
- `GET /api/v1/prayer-times`

`/auth/refresh` and `/auth/logout` read refresh token from `HttpOnly` cookie.

## Swagger Documentation

Swagger UI tersedia setelah server jalan:

- `GET /swagger/index.html`
- Raw OpenAPI JSON: `GET /swagger/doc.json`

Contoh URL lokal (default port `8080`):

- `http://localhost:8080/swagger/index.html`

Generate ulang file Swagger jika endpoint berubah:

```bash
cd server
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g main.go -d ./cmd,./internal,./configs,./pkg -o ./docs --parseDependency --parseInternal
```

## Run Locally

```bash
cd server
cp .env.example .env
go mod tidy
go run ./cmd/migrate -action up
go run ./cmd
```

## Request Examples

### Login with Google ID Token

```http
POST /api/v1/auth/google
Content-Type: application/json

{
  "idToken": "<google-id-token-from-frontend>"
}
```

### Refresh Token

```http
POST /api/v1/auth/refresh
Cookie: bn_refresh_token=<refresh-token>
```

### Get Current User

```http
GET /api/v1/auth/me
Authorization: Bearer <access-token>
```

## Best Practices Applied

- Layered architecture (handler -> service -> repository)
- Context timeout at handler layer
- JWT short-lived access token + rotating refresh token
- Refresh token stored in DB as SHA-256 hash (never store raw token)
- Refresh token delivered via `HttpOnly` cookie
- Transactional refresh-token rotation
- Versioned SQL migration files (no runtime automigrate)
- Auth endpoint rate limiting by client IP
- Fail-fast config validation
- Graceful shutdown for HTTP server
- Centralized CORS and auth middleware
