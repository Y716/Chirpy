# 🐦 Chirpy
**THIS IS A COURSE PROJECT FROM BOOT.DEV. MORE OF IT IN THEIR WEBSITE.**

A lightweight Twitter-like REST API built from scratch in Go. Chirpy lets users register, authenticate, and post short messages ("chirps") with full CRUD support, JWT-based auth, refresh tokens, and a webhook-driven premium upgrade system.

## Features

- **User management** — registration, login, and profile updates with Argon2id password hashing
- **Chirps** — create, read, list, and delete short posts (140-character limit) with automatic profanity filtering
- **JWT authentication** — short-lived access tokens (1 hour) with HS256 signing
- **Refresh tokens** — long-lived tokens (60 days) for seamless session renewal, with revocation support
- **Chirpy Red** — premium tier upgrades via Polka webhook integration with API key verification
- **Admin tools** — metrics dashboard and user reset endpoint (dev-only)
- **Database migrations** — incremental schema evolution using Goose
- **Type-safe SQL** — query generation with sqlc

## Tech Stack

| Layer        | Technology              |
|-------------|------------------------|
| Language     | Go                     |
| Database     | PostgreSQL             |
| SQL Tooling  | sqlc, Goose            |
| Auth         | JWT (golang-jwt/v5), Argon2id |
| Config       | godotenv               |
| IDs          | UUID (google/uuid)     |

## Project Structure

```
.
├── main.go                  # Server setup, routing, middleware
├── handler_chirps.go        # Chirp CRUD handlers
├── handler_users.go         # User registration, login, update, webhooks
├── handler_tokens.go        # Refresh & revoke token handlers
├── json.go                  # JSON response helpers
├── internal/
│   ├── auth/                # JWT, password hashing, refresh tokens, API keys
│   └── database/            # sqlc-generated database layer
├── sql/
│   ├── schema/              # Goose migration files
│   └── queries/             # sqlc query definitions
├── assets/                  # Static assets (logo)
├── index.html               # Served at /app
├── sqlc.yaml                # sqlc configuration
├── go.mod / go.sum
└── .gitignore
```

## API Endpoints

### Health

| Method | Endpoint          | Description        |
|--------|-------------------|--------------------|
| GET    | `/api/healthz`    | Health check       |

### Users

| Method | Endpoint                  | Auth       | Description                      |
|--------|---------------------------|------------|----------------------------------|
| POST   | `/api/users`              | —          | Register a new user              |
| POST   | `/api/login`              | —          | Login (returns JWT + refresh)    |
| PUT    | `/api/users`              | Bearer JWT | Update email and password        |

### Chirps

| Method | Endpoint                  | Auth       | Description                      |
|--------|---------------------------|------------|----------------------------------|
| POST   | `/api/chirps`             | Bearer JWT | Create a chirp (max 140 chars)   |
| GET    | `/api/chirps`             | —          | List all chirps (optional `author_id`, `sort`) |
| GET    | `/api/chirps/{chirpID}`   | —          | Get a single chirp               |
| DELETE | `/api/chirps/{chirpID}`   | Bearer JWT | Delete own chirp                 |

### Tokens

| Method | Endpoint                  | Auth           | Description                 |
|--------|---------------------------|----------------|-----------------------------|
| POST   | `/api/refresh`            | Bearer Refresh | Get a new access token      |
| POST   | `/api/revoke`             | Bearer Refresh | Revoke a refresh token      |

### Webhooks

| Method | Endpoint                  | Auth       | Description                      |
|--------|---------------------------|------------|----------------------------------|
| POST   | `/api/polka/webhooks`     | API Key    | Upgrade user to Chirpy Red       |

### Admin

| Method | Endpoint                  | Description                          |
|--------|---------------------------|--------------------------------------|
| GET    | `/admin/metrics`          | View file server hit count           |
| POST   | `/admin/reset`            | Delete all users (dev only)          |

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev/) (for query generation)

### Setup

1. **Clone the repo**
   ```bash
   git clone https://github.com/Y716/Chirpy.git
   cd Chirpy
   ```

2. **Create a `.env` file**
   ```env
   DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWTSECRET=your-secret-key
   POLKA_KEY=your-polka-api-key
   ```

3. **Run migrations**
   ```bash
   goose -dir sql/schema postgres "$DB_URL" up
   ```

4. **Generate database code** (only if you modify queries)
   ```bash
   sqlc generate
   ```

5. **Run the server**
   ```bash
   go build -o chirpy && ./chirpy
   ```
   The server starts on `http://localhost:8080`.

## Usage Examples

**Register a user:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

**Post a chirp:**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"body": "Hello, Chirpy!"}'
```

## License

This project is unlicensed — feel free to use it however you like.
