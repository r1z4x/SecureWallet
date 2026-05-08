# SecureWallet Repository Structure Audit

**Date:** 2026-05-06  
**Status:** IMPLEMENTED_NOT_VERIFIED  
**Scope:** Full repository structure inspection — Go/Gin, Vue.js, MySQL, Redis, Docker, Swagger/OpenAPI

---

## 1. Project Identity & Stack

| Aspect | Expected | Actual | Status |
|---|---|---|---|
| Language | Go 1.21 | Go 1.21 (`go.mod` line 3) | MATCH |
| HTTP Framework | Gin | `github.com/gin-gonic/gin v1.10.1` | MATCH |
| Frontend | Vue.js | Vue 3.4 + Vite 5 (`frontend/package.json`) | MATCH |
| Database | MySQL 8 | MySQL 8.0 (`docker-compose.yml` + `gorm.io/driver/mysql v1.6.0`) | MATCH |
| Cache | Redis | Redis 7.2-alpine (`docker-compose.yml` + `github.com/go-redis/redis/v8`) | MATCH |
| Containerization | Docker | `docker-compose.yml` + `docker-compose.dev.yml` + `build/docker/Dockerfile.go.dev` | MATCH |
| API Docs | Swagger/OpenAPI | Swagger 2.0 (`docs/swagger.yaml` + `docs/swagger.json`) | MATCH (incomplete) |
| ORM | GORM | `gorm.io/gorm v1.30.2` | MATCH |
| JWT | golang-jwt | `github.com/golang-jwt/jwt/v5 v5.3.0` | MATCH |
| TOTP | pquerna/otp | `github.com/pquerna/otp v1.5.0` | MATCH |
| Hot Reload | air | `.air.toml` + `github.com/cosmtrek/air@v1.49.0` | MATCH |

---

## 2. Repository Directory Map

