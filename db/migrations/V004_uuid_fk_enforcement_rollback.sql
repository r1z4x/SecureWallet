-- V004_uuid_fk_enforcement_rollback.sql
-- Rollback migration: removes FK constraints added by V004.
-- WARNING: This removes referential integrity enforcement at the DB level.
-- Only run in development or after confirming data has been migrated/backed up.

USE securewallet_dev;

-- Remove migration tracking entry
DELETE FROM schema_migrations WHERE version = 4;

-- Drop FK constraints added by V004
SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'wallets'
      AND CONSTRAINT_NAME = 'fk_wallets_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE wallets DROP FOREIGN KEY fk_wallets_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'transactions'
      AND CONSTRAINT_NAME = 'fk_transactions_wallet'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE transactions DROP FOREIGN KEY fk_transactions_wallet', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'transactions'
      AND CONSTRAINT_NAME = 'fk_transactions_counterparty_wallet'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE transactions DROP FOREIGN KEY fk_transactions_counterparty_wallet', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'sessions'
      AND CONSTRAINT_NAME = 'fk_sessions_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE sessions DROP FOREIGN KEY fk_sessions_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'audit_logs'
      AND CONSTRAINT_NAME = 'fk_audit_logs_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE audit_logs DROP FOREIGN KEY fk_audit_logs_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'idempotency_records'
      AND CONSTRAINT_NAME = 'fk_idempotency_records_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE idempotency_records DROP FOREIGN KEY fk_idempotency_records_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'support_tickets'
      AND CONSTRAINT_NAME = 'fk_support_tickets_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE support_tickets DROP FOREIGN KEY fk_support_tickets_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'login_history'
      AND CONSTRAINT_NAME = 'fk_login_history_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE login_history DROP FOREIGN KEY fk_login_history_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'blog_posts'
      AND CONSTRAINT_NAME = 'fk_blog_posts_author'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE blog_posts DROP FOREIGN KEY fk_blog_posts_author', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'blog_comments'
      AND CONSTRAINT_NAME = 'fk_blog_comments_post'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE blog_comments DROP FOREIGN KEY fk_blog_comments_post', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @fk_exists := (
    SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'security_alerts'
      AND CONSTRAINT_NAME = 'fk_security_alerts_user'
);
SET @sql := IF(@fk_exists > 0, 'ALTER TABLE security_alerts DROP FOREIGN KEY fk_security_alerts_user', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Restore original transaction type CHECK constraint
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

ALTER TABLE transactions ADD CONSTRAINT chk_transaction_type CHECK (type IN ('deposit', 'withdrawal', 'transfer_in', 'transfer_out'));

-- Note: V003 recording is not undone to avoid version tracking gaps.
