# Product API — Go + Fiber

---

## Highlights

- **CQRS Repository Pattern** — `QueryRepository` (readDB) + `CommandRepository` (writeDB) — type-enforced DB routing, no replication lag after writes
- **Read/Write DB Separation** — Primary รับ write, Replica รับ read ผ่าน PostgreSQL Streaming Replication
- **Redis Cache-Aside** — cache product/supplier list+detail ตัดรอบ DB queries สำหรับ read-heavy traffic
- **Clean Architecture (DDD)** — แบ่ง 4 layer Domain → Application → Infrastructure → Interface
- **JWT Auth** — Register/Login with bcrypt, protected routes
- **Soft Delete** — Records hidden, not removed; hard delete available separately
- **User Tracking** — `created_by` / `updated_by` stamped on every record
- **Saga Transactions** — Compensating transactions for cross-domain operations
- **1M Seed Data** — suppliers 1,000 records + products 1,000,000 records

---

## Tech Stack

| | |
|---|---|
| Language | Go 1.25 |
| Framework | Fiber v2 |
| Database | PostgreSQL 16 (Primary + Replica) |
| Cache | Redis 7 |
| ORM | GORM |
| Logger | zerolog |
| Validation | go-playground/validator |
| Auth | golang-jwt/jwt/v5 + bcrypt |
| Docs | swaggo/swag + gofiber/swagger |

---

## Architecture

```
internal/{feature}/
  domain/              # Core — entity + interfaces (ไม่ depend ใคร)
    repository.go      # QueryRepository + CommandRepository interfaces
  application/
    queries/           # Read use cases — ใช้ QueryRepository
    commands/          # Write use cases — ใช้ CommandRepository
  infrastructure/
    repository/
      {name}_query_repository.go    # readDB implementation
      {name}_command_repository.go  # writeDB implementation
      {name}_mapper.go              # package-level mappers + fetch helpers
  interface/http/
    handler/           # Function-based Fiber handlers
    dto/               # Request/Response structs + validation
    router.go          # Local DI wiring (wires queryRepo + cmdRepo separately)
```

---

## Connection Info

| Service | Host | Port | User | Password | DB |
|---|---|---|---|---|---|
| PostgreSQL Primary (Write) | localhost | 9191 | postgres | postgres | product_db |
| PostgreSQL Replica (Read) | localhost | 9192 | postgres | postgres | product_db |
| Redis | localhost | 6380 | — | — | 0 |
| App | localhost | 9392 | — | — | — |


---

## Project Structure

```
golang-fiber/
├── cmd/
│   ├── app/main.go                   # Entry point, graceful shutdown
│   └── server/server.go              # Fiber setup, middleware, DI wiring
├── config/
│   └── config.go                     # Config struct, loaded from .env
├── internal/
│   ├── product/
│   │   ├── domain/                   # Interfaces + entities (no external deps)
│   │   │   └── repository.go         # ProductQueryRepository + ProductCommandRepository
│   │   ├── application/
│   │   │   ├── queries/              # Read use cases (QueryRepository)
│   │   │   └── commands/             # Write use cases (CommandRepository)
│   │   ├── infrastructure/repository/
│   │   │   ├── models/               # GORM models
│   │   │   ├── product_query_repository.go    # readDB
│   │   │   ├── product_command_repository.go  # writeDB
│   │   │   └── product_mapper.go              # package-level mappers + fetch helpers
│   │   └── interface/http/
│   │       ├── dto/                  # Request/response + validation
│   │       ├── handler/              # Function-based handlers
│   │       └── router.go             # DI wiring + routes
│   ├── supplier/                     # Same structure as product
│   └── user/
│       ├── domain/
│       │   └── repository.go         # UserQueryRepository + UserCommandRepository
│       ├── application/commands/
│       │   └── user_command.go       # Holds both queryRepo (login) + cmdRepo (register)
│       └── ...
├── pkg/
│   ├── auth/                         # JWT + context helpers
│   ├── cache/                        # Redis cache-aside
│   ├── common/                       # Logger, response, saga transaction, validator
│   ├── constants/                    # Typed enums (status, sort, messages)
│   ├── database/                     # DB connection helper
│   └── middleware/                   # JWT middleware
├── migrations/
│   ├── schema/                       # SQL — auto-run on first docker up
│   └── seed/                         # SQL — run manually
├── docs/                             # Swagger generated files
└── docker-compose.yaml
```
