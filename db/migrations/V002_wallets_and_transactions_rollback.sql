-- V002_wallets_and_transactions_rollback.sql
-- Rollback migration: drops tables created by V002
-- WARNING: This is destructive. All wallet, transaction, and idempotency data will be lost.
-- Only run in development or after confirming data has been migrated/backed up.
-- Order matters: drop dependent tables first (transactions before wallets).

USE securewallet_dev;

-- Remove migration tracking entry
DELETE FROM schema_migrations WHERE version = 2;

-- Drop idempotency records (depends on users)
DROP TABLE IF EXISTS idempotency_records;

-- Drop transactions (depends on wallets)
DROP TABLE IF EXISTS transactions;

-- Drop wallets (depends on users)
DROP TABLE IF EXISTS wallets;

-- Note: schema_migrations table is intentionally kept for future migrations.
-- If this is a full teardown, run: DROP TABLE IF EXISTS schema_migrations;