```
SecureWallet/
├── main.go                         # Entry point — Gin router, middleware, routes, cron jobs
├── go.mod / go.sum                 # Go module: securewallet, Go 1.21
├── .air.toml                       # Hot reload configuration (air)
├── .gitignore                      # Comprehensive — Go, Node, env files, logs, backups
├── README.md                       # Project README with quick start guide
├── env.example                     # Environment variable template (51 lines)
│
├── build/docker/
│   └── Dockerfile.go.dev           # Dev Docker: Go 1.21-alpine + air hot reload
│
├── docker-compose.yml              # Production: mysql + mongodb + redis + backend + frontend
├── docker-compose.dev.yml          # Dev: same + cron service, port mapping, volume mounts
│
├── db/
│   └── init.sql                    # Full schema: 13 tables, UUID PKs, FKs, indexes (232 lines)
│
├── docs/
│   ├── swagger.yaml                # Swagger 2.0 spec — Auth routes only (300 lines, incomplete)
│   ├── swagger.json                # Same in JSON
│   ├── docs.go                     # Go swagger doc generation marker
│   └── TROUBLESHOOTING.md          # Troubleshooting guide
│
├── internal/                       # Go backend (private/internal package)
│   ├── config/
│   │   ├── database.go             # GORM init, schema compat check, AutoMigrate, pool config
│   │   └── redis.go                # go-redis v8 init, ping test
│   │
│   ├── middleware/
│   │   ├── auth.go                 # AuthMiddleware (JWT), AdminMiddleware, SecurityHeaders, OptionalAuth
│   │   ├── rate_limit.go           # In-memory IP-based (5 req/15min, skipped in dev/debug)
│   │   ├── idempotency.go          # Idempotency-Key header → DB record lookup/caching
│   │   ├── validation.go           # Input validation middleware
│   │   └── admin.go                # Admin-specific middleware
│   │
│   ├── models/
│   │   ├── user.go                 # User — UUID PK, bcrypt hash, 2FA fields, GORM soft delete
│   │   ├── wallet.go               # Wallet — UUID PK, UserID FK, decimal(15,2) balance
│   │   ├── transaction.go          # Transaction — UUID PK, WalletID FK, type (deposit/withdrawal/transfer)
│   │   ├── session.go              # Session — UUID PK, UserID FK, token, expires_at
│   │   ├── audit_log.go            # AuditLog — UUID PK, UserID FK, action, resource, details, IP, UA
│   │   ├── idempotency_record.go   # IdempotencyRecord — key (unique), operation, payload_hash, response_body cache
│   │   ├── login_history.go        # LoginHistory — UUID PK, UserID FK, IP, UA, status, location
│   │   ├── support_ticket.go       # SupportTicket — UUID PK, UserID FK, subject, description, status, priority
│   │   └── blog.go                 # BlogPost, BlogComment, BlogCategory, BlogTag models
│   │
│   ├── routes/
│   │   ├── auth.go                 # POST /auth/{register,login,login/2fa,logout,refresh,password-reset,password-verify} + GET /auth/me
│   │   ├── wallets.go              # GET/POST/PUT/DEL /wallets, GET /wallets/balance, POST /wallets/{deposit,transfer}
│   │   ├── transactions.go         # GET/POST/PUT/DEL /transactions
│   │   ├── admin.go                # GET /admin/{dashboard,users,transactions,settings,support/tickets} + POST disable/enable/reply/resolve
│   │   ├── two_factor.go           # POST /2fa/{enable,disable,verify} + GET /2fa/status
│   │   ├── users.go                # User management routes
│   │   ├── support.go              # Support ticket routes
│   │   ├── blog.go                 # Blog routes (public)
│   │   ├── data_management.go      # Data init/reset routes
│   │   ├── login_history.go        # Login history routes
│   │   ├── backup.go               # Backup routes
│   │   ├── security.go             # Security monitoring routes
│   │   └── cron.go                 # Cron job management routes
│   │
│   └── services/
│       ├── auth.go                 # JWT token creation/validation, bcrypt hashing, user auth
│       ├── two_factor.go           # TOTP secret generation/validation, QR code URL
│       ├── backup.go               # Backup service
│       ├── comment.go              # Blog comment auto-approval service
│       ├── cron.go                 # Cron job execution service
│       ├── data_manager.go         # Sample data initialization and DB reset
│       ├── login_history.go        # Login attempt recording
│       └── security_detector.go    # Security event detection (IDOR, brute force, etc.)
│
├── frontend/                       # Vue 3 + Vite application
│   ├── package.json                # Vue 3.4, Vue Router 4.2, Pinia 2.1, Axios 1.6, TailwindCSS 3.4
│   ├── vite.config.js              # Vite 5 configuration
│   ├── tailwind.config.js          # TailwindCSS configuration
│   ├── postcss.config.js           # PostCSS configuration
│   ├── index.html                  # HTML entry point
│   ├── Dockerfile                  # Production Dockerfile
│   ├── Dockerfile.dev              # Development Dockerfile
│   └── src/
│       ├── main.js                 # Vue app bootstrap
│       ├── App.vue                 # Root component
│       ├── style.css               # Global styles
│       ├── components/             # Reusable components (7 files)
│       │   ├── AppHeader.vue
│       │   ├── DeleteAccount.vue
│       │   ├── LanguageSelector.vue
│       │   ├── LoginHistory.vue
│       │   ├── Navigation.vue
│       │   ├── Pagination.vue
│       │   └── TwoFactorAuth.vue
│       ├── views/                  # Page-level components (19 files)
│       │   ├── Admin.vue, Blog.vue, BlogPost.vue, Dashboard.vue, FAQ.vue,
│       │   ├── HelpCenter.vue, Landing.vue, Login.vue, PasswordReset.vue,
│       │   ├── Profile.vue, Register.vue, ResetDatabase.vue, SecurityTips.vue,
│       │   ├── Support.vue, TermsOfService.vue, Transactions.vue, Transfer.vue,
│       │   ├── UserGuide.vue, Wallet.vue
│       ├── services/               # API client modules (10 files)
│       │   ├── auth.js, wallet.js, transaction.js, twoFactor.js, admin.js,
│       │   ├── user.js, support.js, blog.js, loginHistory.js, dataManagement.js
│       ├── stores/
│       │   └── auth.js             # Pinia auth store
│       ├── i18n/
│       │   ├── index.js            # i18n setup
│       │   └── locales/            # en.js, es.js, tr.js
│       └── assets/
│           └── logo.svg
│
├── fix-database.sh                 # Interactive schema fix script
├── start-go-dev.sh                 # Go dev mode start script
├── start-go.sh                     # Go production start script
└── frontend/start-dev.sh           # Frontend dev start script
```

