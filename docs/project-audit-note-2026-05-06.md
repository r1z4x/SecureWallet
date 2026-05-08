# SecureWallet — Project Audit Note

**Date:** 2026-05-06  
**Status:** `IMPLEMENTED_NOT_VERIFIED` — all findings are evidence-backed from direct file reads; runtime verification not yet performed.  
**Audit scope:** full repository structure, schema baseline, module layout, runnable commands, and known blockers.

---

## 1. Schema Version

| Aspect | Finding |
|---|---|
| Migration system | **None.** No `db/migrations/` directory, no `schema_migrations` table, no version numbering. |
| Schema definition | Single monolithic `db/init.sql` (232 lines, 13 `CREATE TABLE IF NOT EXISTS` statements). |
| Schema version | Effectively **V001**. No formal versioning exists. |
| Schema application | Docker entrypoint runs `init.sql`; GORM `AutoMigrate` runs at `main.go:55` (after a UUID compatibility check in `internal/config/database.go`). |
| Tables (13 total) | `users`, `wallets`, `transactions`, `sessions`, `audit_logs`, `idempotency_records`, `support_tickets`, `login_history`, `security_alerts`, `blog_categories`, `blog_tags`, `blog_posts`, `blog_comments` |
| Primary keys | All `CHAR(36)` UUID, except `security_alerts.id` which is `VARCHAR(100)` (inconsistency). |
| Foreign keys | Present on user-owned tables (`wallets`, `transactions`, `sessions`, `audit_logs`, `idempotency_records`, `support_tickets`, `login_history`); all `ON DELETE CASCADE`. |
| AutoMigrate gaps | `IdempotencyRecord` and `SecurityAlert` models are **excluded** from GORM AutoMigrate. `User` struct has 4 columns (`name`, `title`, `avatar`, `bio`) not present in `init.sql`. |
| Financial precision risk | Go models use `float64` for `balance` and `amount`; MySQL stores `DECIMAL(15,2)`. |
| Rollback path | **None.** No down-migration scripts, no backup strategy beyond the optional cron backup job. |

---

## 2. Module Layout

```
securewallet/                           (Go 1.21, module: securewallet)
├── main.go                             Entry point — Gin router, cron dispatcher, middleware chain
├── internal/
│   ├── config/
│   │   ├── database.go                 GORM init, schema compat check, AutoMigrate (11 models), pool
│   │   └── redis.go                    go-redis v8 init, ping
│   ├── middleware/
│   │   ├── auth.go                     AuthMiddleware (JWT), AdminMiddleware, SecurityHeaders, OptionalAuth
│   │   ├── rate_limit.go               In-memory IP rate limiter (5 req/15min, skipped in dev/debug)
│   │   ├── idempotency.go              Idempotency-Key header → DB record lookup; caches response body
│   │   ├── validation.go               Input validation middleware
│   │   └── admin.go                    Admin-specific middleware
│   ├── models/                         (9 files) User, Wallet, Transaction, Session, AuditLog,
│   │                                     IdempotencyRecord, LoginHistory, SupportTicket, Blog (4-in-1)
│   ├── routes/                         (11 files) auth, wallets, transactions, admin, users,
│   │                                     two_factor, support, login_history, backup, security, data_management
│   └── services/                       (7 files) auth, two_factor, login_history, security_detector,
│                                          data_manager, cron_service, comment_service
├── db/
│   └── init.sql                        Monolithic schema (13 tables, UUID PKs, FK constraints)
├── frontend/                           Vue 3.4, Vite 5, Pinia, Vue Router, Axios, TailwindCSS, i18n
├── docs/                               swagger.yaml (Swagger 2.0), swagger.json, TROUBLESHOOTING.md
├── build/docker/                       Dockerfile.go.dev (Go 1.21-alpine + air hot reload)
├── docker-compose.yml                  Production: mysql + mongodb + redis + backend + frontend
├── docker-compose.dev.yml              Dev: same + cron + exposed ports + hot-reload volumes
├── .air.toml                           Go hot reload config
├── env.example                         51 lines — JWT, DB, Redis, CORS, cron config
├── start-go.sh                         Bare-metal start: go mod tidy && go run main.go
└── fix-database.sh                     Interactive DB reset/recreation/check script
```

### Route Group Summary (60+ endpoints, 13 groups)

| Group | Endpoints | Real | Stubs | Key Issues |
|---|---|---|---|---|
| `/api/auth` | 8 | 5 | 3 (logout, refresh, me-minimal) | Refresh returns static "Token refreshed"; logout does not invalidate sessions |
| `/api/wallets` | 8 | 4 | 4 | `GET /:id` has **IDOR** (no ownership check); transfer lacks `SELECT FOR UPDATE` |
| `/api/transactions` | 5 | 1 | 4 | Only list endpoint works |
| `/api/admin` | 10 | 5 | 5 | Dashboard, user enable/disable, system settings are stubs |
| `/api/users` | 7 | 4 | 3 | Create/update use mock data |
| `/api/2fa` | 4 | 4 | 0 | Fully implemented |
| `/api/support` | 5 | 2 | 3 | Reply/resolve are partial |
| `/api/login-history` | 2 | 2 | 0 | Fully implemented |
| `/api/backup` | 4 | 4 | 0 | Fully implemented |
| `/api/security` | 5 | 5 | 0 | Fully implemented |
| `/api/data` | 5 | 5 | 0 | **No auth middleware** — anyone can reset the database |
| `/api/blog` | 7 | 7 | 0 | Public routes on root router |
| `/api/cron` | 6 | 6 | 0 | **No auth middleware** — anyone can trigger cron jobs |

