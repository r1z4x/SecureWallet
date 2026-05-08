-- V005_iam_layer_rollback.sql
-- Rollback: Identity and Access Management Layer

-- Remove user-role assignments for admin users
DELETE FROM user_roles WHERE role_id IN (SELECT id FROM roles WHERE name = 'admin');

-- Remove role-permission assignments
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE is_system = TRUE);

-- Remove system roles
DELETE FROM roles WHERE is_system = TRUE;

-- Remove system permissions
DELETE FROM permissions WHERE resource IN ('wallet', 'transfer', 'user', 'admin', 'audit', 'support');

-- Drop junction tables
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;

-- Drop OAuth tables
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS oauth_providers;

-- Drop RBAC tables
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;

-- Remove enhanced session columns (MySQL doesn't support DROP COLUMN IF EXISTS)
ALTER TABLE sessions DROP COLUMN IF EXISTS device_name;
ALTER TABLE sessions DROP COLUMN IF EXISTS ip_address;
ALTER TABLE sessions DROP COLUMN IF EXISTS user_agent;
ALTER TABLE sessions DROP COLUMN IF EXISTS is_active;
ALTER TABLE sessions DROP COLUMN IF EXISTS last_accessed;
ALTER TABLE sessions DROP COLUMN IF EXISTS token_last_8;

-- Remove session indexes
DROP INDEX IF EXISTS idx_sessions_user_active ON sessions;
DROP INDEX IF EXISTS idx_sessions_token ON sessions;
