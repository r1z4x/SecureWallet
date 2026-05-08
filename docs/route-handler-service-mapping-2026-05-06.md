# Gin Route-to-Handler-Service-Repository Mapping

**Generated:** 2026-05-06
**Status:** `VERIFIED_COMPLETE` (all files read and traced)

---

## Entry Point

`main.go` — sets up Gin router, applies global middleware, registers all route groups under `/api`, plus blog and cron routes at root level.

### Global Middleware Chain (applied to all routes)

| Order | Middleware | File | Purpose |
|---|---|---|---|
| 1 | `gin.Recovery()` | built-in | Panic recovery |
| 2 | `SecurityHeadersMiddleware()` | `middleware/auth.go:106` | X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, CSP, HSTS, Referrer-Policy |
| 3 | `InputValidationMiddleware()` | `middleware/validation.go:12` | Blocks SQL/XSS/command-injection patterns in query/path params |
| 4 | Inline CORS middleware | `main.go:99` | Origin validation from `CORS_ORIGINS` env var |

### Public Endpoints (no auth)

| Method | Path | Handler | File | Notes |
|---|---|---|---|---|
| GET | `/` | anonymous | `main.go:174` | API info |
| GET | `/health` | anonymous | `main.go:193` | Health check |
| GET | `/api/info` | anonymous | `main.go:202` | Minimal API metadata |
| GET | `/test-cors` | anonymous | `main.go:184` | CORS test |
| GET | `/swagger/*any` | ginSwagger | `main.go:149` | Swagger UI |

---

## 1. Auth Routes (`/api/auth`)