---

## 3. Database Schema (from `db/init.sql`)

| Table | Primary Key | Key Foreign Keys | Notable Columns |
|---|---|---|---|
| `users` | CHAR(36) | — | username, email, password_hash, two_factor_secret, two_factor_enabled, is_active, is_admin, deleted_at |
| `wallets` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | balance DECIMAL(15,2), currency VARCHAR(3) |
| `transactions` | CHAR(36) | wallet_id → wallets(id) ON DELETE CASCADE | type VARCHAR(20), amount DECIMAL(15,2), currency, description, status |
| `sessions` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | token VARCHAR(255), expires_at TIMESTAMP |
| `audit_logs` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | action VARCHAR(100), resource, details TEXT, ip_address, user_agent |
| `idempotency_records` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | key CHAR(36) UNIQUE, operation, payload_hash, status, http_status, response_body, expires_at |
| `support_tickets` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | subject, description, status, priority |
| `login_history` | CHAR(36) | user_id → users(id) ON DELETE CASCADE | ip_address, user_agent, status, location |
| `blog_categories` | CHAR(36) | — | name, slug (unique), description, color |
| `blog_tags` | CHAR(36) | — | name, slug (unique) |
| `blog_posts` | CHAR(36) | author_id → users(id) ON DELETE CASCADE | title, slug (unique), excerpt, content LONGTEXT, category, tags (JSON), status, view_count, published_at |
| `blog_comments` | CHAR(36) | post_id → blog_posts(id) ON DELETE CASCADE | name, email, content, status, ip_address, user_agent |
| `security_alerts` | VARCHAR(100) | user_id → users(id) ON DELETE CASCADE | type, severity, ip_address, details JSON, status, resolved_by, resolved_at |

**Key observations:**
- All tables use CHAR(36) UUID primary keys
- Foreign keys are properly declared with ON DELETE CASCADE
- Every table has appropriate indexes
- GORM soft delete (deleted_at) on core tables
- `security_alerts` uses VARCHAR(100) PK instead of CHAR(36) — inconsistency noted
- No migration versioning (not using golang-migrate or similar) — GORM AutoMigrate at startup

---

## 4. API Route Map

### Route Registration Order (from `main.go`)

| Route Group | Prefix | Middleware | Key Endpoints |
|---|---|---|---|
| Swagger | `/swagger/*any` | None | GET (docs UI) |
| Auth | `/api/auth` | RateLimitMiddleware (register, login, 2fa, pwd-reset, pwd-verify) | POST register, login, login/2fa, logout, refresh, password-reset, password-verify; GET me |
| Users | `/api/users` | AuthMiddleware | User CRUD |
| Wallets | `/api/wallets` | AuthMiddleware + IdempotencyMiddleware (transfer only) | GET /, /balance, /:id; POST /, /deposit, /transfer; PUT /:id; DELETE /:id |
| Transactions | `/api/transactions` | AuthMiddleware | GET /, /:id; POST /; PUT /:id; DELETE /:id |
| Admin | `/api/admin` | AuthMiddleware + AdminMiddleware | GET dashboard, users, transactions, settings, support/tickets; POST users/:id/disable, users/:id/enable, support/tickets/:id/reply, support/tickets/:id/resolve |
| Support | `/api/support` | AuthMiddleware | Support ticket CRUD |
| Data Management | `/api/data` | AuthMiddleware | Data init/reset |
| 2FA | `/api/2fa` | AuthMiddleware | POST enable, disable, verify; GET status |
| Login History | `/api/login-history` | AuthMiddleware | Login history access |
| Backup | `/api/backup` | AuthMiddleware | Backup operations |
| Security | `/api/security` | AuthMiddleware | Security monitoring |
| Blog | `/blog` (top-level) | OptionalAuthMiddleware | Public blog access |
| Cron | `/cron` (top-level) | None | Cron job triggers |

