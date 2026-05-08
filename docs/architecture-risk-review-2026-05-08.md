# SecureWallet — Architecture Risk Review

**Date:** 2026-05-08
**Assessed against:** OWASP ASVS v4.0.3, PCI-DSS v4.0, Financial Regulatory Baselines
**Status:** `IMPLEMENTED_NOT_VERIFIED` — findings from direct file inspection; runtime verification not performed.

---

## 1. Executive Summary

SecureWallet implements a Go 1.21 / Gin backend with MySQL 8, Redis, and Vue.js frontend for digital banking. The codebase has made meaningful progress on JWT authentication, TOTP 2FA, transfer service with row-level locking, idempotency records, audit logging, and security detection. However, **critical gaps remain** across authentication hardening, authorization enforcement, data management endpoint exposure, financial precision, and operational security controls.

### Risk Distribution

| Severity | Count | Summary |
|---|---|---|
| Critical | 6 | Unauthenticated destructive endpoints, IDOR, float64 for money, missing auth on admin mutation, password reset token exposure, no CSRF protection |
| High | 9 | In-memory rate limiting, stub admin endpoints, bcrypt.DefaultCost, no JWT revocation on access tokens, weak input validation, CORS debug logging, mock data in production paths |
| Medium | 7 | No structured logging, no log rotation, no request size limits, no CSRF tokens, TOTP rate limit bypass in dev, no account lockout on login, missing audit on admin actions |
| Low | 4 | Unused MongoDB container, frontend lint broken, swagger incomplete, version header leaks |

---

## 2. OWASP ASVS v4.0.3 Assessment

### V1 — Architecture, Design and Threat Modeling

| Requirement | Status | Finding |
|---|---|---|
| 1.1.1 — Security requirements documented | FAIL | No threat model, no security requirements document |
| 1.2.1 — Tiered architecture | PASS | Clear separation: routes → services → models → config |
| 1.2.4 — Centralized security controls | PARTIAL | Middleware exists but not applied consistently (see `/api/data/*`, `/api/cron/*`) |
| 1.3.1 — Secure defaults | FAIL | `RESET_DATABASE_ON_STARTUP` defaults to `false` in env.example but is a production risk; `FORCE_DATABASE_RECREATION` exposed |

### V2 — Authentication

| Requirement | Status | Finding |
|---|---|---|
| 2.1.1 — Password policy | PASS | 12+ chars, upper+lower+digit+special enforced (`auth.go:577-595`) |
| 2.1.3 — Password hashing | PARTIAL | bcrypt used but `DefaultCost` (10) — should be `MinCost` or higher (12+) for financial apps |
| 2.1.5 — Password reset | FAIL | Reset token stored in Redis with 24h TTL but **no rate limit on reset requests per email**; fallback to in-memory `secondOrderStorage` map is not thread-safe |
| 2.2.1 — JWT algorithm | PASS | HS256 explicitly validated (`auth.go:264`) |
| 2.2.3 — JWT expiry | PASS | Access token: 30min configurable; refresh token: 7 days configurable |
| 2.2.6 — JWT revocation | PARTIAL | Refresh tokens stored in DB and revocable; **access tokens NOT revocable** (no jti blacklist check in `GetCurrentUser`) |
| 2.2.8 — Token rotation | PASS | `RotateRefreshToken` invalidates old, issues new in DB transaction (`auth.go:160-195`) |
| 2.3.1 — TOTP implementation | PASS | `pquerna/otp` with crypto/rand secret generation |
| 2.3.2 — TOTP rate limiting | PARTIAL | Redis-backed rate limit exists (`rate_limit_2fa.go`) but **bypassed in dev mode**; no lockout after recovery code exhaustion |
| 2.4.1 — Brute force protection | FAIL | Login rate limit is global IP-based (5/15min, in-memory), **no per-account lockout**; `rate_limit.go` skipped in dev/debug |

### V3 — Session Management

| Requirement | Status | Finding |
|---|---|---|
| 3.1.1 — Session ID generation | PASS | Refresh tokens: 32 bytes from `crypto/rand` |
| 3.2.1 — Session expiration | PASS | Refresh tokens have configurable expiry; expired tokens deleted on validation |
| 3.3.1 — Session termination | PARTIAL | `RevokeToken` and `RevokeAllUserTokens` exist; logout endpoint revokes only if `refresh_token` provided in body |
| 3.4.1 — Concurrent sessions | FAIL | No limit on concurrent sessions per user |
| 3.5.1 — Remember me | N/A | Not implemented |

