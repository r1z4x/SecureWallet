-- V004_uuid_fk_enforcement.sql
-- Forward migration: enforce foreign key constraints on tables created via GORM AutoMigrate,
-- fix transaction type CHECK constraint to match actual Go code values,
-- and record V003 in schema_migrations if not already present.
-- Idempotent: safe to run multiple times.
-- Rollback: see V004_uuid_fk_enforcement_rollback.sql

USE securewallet_dev;

-- ─────────────────────────────────────────────
-- 1. Record V003 if missing (it didn't self-record)
-- ─────────────────────────────────────────────
INSERT IGNORE INTO schema_migrations (version, name, checksum)
VALUES (3, 'V003_audit_log_enhancements', SHA2('V003_audit_log_enhancements.sql', 256));

-- ─────────────────────────────────────────────
-- 2. Fix transaction type CHECK constraint
--    Go code uses: deposit, withdrawal, transfer_in, transfer_out (lowercase with underscores)
--    The original V002 constraint already uses these correct values.
--    However, sample data and transfer.go used UPPERCASE values which would fail.
--    This step drops and recreates the constraint to ensure it matches.
-- ─────────────────────────────────────────────
SET @check_exists := (
    SELECT COUNT(*)
    FROM information_schema.CHECK_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND CONSTRAINT_NAME = 'chk_transaction_type'
);
SET @sql := IF(@check_exists > 0,
    'ALTER TABLE transactions DROP CHECK chk_transaction_type',
    'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

ALTER TABLE transactions
    ADD CONSTRAINT chk_transaction_type CHECK (type IN ('deposit', 'withdrawal', 'transfer_in', 'transfer_out'));

-- ─────────────────────────────────────────────
-- 3. Enforce foreign key constraints on all tables
--    GORM AutoMigrate does NOT create FK constraints.
--    The V002 SQL migration already creates some FKs with older names, so each
--    add is guarded by the actual table/column relationship instead of only the
--    desired constraint name.
-- ─────────────────────────────────────────────

-- wallets -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'wallets'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE wallets ADD CONSTRAINT fk_wallets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- transactions -> wallets
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'transactions'
      AND COLUMN_NAME = 'wallet_id' AND REFERENCED_TABLE_NAME = 'wallets'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE transactions ADD CONSTRAINT fk_transactions_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- transactions -> wallets (counterparty)
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'transactions'
      AND COLUMN_NAME = 'counterparty_wallet_id' AND REFERENCED_TABLE_NAME = 'wallets'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE transactions ADD CONSTRAINT fk_transactions_counterparty_wallet FOREIGN KEY (counterparty_wallet_id) REFERENCES wallets(id) ON DELETE SET NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- sessions -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sessions'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE sessions ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- audit_logs -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'audit_logs'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE audit_logs ADD CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- idempotency_records -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'idempotency_records'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE idempotency_records ADD CONSTRAINT fk_idempotency_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- support_tickets -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'support_tickets'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE support_tickets ADD CONSTRAINT fk_support_tickets_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- login_history -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'login_history'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE login_history ADD CONSTRAINT fk_login_history_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- blog_posts -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'blog_posts'
      AND COLUMN_NAME = 'author_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE blog_posts ADD CONSTRAINT fk_blog_posts_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- blog_comments -> blog_posts
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'blog_comments'
      AND COLUMN_NAME = 'post_id' AND REFERENCED_TABLE_NAME = 'blog_posts'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE blog_comments ADD CONSTRAINT fk_blog_comments_post FOREIGN KEY (post_id) REFERENCES blog_posts(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- security_alerts -> users
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'security_alerts'
      AND COLUMN_NAME = 'user_id' AND REFERENCED_TABLE_NAME = 'users'
      AND REFERENCED_COLUMN_NAME = 'id'
);
SET @sql := IF(@fk_exists = 0, 'ALTER TABLE security_alerts ADD CONSTRAINT fk_security_alerts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ─────────────────────────────────────────────
-- 4. Record this migration
-- ─────────────────────────────────────────────
INSERT IGNORE INTO schema_migrations (version, name, checksum)
VALUES (4, 'V004_uuid_fk_enforcement', SHA2('V004_uuid_fk_enforcement.sql', 256));
