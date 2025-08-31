"""
Weak Authentication Vulnerabilities
This module contains intentionally weak authentication mechanisms for testing purposes.
"""

import hashlib
import base64
import json
from datetime import datetime, timedelta
from typing import Optional, Dict, Any
from jose import jwt
from src.config.settings import settings


class WeakAuthenticationVulnerabilities:
    """Class containing weak authentication vulnerabilities for testing"""
    
    def basic_weak_password_storage(self, password: str) -> str:
        """
        Basic Weak Password Storage - Plain text
        Vulnerability Level: Basic
        """
        # VULNERABLE: Plain text password storage
        return password
    
    def medium_weak_password_storage(self, password: str) -> str:
        """
        Medium Weak Password Storage - MD5 hash
        Vulnerability Level: Medium
        """
        # VULNERABLE: MD5 hash (cryptographically broken)
        return hashlib.md5(password.encode()).hexdigest()
    
    def hard_weak_password_storage(self, password: str) -> str:
        """
        Hard Weak Password Storage - SHA1 hash
        Vulnerability Level: Hard
        """
        # VULNERABLE: SHA1 hash (cryptographically broken)
        return hashlib.sha1(password.encode()).hexdigest()
    
    def expert_weak_password_storage(self, password: str) -> str:
        """
        Expert Weak Password Storage - Base64 encoding
        Vulnerability Level: Expert
        """
        # VULNERABLE: Base64 encoding (not hashing)
        return base64.b64encode(password.encode()).decode()
    
    def basic_weak_jwt_secret(self, payload: Dict[str, Any]) -> str:
        """
        Basic Weak JWT - Hardcoded weak secret
        Vulnerability Level: Basic
        """
        # VULNERABLE: Hardcoded weak secret
        secret = "secret123"
        return jwt.encode(payload, secret, algorithm="HS256")
    
    def medium_weak_jwt_secret(self, payload: Dict[str, Any]) -> str:
        """
        Medium Weak JWT - Predictable secret
        Vulnerability Level: Medium
        """
        # VULNERABLE: Predictable secret
        secret = "jwt_secret_key_2024"
        return jwt.encode(payload, secret, algorithm="HS256")
    
    def hard_weak_jwt_secret(self, payload: Dict[str, Any]) -> str:
        """
        Hard Weak JWT - Short secret
        Vulnerability Level: Hard
        """
        # VULNERABLE: Short secret (brute forceable)
        secret = "abc123"
        return jwt.encode(payload, secret, algorithm="HS256")
    
    def expert_weak_jwt_secret(self, payload: Dict[str, Any]) -> str:
        """
        Expert Weak JWT - Algorithm confusion
        Vulnerability Level: Expert
        """
        # VULNERABLE: Algorithm confusion attack
        secret = "public_key_for_rsa"
        return jwt.encode(payload, secret, algorithm="HS256")
    
    def basic_weak_session_management(self, user_id: int) -> str:
        """
        Basic Weak Session Management - Predictable session ID
        Vulnerability Level: Basic
        """
        # VULNERABLE: Predictable session ID
        session_id = f"session_{user_id}_{datetime.now().strftime('%Y%m%d')}"
        return session_id
    
    def medium_weak_session_management(self, user_id: int) -> str:
        """
        Medium Weak Session Management - Short session ID
        Vulnerability Level: Medium
        """
        # VULNERABLE: Short session ID
        import random
        session_id = f"s{user_id}{random.randint(100, 999)}"
        return session_id
    
    def hard_weak_session_management(self, user_id: int) -> str:
        """
        Hard Weak Session Management - Time-based session ID
        Vulnerability Level: Hard
        """
        # VULNERABLE: Time-based session ID
        timestamp = int(datetime.now().timestamp())
        session_id = f"{user_id}_{timestamp}"
        return session_id
    
    def expert_weak_session_management(self, user_id: int) -> str:
        """
        Expert Weak Session Management - Encoded session ID
        Vulnerability Level: Expert
        """
        # VULNERABLE: Encoded session ID with weak encoding
        session_data = {"user_id": user_id, "timestamp": datetime.now().isoformat()}
        session_id = base64.b64encode(json.dumps(session_data).encode()).decode()
        return session_id
    
    def basic_weak_password_validation(self, password: str) -> bool:
        """
        Basic Weak Password Validation - No validation
        Vulnerability Level: Basic
        """
        # VULNERABLE: No password validation
        return len(password) > 0
    
    def medium_weak_password_validation(self, password: str) -> bool:
        """
        Medium Weak Password Validation - Weak requirements
        Vulnerability Level: Medium
        """
        # VULNERABLE: Weak password requirements
        return len(password) >= 4
    
    def hard_weak_password_validation(self, password: str) -> bool:
        """
        Hard Weak Password Validation - Predictable patterns
        Vulnerability Level: Hard
        """
        # VULNERABLE: Predictable password patterns
        common_passwords = ["password", "123456", "admin", "qwerty", "letmein"]
        return password not in common_passwords and len(password) >= 6
    
    def expert_weak_password_validation(self, password: str) -> bool:
        """
        Expert Weak Password Validation - Complex but bypassable
        Vulnerability Level: Expert
        """
        # VULNERABLE: Complex but bypassable validation
        has_upper = any(c.isupper() for c in password)
        has_lower = any(c.islower() for c in password)
        has_digit = any(c.isdigit() for c in password)
        has_special = any(c in "!@#$%^&*" for c in password)
        
        return (len(password) >= 8 and has_upper and has_lower and 
                has_digit and has_special)
    
    def basic_weak_token_expiration(self, payload: Dict[str, Any]) -> str:
        """
        Basic Weak Token Expiration - No expiration
        Vulnerability Level: Basic
        """
        # VULNERABLE: No expiration
        return jwt.encode(payload, settings.jwt_secret_key, algorithm="HS256")
    
    def medium_weak_token_expiration(self, payload: Dict[str, Any]) -> str:
        """
        Medium Weak Token Expiration - Long expiration
        Vulnerability Level: Medium
        """
        # VULNERABLE: Long expiration time
        payload["exp"] = datetime.utcnow() + timedelta(days=365)
        return jwt.encode(payload, settings.jwt_secret_key, algorithm="HS256")
    
    def hard_weak_token_expiration(self, payload: Dict[str, Any]) -> str:
        """
        Hard Weak Token Expiration - Predictable expiration
        Vulnerability Level: Hard
        """
        # VULNERABLE: Predictable expiration
        payload["exp"] = datetime.utcnow() + timedelta(hours=24)
        return jwt.encode(payload, settings.jwt_secret_key, algorithm="HS256")
    
    def expert_weak_token_expiration(self, payload: Dict[str, Any]) -> str:
        """
        Expert Weak Token Expiration - Renewal without validation
        Vulnerability Level: Expert
        """
        # VULNERABLE: Renewal without proper validation
        payload["exp"] = datetime.utcnow() + timedelta(minutes=30)
        payload["renewable"] = True
        return jwt.encode(payload, settings.jwt_secret_key, algorithm="HS256")