### V4 — Access Control

| Requirement | Status | Finding |
|---|---|---|
| 4.1.1 — Explicit deny by default | PASS | Routes require `AuthMiddleware()` by default |
| 4.1.2 — Unique resource IDs | PASS | UUIDs used throughout |
| 4.1.3 — Horizontal access control | **FAIL** | `getWallet` in `wallets.go:317-347` checks `id AND user_id` — this is correct. However, `getTransaction` (`transactions.go:87-90`) is a **stub** returning mock data with no ownership check |
| 4.1.4 — Vertical access control | PARTIAL | `AdminMiddleware` checks `IsAdmin` flag; however, `disableUser`/`enableUser`/`saveSystemSettings` are **stubs** with no implementation |
| 4.2.1 — Directory listing | PASS | Not applicable (API-only) |
| 4.3.1 — CSRF protection | **FAIL** | No CSRF tokens; CORS allows credentials; state-changing operations rely solely on JWT bearer auth |

### V5 — Validation, Sanitization, Encoding

| Requirement | Status | Finding |
|---|---|---|
| 5.1.1 — Input validation | PARTIAL | `InputValidationMiddleware` blocks SQL/XSS patterns but **does not validate request body**; blacklist approach is bypassable |
| 5.1.3 — Output encoding | PASS | Gin auto-encodes JSON responses |
| 5.2.1 — SQL injection | PASS | GORM parameterized queries used throughout |
| 5.2.3 — Command injection | PASS | No shell execution in code paths |
| 5.3.1 — File upload validation | N/A | No file upload endpoints |
| 5.4.1 — SSRF | PASS | No outbound HTTP requests from user input |
| 5.5.1 — Request size limits | **FAIL** | No `Content-Length` or body size limit middleware |

### V6 — Cryptography

| Requirement | Status | Finding |
|---|---|---|
| 6.1.1 — Approved algorithms | PASS | bcrypt, HS256, SHA-256, TOTP (RFC 6238) |
| 6.2.1 — Random value generation | PASS | `crypto/rand` used for tokens, TOTP secrets, recovery codes |
| 6.3.1 — Secrets management | FAIL | JWT secret, DB password, Redis password in env vars with `CHANGE_THIS` placeholders; no secret rotation mechanism |

### V7 — Error Handling and Logging

| Requirement | Status | Finding |
|---|---|---|
| 7.1.1 — Generic error messages | PARTIAL | Auth errors use generic "Incorrect username or password"; but transfer errors leak details (fee, total_amount) |
| 7.1.2 — No stack traces | PASS | No stack traces in error responses |
| 7.2.1 — Structured logging | **FAIL** | `log.Printf` and `log.Println` used; no JSON structured logging, no log levels |
| 7.2.2 — Sensitive data in logs | PARTIAL | `sanitizeDetails` and `sanitizeMetadata` redact known patterns; but `log.Printf("CORS Debug - Origin: %s")` in `main.go:112` may leak origin headers in production |
| 7.3.1 — Log protection | FAIL | No log file permissions set; no log rotation configured |

### V8 — Data Protection

| Requirement | Status | Finding |
|---|---|---|
| 8.1.1 — Sensitive data in transit | PARTIAL | HSTS header set; but no TLS enforcement at application layer |
| 8.2.1 — Sensitive data at rest | PARTIAL | Passwords hashed; TOTP secrets stored in DB column (should be encrypted) |
| 8.3.1 — Financial data precision | **FAIL** | Go models use `float64` for `Balance` and `Amount`; MySQL uses `DECIMAL(15,2)`. Floating-point arithmetic can cause rounding errors in balance calculations |
| 8.4.1 — PII minimization | PARTIAL | User responses exclude `PasswordHash` and `TwoFactorSecret` via `json:"-"`; but `Name`, `Title`, `Avatar`, `Bio` fields returned in `/auth/me` |

### V9 — Communications

| Requirement | Status | Finding |
|---|---|---|
| 9.1.1 — TLS | PASS | HSTS header configured; TLS expected at reverse proxy |
| 9.1.3 — Certificate pinning | N/A | Not applicable for web API |

### V10 — Malicious Code

| Requirement | Status | Finding |
|---|---|---|
| 10.1.1 — Dependency review | PARTIAL | Dependencies appear reasonable; no known vulnerable versions in go.mod |
| 10.2.1 — Code signing | N/A | Not applicable |

