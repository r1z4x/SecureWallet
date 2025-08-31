#!/usr/bin/env python3

import sys
import os
from datetime import datetime

# Add the src directory to Python path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'src'))

from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from src.config.settings import settings
from src.services.data_manager import DataManager
from src.models.audit_log import AuditLog
from src.models.session import Session
from src.models.user import User

class DataManagerCLI:
    def __init__(self):
        self.db_url = settings.database_url
        self.engine = create_engine(self.db_url)
        self.SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=self.engine)
        self.db = self.SessionLocal()
        self.data_manager = DataManager(self.db)
    
    def __del__(self):
        if hasattr(self, 'db'):
            self.db.close()
    
    def create_demo_data(self):
        """Create demo data"""
        try:
            result = self.data_manager.create_demo_data()
            print("✅ Demo data created successfully!")
            print(f"📊 Users created: {result['users_created']}")
            print(f"💰 Wallets created: {result['wallets_created']}")
            print(f"💳 Transactions created: {result['transactions_created']}")
            print(f"🎫 Support tickets created: {result['tickets_created']}")
            print("\n🔑 Login Credentials:")
            for user, creds in result['credentials'].items():
                print(f"  {user}: {creds}")
            return True
        except Exception as e:
            print(f"❌ Error creating demo data: {e}")
            return False
    
    def list_snapshots(self):
        """List all snapshots"""
        try:
            snapshots = self.data_manager.list_snapshots()
            
            if not snapshots:
                print("📋 No snapshots found")
                return True
            
            print("📋 Available Snapshots:")
            for snapshot in snapshots:
                print(f"  📸 Version: {snapshot['version']}")
                print(f"     Description: {snapshot['description']}")
                print(f"     Created: {snapshot['created_at']}")
                print(f"     Users: {snapshot['user_count']}")
                print(f"     Wallets: {snapshot['wallet_count']}")
                print(f"     Transactions: {snapshot['transaction_count']}")
                print()
            return True
        except Exception as e:
            print(f"❌ Error listing snapshots: {e}")
            return False
    
    def create_snapshot(self, version, description=""):
        """Create a snapshot"""
        try:
            snapshot = self.data_manager.create_snapshot(version, description)
            print(f"✅ Snapshot '{version}' created successfully!")
            print(f"📸 Description: {description}")
            print(f"🕒 Created: {snapshot['created_at']}")
            return True
        except Exception as e:
            print(f"❌ Error creating snapshot: {e}")
            return False
    
    def restore_snapshot(self, version, clear_existing=True):
        """Restore a snapshot"""
        try:
            success = self.data_manager.restore_snapshot(version, clear_existing)
            if success:
                print(f"✅ Snapshot '{version}' restored successfully!")
                return True
            else:
                print(f"❌ Failed to restore snapshot '{version}'")
                return False
        except Exception as e:
            print(f"❌ Error restoring snapshot: {e}")
            return False
    
    def clear_data(self):
        """Clear all data"""
        try:
            self.data_manager._clear_all_data()
            print("✅ All data cleared successfully!")
            return True
        except Exception as e:
            print(f"❌ Error clearing data: {e}")
            return False
    
    def get_credentials(self):
        """Get demo credentials"""
        try:
            result = self.data_manager.get_all_credentials()
            print("🔑 Demo Credentials:")
            for user, creds in result["credentials"].items():
                print(f"  {user}: {creds}")
            print(f"📊 Total users: {result['total_users']}")
            return True
        except Exception as e:
            print(f"❌ Error getting credentials: {e}")
            return False
    
    def create_admin(self):
        """Create admin user"""
        try:
            result = self.data_manager.create_admin_only()
            print("✅ Admin user created successfully!")
            print(f"🔑 {result['admin_credentials']}")
            return True
        except Exception as e:
            print(f"❌ Error creating admin: {e}")
            return False
    
    def setup_fresh(self):
        """Setup fresh database"""
        try:
            result = self.data_manager.setup_fresh_database()
            print("✅ Fresh database setup completed!")
            print(f"📊 Users created: {result['data']['users_created']}")
            print(f"💰 Wallets created: {result['data']['wallets_created']}")
            print(f"💳 Transactions created: {result['data']['transactions_created']}")
            print(f"🎫 Support tickets created: {result['data']['tickets_created']}")
            print("\n🔑 Login Credentials:")
            for user, creds in result['data']['credentials'].items():
                print(f"  {user}: {creds}")
            return True
        except Exception as e:
            print(f"❌ Error setting up fresh database: {e}")
            return False
    
    def reset_demo(self):
        """Reset to demo state"""
        try:
            result = self.data_manager.reset_to_demo()
            print("✅ Database reset to demo state!")
            print(f"📊 Users created: {result['data']['users_created']}")
            print(f"💰 Wallets created: {result['data']['wallets_created']}")
            print(f"💳 Transactions created: {result['data']['transactions_created']}")
            print(f"🎫 Support tickets created: {result['data']['tickets_created']}")
            print("\n🔑 Login Credentials:")
            for user, creds in result['data']['credentials'].items():
                print(f"  {user}: {creds}")
            return True
        except Exception as e:
            print(f"❌ Error resetting demo: {e}")
            return False
    
    def make_admin(self, username):
        """Make user admin"""
        try:
            result = self.data_manager.make_user_admin(username)
            if result['success']:
                print(f"✅ {result['message']}")
                return True
            else:
                print(f"❌ {result['message']}")
                return False
        except Exception as e:
            print(f"❌ Error making admin: {e}")
            return False
    
    def show_database_info(self):
        """Show database connection info"""
        print("🗄️ Database Information:")
        print(f"  URL: {self.db_url}")
        print(f"  Engine: {self.engine}")
        print(f"  Session: {self.db}")
        print()
    
    def show_login_activity(self, limit=20):
        """Show recent login activity"""
        try:
            # Get recent login audit logs
            login_logs = self.db.query(AuditLog).filter(
                AuditLog.action == "login"
            ).order_by(AuditLog.created_at.desc()).limit(limit).all()
            
            if not login_logs:
                print("📋 No login activity found")
                return True
            
            print(f"📋 Recent Login Activity (Last {len(login_logs)} entries):")
            print("-" * 80)
            
            for log in login_logs:
                user = self.db.query(User).filter(User.id == log.user_id).first()
                username = user.username if user else "Unknown"
                
                print(f"👤 User: {username}")
                print(f"   📅 Time: {log.created_at}")
                print(f"   🌐 IP: {log.ip_address or 'N/A'}")
                print(f"   🔍 Details: {log.details or 'N/A'}")
                print()
            
            return True
        except Exception as e:
            print(f"❌ Error showing login activity: {e}")
            return False
    
    def show_active_sessions(self):
        """Show currently active sessions"""
        try:
            from datetime import datetime, timezone
            
            # Get active sessions (not expired)
            active_sessions = self.db.query(Session).filter(
                Session.expires_at > datetime.now(timezone.utc)
            ).order_by(Session.created_at.desc()).all()
            
            if not active_sessions:
                print("📋 No active sessions found")
                return True
            
            print(f"📋 Active Sessions ({len(active_sessions)} sessions):")
            print("-" * 80)
            
            for session in active_sessions:
                user = self.db.query(User).filter(User.id == session.user_id).first()
                username = user.username if user else "Unknown"
                
                print(f"👤 User: {username}")
                print(f"   📅 Created: {session.created_at}")
                print(f"   ⏰ Expires: {session.expires_at}")
                print(f"   🔑 Token: {session.session_token[:20]}...")
                print()
            
            return True
        except Exception as e:
            print(f"❌ Error showing active sessions: {e}")
            return False
    
    def show_user_activity(self, username=None, limit=10):
        """Show user activity (login, actions, etc.)"""
        try:
            query = self.db.query(AuditLog)
            
            if username:
                user = self.db.query(User).filter(User.username == username).first()
                if not user:
                    print(f"❌ User '{username}' not found")
                    return False
                query = query.filter(AuditLog.user_id == user.id)
            
            activities = query.order_by(AuditLog.created_at.desc()).limit(limit).all()
            
            if not activities:
                if username:
                    print(f"📋 No activity found for user '{username}'")
                else:
                    print("📋 No user activity found")
                return True
            
            print(f"📋 User Activity ({'for ' + username if username else 'all users'}):")
            print("-" * 80)
            
            for activity in activities:
                user = self.db.query(User).filter(User.id == activity.user_id).first()
                username = user.username if user else "Unknown"
                
                print(f"👤 User: {username}")
                print(f"   🔧 Action: {activity.action}")
                print(f"   📅 Time: {activity.created_at}")
                print(f"   🌐 IP: {activity.ip_address or 'N/A'}")
                if activity.resource_type:
                    print(f"   📁 Resource: {activity.resource_type} (ID: {activity.resource_id})")
                if activity.details:
                    print(f"   📝 Details: {activity.details}")
                print()
            
            return True
        except Exception as e:
            print(f"❌ Error showing user activity: {e}")
            return False
    
    def clear_expired_sessions(self):
        """Clear expired sessions"""
        try:
            from datetime import datetime, timezone
            
            # Delete expired sessions
            expired_count = self.db.query(Session).filter(
                Session.expires_at <= datetime.now(timezone.utc)
            ).delete()
            
            self.db.commit()
            print(f"✅ Cleared {expired_count} expired sessions")
            return True
        except Exception as e:
            print(f"❌ Error clearing expired sessions: {e}")
            return False
    
    def show_login_stats(self):
        """Show login statistics"""
        try:
            # Get total login count
            total_logins = self.db.query(AuditLog).filter(
                AuditLog.action == "login"
            ).count()
            
            # Get unique users who logged in
            unique_users = self.db.query(AuditLog.user_id).filter(
                AuditLog.action == "login"
            ).distinct().count()
            
            # Get today's logins
            from datetime import datetime, timezone, timedelta
            today = datetime.now(timezone.utc).date()
            today_logins = self.db.query(AuditLog).filter(
                AuditLog.action == "login",
                AuditLog.created_at >= today
            ).count()
            
            # Get this week's logins
            week_ago = datetime.now(timezone.utc) - timedelta(days=7)
            week_logins = self.db.query(AuditLog).filter(
                AuditLog.action == "login",
                AuditLog.created_at >= week_ago
            ).count()
            
            print("📊 Login Statistics:")
            print("-" * 40)
            print(f"   Total Logins: {total_logins}")
            print(f"   Unique Users: {unique_users}")
            print(f"   Today's Logins: {today_logins}")
            print(f"   This Week's Logins: {week_logins}")
            print()
            
            return True
        except Exception as e:
            print(f"❌ Error showing login stats: {e}")
            return False

