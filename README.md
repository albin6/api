# Enterprise Role-Based Auth System

A production-ready, highly scalable Role-Based Authentication System built in Go using Clean Architecture.

## Features

- **Clean Architecture**: `Core (Domain/Ports)` -> `Adapters (Handler/Repo)`.
- **Manual Dependency Injection**: Explicit wiring via `bootstrap.go`.
- **RBAC**: Roles `ADMIN` (System), `HEAD`, `LEAD`, `MEMBER`.
- **JWT Authentication**: Access (15m) + Refresh (7d) with Rotation (Redis).
- **Admin Security**: Segregated `admins` table, protected by `x-admin-auth-key`.
- **Database Pooling**: Configured GORM connection limits.
- **Graceful Shutdown**: Handles OS signals to close DB/Redis connections.
- **Structured Logging**: JSON logs with `slog`.

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin
- **Database**: PostgreSQL (GORM)
- **Cache**: Redis
- **Config**: Viper
- **Container**: Docker & Docker Compose

## Architecture

```text
cmd/api/
  bootstrap.go    # DI & Server Composition
  main.go         # Entry Point
config/           # Environment Config
internal/
  core/
    domain/       # Entities (User, Admin, Role)
    port/         # Interfaces (Repository, Service)
    service/      # Use Cases (Auth, Admin)
  adapters/
    handler/      # HTTP Handlers
    repo/         # Postgres Repositories
    storage/      # Redis Repository
    api/          # Middlewares
pkg/              # Shared Utilities (Logger, Password, Token)
```

## API Documentation

### Public

- `GET /health`: Health check (Postgres & Redis status).
- `POST /auth/signup`: Create a new user (Default Role: MEMBER).
- `POST /auth/login`: Login user. Returns Access & Refresh tokens.
- `POST /auth/refresh`: Rotate refresh token.
- `POST /auth/logout`: Logout user (Revoke refresh token).

### Admin (Protected by `x-admin-auth-key`)

- `POST /admin/create`: Create a new System Admin.

### Private (Protected by `Authorization: Bearer <token>`)

- `GET /api/profile`: Get user profile (Example protected route).

## How to Run

1. **Start Infrastructure**:
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Run Application**:
   ```bash
   go run ./cmd/api
   ```
   Or with Docker:
   ```bash
   docker-compose up --build
   ```

## Configuration

Set the following environment variables (defaults in `config/config.go`):
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- `REDIS_HOST`, `REDIS_PORT`
- `JWT_SECRET`, `ADMIN_SECRET_KEY`