### V11 — Business Logic

| Requirement | Status | Finding |
|---|---|---|
| 11.1.1 — Business logic flows documented | FAIL | No documented business logic flows |
| 11.2.1 — Transaction integrity | PARTIAL | Transfer service uses DB transaction with `SELECT FOR UPDATE` locking; but **no daily/monthly transfer limits enforced** |
| 11.3.1 — Idempotency | PARTIAL | Idempotency middleware and service-layer checks exist; but middleware has race condition (check-then-act without DB lock) |

### V12 — Files and Resources

| Requirement | Status | Finding |
|---|---|---|
| 12.1.1 — File access control | N/A | No file serving endpoints |

### V13 — API and Web Service

| Requirement | Status | Finding |
|---|---|---|
| 13.1.1 — API schema validation | PARTIAL | Gin binding validates required fields; but no JSON schema validation |
| 13.2.1 — API versioning | FAIL | No API versioning (all routes at `/api/`) |
| 13.3.1 — Rate limiting per client | PARTIAL | IP-based rate limiting exists; no per-user API rate limits |

### V14 — Configuration

| Requirement | Status | Finding |
|---|---|---|
| 14.1.1 — Secure configuration | FAIL | `env.example` contains `CHANGE_THIS` placeholders; no config validation at startup |
| 14.2.1 — Dependency versions | PASS | Go 1.21, pinned dependency versions |
| 14.3.1 — HTTP security headers | PASS | `SecurityHeadersMiddleware` sets X-Content-Type-Options, X-Frame-Options, CSP, HSTS, Referrer-Policy |

---

## 3. PCI-DSS v4.0 Assessment

### Requirement 1 — Network Security Controls

| Control | Status | Finding |
|---|---|---|
| Firewall rules | PARTIAL | Docker network isolation configured; no application-level firewall rules |
| Network segmentation | PARTIAL | Services on `securewallet_network`; but MongoDB is unused and should be removed |

### Requirement 2 — Secure Configuration

| Control | Status | Finding |
|---|---|---|
| Default passwords | **FAIL** | `env.example` and `docker-compose.yml` contain `CHANGE_THIS_*` placeholders; no startup validation rejects defaults |
| Unnecessary services | FAIL | MongoDB container defined but unused |
| Security parameters | PARTIAL | CORS, JWT expiry configurable; no validation of secure ranges |

### Requirement 3 — Protect Stored Account Data

| Control | Status | Finding |
|---|---|---|
| PAN storage | N/A | No card data stored |
| Cryptographic protection | PARTIAL | Passwords hashed with bcrypt; TOTP secrets stored in plaintext in DB |
| Key management | FAIL | No key rotation for JWT secret; single shared secret |

### Requirement 4 — Protect Data in Transit

| Control | Status | Finding |
|---|---|---|
| TLS for transmission | PARTIAL | HSTS header set; TLS expected at proxy layer; no application-level TLS |

### Requirement 6 — Secure Development

| Control | Status | Finding |
|---|---|---|
| Secure SDLC | FAIL | No CI/CD, no automated testing, no code review process |
| Vulnerability management | FAIL | No dependency vulnerability scanning |
| Code review | FAIL | No evidence of peer review process |

### Requirement 7 — Access Control

| Control | Status | Finding |
|---|---|---|
| Least privilege | PARTIAL | Admin/non-admin roles exist; but `IsAdmin` is a simple boolean with no granular permissions |
| Access assignment | FAIL | No approval workflow for admin role assignment |

### Requirement 8 — Authentication

| Control | Status | Finding |
|---|---|---|
| Unique IDs | PASS | UUID-based user IDs |
| MFA | PASS | TOTP 2FA implemented with recovery codes |
| Password complexity | PASS | 12+ chars with character class requirements |
| Account lockout | FAIL | No account lockout after failed login attempts |
| Session timeout | PASS | 30-min access token, 7-day refresh token |

### Requirement 10 — Logging and Monitoring

| Control | Status | Finding |
|---|---|---|
| Audit trails | PARTIAL | `AuditLog` model captures action, resource, result, IP, user agent, correlation ID |
| Log integrity | FAIL | Audit logs stored in same DB with `ON DELETE CASCADE` — can be deleted with user |
| Log retention | FAIL | No retention policy; no log archival |
| Real-time alerting | PARTIAL | `SecurityDetector` with Redis-based detection; but `sendRealTimeAlert` is TODO |
| Time synchronization | FAIL | No NTP configuration; relies on host clock |

