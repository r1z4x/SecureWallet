#!/bin/bash

# SecureWallet Database Fix Script
# This script helps fix database schema compatibility issues

echo "🔧 SecureWallet Database Fix Script"
echo "=================================="
echo ""

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: Please run this script from the SecureWallet project root directory"
    exit 1
fi

echo "📋 Available options:"
echo "1. Quick fix - Reset database with sample data"
echo "2. Force recreation - Completely drop and recreate database"
echo "3. Manual cleanup - Drop all tables manually"
echo "4. Check current schema"
echo ""

read -p "Choose an option (1-4): " choice

case $choice in
    1)
        echo "🔄 Option 1: Quick database reset"
        echo "This will clear all data and recreate tables with the correct schema."
        read -p "Are you sure? This will delete all existing data! (y/N): " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            export RESET_DATABASE_ON_STARTUP=true
            echo "✅ Environment variable set. Now restart your application:"
            echo "   ./start-go-dev.sh"
        else
            echo "❌ Operation cancelled"
        fi
        ;;
    2)
        echo "💥 Option 2: Force database recreation"
        echo "This will completely drop and recreate the database."
        read -p "Are you sure? This will delete the entire database! (y/N): " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            export FORCE_DATABASE_RECREATION=true
            export RESET_DATABASE_ON_STARTUP=true
            echo "✅ Environment variables set. Now restart your application:"
            echo "   ./start-go-dev.sh"
        else
            echo "❌ Operation cancelled"
        fi
        ;;
    3)
        echo "🗑️  Option 3: Manual table cleanup"
        echo "This will drop all tables manually. You'll need to restart the app after."
        read -p "Are you sure? This will delete all tables! (y/N): " confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            echo "📝 Run these SQL commands in your MySQL client:"
            echo ""
            echo "USE securewallet_dev;"
            echo "DROP TABLE IF EXISTS transactions;"
            echo "DROP TABLE IF EXISTS login_history;"
            echo "DROP TABLE IF EXISTS support_tickets;"
            echo "DROP TABLE IF EXISTS audit_logs;"
            echo "DROP TABLE IF EXISTS sessions;"
            echo "DROP TABLE IF EXISTS wallets;"
            echo "DROP TABLE IF EXISTS users;"
            echo ""
            echo "After running these commands, restart your application."
        else
            echo "❌ Operation cancelled"
        fi
        ;;
    4)
        echo "🔍 Option 4: Check current schema"
        echo "This will show the current database schema status."
        if command -v mysql &> /dev/null; then
            echo "📊 Current schema status:"
            mysql -u"$DB_USER" -p"$DB_PASSWORD" -h"$DB_HOST" -P"$DB_PORT" securewallet_dev -e "
                SELECT 
                    TABLE_NAME,
                    COLUMN_NAME,
                    COLUMN_TYPE
                FROM INFORMATION_SCHEMA.COLUMNS 
                WHERE TABLE_SCHEMA = 'securewallet_dev' 
                    AND COLUMN_NAME IN ('id', 'user_id', 'wallet_id')
                ORDER BY TABLE_NAME, COLUMN_NAME;
            " 2>/dev/null || echo "❌ Could not connect to database. Check your environment variables."
        else
            echo "❌ MySQL client not found. Please check your schema manually."
        fi
        ;;
    *)
        echo "❌ Invalid option. Please choose 1-4."
        exit 1
        ;;
esac

echo ""
echo "📚 For more help, see docs/TROUBLESHOOTING.md"
echo "🌐 API endpoint for force recreation: POST /api/data/force-recreate"
