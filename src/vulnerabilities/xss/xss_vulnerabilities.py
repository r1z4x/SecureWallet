"""
XSS (Cross-Site Scripting) Vulnerabilities
This module contains intentionally vulnerable XSS code for testing purposes.
"""

import re
from typing import List, Dict, Any
from sqlalchemy.orm import Session
from src.models.user import User


class XSSVulnerabilities:
    """Class containing XSS vulnerabilities for testing"""
    
    def basic_reflected_xss(self, user_input: str) -> str:
        """
        Basic Reflected XSS - Direct output
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct output without sanitization
        return f"<h1>Search Results for: {user_input}</h1>"
    
    def medium_reflected_xss(self, user_input: str) -> str:
        """
        Medium Reflected XSS - Partial sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization that can be bypassed
        user_input = user_input.replace("<script>", "").replace("</script>", "")
        return f"<h1>Search Results for: {user_input}</h1>"
    
    def hard_reflected_xss(self, user_input: str) -> str:
        """
        Hard Reflected XSS - Complex bypass scenarios
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex bypass scenarios
        user_input = user_input.replace("javascript:", "").replace("onerror=", "")
        return f"<h1>Search Results for: {user_input}</h1>"
    
    def expert_reflected_xss(self, user_input: str) -> str:
        """
        Expert Reflected XSS - Advanced bypass techniques
        Vulnerability Level: Expert
        """
        # VULNERABLE: Advanced bypass techniques
        user_input = re.sub(r'<script[^>]*>.*?</script>', '', user_input, flags=re.IGNORECASE)
        return f"<h1>Search Results for: {user_input}</h1>"
    
    def basic_stored_xss(self, db: Session, user_id: int, comment: str) -> str:
        """
        Basic Stored XSS - Direct storage
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct storage without sanitization
        user = db.query(User).filter(User.id == user_id).first()
        if user:
            # In a real scenario, this would be stored in a comments table
            return f"<div class='comment'>User {user.username} says: {comment}</div>"
        return "User not found"
    
    def medium_stored_xss(self, db: Session, user_id: int, comment: str) -> str:
        """
        Medium Stored XSS - Partial sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization
        comment = comment.replace("<script>", "").replace("</script>", "")
        user = db.query(User).filter(User.id == user_id).first()
        if user:
            return f"<div class='comment'>User {user.username} says: {comment}</div>"
        return "User not found"
    
    def hard_stored_xss(self, db: Session, user_id: int, comment: str) -> str:
        """
        Hard Stored XSS - Complex sanitization bypass
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex sanitization bypass
        comment = re.sub(r'<script[^>]*>.*?</script>', '', comment, flags=re.IGNORECASE)
        user = db.query(User).filter(User.id == user_id).first()
        if user:
            return f"<div class='comment'>User {user.username} says: {comment}</div>"
        return "User not found"
    
    def expert_stored_xss(self, db: Session, user_id: int, comment: str) -> str:
        """
        Expert Stored XSS - Advanced techniques
        Vulnerability Level: Expert
        """
        # VULNERABLE: Advanced techniques
        comment = re.sub(r'<script[^>]*>.*?</script>', '', comment, flags=re.IGNORECASE)
        comment = comment.replace("javascript:", "").replace("onerror=", "")
        user = db.query(User).filter(User.id == user_id).first()
        if user:
            return f"<div class='comment'>User {user.username} says: {comment}</div>"
        return "User not found"
    
    def basic_dom_xss(self, user_input: str) -> str:
        """
        Basic DOM XSS - Direct DOM manipulation
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct DOM manipulation
        return f"""
        <script>
            document.getElementById('result').innerHTML = '{user_input}';
        </script>
        """
    
    def medium_dom_xss(self, user_input: str) -> str:
        """
        Medium DOM XSS - Partial sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization
        user_input = user_input.replace("<script>", "").replace("</script>", "")
        return f"""
        <script>
            document.getElementById('result').innerHTML = '{user_input}';
        </script>
        """
    
    def hard_dom_xss(self, user_input: str) -> str:
        """
        Hard DOM XSS - Complex bypass
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex bypass
        user_input = user_input.replace("javascript:", "").replace("onerror=", "")
        return f"""
        <script>
            document.getElementById('result').innerHTML = '{user_input}';
        </script>
        """
    
    def expert_dom_xss(self, user_input: str) -> str:
        """
        Expert DOM XSS - Advanced bypass
        Vulnerability Level: Expert
        """
        # VULNERABLE: Advanced bypass
        user_input = re.sub(r'<script[^>]*>.*?</script>', '', user_input, flags=re.IGNORECASE)
        return f"""
        <script>
            document.getElementById('result').innerHTML = '{user_input}';
        </script>
        """
    
    def basic_xss_in_json(self, user_input: str) -> str:
        """
        Basic XSS in JSON - Direct JSON output
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct JSON output
        data = {
            "message": user_input,
            "status": "success"
        }
        return str(data)
    
    def medium_xss_in_json(self, user_input: str) -> str:
        """
        Medium XSS in JSON - Partial sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization
        user_input = user_input.replace("<script>", "").replace("</script>", "")
        data = {
            "message": user_input,
            "status": "success"
        }
        return str(data)
    
    def hard_xss_in_json(self, user_input: str) -> str:
        """
        Hard XSS in JSON - Complex sanitization
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex sanitization
        user_input = re.sub(r'<script[^>]*>.*?</script>', '', user_input, flags=re.IGNORECASE)
        data = {
            "message": user_input,
            "status": "success"
        }
        return str(data)
    
    def expert_xss_in_json(self, user_input: str) -> str:
        """
        Expert XSS in JSON - Advanced techniques
        Vulnerability Level: Expert
        """
        # VULNERABLE: Advanced techniques
        user_input = re.sub(r'<script[^>]*>.*?</script>', '', user_input, flags=re.IGNORECASE)
        user_input = user_input.replace("javascript:", "").replace("onerror=", "")
        data = {
            "message": user_input,
            "status": "success"
        }
        return str(data)