---

## 5. Middleware Chain (in `main.go`)

```
Global middleware (applied to all routes):
  1. gin.Recovery()                    — panic recovery
  2. SecurityHeadersMiddleware()       — X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy, CSP, HSTS
  3. InputValidationMiddleware()       — request validation
  4. Custom CORS middleware (inline)   — origin validation from CORS_ORIGINS env var

Route-level middleware:
  - AuthMiddleware()                   — JWT Bearer token validation, sets user in context
  - AdminMiddleware()                  — checks user.IsAdmin from context
  - RateLimitMiddleware()              — IP-based, 5 req/15min, in-memory, skipped in dev/debug mode
  - IdempotencyMiddleware("transfer")  — Idempotency-Key header validation, UUID check, DB record lookup
  - OptionalAuthMiddleware()           — attempts auth but doesn't require it
```

---

## 6. Service Layer Map

| Service | File | Responsibilities |
|---|---|---|
| AuthService (functions) | `services/auth.go` | GetJWTSecret, AuthenticateUser, CreateAccessToken, GetCurrentUser, GetPasswordHash |
| TwoFactorService | `services/two_factor.go` | GenerateSecret, ValidateCode, GenerateQRCodeURL, GetCurrentCode |
| LoginHistoryService | `services/login_history.go` | RecordLoginAttempt |
| SecurityDetector | `services/security_detector.go` | DetectIDOR, DetectBruteForce, etc. |
| CronService | `services/cron.go` | ExecuteCronJob (comment-approval, backup, log-cleanup, security-monitor) |
| CommentService | `services/comment.go` | Auto-approval scheduler for blog comments |
| DataManager | `services/data_manager.go` | ResetDatabase, CompleteDatabaseRecreation, sample data init |

---

## 7. Go Model Relationships

```
User (1) ──→ (*) Wallet       [foreignKey:UserID]
User (1) ──→ (*) Session      [foreignKey:UserID]
User (1) ──→ (*) AuditLog     [foreignKey:UserID]
User (1) ──→ (*) SupportTicket [foreignKey:UserID]
User (1) ──→ (*) BlogPost     [foreignKey:AuthorID]

Wallet (1) ──→ (*) Transaction [foreignKey:WalletID]
```

All models use `BeforeCreate` hook to auto-generate UUID via `uuid.New()`.

---

## 8. Docker Infrastructure

### Services (Production — `docker-compose.yml`)

| Service | Image | Port Mapping | Health Check |
|---|---|---|---|
| mysql | mysql:8.0 | (none, secure default) | mysqladmin ping, 10s interval |
| mongodb | mongo:7.0 | (none, secure default) | None |
| redis | redis:7.2-alpine | (none, secure default) | None |
| backend | Go app (build/docker/Dockerfile.go) | 8081:8080 | wget /health, 30s |
| frontend | Node.js (frontend/Dockerfile) | 3001:3000 | None |

### Services (Development — `docker-compose.dev.yml`)

Same as production plus:
- Ports exposed for direct access (mysql 3307, mongodb 27018, redis 6380)
- Volume mounts for hot reload (`. :/app` for backend, `./frontend/src` for frontend)
- Cron service running dcron with 4 scheduled jobs
- Backend uses `air` for hot reload
- Frontend uses `npm run dev`

---

## 9. Swagger/OpenAPI Status