### Requirement 11 — Security Testing

| Control | Status | Finding |
|---|---|---|
| Penetration testing | FAIL | No pen testing evidence |
| Vulnerability scanning | FAIL | No automated scanning |
| Intrusion detection | PARTIAL | `SecurityDetector` provides basic IDOR/brute-force detection |

### Requirement 12 — Security Policy

| Control | Status | Finding |
|---|---|---|
| Security policy | FAIL | No security policy document |
| Incident response | FAIL | No incident response plan |
| Risk assessment | FAIL | This is the first formal risk assessment |

---

## 4. Financial Regulatory Baseline Assessment

### 4.1 Transaction Integrity

| Requirement | Status | Finding |
|---|---|---|
| Atomic transfers | PASS | `Transfer()` uses `db.Transaction()` with `SELECT FOR UPDATE` row locking |
| Double-entry bookkeeping | FAIL | Transfer creates separate `transfer_out` and `transfer_in` records but **no linking reference** between them; no counterparty_wallet_id populated |
| Audit trail | PARTIAL | Transaction records exist; but no immutable audit log (soft deletes enabled) |
| Reconciliation | FAIL | No reconciliation mechanism; no daily balance verification |
| Idempotency | PARTIAL | Idempotency key supported but middleware has race condition |

### 4.2 Financial Data Precision

| Requirement | Status | Finding |
|---|---|---|
| Decimal precision | **FAIL** | Go `float64` for balance/amount — must use `decimal.Decimal` or `int64` (cents) |
| Rounding rules | FAIL | No explicit rounding strategy; floating-point arithmetic may produce non-deterministic results |
| Fee calculation | PASS | Fee calculated as percentage with min/max bounds |

### 4.3 Regulatory Reporting

| Requirement | Status | Finding |
|---|---|---|
| Transaction reporting | FAIL | No export or reporting mechanism |
| Suspicious activity reporting | PARTIAL | `SecurityDetector` detects patterns but no SAR workflow |
| Data retention | FAIL | No retention policy defined |

---

## 5. Audit Requirements Specification

### 5.1 Audit Event Catalog

| Event Category | Event Name | Trigger | Required Fields |
|---|---|---|---|
| Authentication | `LOGIN_ATTEMPT` | Login request | user_id (or nil), ip_address, user_agent, result, correlation_id |
| Authentication | `LOGIN_SUCCESS` | Successful login | user_id, ip_address, user_agent, result, correlation_id |
| Authentication | `LOGIN_FAILURE` | Failed login | user_id (if found), ip_address, user_agent, result, correlation_id |
| Authentication | `LOGOUT` | Logout request | user_id, result, correlation_id |
| Authentication | `TOKEN_REFRESH` | Token refresh | user_id, result, correlation_id |
| Authentication | `2FA_ENROLL` | 2FA setup | user_id, result, correlation_id |
| Authentication | `2FA_VERIFY` | 2FA code check | user_id, result, correlation_id |
| Authentication | `PASSWORD_RESET_REQUEST` | Reset requested | user_id, ip_address, result, correlation_id |
| Authentication | `PASSWORD_RESET_COMPLETE` | Password changed | user_id, result, correlation_id |
| Wallet | `WALLET_CREATE` | New wallet | user_id, wallet_id, currency, result, correlation_id |
| Wallet | `WALLET_READ` | Wallet accessed | user_id, wallet_id, result, correlation_id |
| Wallet | `DEPOSIT` | Deposit made | user_id, wallet_id, amount, currency, transaction_id, result, correlation_id |
| Wallet | `TRANSFER_INITIATED` | Transfer started | user_id, recipient_email, amount, currency, idempotency_key, correlation_id |
| Wallet | `TRANSFER_COMPLETED` | Transfer succeeded | user_id, recipient_id, sender_wallet_id, recipient_wallet_id, amount, fee, transaction_id, correlation_id |
| Wallet | `TRANSFER_FAILED` | Transfer failed | user_id, recipient_email, amount, error_code, correlation_id |
| Admin | `USER_DISABLED` | Admin disables user | admin_id, target_user_id, result, correlation_id |
| Admin | `USER_ENABLED` | Admin enables user | admin_id, target_user_id, result, correlation_id |
| Admin | `ADMIN_LOGIN` | Admin logs in | user_id, ip_address, result, correlation_id |
| Security | `BRUTE_FORCE_DETECTED` | Rate limit exceeded | user_id/ip, attempt_count, result, correlation_id |
| Security | `IDOR_ATTEMPT` | Unauthorized resource access | user_id, resource_id, resource_owner, result, correlation_id |
| System | `DATABASE_RESET` | Database reset | user_id (or system), result, correlation_id |
| System | `CONFIG_CHANGE` | Settings updated | admin_id, setting_name, result, correlation_id |