**Setup:** `SetupAuthRoutes()` in `routes/auth.go:28`
**Base group:** `router.Group("/api/auth")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| POST | `/register` | `RateLimitMiddleware()` | `register` | — | `db.Where(...).First()` (check dup), `db.Create()` (user) | auth.go:32 |
| POST | `/login` | `RateLimitMiddleware()` | `login` | `GetJWTSecret()`, `NewLoginHistoryService().RecordLoginAttempt()` | `db.Where(...).First()` (user), bcrypt compare | auth.go:33 |
| POST | `/login/2fa` | `RateLimitMiddleware()` | `login2FA` | `NewTwoFactorService().ValidateCode()`, `GetJWTSecret()`, `NewLoginHistoryService().RecordLoginAttempt()` | `db.First()` (user) | auth.go:34 |
| POST | `/logout` | none | `logout` | — | none (no-op stub) | auth.go:35 |
| GET | `/me` | `AuthMiddleware()` | `getCurrentUser` | — | reads `c.Get("user")` from context | auth.go:36 |
| POST | `/refresh` | none | `refreshToken` | — | none (stub) | auth.go:37 |
| POST | `/password-reset` | `RateLimitMiddleware()` | `passwordReset` | — | `db.Where(...).First()` (user), Redis SET or in-memory map | auth.go:38 |
| POST | `/password-verify` | `RateLimitMiddleware()` | `passwordVerify` | — | `db.Where(...).First()` (user), Redis GET, `db.Model().Update()` (password_hash) | auth.go:39 |

**Services used:**
- `services.GetJWTSecret()` → reads `JWT_SECRET_KEY` env var
- `services.NewLoginHistoryService()` → `services/login_history.go`
- `services.NewTwoFactorService()` → `services/two_factor.go`

**Known gaps:**
- `logout` is a no-op (line 291) — does not invalidate sessions or tokens
- `refreshToken` is a stub (line 321) — returns static message, no token rotation
- JWT claims only contain `sub`, `exp`, `iat` — missing `jti`, `role`, `admin` claims

---

## 2. Wallet Routes (`/api/wallets`)

**Setup:** `SetupWalletRoutes()` in `routes/wallets.go:20`
**Base group:** `router.Group("/api/wallets")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/` | `AuthMiddleware()` | `getWallets` | — | `db.Where("user_id = ?", currentUser.ID).Find()` | wallets.go:23 |
| GET | `/balance` | `AuthMiddleware()` | `getBalance` | — | `db.Where("user_id = ?").First()` (wallet), `db.Model().Count()` (transactions) | wallets.go:24 |
| POST | `/deposit` | `AuthMiddleware()` | `deposit` | — | `db.Begin()`, `tx.Model().Update()` (balance), `tx.Create()` (transaction), `tx.Commit()` | wallets.go:25 |
| POST | `/transfer` | `AuthMiddleware()`, `IdempotencyMiddleware("transfer")` | `transfer` | `NewSecurityDetector().DetectIDOR()` (on IDOR) | `db.Where(...).First()` (sender, recipient wallets), `db.Begin()`, `tx.Model().Update()` (both balances), `tx.Create()` (2 transactions), `tx.Commit()` | wallets.go:26 |
| GET | `/:id` | `AuthMiddleware()` | `getWallet` | `NewSecurityDetector().DetectIDOR()` | `db.Preload("User").Where("id = ?").First()` — **NO ownership check** | wallets.go:27 |
| POST | `/` | `AuthMiddleware()` | `createWallet` | — | none (stub) | wallets.go:28 |
| PUT | `/:id` | `AuthMiddleware()` | `updateWallet` | — | none (stub) | wallets.go:29 |
| DELETE | `/:id` | `AuthMiddleware()` | `deleteWallet` | — | none (stub) | wallets.go:30 |

**Services used:**
- `services.NewSecurityDetector()` → `services/security_detector.go`

**Known gaps:**
- `getWallet` (line 348): **IDOR vulnerability** — no ownership validation; any authenticated user can read any wallet by ID. Security detector fires but does NOT block the response.
- `transfer` (line 217): **No `SELECT FOR UPDATE`** — uses plain `db.Begin()` without row-level locking; race condition possible under concurrent transfers.
- `transfer`: **Idempotency record never created** — `IdempotencyMiddleware` checks for existing records but the handler never calls `CreateIdempotencyRecord()` after success.
- `createWallet`, `updateWallet`, `deleteWallet` are stubs returning placeholder JSON.
- Balance type: `float64` in Go model vs `DECIMAL(15,2)` in MySQL — precision loss risk.

---

## 3. Transaction Routes (`/api/transactions`)

**Setup:** `SetupTransactionRoutes()` in `routes/transactions.go:15`
**Base group:** `router.Group("/api/transactions")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/` | `AuthMiddleware()` | `getTransactions` | — | `db.Where("user_id = ?").First()` (wallet), `db.Where("wallet_id = ?").Order().Limit().Find()` (transactions) | transactions.go:18 |
| GET | `/:id` | `AuthMiddleware()` | `getTransaction` | — | none (stub) | transactions.go:19 |
| POST | `/` | `AuthMiddleware()` | `createTransaction` | — | none (stub) | transactions.go:20 |
| PUT | `/:id` | `AuthMiddleware()` | `updateTransaction` | — | none (stub) | transactions.go:21 |
| DELETE | `/:id` | `AuthMiddleware()` | `deleteTransaction` | — | none (stub) | transactions.go:22 |

**Known gaps:** Only `getTransactions` is implemented. All other CRUD endpoints are stubs.

---

## 4. Admin Routes (`/api/admin`)

**Setup:** `SetupAdminRoutes()` in `routes/admin.go:16`
**Base group:** `router.Group("/api/admin")` with `AuthMiddleware()` + `AdminMiddleware()` on all routes.

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/dashboard` | `AuthMiddleware()`, `AdminMiddleware()` | `getDashboard` | — | none (stub) | admin.go:23 |
| GET | `/users` | `AuthMiddleware()`, `AdminMiddleware()` | `getAdminUsers` | — | none (stub) | admin.go:24 |
| GET | `/transactions` | `AuthMiddleware()`, `AdminMiddleware()` | `getAdminTransactions` | — | `db.Preload("Wallet.User").Order().Limit().Find()` (all transactions) | admin.go:25 |
| POST | `/users/:id/disable` | `AuthMiddleware()`, `AdminMiddleware()` | `disableUser` | — | none (stub) | admin.go:27 |
| POST | `/users/:id/enable` | `AuthMiddleware()`, `AdminMiddleware()` | `enableUser` | — | none (stub) | admin.go:28 |
| GET | `/settings` | `AuthMiddleware()`, `AdminMiddleware()` | `getSystemSettings` | — | none (hardcoded settings) | admin.go:29 |
| POST | `/settings` | `AuthMiddleware()`, `AdminMiddleware()` | `saveSystemSettings` | — | none (stub) | admin.go:30 |
| GET | `/support/tickets` | `AuthMiddleware()`, `AdminMiddleware()` | `getAdminSupportTickets` | — | `db.Preload("User").Order().Find()` (all tickets) | admin.go:32 |
| POST | `/support/tickets/:id/reply` | `AuthMiddleware()`, `AdminMiddleware()` | `replyToTicket` | — | `db.First()` (ticket), `db.Save()` (ticket status) | admin.go:33 |
| POST | `/support/tickets/:id/resolve` | `AuthMiddleware()`, `AdminMiddleware()` | `resolveTicket` | — | `db.First()` (ticket), `db.Save()` (ticket status) | admin.go:34 |

