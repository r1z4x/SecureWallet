-- V005_iam_layer.sql
-- Migration: Identity and Access Management Layer
-- Adds RBAC (roles, permissions), OAuth2 support, and enhanced session management

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- User-Role junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id CHAR(36) NOT NULL,
    role_id CHAR(36) NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Role-Permission junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id CHAR(36) NOT NULL,
    permission_id CHAR(36) NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- OAuth Providers table
CREATE TABLE IF NOT EXISTS oauth_providers (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    auth_url VARCHAR(500) NOT NULL,
    token_url VARCHAR(500) NOT NULL,
    user_info_url VARCHAR(500),
    scopes VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- OAuth Accounts table
CREATE TABLE IF NOT EXISTS oauth_accounts (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    provider_id CHAR(36) NOT NULL,
    provider_name VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    access_token VARCHAR(1000),
    refresh_token VARCHAR(1000),
    token_expiry TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES oauth_providers(id) ON DELETE CASCADE,
    UNIQUE KEY unique_provider_user (provider_id, provider_user_id)
);

-- Enhanced sessions table (add columns if not exists)
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS device_name VARCHAR(100);
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45);
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS user_agent VARCHAR(500);
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS token_last_8 VARCHAR(8);

-- Add index for session lookups
CREATE INDEX IF NOT EXISTS idx_sessions_user_active ON sessions(user_id, is_active, expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);

-- Insert system permissions
INSERT IGNORE INTO permissions (id, name, description, resource, action) VALUES
    (UUID(), 'wallet:read', 'Read wallet information', 'wallet', 'read'),
    (UUID(), 'wallet:write', 'Modify wallet information', 'wallet', 'write'),
    (UUID(), 'wallet:delete', 'Delete wallets', 'wallet', 'delete'),
    (UUID(), 'transfer:read', 'Read transfer history', 'transfer', 'read'),
    (UUID(), 'transfer:write', 'Initiate transfers', 'transfer', 'write'),
    (UUID(), 'user:read', 'Read user information', 'user', 'read'),
    (UUID(), 'user:write', 'Modify user information', 'user', 'write'),
    (UUID(), 'user:delete', 'Delete users', 'user', 'delete'),
    (UUID(), 'admin:*', 'Full administrative access', 'admin', '*'),
    (UUID(), 'audit:read', 'Read audit logs', 'audit', 'read'),
    (UUID(), 'support:read', 'Read support tickets', 'support', 'read'),
    (UUID(), 'support:write', 'Manage support tickets', 'support', 'write');

-- Insert system roles
INSERT IGNORE INTO roles (id, name, description, is_system) VALUES
    (UUID(), 'admin', 'Full system administrator access', TRUE),
    (UUID(), 'user', 'Standard user access', TRUE),
    (UUID(), 'auditor', 'Read-only access to audit logs and reports', TRUE),
    (UUID(), 'support', 'Customer support access', TRUE);

-- Assign permissions to admin role (admin:*)
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'admin' AND p.name = 'admin:*';

-- Assign permissions to user role
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'user' AND p.name IN ('wallet:read', 'wallet:write', 'transfer:read', 'transfer:write');

-- Assign permissions to auditor role
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'auditor' AND p.name IN ('audit:read', 'wallet:read', 'transfer:read');

-- Assign permissions to support role
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'support' AND p.name IN ('support:read', 'support:write', 'user:read');

-- Assign admin role to existing admin users
INSERT IGNORE INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r
WHERE u.is_admin = TRUE AND r.name = 'admin';