**Spec version:** Swagger 2.0  
**Location:** `docs/swagger.yaml`, `docs/swagger.json`  
**Swagger UI route:** `GET /swagger/*any` (via gin-swagger)  
**Swagger annotations:** Present on Go route handlers (e.g., `@Summary`, `@Tags`, `@Router`, `@Security`)

### Documented Endpoints (5 total — incomplete):
- `POST /auth/register` — Register user
- `POST /auth/login` — Login user
- `POST /auth/logout` — Logout
- `GET /auth/me` — Get current user (secured)
- `POST /auth/refresh` — Refresh token

### Missing from Swagger:
- Wallet routes (CRUD, balance, deposit, transfer)
- Transaction routes
- Admin routes
- 2FA routes (enable, disable, verify, status)
- Support routes
- Login history routes
- Data management routes
- Backup routes
- Security routes
- Blog routes

**Swagger type inconsistency:** Definitions use `type: integer` for IDs, but actual models use UUIDs (`type: char(36)`). The Swagger spec is out of date with the actual UUID-based data model.

---

## 10. Frontend Architecture

**Framework:** Vue 3.4 (Composition API) + Vite 5  
**State Management:** Pinia 2.1 (`stores/auth.js`)  
**Routing:** Vue Router 4.2  
**HTTP Client:** Axios 1.6  
**Styling:** TailwindCSS 3.4  
**Internationalization:** vue-i18n 9.x (en, es, tr)

### Frontend Service Modules → Backend Route Mapping

| Frontend Service | Expected Backend Routes |
|---|---|
| `auth.js` | `/api/auth/*` |
| `wallet.js` | `/api/wallets/*` |
| `transaction.js` | `/api/transactions/*` |
| `twoFactor.js` | `/api/2fa/*` |
| `admin.js` | `/api/admin/*` |
| `user.js` | `/api/users/*` |
| `support.js` | `/api/support/*` |
| `blog.js` | `/blog/*` |
| `loginHistory.js` | `/api/login-history/*` |
| `dataManagement.js` | `/api/data/*` |

### Views → Purpose Mapping

| View | Purpose |
|---|---|
| Landing.vue | Public landing/marketing page |
| Login.vue | User login form |
| Register.vue | User registration form |
| Dashboard.vue | Main user dashboard |
| Wallet.vue | Wallet overview and management |
| Transfer.vue | Money transfer between users |
| Transactions.vue | Transaction history |
| Profile.vue | User profile management |
| TwoFactorAuth.vue | 2FA setup and management |
| Admin.vue | Full admin panel |
| Support.vue | Support ticket system |
| Blog.vue / BlogPost.vue | Public blog |
| FAQ.vue, HelpCenter.vue, UserGuide.vue | Help/documentation pages |
| SecurityTips.vue, TermsOfService.vue | Static content pages |
| PasswordReset.vue | Password reset flow |
| ResetDatabase.vue | Database reset admin tool |

---

## 11. Test Coverage

**Test files found:** **0** (`**/*_test.go` glob returned nothing)

No Go tests exist in the repository. No frontend test configuration found (no vitest, jest, or testing-library dependencies).

---

## 12. Key Architectural Observations

### Strengths
1. **UUID-based data model** — All entities use UUID PKs with proper FK relationships and CASCADE deletes
2. **Security headers middleware** — CSP, HSTS, X-Frame-Options, X-Content-Type-Options applied globally
3. **bcrypt password hashing** — Used consistently across auth routes and service layer
4. **Structured Go packages** — Clean separation: config, models, middleware, routes, services
5. **Environment variable configuration** — No hardcoded secrets in source (env.example provided)
6. **GORM soft deletes** — Standard soft-delete pattern on all core entities
7. **Transfer transaction wrapping** — DB transactions used for balance changes and transaction record creation
8. **Dual docker-compose files** — Separate dev and production configurations
9. **Hot reload in dev** — `air` for Go, Vite HMR for Vue

### Concerns & Gaps