### 5.2 Audit Log Schema Requirements

```sql
-- Recommended audit_logs table enhancements
ALTER TABLE audit_logs
    ADD COLUMN correlation_id CHAR(36) COMMENT 'Request correlation ID',
    ADD COLUMN resource_id CHAR(36) COMMENT 'Target resource UUID',
    ADD COLUMN metadata JSON COMMENT 'Sanitized event metadata',
    ADD COLUMN event_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD INDEX idx_correlation_id (correlation_id),
    ADD INDEX idx_event_timestamp (event_timestamp),
    ADD INDEX idx_result (result);

-- Remove CASCADE delete to preserve audit trail
-- Audit logs must be IMMUTABLE after creation
```

### 5.3 Audit Log Immutability

- Audit logs **MUST NOT** be soft-deletable (`DeletedAt` field prohibited on audit_logs table)
- Audit logs **MUST NOT** have `ON DELETE CASCADE` from user deletion
- Audit log updates **MUST** be append-only (no UPDATE allowed after INSERT)
- Audit log archival **MUST** occur before deletion from primary store

---

## 6. Data Retention Requirements

| Data Type | Retention Period | Rationale | Deletion Method |
|---|---|---|---|
| Audit logs (auth events) | 7 years | Financial regulatory requirement | Archive to cold storage, then purge |
| Audit logs (financial events) | 7 years | PCI-DSS + financial regulation | Archive to cold storage, then purge |
| Audit logs (security events) | 3 years | Security investigation window | Archive to cold storage, then purge |
| Transaction records | 7 years | Financial regulatory requirement | Archive to cold storage, then purge |
| Session records | 30 days after expiry | Session management | Soft delete, then hard delete |
| Login history | 2 years | Security investigation | Soft delete, then hard delete |
| Idempotency records | 24 hours after expiry | Duplicate prevention | Hard delete via TTL |
| Security alerts | 90 days | Incident response | Hard delete after resolution |
| Support tickets | 3 years | Customer service records | Soft delete, then hard delete |
| User personal data | 30 days after account deletion | GDPR right to erasure | Hard delete with audit trail preservation |
| Recovery codes | 90 days | 2FA recovery window | Redis TTL expiry |

---

## 7. Monitoring Requirements

### 7.1 Real-Time Monitoring

| Metric | Threshold | Alert Channel | Response |
|---|---|---|---|
| Failed login attempts per user | >5 in 15 min | Admin dashboard, email | Account lockout for 30 min |
| Failed 2FA attempts per user | >5 in 15 min | Admin dashboard, email | Account lockout for 30 min |
| Transfer failures per user | >10 in 1 hour | Admin dashboard | Manual review trigger |
| IDOR attempts per user | >3 in 5 min | Admin dashboard, security log | IP block consideration |
| Database reset attempts | Any | Admin dashboard, email | Immediate investigation |
| Admin privilege escalation | Any | Admin dashboard, email | Immediate investigation |
| Concurrent session count per user | >5 | Admin dashboard | Review and notify user |

### 7.2 Periodic Monitoring

| Report | Frequency | Recipients | Content |
|---|---|---|---|
| Daily transaction summary | Daily 00:00 UTC | Admin team | Total transfers, deposits, failures, fees |
| Weekly security summary | Weekly Monday | Security team | Failed auth attempts, IDOR detections, lockouts |
| Monthly audit log review | Monthly 1st | Compliance team | Audit log completeness, anomalies |
| Quarterly access review | Quarterly | Admin team | Admin role assignments, privilege changes |
| Annual penetration test | Annually | External auditor | Full security assessment |

### 7.3 Health Monitoring

| Check | Interval | Failure Action |
|---|---|---|---|
| Database connectivity | 30s | Alert, failover consideration |
| Redis connectivity | 30s | Fall back to in-memory rate limiting |
| API response time (p95) | 1 min | Alert if >500ms |
| Error rate (5xx) | 1 min | Alert if >1% of requests |
| Disk space (logs) | 5 min | Alert if >80% used |
| JWT secret validity | Startup | Fail to start if not set |

