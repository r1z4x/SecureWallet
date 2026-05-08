-- V002_wallets_and_transactions.sql
-- Forward migration: wallets, transactions, and idempotency support
-- Dependencies: V001 (users table with UUID primary key)
-- Idempotent: safe to run multiple times via IF NOT EXISTS guards
-- Rollback: see V002_wallets_and_transactions_rollback.sql

USE securewallet_dev;

-- ─────────────────────────────────────────────
-- 1. Wallets table
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wallets (
    id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    version INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Optimistic locking version',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    PRIMARY KEY (id),
    CONSTRAINT uq_wallet_user_currency UNIQUE (user_id, currency, deleted_at),
    CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_wallet_balance_nonnegative CHECK (balance >= 0),
    CONSTRAINT chk_wallet_currency_length CHECK (CHAR_LENGTH(currency) = 3),

    INDEX idx_wallet_user_id (user_id),
    INDEX idx_wallet_currency (currency),
    INDEX idx_wallet_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ─────────────────────────────────────────────
-- 2. Transactions table
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS transactions (
    id CHAR(36) NOT NULL,
    wallet_id CHAR(36) NOT NULL,
    counterparty_wallet_id CHAR(36) NULL COMMENT 'Destination wallet for transfers; NULL for deposits/withdrawals',
    type VARCHAR(20) NOT NULL COMMENT 'deposit, withdrawal, transfer_in, transfer_out',
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    description VARCHAR(255) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT 'pending, completed, failed, reversed',
    idempotency_key CHAR(36) NULL COMMENT 'Client-supplied key for deduplication',
    reference_id VARCHAR(100) NULL COMMENT 'External reference (e.g. bank transfer ID)',
    metadata JSON NULL COMMENT 'Sanitized audit metadata; never store secrets',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    PRIMARY KEY (id),
    CONSTRAINT uq_transaction_idempotency UNIQUE (idempotency_key),
    CONSTRAINT fk_transaction_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_counterparty_wallet FOREIGN KEY (counterparty_wallet_id) REFERENCES wallets(id) ON DELETE SET NULL,
    CONSTRAINT chk_transaction_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_transaction_currency_length CHECK (CHAR_LENGTH(currency) = 3),
    CONSTRAINT chk_transaction_type CHECK (type IN ('deposit', 'withdrawal', 'transfer_in', 'transfer_out')),
    CONSTRAINT chk_transaction_status CHECK (status IN ('pending', 'completed', 'failed', 'reversed')),

    INDEX idx_transaction_wallet_id (wallet_id),
    INDEX idx_transaction_type (type),
    INDEX idx_transaction_status (status),
    INDEX idx_transaction_created_at (created_at),
    INDEX idx_transaction_idempotency_key (idempotency_key),
    INDEX idx_transaction_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ─────────────────────────────────────────────
-- 3. Idempotency records table
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS idempotency_records (
    id CHAR(36) NOT NULL,
    key CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    operation VARCHAR(100) NOT NULL COMMENT 'transfer, deposit, withdrawal',
    payload_hash CHAR(64) NOT NULL COMMENT 'SHA-256 of request body',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT 'pending, completed, failed',
    http_status INT NOT NULL DEFAULT 200,
    response_body TEXT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    PRIMARY KEY (id),
    CONSTRAINT uq_idempotency_key UNIQUE (key),
    CONSTRAINT fk_idempotency_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_idempotency_status CHECK (status IN ('pending', 'completed', 'failed')),

    INDEX idx_idempotency_user_id (user_id),
    INDEX idx_idempotency_operation (operation),
    INDEX idx_idempotency_expires_at (expires_at),
    INDEX idx_idempotency_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ─────────────────────────────────────────────
-- 4. Migration tracking
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INT UNSIGNED NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    checksum CHAR(64) NULL COMMENT 'SHA-256 of migration file contents'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert migration record (idempotent via INSERT IGNORE)
INSERT IGNORE INTO schema_migrations (version, name, checksum)
VALUES (2, 'V002_wallets_and_transactions', SHA2('V002_wallets_and_transactions.sql', 256));
