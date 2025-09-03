# Troubleshooting Guide

## Database Schema Issues

### Common Error: Foreign Key Constraint Incompatibility

If you encounter this error:
```
Error 3780 (HY000): Referencing column 'user_id' and referenced column 'id' in foreign key constraint 'fk_users_wallets' are incompatible.
```

This indicates a schema mismatch between your existing database and the new UUID-based schema.

### Root Cause

The application was updated to use UUIDs (`CHAR(36)`) for all primary keys and foreign keys, but your existing database still uses the old schema with `BIGINT UNSIGNED` columns.

### Solutions

#### Option 1: Force Database Recreation (Recommended for Development)

This completely drops and recreates the database with the correct schema:

```bash
# Set environment variables
export FORCE_DATABASE_RECREATION=true
export RESET_DATABASE_ON_STARTUP=true

# Restart the application
./start-go-dev.sh
```

#### Option 2: Use the Force Recreation API Endpoint

If the application is already running, you can call the API endpoint:

```bash
curl -X POST http://localhost:8080/api/data/force-recreate
```

#### Option 3: Manual Database Reset

For a standard reset (may not work if schema is incompatible):

```bash
# Set environment variable
export RESET_DATABASE_ON_STARTUP=true

# Restart the application
./start-go-dev.sh
```

#### Option 4: Manual Database Cleanup

If the above options don't work, you can manually clean up:

```sql
-- Connect to MySQL and run these commands
USE securewallet_dev;

-- Drop all tables (in correct order to avoid foreign key issues)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS login_history;
DROP TABLE IF EXISTS support_tickets;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS users;

-- Then restart the application - it will recreate tables with correct schema
```

### Environment Variables

Add these to your `.env` file for easier management:

```bash
# Database reset options
RESET_DATABASE_ON_STARTUP=false
FORCE_DATABASE_RECREATION=false

# Set to 'true' when you need to reset the database
```

### Verification

After fixing the issue, verify the schema is correct:

```sql
-- Check that all ID columns are CHAR(36)
SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'securewallet_dev' 
    AND COLUMN_NAME IN ('id', 'user_id', 'wallet_id')
ORDER BY TABLE_NAME, COLUMN_NAME;
```

Expected output:
```
+----------------+-------------+-------------+
| TABLE_NAME     | COLUMN_NAME | COLUMN_TYPE |
+----------------+-------------+-------------+
| audit_logs     | id          | char(36)    |
| audit_logs     | user_id     | char(36)    |
| login_history  | id          | char(36)    |
| login_history  | user_id     | char(36)    |
| sessions       | id          | char(36)    |
| sessions       | user_id     | char(36)    |
| support_tickets| id          | char(36)    |
| support_tickets| user_id     | char(36)    |
| transactions   | id          | char(36)    |
| transactions   | wallet_id   | char(36)    |
| users          | id          | char(36)    |
| wallets        | id          | char(36)    |
| wallets        | user_id     | char(36)    |
+----------------+-------------+-------------+
```

### Prevention

To avoid this issue in the future:

1. **Always use environment variables** for database resets
2. **Test schema changes** in development first
3. **Use the force recreation option** when switching between major schema versions
4. **Backup your data** before major schema changes

### Getting Help

If you continue to have issues:

1. Check the application logs for detailed error messages
2. Verify your MySQL version is 8.0+
3. Ensure you have proper permissions to drop/create databases
4. Check that all environment variables are set correctly
