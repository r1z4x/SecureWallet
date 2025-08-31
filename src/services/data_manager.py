from datetime import datetime
from typing import List, Dict, Any, Optional
from sqlalchemy.orm import Session
from sqlalchemy import text
from src.models.user import User
from src.models.wallet import Wallet
from src.models.transaction import Transaction
from src.models.support_ticket import SupportTicket
from src.services.auth import get_password_hash
import json
import os

class DataManager:
    """Data management service for SecureWallet"""
    
    def __init__(self, db: Session):
        self.db = db
        self.data_dir = "data"
        self.ensure_data_directory()
    
    def ensure_data_directory(self):
        """Ensure data directory exists"""
        if not os.path.exists(self.data_dir):
            os.makedirs(self.data_dir)
    
    def create_snapshot(self, version: str, description: str = "") -> Dict[str, Any]:
        """Create a data snapshot"""
        snapshot = {
            "version": version,
            "description": description,
            "created_at": datetime.utcnow().isoformat(),
            "data": {
                "users": self._get_users_data(),
                "wallets": self._get_wallets_data(),
                "transactions": self._get_transactions_data(),
                "support_tickets": self._get_support_tickets_data()
            }
        }
        
        # Save snapshot to file
        filename = f"{self.data_dir}/snapshot_{version}.json"
        with open(filename, 'w') as f:
            json.dump(snapshot, f, indent=2)
        
        return snapshot
    
    def load_snapshot(self, version: str) -> Optional[Dict[str, Any]]:
        """Load a data snapshot"""
        filename = f"{self.data_dir}/snapshot_{version}.json"
        if not os.path.exists(filename):
            return None
        
        with open(filename, 'r') as f:
            return json.load(f)
    
    def list_snapshots(self) -> List[Dict[str, Any]]:
        """List all available snapshots"""
        snapshots = []
        for filename in os.listdir(self.data_dir):
            if filename.startswith("snapshot_") and filename.endswith(".json"):
                version = filename.replace("snapshot_", "").replace(".json", "")
                snapshot = self.load_snapshot(version)
                if snapshot:
                    snapshots.append({
                        "version": version,
                        "description": snapshot.get("description", ""),
                        "created_at": snapshot.get("created_at", ""),
                        "user_count": len(snapshot.get("data", {}).get("users", [])),
                        "wallet_count": len(snapshot.get("data", {}).get("wallets", [])),
                        "transaction_count": len(snapshot.get("data", {}).get("transactions", []))
                    })
        
        return sorted(snapshots, key=lambda x: x["created_at"], reverse=True)
    
    def restore_snapshot(self, version: str, clear_existing: bool = True) -> bool:
        """Restore data from a snapshot"""
        snapshot = self.load_snapshot(version)
        if not snapshot:
            return False
        
        try:
            if clear_existing:
                self._clear_all_data()
            
            # Restore users
            for user_data in snapshot["data"]["users"]:
                self._create_user_from_data(user_data)
            
            # Restore wallets
            for wallet_data in snapshot["data"]["wallets"]:
                self._create_wallet_from_data(wallet_data)
            
            # Restore transactions
            for transaction_data in snapshot["data"]["transactions"]:
                self._create_transaction_from_data(transaction_data)
            
            # Restore support tickets
            for ticket_data in snapshot["data"]["support_tickets"]:
                self._create_support_ticket_from_data(ticket_data)
            
            self.db.commit()
            return True
            
        except Exception as e:
            self.db.rollback()
            print(f"Error restoring snapshot: {e}")
            return False
    
    def create_demo_data(self) -> Dict[str, Any]:
        """Create demo data with sample users and transactions"""
        # Clear existing data
        self._clear_all_data()
        
        # Create demo users
        demo_users = [
            {
                "username": "admin",
                "email": "admin@securewallet.com",
                "password": "admin123",
                "is_admin": True,
                "is_active": True
            },
            {
                "username": "john",
                "email": "john@example.com",
                "password": "password123",
                "is_admin": False,
                "is_active": True
            },
            {
                "username": "jane",
                "email": "jane@example.com",
                "password": "password123",
                "is_admin": False,
                "is_active": True
            },
            {
                "username": "bob",
                "email": "bob@example.com",
                "password": "password123",
                "is_admin": False,
                "is_active": True
            },
            {
                "username": "demo",
                "email": "demo@securewallet.com",
                "password": "demo123",
                "is_admin": False,
                "is_active": True
            }
        ]
        
        created_users = []
        for user_data in demo_users:
            user = self._create_user_from_data(user_data)
            created_users.append(user)
        
        # Create wallets for users
        demo_wallets = [
            {"user_id": 1, "wallet_name": "Admin Wallet", "balance": 10000.00, "currency": "USD"},
            {"user_id": 2, "wallet_name": "John's Wallet", "balance": 5000.00, "currency": "USD"},
            {"user_id": 3, "wallet_name": "Jane's Wallet", "balance": 7500.00, "currency": "USD"},
            {"user_id": 4, "wallet_name": "Bob's Wallet", "balance": 3000.00, "currency": "USD"},
            {"user_id": 5, "wallet_name": "Demo Wallet", "balance": 1000.00, "currency": "USD"}
        ]
        
        created_wallets = []
        for wallet_data in demo_wallets:
            wallet = self._create_wallet_from_data(wallet_data)
            created_wallets.append(wallet)
        
        # Create sample transactions
        demo_transactions = [
            {"from_wallet_id": 1, "to_wallet_id": 2, "amount": 1000.00, "transaction_type": "TRANSFER", "status": "COMPLETED", "description": "Welcome bonus"},
            {"from_wallet_id": 1, "to_wallet_id": 3, "amount": 1000.00, "transaction_type": "TRANSFER", "status": "COMPLETED", "description": "Welcome bonus"},
            {"from_wallet_id": 1, "to_wallet_id": 4, "amount": 1000.00, "transaction_type": "TRANSFER", "status": "COMPLETED", "description": "Welcome bonus"},
            {"from_wallet_id": 2, "to_wallet_id": 3, "amount": 500.00, "transaction_type": "TRANSFER", "status": "COMPLETED", "description": "Payment for services"},
            {"from_wallet_id": 3, "to_wallet_id": 4, "amount": 250.00, "transaction_type": "TRANSFER", "status": "COMPLETED", "description": "Shared expenses"},
            {"from_wallet_id": None, "to_wallet_id": 5, "amount": 1000.00, "transaction_type": "DEPOSIT", "status": "COMPLETED", "description": "Initial deposit"}
        ]
        
        created_transactions = []
        for transaction_data in demo_transactions:
            transaction = self._create_transaction_from_data(transaction_data)
            created_transactions.append(transaction)
        
        # Create sample support tickets
        demo_tickets = [
            {"user_id": 2, "subject": "Account verification", "message": "I need help with account verification", "status": "open", "priority": "medium"},
            {"user_id": 3, "subject": "Transaction issue", "message": "My transaction is pending for too long", "status": "in_progress", "priority": "high"},
            {"user_id": 4, "subject": "Password reset", "message": "I forgot my password", "status": "closed", "priority": "low"}
        ]
        
        created_tickets = []
        for ticket_data in demo_tickets:
            ticket = self._create_support_ticket_from_data(ticket_data)
            created_tickets.append(ticket)
        
        self.db.commit()
        
        return {
            "users_created": len(created_users),
            "wallets_created": len(created_wallets),
            "transactions_created": len(created_transactions),
            "tickets_created": len(created_tickets),
            "credentials": {
                "admin": "admin / admin123",
                "john": "john / password123",
                "jane": "jane / password123",
                "bob": "bob / password123",
                "demo": "demo / demo123"
            }
        }
    
    def _clear_all_data(self):
        """Clear all data from database"""
        self.db.execute(text("DELETE FROM support_tickets"))
        self.db.execute(text("DELETE FROM transactions"))
        self.db.execute(text("DELETE FROM wallets"))
        self.db.execute(text("DELETE FROM audit_logs"))
        self.db.execute(text("DELETE FROM sessions"))
        self.db.execute(text("DELETE FROM users"))
        # Reset auto increment
        self.db.execute(text("ALTER TABLE users AUTO_INCREMENT = 1"))
        self.db.execute(text("ALTER TABLE wallets AUTO_INCREMENT = 1"))
        self.db.execute(text("ALTER TABLE transactions AUTO_INCREMENT = 1"))
        self.db.commit()
    
    def _get_users_data(self) -> List[Dict[str, Any]]:
        """Get users data for snapshot"""
        users = self.db.query(User).all()
        return [
            {
                "username": user.username,
                "email": user.email,
                "password_hash": user.password_hash,
                "is_admin": user.is_admin,
                "is_active": user.is_active,
                "created_at": user.created_at.isoformat() if user.created_at else None
            }
            for user in users
        ]
    
    def _get_wallets_data(self) -> List[Dict[str, Any]]:
        """Get wallets data for snapshot"""
        wallets = self.db.query(Wallet).all()
        return [
            {
                "user_id": wallet.user_id,
                "wallet_name": wallet.wallet_name,
                "balance": float(wallet.balance),
                "currency": wallet.currency,
                "is_active": wallet.is_active,
                "created_at": wallet.created_at.isoformat() if wallet.created_at else None
            }
            for wallet in wallets
        ]
    
    def _get_transactions_data(self) -> List[Dict[str, Any]]:
        """Get transactions data for snapshot"""
        transactions = self.db.query(Transaction).all()
        return [
            {
                "from_wallet_id": transaction.from_wallet_id,
                "to_wallet_id": transaction.to_wallet_id,
                "amount": float(transaction.amount),
                "transaction_type": transaction.transaction_type,
                "status": transaction.status,
                "description": transaction.description,
                "created_at": transaction.created_at.isoformat() if transaction.created_at else None
            }
            for transaction in transactions
        ]
    
    def _get_support_tickets_data(self) -> List[Dict[str, Any]]:
        """Get support tickets data for snapshot"""
        tickets = self.db.query(SupportTicket).all()
        return [
            {
                "user_id": ticket.user_id,
                "subject": ticket.subject,
                "message": ticket.message,
                "status": ticket.status,
                "priority": ticket.priority,
                "created_at": ticket.created_at.isoformat() if ticket.created_at else None
            }
            for ticket in tickets
        ]
    
    def _create_user_from_data(self, user_data: Dict[str, Any]) -> User:
        """Create user from data"""
        user = User(
            username=user_data["username"],
            email=user_data["email"],
            password_hash=user_data.get("password_hash") or get_password_hash(user_data["password"]),
            is_admin=user_data.get("is_admin", False),
            is_active=user_data.get("is_active", True)
        )
        self.db.add(user)
        self.db.flush()  # Get the ID
        return user
    
    def _create_wallet_from_data(self, wallet_data: Dict[str, Any]) -> Wallet:
        """Create wallet from data"""
        wallet = Wallet(
            user_id=wallet_data["user_id"],
            wallet_name=wallet_data["wallet_name"],
            balance=wallet_data["balance"],
            currency=wallet_data["currency"],
            is_active=wallet_data.get("is_active", True)
        )
        self.db.add(wallet)
        self.db.flush()  # Get the ID
        return wallet
    
    def _create_transaction_from_data(self, transaction_data: Dict[str, Any]) -> Transaction:
        """Create transaction from data"""
        transaction = Transaction(
            from_wallet_id=transaction_data.get("from_wallet_id"),
            to_wallet_id=transaction_data["to_wallet_id"],
            amount=transaction_data["amount"],
            transaction_type=transaction_data["transaction_type"],
            status=transaction_data["status"],
            description=transaction_data.get("description", "")
        )
        self.db.add(transaction)
        self.db.flush()  # Get the ID
        return transaction
    
    def _create_support_ticket_from_data(self, ticket_data: Dict[str, Any]) -> SupportTicket:
        """Create support ticket from data"""
        ticket = SupportTicket(
            user_id=ticket_data["user_id"],
            subject=ticket_data["subject"],
            message=ticket_data["message"],
            status=ticket_data["status"],
            priority=ticket_data["priority"]
        )
        self.db.add(ticket)
        self.db.flush()  # Get the ID
        return ticket
    
    def setup_fresh_database(self) -> Dict[str, Any]:
        """Setup fresh database with admin user and demo data"""
        # Clear all data first
        self._clear_all_data()
        
        # Create demo data
        result = self.create_demo_data()
        
        return {
            "message": "Fresh database setup completed",
            "data": result
        }
    
    def create_admin_only(self) -> Dict[str, Any]:
        """Create only admin user"""
        # Check if admin already exists
        existing_admin = self.db.query(User).filter(User.username == "admin").first()
        if existing_admin:
            return {
                "message": "Admin user already exists",
                "admin_credentials": "admin / admin123"
            }
        
        # Create admin user
        admin_user = self._create_user_from_data({
            "username": "admin",
            "email": "admin@securewallet.com",
            "password": "admin123",
            "is_admin": True,
            "is_active": True
        })
        
        self.db.commit()
        
        return {
            "message": "Admin user created successfully",
            "admin_credentials": "admin / admin123"
        }
    
    def get_all_credentials(self) -> Dict[str, Any]:
        """Get all user credentials"""
        users = self.db.query(User).all()
        credentials = {}
        
        for user in users:
            if user.username == "admin":
                credentials["admin"] = f"{user.username} / admin123"
            elif user.username in ["john", "jane", "bob"]:
                credentials[user.username] = f"{user.username} / password123"
            elif user.username == "demo":
                credentials[user.username] = f"{user.username} / demo123"
        
        return {
            "credentials": credentials,
            "total_users": len(users)
        }
    
    def reset_to_demo(self) -> Dict[str, Any]:
        """Reset database to demo state"""
        return self.setup_fresh_database()
    
    def make_user_admin(self, username: str) -> Dict[str, Any]:
        """Make a user admin"""
        user = self.db.query(User).filter(User.username == username).first()
        if not user:
            return {
                "success": False,
                "message": f"User {username} not found"
            }
        
        user.is_admin = True
        self.db.commit()
        
        return {
            "success": True,
            "message": f"User {username} is now admin"
        }
