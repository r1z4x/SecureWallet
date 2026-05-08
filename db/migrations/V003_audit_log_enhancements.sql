-- V003_audit_log_enhancements.sql
-- Add result and correlation_id columns to audit_logs table
-- Result: tracks success/failure/denied/unknown outcome
-- CorrelationID: links related audit entries across request lifecycle

ALTER TABLE audit_logs
    ADD COLUMN IF NOT EXISTS result VARCHAR(20) NOT NULL DEFAULT 'unknown' AFTER details,
    ADD COLUMN IF NOT EXISTS correlation_id CHAR(36) NULL AFTER result,
    ADD INDEX idx_audit_logs_user_id (user_id),
    ADD INDEX idx_audit_logs_action (action),
    ADD INDEX idx_audit_logs_resource (resource);

-- Backfill existing rows with 'unknown' result
UPDATE audit_logs SET result = 'unknown' WHERE result = '' OR result IS NULL;