**Known gaps:**
- `getDashboard`, `getAdminUsers`, `disableUser`, `enableUser`, `getSystemSettings`, `saveSystemSettings` are stubs or return hardcoded data.
- `getAdminTransactions` is fully implemented with proper `Preload("Wallet.User")` and response transformation.
- Support ticket reply/resolve have mock ticket ID handling hardcoded.

---

## 5. User Routes (`/api/users`)

**Setup:** `SetupUserRoutes()` in `routes/users.go:53`
**Base group:** `router.Group("/api/users")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/` | `AuthMiddleware()` | `getUsers` | — | `db.Limit(100).Find()` (users) — but **overridden by mock data** | users.go:56 |
| GET | `/search` | `AuthMiddleware()` | `searchUsers` | — | `db.Where("username LIKE ? OR email LIKE ?").Where("is_active = ?").Limit(10).Find()` | users.go:57 |
| GET | `/:id` | `AuthMiddleware()` | `getUser` | — | `db.First()` (user) — admin-only check in handler | users.go:58 |
| POST | `/` | `AuthMiddleware()` | `createUser` | — | **mock data only** — no DB write | users.go:59 |
| PUT | `/:id` | `AuthMiddleware()` | `updateUser` | — | Mock data path + `db.First()` + `db.Save()` | users.go:60 |
| DELETE | `/:id` | `AuthMiddleware()` | `deleteUser` | — | Mock data path + `db.Begin()` cascade delete + `tx.Delete()` | users.go:61 |
| DELETE | `/account` | `AuthMiddleware()` | `deleteCurrentUserAccount` | — | `db.Begin()` cascade delete (login_history, sessions, tickets, audit_logs, transactions, wallets, user) + `tx.Delete()` | users.go:62 |

**Known gaps:**
- `getUsers` and `createUser` use in-memory mock data instead of the database.
- `updateUser` and `deleteUser` have dual paths: mock data first, then DB fallback.
- `deleteCurrentUserAccount` and `deleteUser` are the most fully-implemented handlers — proper cascade delete with transaction.

---

## 6. Two-Factor Routes (`/api/2fa`)

