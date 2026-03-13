# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Commands

```bash
# Build
go build ./...

# Test
go test -v -race ./internal/... ./pkg/...
go test -v -race -coverprofile=coverage.out ./...
go test -v -tags=integration ./tests/integration/...

# Single package test
go test -v -race ./pkg/mikrotik/hotspot/...

# Lint & tidy
make lint
go mod tidy

# Database migrations (requires goose installed)
make migrate-up
make migrate-down
make migrate-status
make migrate-reset

# Infrastructure (Docker)
make docker-up    # start postgres, redis, rabbitmq
make docker-down
```

### Scaffolding (Makefile generators)

These generators create boilerplate files and auto-update the corresponding `registry.go`:

```bash
make model VAL=customer          # internal/model/
make domain VAL=billing          # internal/domain/billing/
make service VAL=billing         # internal/service/
make handler VAL=billing         # internal/handler/
make repository VAL=customer     # internal/repository/ + postgres/
make migration VAL=voucher_sales # internal/migration/ (next number auto-assigned)
make queue-producer VAL=suspend  # internal/queue/producer/
make queue-consumer VAL=suspend  # internal/queue/consumer/
make scheduler VAL=billing       # internal/scheduler/
make mikrotik-domain VAL=ppp     # pkg/mikrotik/domain/
make mikrotik-module VAL=ppp     # pkg/mikrotik/<name>/ repo+service + mikrotik.go facade
```

After scaffolding, manually wire the new component into the relevant `registry.go`.

---

## Architecture

The project is split into two completely independent concerns:

```
internal/   ← ISP business logic (customers, billing, payments)
pkg/        ← RouterOS API client (hardware abstraction)
```

### `internal/` — ISP Business Logic

Every layer uses a **Registry pattern**: a `Registry` struct in each package wires all components together. Dependency flow:

```
cmd/server/main.go
  └─ postgres.NewRepository(db)        → repository.Registry
  └─ domain.NewRegistry(...)           → domain.Registry  (pure business rules, no DB)
  └─ queue.NewRegistry(rabbitClient)   → queue.Registry   (RabbitMQ producers/consumers)
  └─ service.NewRegistry(repo, domain, queue, mikrotik)
  └─ handler.NewRegistry(services)
  └─ router.New(handlers, middleware)
```

| Directory | Role |
|---|---|
| `internal/model/` | GORM DB structs |
| `internal/domain/<name>/` | Pure business logic/rules (no DB access) |
| `internal/repository/` | Interfaces + `postgres/` GORM implementations |
| `internal/service/` | Use cases — orchestrate repo + domain + queue |
| `internal/handler/` | Gin HTTP handlers |
| `internal/router/` | Route registration (`/api/v1/...`) |
| `internal/middleware/` | JWT auth, request logging, request ID |
| `internal/migration/` | Goose Go-code migrations, numbered `001_`, `002_`, ... |
| `internal/queue/` | RabbitMQ producers/consumers (billing, suspend, notification) |
| `internal/scheduler/` | Cron jobs via robfig/cron |

### `pkg/mikrotik/` — RouterOS API Client

Each RouterOS subsystem is its own sub-package with the same structure:

```
pkg/mikrotik/
  client/        ← TCP connection wrapper (async mode, auto-reconnect)
  domain/        ← RouterOS domain types (no DB)
  <subsystem>/
    repository.go  ← direct RouterOS API calls
    service.go     ← thin wrapper; exposed on the facade
  mikrotik.go    ← Client facade (PPP, Hotspot, Queue, Firewall, IPPool, IPAddress, ...)
```

Current subsystem packages: `hotspot`, `ppp`, `queue`, `firewall`, `ippool`, `ipaddress`, `monitor`, `report`, `script`, `voucher`.

**RouterOS client patterns:**
- `RunContext(ctx, "/path/to/command", "?filter=val", "=param=val")` — one-shot (print/add/set/remove/enable/disable)
- `ListenArgsContext(ctx, args)` — streaming (follow, interval); returns `*ListenReply` with `.Chan()`
- CLI path → API: replace spaces with `/` → `/ip/hotspot/user/print`
- Filter by field: `?name=foo`, filter by ID: `?.id=*3A`
- Set params: `=name=foo`, reference item: `=.id=*3A`

**Per-package helpers** (repeated in each mikrotik package, do not import cross-package):
```go
func parseInt(s string) int64 { ... }
func parseBool(s string) bool  { return s == "true" || s == "yes" }
```

### Migrations

Migrations are Go code (not SQL files), registered with goose via `init()`:

```go
func init() { goose.AddMigrationContext(up021, down021) }

func up021(ctx context.Context, tx *sql.Tx) error {
    _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS ...`)
    return err
}
```

Number files sequentially. The `users` table is `001`, `mikrotik_routers` is `004`. Always check the last migration number before adding a new one.

### Mikhmon-style Voucher Reports

`pkg/mikrotik/report/` stores sales reports as RouterOS `/system/script` entries with a special name format:

```
date-|-time-|-username-|-price-|-ip-|-mac-|-validity-|-profile-|-comment
```

**This is the mikhmon format and must not be changed.** New DB-backed reports use `pkg/mikrotik/domain.VoucherSaleRecord` + `internal/migration/021_create_voucher_sales.go` alongside it — the two paths are independent.

### Hotspot IP Binding (Static IP for permanent subscribers)

MikroTik IP bindings (`/ip/hotspot/ip-binding`) are managed via `pkg/mikrotik/hotspot.Service` (GetIPBindings, AddIPBinding, RemoveIPBinding, EnableIPBinding, DisableIPBinding). When a hotspot subscription with `static_ip` is synced to the router, store the returned MikroTik entry ID in `subscriptions.mt_ip_binding_id` (added in migration 022).

---

## Key Conventions

- **`internal/queue/`** = RabbitMQ message queue (ISP billing/suspend/notification). **`pkg/mikrotik/queue/`** = RouterOS `/queue/simple` management. These are completely different things.
- Mikrotik `Client` in `service.NewRegistry` is `nil` at startup — it is initialized per-router at request time.
- All response shapes use `pkg/response` helpers.
- Validation uses `go-playground/validator/v10`; struct tags are `validate:"..."`.
- Logger is `go.uber.org/zap`; use structured fields, not `fmt.Sprintf`.
