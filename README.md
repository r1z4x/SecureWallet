# SecureWallet - Digital Banking Platform (Vulnerable)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Python](https://img.shields.io/badge/Python-3.12+-blue.svg)](https://www.python.org/downloads/)
[![Node.js](https://img.shields.io/badge/Node.js-18+-green.svg)](https://nodejs.org/)
[![Go](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://www.docker.com/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.104+-green.svg)](https://fastapi.tiangolo.com/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.4+-green.svg)](https://vuejs.org/)
[![OWASP](https://img.shields.io/badge/OWASP-Top%2010-orange.svg)](https://owasp.org/www-project-top-ten/)
[![Vulnerable](https://img.shields.io/badge/Vulnerable-For%20Testing-red.svg)](https://owasp.org/)

A comprehensive vulnerable application designed for OWASP Top 10 training and educational purposes. This application is dedicated to **[OWASPAttackSimulator](https://github.com/r1z4x/OWASPAttackSimulator)** and has been specifically developed to work seamlessly with this security testing platform.

## Purpose

This application is designed to work stably with [OWASPAttackSimulator](https://github.com/r1z4x/OWASPAttackSimulator), enabling successful testing of security products and services. It provides a realistic digital banking environment with intentionally implemented vulnerabilities across all OWASP Top 10 categories, making it an ideal platform for:

- **Security Product Testing**: Comprehensive evaluation of security tools and solutions
- **Service Validation**: Testing security services and monitoring capabilities
- **Training & Education**: Hands-on learning for security professionals
- **Research & Development**: Development and testing of new security technologies

## 🚨 **IMPORTANT DISCLAIMER**

**This application is intentionally vulnerable and designed for educational and testing purposes in controlled environments. It demonstrates various security vulnerabilities and should NEVER be used in production environments or deployed on public networks.**

## Vulnerability Types

### 1. Injection Vulnerabilities
- **SQL Injection** (MySQL) - Direct concatenation, partial sanitization, complex scenarios
- **NoSQL Injection** (MongoDB) - JSON injection, operator injection, aggregation
- **Command Injection** - System command execution, file operations

### 2. Broken Authentication
- **Weak Password Storage** - Plain text, MD5, SHA1, Base64
- **Insecure JWT** - Weak secrets, long expiration, algorithm confusion
- **Session Management** - Predictable IDs, weak validation

### 3. Sensitive Data Exposure
- **Plain Text Passwords** - Unencrypted storage
- **Weak Encryption** - Outdated algorithms
- **Information Disclosure** - Verbose error messages

### 4. XSS (Cross-Site Scripting)
- **Reflected XSS** - Direct output, partial sanitization
- **Stored XSS** - Database storage, complex bypass
- **DOM XSS** - Client-side manipulation
- **XSS in JSON** - API responses

### 5. Broken Access Control
- **IDOR** - Insecure direct object references
- **Privilege Escalation** - Admin bypass, role manipulation
- **Missing Authorization** - Unprotected endpoints

### 6. Security Misconfiguration
- **Verbose Error Messages** - Detailed stack traces
- **Misconfigured CORS** - Overly permissive settings
- **Default Configurations** - Weak defaults

### 7. Insecure Deserialization
- **Pickle Injection** - Unsafe deserialization with command execution
- **JSON Injection** - Malicious payloads

### 8. XXE (XML External Entity)
- **File Read** - Direct file system access
- **Remote Entity** - External resource inclusion
- **Parameter Entity** - Complex injection scenarios
- **Chained Attacks** - Multiple attack vectors

### 9. Vulnerable Components
- **Outdated Dependencies** - Known vulnerabilities
- **Unpatched Libraries** - Security issues

### 10. Insufficient Logging/Monitoring
- **Missing Audit Logs** - No security events
- **Unmonitored Actions** - Admin operations

## Quick Start

### Prerequisites
- Docker and Docker Compose

### Setup
1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/securewallet-vulnerable-app.git
   cd securewallet-vulnerable-app
   ```

2. **Setup environment variables**
   ```bash
   cp env.example .env
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - **Frontend**: http://localhost:3000
   - **Backend API**: http://localhost:8000
   - **API Documentation**: http://localhost:8000/docs

## Default Credentials

- **Admin User**: `admin@vulnerable-app.com` / `admin123`

## Vulnerability Levels

Set `VULNERABILITY_LEVEL` in your environment:

- **basic**: Simple vulnerabilities for beginners
- **medium**: Moderate complexity with some protection
- **hard**: Advanced scenarios requiring bypass techniques
- **expert**: Complex chained attacks and advanced techniques

## License

This project is licensed under the MIT License.

## Acknowledgments

- **[OWASPAttackSimulator](https://github.com/r1z4x/OWASPAttackSimulator)** - This application is dedicated to and designed for OWASPAttackSimulator
- OWASP Foundation for the OWASP Top 10
- FastAPI for the excellent web framework
- All contributors and security researchers

---

**Remember: This application is intentionally vulnerable. Use responsibly and only in controlled environments.**