def main():
    if len(sys.argv) < 2:
        print("🚀 SecureWallet Data Manager (Direct Database)")
        print("=" * 50)
        print("🗄️ Direct database operations - no API/HTTP calls")
        print()
        print("Usage:")
        print("  python manage-data.py create-admin")
        print("  python manage-data.py setup-fresh")
        print("  python manage-data.py demo-data")
        print("  python manage-data.py create-snapshot [version] [description]")
        print("  python manage-data.py list-snapshots")
        print("  python manage-data.py restore-snapshot [version]")
        print("  python manage-data.py clear-data")
        print("  python manage-data.py credentials")
        print("  python manage-data.py reset-demo")
        print("  python manage-data.py make-admin [username]")
        print("  python manage-data.py db-info")
        print("  python manage-data.py login-activity [limit]")
        print("  python manage-data.py active-sessions")
        print("  python manage-data.py user-activity [username] [limit]")
        print("  python manage-data.py clear-expired-sessions")
        print("  python manage-data.py login-stats")
        print()
        print("Examples:")
        print("  python manage-data.py create-admin")
        print("  python manage-data.py make-admin admin")
        print("  python manage-data.py setup-fresh")
        print("  python manage-data.py demo-data")
        print("  python manage-data.py create-snapshot v1.0 'Initial demo data'")
        print("  python manage-data.py list-snapshots")
        print("  python manage-data.py restore-snapshot v1.0")
        print("  python manage-data.py login-activity 10")
        print("  python manage-data.py user-activity admin 5")
        print()
        print("🔧 All operations are performed directly on the database")
        return
    
    try:
        cli = DataManagerCLI()
    except Exception as e:
        print(f"❌ Failed to connect to database: {e}")
        print("   Make sure the database is running and accessible")
        return
    
    command = sys.argv[1]
    
    if command == "create-admin":
        cli.create_admin()
    
    elif command == "setup-fresh":
        cli.setup_fresh()
    
    elif command == "demo-data":
        cli.create_demo_data()
    
    elif command == "create-snapshot":
        if len(sys.argv) < 3:
            print("❌ Usage: python manage-data.py create-snapshot [version] [description]")
            return
        
        version = sys.argv[2]
        description = sys.argv[3] if len(sys.argv) > 3 else ""
        cli.create_snapshot(version, description)
    
    elif command == "list-snapshots":
        cli.list_snapshots()
    
    elif command == "restore-snapshot":
        if len(sys.argv) < 3:
            print("❌ Usage: python manage-data.py restore-snapshot [version]")
            return
        
        version = sys.argv[2]
        cli.restore_snapshot(version)
    
    elif command == "clear-data":
        cli.clear_data()
    
    elif command == "credentials":
        cli.get_credentials()
    
    elif command == "reset-demo":
        cli.reset_demo()
    
    elif command == "make-admin":
        if len(sys.argv) < 3:
            print("❌ Usage: python manage-data.py make-admin [username]")
            return
        
        username = sys.argv[2]
        cli.make_admin(username)
    
    elif command == "db-info":
        cli.show_database_info()
    
    elif command == "login-activity":
        limit = int(sys.argv[2]) if len(sys.argv) > 2 else 20
        cli.show_login_activity(limit)
    
    elif command == "active-sessions":
        cli.show_active_sessions()
    
    elif command == "user-activity":
        username = sys.argv[2] if len(sys.argv) > 2 else None
        limit = int(sys.argv[3]) if len(sys.argv) > 3 else 10
        cli.show_user_activity(username, limit)
    
    elif command == "clear-expired-sessions":
        cli.clear_expired_sessions()
    
    elif command == "login-stats":
        cli.show_login_stats()
    
    else:
        print(f"❌ Unknown command: {command}")

if __name__ == "__main__":
    main()