| # | Finding | Severity | Context |
|---|---|---|---|
| 1 | **No tests exist** | CRITICAL | Zero Go tests, zero frontend tests. Cannot verify security, atomicity, or correctness. |
| 2 | **IDOR in getWallet** | CRITICAL | `routes/wallets.go:348` — intentionally vulnerable: no ownership check, returns any wallet by ID |
| 3 | **No SELECT FOR UPDATE on transfers** | HIGH | `routes/wallets.go:217` — tx.Begin() but no row-level locking, race condition possible between balance read and write |
| 4 | **Idempotency record not created after transfer** | HIGH | Middleware checks for existing record but `CreateIdempotencyRecord` is never called in the transfer handler |
| 5 | **Rate limiting is in-memory only** | HIGH | `middleware/rate_limit.go` — loses all state on restart, not shared across instances |
| 6 | **No refresh token rotation** | HIGH | `routes/auth.go:321` — refreshToken returns stub: `"Token refreshed"` with no implementation |
| 7 | **Logout is a no-op** | HIGH | `routes/auth.go:291` — logout just returns 200, no token invalidation |
| 8 | **Swagger spec is incomplete** | HIGH | Only 5 auth endpoints documented, missing all wallet/transaction/admin/2FA routes. IDs typed as `integer` instead of UUID strings. |
| 9 | **JWT claims lack user ID** | MEDIUM | Token carries `sub: username` only — no `user_id` claim, causing extra DB lookup on every auth check |
| 10 | **No JWT refresh token** | MEDIUM | Only access tokens are issued; no refresh token mechanism for silent renewal |
| 11 | **Password policy is inconsistent** | MEDIUM | `routes/auth.go:474` enforces 12-char password, but admin settings UI shows minLength: 8 |
| 12 | **MongoDB included but unused** | LOW | `docker-compose.yml` includes MongoDB 7.0, but no Go code references MongoDB |
| 13 | **Login uses `user.Username` in `sub` claim** | MEDIUM | 2FA login hardcodes 60-min expiry (line 263) ignoring `ACCESS_TOKEN_EXPIRE_MINUTES` env var |
| 14 | **CORS origin parsing is fragile** | LOW | `main.go:111-113` — manual JSON array parsing could break with unusual formatting |
| 15 | **No structured logging** | MEDIUM | Only `log.Println`/`log.Printf` — no log levels, no JSON format capability |
| 16 | **No health check endpoint for readiness** | LOW | Only `/health` returns static JSON, no DB/Redis connectivity check |
| 17 | **No input sanitization for XSS** | MEDIUM | No HTML/script escaping in content that gets rendered back (blog posts, comments, descriptions) |
| 18 | **API returns full user objects** | MEDIUM | `getCurrentUser` returns the full User model including all relationships |

---

## 13. Baseline Commands

Based on the repository structure, the following commands should work:

```bash
# Go backend
cd /workspace/.../SecureWallet
go build -o ./tmp/main .           # Build
go vet ./...                       # Lint
go test ./...                      # Test (no tests exist — expected to pass with no test files)

# Frontend
cd frontend
npm install                        # Install dependencies
npm run build                      # Production build
npm run dev                        # Development server

# Docker
docker-compose -f docker-compose.dev.yml up    # Dev
docker-compose up                               # Production
```

---

## 14. Migration Status

- **Versioning:** None — no migration framework detected
- **Approach:** GORM AutoMigrate at startup (`config/database.go:130`)
- **Rollback:** Not supported — `FORCE_DATABASE_RECREATION` env var drops and recreates entire database
- **Schema file:** Single `db/init.sql` with full DDL (not versioned migrations)
- **Vulnerability:** In production, AutoMigrate is skipped if schema incompatibility is detected, requiring manual intervention

---

*Generated: 2026-05-06 | Repository: r1z4x/SecureWallet | Task: bd0841f9-c59e-4d25-a9e6-686a0146bd6a (subtask 1: architect)*
