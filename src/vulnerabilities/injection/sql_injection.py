"""
SQL Injection Vulnerabilities
This module contains intentionally vulnerable SQL queries for testing purposes.
"""

from sqlalchemy.orm import Session
from sqlalchemy import text
from typing import List, Optional
from src.models.user import User
from src.models.wallet import Wallet
from src.models.transaction import Transaction


class SQLInjectionVulnerabilities:
    """Class containing SQL injection vulnerabilities for testing"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def basic_sql_injection_user_search(self, username: str) -> List[User]:
        """
        Basic SQL Injection - Direct string concatenation
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct string concatenation
        query = f"SELECT * FROM users WHERE username = '{username}'"
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def medium_sql_injection_user_search(self, username: str) -> List[User]:
        """
        Medium SQL Injection - Partial sanitization bypass
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization that can be bypassed
        username = username.replace("'", "''")  # Weak sanitization
        query = f"SELECT * FROM users WHERE username = '{username}'"
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def hard_sql_injection_user_search(self, username: str) -> List[User]:
        """
        Hard SQL Injection - Complex injection scenarios
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex injection with multiple conditions
        query = f"""
        SELECT u.*, w.balance 
        FROM users u 
        LEFT JOIN wallets w ON u.id = w.user_id 
        WHERE u.username = '{username}' OR u.email LIKE '%{username}%'
        """
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def expert_sql_injection_user_search(self, username: str) -> List[User]:
        """
        Expert SQL Injection - Union-based injection
        Vulnerability Level: Expert
        """
        # VULNERABLE: Union-based injection
        query = f"""
        SELECT id, username, email, password_hash, is_admin 
        FROM users 
        WHERE username = '{username}'
        UNION ALL
        SELECT id, username, email, password_hash, is_admin 
        FROM users 
        WHERE is_admin = 1
        """
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def basic_sql_injection_wallet_search(self, user_id: str) -> List[Wallet]:
        """
        Basic SQL Injection - Wallet search
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct string concatenation
        query = f"SELECT * FROM wallets WHERE user_id = {user_id}"
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def medium_sql_injection_transaction_search(self, wallet_id: str, amount: str) -> List[Transaction]:
        """
        Medium SQL Injection - Transaction search with multiple parameters
        Vulnerability Level: Medium
        """
        # VULNERABLE: Multiple parameter injection
        query = f"""
        SELECT * FROM transactions 
        WHERE to_wallet_id = {wallet_id} 
        AND amount >= {amount}
        """
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def hard_sql_injection_admin_bypass(self, user_id: str) -> List[User]:
        """
        Hard SQL Injection - Admin privilege bypass
        Vulnerability Level: Hard
        """
        # VULNERABLE: Admin privilege bypass
        query = f"""
        SELECT * FROM users 
        WHERE id = {user_id} 
        OR is_admin = 1
        """
        result = self.db.execute(text(query))
        return result.fetchall()
    
    def expert_sql_injection_data_exfiltration(self, condition: str) -> List[dict]:
        """
        Expert SQL Injection - Data exfiltration
        Vulnerability Level: Expert
        """
        # VULNERABLE: Data exfiltration through complex injection
        query = f"""
        SELECT 
            u.username,
            u.email,
            u.password_hash,
            w.balance,
            t.amount,
            t.description
        FROM users u
        LEFT JOIN wallets w ON u.id = w.user_id
        LEFT JOIN transactions t ON w.id = t.to_wallet_id
        WHERE {condition}
        """
        result = self.db.execute(text(query))
        return [dict(row) for row in result.fetchall()]
