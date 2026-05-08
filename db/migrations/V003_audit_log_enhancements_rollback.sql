-- V003_audit_log_enhancements_rollback.sql
-- Remove result and correlation_id columns from audit_logs table
-- Note: This drops data in these columns; only use if reverting the enhancement

ALTER TABLE audit_logs
    DROP COLUMN IF EXISTS result,
    DROP COLUMN IF EXISTS correlation_id;

-- Note: indexes are not dropped as they may be used by other queries
-- If needed, manually drop: idx_audit_logs_user_id, idx_audit_logs_action, idx_audit_logs_resource