---

## 8. Critical Findings — Remediation Priority

### P0 — Immediate (Security / Financial Integrity)

| ID | Finding | Location | Remediation |
|---|---|---|---|
| F01 | `/api/data/*` and `/api/cron/*` unauthenticated | `data_management.go`, `cron.go` | Add `AuthMiddleware()` + `AdminOnlyMiddleware()` |
| F02 | `float64` for financial amounts | `wallet.go`, `transaction.go`, `transfer.go` | Replace with `shopspring/decimal` or `int64` cents |
| F03 | `getTransaction` is stub with no ownership check | `transactions.go:87-90` | Implement with wallet ownership verification |
| F04 | No CSRF protection for state-changing operations | `main.go` CORS, all POST/PUT/DELETE | Add CSRF middleware for browser clients |
| F05 | Password reset has no per-email rate limit | `auth.go:431-482` | Add Redis-backed rate limit per email |
| F06 | Audit logs have `ON DELETE CASCADE` from users | `init.sql:84` | Remove CASCADE; audit logs must outlive users |

### P1 — High Priority

| ID | Finding | Location | Remediation |
|---|---|---|---|
| F07 | Rate limiting is in-memory, skipped in dev | `rate_limit.go` | Migrate to Redis; enforce in all environments |
| F08 | bcrypt.DefaultCost (10) too low for financial app | `auth.go:95,299` | Use cost 12+ or `bcrypt.MinCost` |
| F09 | Admin mutation endpoints are stubs | `admin.go:152-228` | Implement disableUser, enableUser, saveSystemSettings |
| F10 | Access tokens not revocable | `auth.go:257-295` | Add jti blacklist check in `GetCurrentUser` |
| F11 | Input validation uses blacklist approach | `validation.go:43-84` | Add whitelist validation for known fields |
| F12 | CORS debug logging in production | `main.go:112-138` | Remove or gate behind debug flag |
| F13 | Mock data paths in production routes | `users.go` | Remove mock data; use DB-only |
| F14 | No request body size limit | `main.go` middleware chain | Add `maxBodySize` middleware (e.g., 1MB) |
| F15 | TOTP secrets stored in plaintext | `user.go:20`, `init.sql:14` | Encrypt TOTP secrets at rest |

### P2 — Medium Priority

| ID | Finding | Location | Remediation |
|---|---|---|---|
| F16 | No structured logging | All packages | Replace `log.Printf` with structured JSON logger |
| F17 | No log rotation | `security_detector.go:324-345` | Implement log rotation (e.g., lumberjack) |
| F18 | No account lockout on login failures | `auth.go:135-194` | Add Redis-backed account lockout |
| F19 | No audit logging on admin actions | `admin.go` | Add audit.Log calls to all admin handlers |
| F20 | Transfer has no daily/monthly limits | `transfer.go:94-210` | Add configurable transfer limits |
| F21 | Idempotency middleware race condition | `idempotency.go:50-76` | Use DB-level unique constraint + INSERT IGNORE |
| F22 | No API versioning | All routes | Add `/api/v1/` prefix |

---

## 9. Compliance Gap Summary

| Framework | Coverage | Critical Gaps |
|---|---|---|
| OWASP ASVS v4.0.3 | ~45% | CSRF, input validation, brute force protection, error handling, configuration |
| PCI-DSS v4.0 | ~30% | Default passwords, logging integrity, vulnerability management, security testing |
| Financial Regulations | ~35% | Decimal precision, double-entry linking, reconciliation, data retention, reporting |

---

## 10. Recommended Next Steps

1. **P0 Remediation Sprint**: Address F01-F06 immediately — these represent active security and financial integrity risks
2. **Decimal Migration**: Replace `float64` with `decimal.Decimal` across all financial models and services
3. **Authentication Hardening**: Implement account lockout, increase bcrypt cost, add CSRF protection
4. **Audit Log Immutability**: Remove CASCADE deletes, add append-only enforcement, implement archival
5. **Monitoring Setup**: Implement real-time alerting channels, configure health checks, set up log aggregation
6. **Testing Framework**: Add integration tests for auth, wallet ownership, transfer atomicity, and idempotency
7. **CI/CD Pipeline**: Set up automated build, test, lint, and security scanning
8. **Compliance Documentation**: Document security policies, incident response procedures, and data retention schedules

---

*Assessment completed: 2026-05-08 | Next review: After P0 remediation*
