# Product API — Go + Fiber
---
## Highlights

- **Read/Write DB Separation** — Primary รับ write, Replica รับ read ผ่าน PostgreSQL Streaming Replication
- **Redis Cache-Aside** — cache product list/detail ตัดรอบ DB queries สำหรับ read-heavy traffic
- **Clean Architecture (DDD)** — แบ่ง 4 layer Domain → Application → Infrastructure → Interface
- **1M Seed Data** — suppliers 1,000 records + products 1,000,000 records (มีและไม่มี supplier relation ปนกัน)

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

---

## Architecture

```
internal/{feature}/
  domain/              # Core — entity + interfaces (ไม่ depend ใคร)
  application/
    queries/           # Read use cases + cache-aside
    commands/          # Write use cases + cache invalidation
  infrastructure/
    repository/        # GORM implementation + entity↔model mapper
  interface/http/
    handler/           # Function-based Fiber handlers
    dto/               # Request/Response structs + validation
    router.go          # Local DI wiring
```

---


## Connection Info

| Service | Host | Port | User | Password | DB |
|---|---|---|---|---|---|
| PostgreSQL Primary | localhost | 9191 | postgres | postgres | product_db |
| PostgreSQL Replica | localhost | 9192 | postgres | postgres | product_db |
| Redis | localhost | 6380 | — | — | 0 |


---

## Project Structure

```
golang-fiber/
├── cmd/
│   ├── app/main.go          # Entry point, graceful shutdown
│   └── server/server.go     # Fiber setup, middleware, DI
├── config/
│   └── config.go            # Config struct, loaded from .env
├── internal/
│   └── product/             # Product feature
│       ├── domain/          # Interfaces + entities
│       ├── application/     # Business logic + cache
│       ├── infrastructure/  # GORM + DB implementation
│       └── interface/http/  # Handlers + DTOs + Router
├── pkg/
│   ├── cache/               # Redis helper (Get/Set/DelPattern)
│   ├── common/              # Logger, ResponseHelper, Validator
│   ├── constants/           # Response codes, Sort enums
│   └── database/            # DB connection helper
├── migrations/
│   ├── schema/              # Auto-run on first docker up
│   └── seed/                # Run manually
└── docker-compose.yaml
```