**Setup:** `SetupTwoFactorRoutes()` in `routes/two_factor.go:15`
**Base group:** `router.Group("/api/2fa")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| POST | `/enable` | `AuthMiddleware()` | `enable2FA` | `NewTwoFactorService().ValidateCode()` | `db.First()` (user), `db.Model().Update()` (two_factor_enabled) | two_factor.go:18 |
| POST | `/disable` | `AuthMiddleware()` | `disable2FA` | `NewTwoFactorService().ValidateCode()` | `db.First()` (user), `db.Model().Updates()` (two_factor_enabled=false, two_factor_secret="") | two_factor.go:19 |
| POST | `/verify` | `AuthMiddleware()` | `verify2FA` | `NewTwoFactorService().ValidateCode()` | `db.First()` (user) | two_factor.go:20 |
| GET | `/status` | `AuthMiddleware()` | `get2FAStatus` | `NewTwoFactorService().GenerateSecret()` | `db.First()` (user), `db.Model().Update()` (two_factor_secret) | two_factor.go:21 |

**Services used:**
- `services.NewTwoFactorService()` → `services/two_factor.go` (TOTP via `pquerna/otp`)

**Known gaps:**
- No recovery codes mechanism — if user loses TOTP device, account is locked.
- `get2FAStatus` generates and saves a new secret every time 2FA is not enabled (side effect on GET).

---

## 7. Support Routes (`/api/support`)

**Setup:** `SetupSupportRoutes()` in `routes/support.go:13`
**Base group:** `router.Group("/api/support")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/tickets` | `AuthMiddleware()` | `getTickets` | — | `db.Where("user_id = ?", currentUser.ID).Order().Find()` | support.go:16 |
| GET | `/tickets/:id` | `AuthMiddleware()` | `getTicket` | — | none (stub) | support.go:17 |
| POST | `/tickets` | `AuthMiddleware()` | `createTicket` | — | `db.Create()` (ticket) | support.go:18 |
| PUT | `/tickets/:id` | `AuthMiddleware()` | `updateTicket` | — | none (stub) | support.go:19 |
| DELETE | `/tickets/:id` | `AuthMiddleware()` | `deleteTicket` | — | none (stub) | support.go:20 |

**Known gaps:** Only `getTickets` and `createTicket` are implemented.

---

## 8. Login History Routes (`/api/login-history`)

**Setup:** `SetupLoginHistoryRoutes()` in `routes/login_history.go:15`
**Base group:** `router.Group("/api/login-history")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/` | `AuthMiddleware()` | `getLoginHistory` | `NewLoginHistoryService().GetLoginHistory()` | DB query via service | login_history.go:18 |
| GET | `/recent` | `AuthMiddleware()` | `getRecentLoginHistory` | `NewLoginHistoryService().GetLoginHistory()` | DB query via service | login_history.go:19 |

**Services used:**
- `services.NewLoginHistoryService()` → `services/login_history.go`

---

## 9. Backup Routes (`/api/backup`)

**Setup:** `SetupBackupRoutes()` in `routes/backup.go:12`
**Base group:** `router.Group("/api/backup")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/` | `AuthMiddleware()` | `listBackups` | `NewBackupService().ListBackups()` | filesystem | backup.go:17 |
| GET | `/:filename` | `AuthMiddleware()` | `getBackupInfo` | `NewBackupService().GetBackupInfo()` | filesystem | backup.go:18 |
| GET | `/stats` | `AuthMiddleware()` | `getBackupStats` | `NewBackupService().GetBackupStats()` | filesystem | backup.go:19 |
| GET | `/config` | `AuthMiddleware()` | `getBackupConfig` | — | returns `DefaultBackupConfig` | backup.go:20 |

**Services used:**
- `services.NewBackupService()` → `services/backup.go`

**Known gaps:** Backup endpoints accessible to any authenticated user (not admin-only).

---

## 10. Security Routes (`/api/security`)

**Setup:** `SetupSecurityRoutes()` in `routes/security.go:13`
**Base group:** `router.Group("/api/security")`

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| GET | `/idor/stats` | `AuthMiddleware()`, `AdminOnlyMiddleware()` | `getIDORStats` | `NewSecurityDetector().GetIDORStats()` | Redis + DB (alerts) | security.go:18 |
| GET | `/alerts` | `AuthMiddleware()`, `AdminOnlyMiddleware()` | `getSecurityAlerts` | `NewSecurityDetector().GetSecurityAlerts()` | DB (alerts) | security.go:19 |
| PUT | `/alerts/:id/status` | `AuthMiddleware()`, `AdminOnlyMiddleware()` | `updateAlertStatus` | `NewSecurityDetector().UpdateAlertStatus()` | DB (alerts) | security.go:20 |
| POST | `/users/:id/reset-attempts` | `AuthMiddleware()`, `AdminOnlyMiddleware()` | `resetUserAttempts` | `NewSecurityDetector().ResetUserAttempts()` | Redis (del keys) | security.go:21 |
| POST | `/cleanup` | `AuthMiddleware()`, `AdminOnlyMiddleware()` | `cleanupSecurityData` | `NewSecurityDetector().CleanupOldData()` | DB (delete old alerts) | security.go:22 |

**Services used:**
- `services.NewSecurityDetector()` → `services/security_detector.go`

---

## 11. Data Management Routes (`/api/data`)

**Setup:** `SetupDataManagementRoutes()` in `routes/data_management.go:11`
**Base group:** `router.Group("/api/data")` — **NO AUTH MIDDLEWARE**

| Method | Path | Middleware | Handler | Service Calls | Data Access | Line |
|---|---|---|---|---|---|---|
| POST | `/init-db` | none | `initDatabase` | `NewSampleDataManager().InitializeDatabase()` | GORM AutoMigrate | data_management.go:14 |
| DELETE | `/clear-sample` | none | `clearSampleData` | `NewSampleDataManager().ClearSampleData()` | DB DELETE | data_management.go:15 |
| GET | `/stats` | none | `getDataStats` | `NewSampleDataManager().GetSampleDataStats()` | DB COUNTs | data_management.go:16 |
| POST | `/reset-database` | none | `resetDatabase` | `NewSampleDataManager().ResetDatabase()` | cascade DELETE + reseed | data_management.go:17 |
| POST | `/force-recreate` | none | `forceDatabaseRecreation` | `NewSampleDataManager().CompleteDatabaseRecreation()` + reseed | DROP + CREATE + seed | data_management.go:18 |

**Services used:**
- `services.NewSampleDataManager()` → `services/data_manager.go`

**Known gaps:**
- **CRITICAL:** All data management endpoints are **unauthenticated and unprotected**. Anyone can reset or recreate the database.

---

## 12. Blog Routes (`/api/blog`) — Public

**Setup:** `BlogRoutes()` in `routes/blog.go:18`
**Base group:** `router.Group("/api/blog")` — registered on `*gin.Engine` directly, no auth.

| Method | Path | Middleware | Handler | Data Access | Line |
|---|---|---|---|---|---|
| GET | `/posts` | none | anonymous inline | `db.Table("blog_posts").Joins("LEFT JOIN users ...").Where("status = published").Offset().Limit().Find()` | blog.go:22 |
| GET | `/posts/:slug` | none | anonymous inline | `db.Table("blog_posts").Joins(...).Where("slug = ? AND status = published").First()` + comments + related posts | blog.go:103 |
| GET | `/posts/:slug/comments` | none | anonymous inline | `db.Table("blog_comments").Where("post_id = ? AND status = approved").Offset().Limit().Find()` | blog.go:268 |
| POST | `/posts/:slug/comments` | none | anonymous inline | `db.Table("blog_comments").Create()` (status=pending) | blog.go:324 |
| GET | `/categories` | none | anonymous inline | `db.Table("blog_categories").Find()` | blog.go:370 |
| GET | `/tags` | none | anonymous inline | `db.Table("blog_tags").Find()` | blog.go:388 |
| GET | `/comments/stats` | none | anonymous inline | `db.Table("blog_comments").Count()` (total, pending, approved) | blog.go:404 |

---

## 13. Cron Routes (`/api/cron`) — Unauthenticated

**Setup:** `CronRoutes()` in `routes/cron.go:11`
**Base group:** `router.Group("/api/cron")` — **NO AUTH MIDDLEWARE**

| Method | Path | Middleware | Handler | Service Calls | Line |
|---|---|---|---|---|---|
| GET | `/status` | none | anonymous inline | `NewCronService().GetCronStatus()` | cron.go:15 |
| POST | `/execute/:job` | none | anonymous inline | `NewCronService().ExecuteCronJob()` | cron.go:22 |
| POST | `/setup` | none | anonymous inline | `NewCronService().SetupCronJobs()` | cron.go:41 |
| DELETE | `/remove` | none | anonymous inline | `NewCronService().RemoveCronJobs()` | cron.go:51 |
| GET | `/comments/stats` | none | anonymous inline | `NewCommentService().GetCommentStats()` | cron.go:67 |
| GET | `/backup/stats` | none | anonymous inline | `NewBackupService().ListBackups()` | cron.go:74 |

**Known gaps:**
- **CRITICAL:** All cron endpoints are **unauthenticated**. Anyone can execute cron jobs, setup/remove cron schedules.

---

## Middleware Inventory

| Middleware | File | Function | Behavior |
|---|---|---|---|
| `AuthMiddleware()` | `middleware/auth.go:14` | JWT validation | Extracts `Bearer` token, calls `services.GetCurrentUser()`, sets `c.Set("user", user)` |
| `OptionalAuthMiddleware()` | `middleware/auth.go:48` | Optional JWT | Same as above but does not abort on failure |
| `AdminMiddleware()` | `middleware/auth.go:79` | Admin check | Checks `user.(*models.User).IsAdmin`, returns 403 if false |
| `AdminOnlyMiddleware()` | `middleware/admin.go:11` | Admin check (alt) | Same as AdminMiddleware, slightly different error message |
| `RateLimitMiddleware()` | `middleware/rate_limit.go:19` | IP-based rate limit | In-memory map, 5 req/15min per IP. **Skipped in development mode** |
| `IdempotencyMiddleware(op)` | `middleware/idempotency.go:20` | Idempotency check | Requires `Idempotency-Key` header (UUID), checks `idempotency_records` table, returns cached response if found |
| `SecurityHeadersMiddleware()` | `middleware/auth.go:106` | Security headers | Sets 6 security headers on every response |
| `InputValidationMiddleware()` | `middleware/validation.go:12` | Input sanitization | Blocks SQL/XSS/command patterns in query and path params |

---

## Service Inventory

| Service | File | Key Methods | Used By |
|---|---|---|---|
| `GetJWTSecret()` | `services/auth.go:18` | Returns `JWT_SECRET_KEY` env var | auth routes, auth middleware |
| `AuthenticateUser()` | `services/auth.go:32` | DB lookup + bcrypt | not called by routes (logic inlined in `login`) |
| `CreateAccessToken()` | `services/auth.go:49` | JWT creation with claims | not called by routes (logic inlined in `login`) |
| `GetCurrentUser()` | `services/auth.go:78` | JWT parse + DB lookup | `AuthMiddleware()` |
| `NewTwoFactorService()` | `services/two_factor.go:16` | `GenerateSecret`, `ValidateCode`, `GenerateQRCodeURL` | 2FA routes, login2FA |
| `NewLoginHistoryService()` | `services/login_history.go:17` | `RecordLoginAttempt`, `GetLoginHistory`, `GetFailedLoginAttempts` | login, login2FA, login-history routes |
| `NewSecurityDetector()` | `services/security_detector.go:74` | `DetectIDOR`, `DetectSecurityEvent`, `GetIDORStats`, `GetSecurityAlerts`, `UpdateAlertStatus`, `ResetUserAttempts`, `CleanupOldData` | wallet getWallet, security routes |
| `NewBackupService()` | `services/backup.go` | `ListBackups`, `GetBackupInfo`, `GetBackupStats` | backup routes, cron routes |
| `NewCronService()` | `services/cron.go` | `GetCronStatus`, `ExecuteCronJob`, `SetupCronJobs`, `RemoveCronJobs` | cron routes |
| `NewCommentService()` | `services/comment.go` | `GetCommentStats` | cron routes |
| `NewSampleDataManager()` | `services/data_manager.go` | `InitializeDatabase`, `ClearSampleData`, `GetSampleDataStats`, `ResetDatabase`, `CompleteDatabaseRecreation`, `CreateSampleUsers/Wallets/Transactions/LoginHistory` | data management routes |

---

## Model-to-Table Mapping

| Go Model | File | Table | PK | FK | Notes |
|---|---|---|---|---|---|
| `User` | `models/user.go` | `users` | `id` CHAR(36) | — | `BeforeCreate` UUID hook, soft delete |
| `Wallet` | `models/wallet.go` | `wallets` | `id` CHAR(36) | `user_id` → users(id) | Balance as `float64` (precision risk) |
| `Transaction` | `models/transaction.go` | `transactions` | `id` CHAR(36) | `wallet_id` → wallets(id) | Amount as `float64` (precision risk) |
| `Session` | `models/session.go` | `sessions` | `id` CHAR(36) | `user_id` → users(id) | Not used by auth routes |
| `AuditLog` | `models/audit_log.go` | `audit_logs` | `id` CHAR(36) | `user_id` → users(id) | Created manually in getWallet |
| `IdempotencyRecord` | `models/idempotency_record.go` | `idempotency_records` | `id` CHAR(36) | `user_id` → users(id) | **Not in AutoMigrate** |
| `LoginHistory` | `models/login_history.go` | `login_history` | `id` CHAR(36) | `user_id` → users(id) | — |
| `SupportTicket` | `models/support_ticket.go` | `support_tickets` | `id` CHAR(36) | `user_id` → users(id) | — |
| `SecurityAlert` | (no model file) | `security_alerts` | `id` VARCHAR(100) | `user_id` → users(id) | Struct defined in security_detector.go only |

---

## Critical Security Findings

| Severity | Location | Issue |
|---|---|---|
| **CRITICAL** | `/api/data/*` | All 5 endpoints unauthenticated — anyone can reset/recreate database |
| **CRITICAL** | `/api/cron/*` | All 6 endpoints unauthenticated — anyone can execute/setup/remove cron jobs |
| **HIGH** | `wallets.go:348` | IDOR in `getWallet` — any authenticated user can read any wallet |
| **HIGH** | `wallets.go:217` | Transfer has no `SELECT FOR UPDATE` — race condition under concurrency |
| **HIGH** | `wallets.go:26` | Idempotency middleware checks but handler never creates records |
| **HIGH** | `auth.go:291` | Logout is a no-op — tokens remain valid |
| **MEDIUM** | `auth.go:321` | Refresh token is a stub — no token rotation |
| **MEDIUM** | `rate_limit.go` | In-memory only, skipped in development, no Redis backing |
| **MEDIUM** | `models/wallet.go` | `float64` for money — precision loss vs MySQL DECIMAL |
| **LOW** | `two_factor.go:203` | GET `/2fa/status` has side effect (generates/saves secret) |
