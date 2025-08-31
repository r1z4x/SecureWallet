"""
NoSQL Injection Vulnerabilities
This module contains intentionally vulnerable MongoDB queries for testing purposes.
"""

from pymongo.database import Database
from typing import List, Optional, Dict, Any
import json


class NoSQLInjectionVulnerabilities:
    """Class containing NoSQL injection vulnerabilities for testing"""
    
    def __init__(self, db: Database):
        self.db = db
    
    def basic_nosql_injection_user_search(self, username: str) -> List[Dict[str, Any]]:
        """
        Basic NoSQL Injection - Direct string injection
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct string injection
        query = {"username": username}
        return list(self.db.users.find(query))
    
    def medium_nosql_injection_user_search(self, username: str) -> List[Dict[str, Any]]:
        """
        Medium NoSQL Injection - JSON injection
        Vulnerability Level: Medium
        """
        # VULNERABLE: JSON injection
        try:
            query = json.loads(username)
        except json.JSONDecodeError:
            query = {"username": username}
        return list(self.db.users.find(query))
    
    def hard_nosql_injection_user_search(self, username: str) -> List[Dict[str, Any]]:
        """
        Hard NoSQL Injection - Operator injection
        Vulnerability Level: Hard
        """
        # VULNERABLE: Operator injection
        query = {"username": {"$ne": None}}
        if username:
            query["username"] = username
        return list(self.db.users.find(query))
    
    def expert_nosql_injection_user_search(self, username: str) -> List[Dict[str, Any]]:
        """
        Expert NoSQL Injection - Complex operator injection
        Vulnerability Level: Expert
        """
        # VULNERABLE: Complex operator injection
        query = {
            "$or": [
                {"username": username},
                {"is_admin": True}
            ]
        }
        return list(self.db.users.find(query))
    
    def basic_nosql_injection_wallet_search(self, user_id: str) -> List[Dict[str, Any]]:
        """
        Basic NoSQL Injection - Wallet search
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct injection
        query = {"user_id": int(user_id) if user_id.isdigit() else user_id}
        return list(self.db.wallets.find(query))
    
    def medium_nosql_injection_transaction_search(self, wallet_id: str, amount: str) -> List[Dict[str, Any]]:
        """
        Medium NoSQL Injection - Transaction search
        Vulnerability Level: Medium
        """
        # VULNERABLE: Multiple parameter injection
        query = {
            "to_wallet_id": int(wallet_id) if wallet_id.isdigit() else wallet_id,
            "amount": {"$gte": float(amount) if amount.replace('.', '').isdigit() else amount}
        }
        return list(self.db.transactions.find(query))
    
    def hard_nosql_injection_admin_bypass(self, user_id: str) -> List[Dict[str, Any]]:
        """
        Hard NoSQL Injection - Admin bypass
        Vulnerability Level: Hard
        """
        # VULNERABLE: Admin bypass through injection
        query = {
            "$or": [
                {"_id": int(user_id) if user_id.isdigit() else user_id},
                {"is_admin": True}
            ]
        }
        return list(self.db.users.find(query))
    
    def expert_nosql_injection_data_exfiltration(self, condition: str) -> List[Dict[str, Any]]:
        """
        Expert NoSQL Injection - Data exfiltration
        Vulnerability Level: Expert
        """
        # VULNERABLE: Data exfiltration through complex injection
        try:
            query = json.loads(condition)
        except json.JSONDecodeError:
            query = {}
        
        pipeline = [
            {"$match": query},
            {"$lookup": {
                "from": "wallets",
                "localField": "_id",
                "foreignField": "user_id",
                "as": "wallets"
            }},
            {"$lookup": {
                "from": "transactions",
                "localField": "wallets._id",
                "foreignField": "to_wallet_id",
                "as": "transactions"
            }}
        ]
        return list(self.db.users.aggregate(pipeline))