---

## 3. Runnable Commands

### Backend (Go)
```bash
# Build — PASSES (zero compilation errors)
go build ./...

# Vet — PASSES (zero static analysis warnings)
go vet ./...

# Test — PASSES vacuously (0 test files across 7 packages)
go test ./...

# Run standalone (needs MySQL + Redis)
go run main.go                # or: ./start-go.sh

# Hot reload (needs air CLI)
air
```

### Frontend (Vue.js)
```bash
cd frontend/

# Install dependencies — PASSES (11 audit vulnerabilities)
npm install

# Build — PASSES (574 KB JS bundle, chunk >500 KB warning)
npm run build

# Dev server
npm run dev

# Lint — FAILS (eslint not in devDependencies)
npm run lint
```

### Docker
```bash
# Development (exposed ports, hot reload, cron service)
docker compose -f docker-compose.dev.yml up -d

# Production (ports locked down, health checks)
docker compose up -d

# Stop
docker compose down
```

### Maintenance
```bash
# Interactive database fix (reset / recreate / manual cleanup / schema check)
./fix-database.sh

# Force database recreation via env var
FORCE_DATABASE_RECREATION=true RESET_DATABASE_ON_STARTUP=true go run main.go

# Cron job execution
go run main.go --cron=backup
go run main.go --cron=security-monitor
go run main.go --cron=comment-approval
go run main.go --cron=log-cleanup
```

---

## 4. Known Blockers

### Critical (security / financial integrity)

| # | Blocker | Location | Detail |
|---|---|---|---|
| B1 | **Zero automated tests** | Entire codebase | 0 Go test files, 0 frontend test files, no test framework configured. No falsifiable verification of any behavior exists. |
| B2 | **IDOR in wallet endpoint** | `internal/routes/wallets.go:348` | `getWallet()` reads any wallet by ID without ownership check. Security detector fires a warning but does **not block** the response. |
| B3 | **Transfer race condition** | `internal/routes/wallets.go:217` | Transfer handler uses `db.Begin()` + `db.Commit()` but **omits `SELECT FOR UPDATE`** on wallet rows. Concurrent transfers can produce incorrect balances. |
| B4 | **Idempotency is broken** | `internal/routes/wallets.go` + `internal/middleware/idempotency.go` | The idempotency middleware checks for existing records, but the transfer handler **never calls `CreateIdempotencyRecord()`** after success. Duplicate requests always pass through. |
| B5 | **Unauthenticated destructive endpoints** | `internal/routes/data_management.go` + `internal/routes/cron.go` | `/api/data/*` and `/api/cron/*` have **zero auth middleware**. Anyone can reset the database or execute cron jobs. |

### High (infrastructure / contract)

| # | Blocker | Location | Detail |
|---|---|---|---|
| B6 | **No CI/CD pipeline** | `.github/workflows/` missing | No automated build, test, or deploy pipeline. |
| B7 | **Rate limiting is in-memory** | `internal/middleware/rate_limit.go` | State lost on restart; not shared across instances; **skipped entirely** in dev/debug mode. |
| B8 | **Float64 for money** | `internal/models/wallet.go`, `transaction.go` | Go `float64` cannot exactly represent decimal values. MySQL uses `DECIMAL(15,2)`. |
| B9 | **Auth stubs** | `internal/routes/auth.go` | Logout returns 200 without invalidating sessions; refresh returns static "Token refreshed" without rotating tokens. |
| B10 | **Swagger incomplete** | `docs/swagger.yaml` | Only 5 of 60+ endpoints documented. ID type is `integer` but all real IDs are UUIDs. |

### Medium (operational / observability)

| # | Blocker | Location | Detail |
|---|---|---|---|
| B11 | **No structured logging** | all packages | `log.Println` / `fmt.Printf` used everywhere; no log levels, no request correlation IDs, no JSON format in practice. |
| B12 | **Full user objects returned** | `routes/auth.go`, `routes/users.go` | API responses include password hash, 2FA secret, and other sensitive fields. |
| B13 | **MongoDB unused** | `docker-compose.yml` | MongoDB container is defined and started but no Go code references it. Wasteful in production. |
| B14 | **Frontend lint broken** | `frontend/package.json` | `eslint` referenced in `npm run lint` script but not listed in `devDependencies`. |
| B15 | **No migration rollback** | `db/` | No down-migration scripts. `FORCE_DATABASE_RECREATION=true` is the only reset path (destructive). |

---

## 5. Audit Trail

This note synthesizes findings from three prior evidence-backed subtask outputs:
- `docs/repository-audit-2026-05-06.md` (426 lines) — full directory inventory, stack verification, component wiring
- `docs/repository-health-check-baseline-2026-05-06.md` (50 lines) — build/test/lint command results
- `docs/route-handler-service-mapping-2026-05-06.md` (362 lines) — every route → handler → service → DB trace

All blockers above were confirmed by direct file reads during this audit session. No claim is based on memory or inference.
