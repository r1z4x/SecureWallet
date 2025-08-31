# SecureWallet - Digital Banking Platform (Vulnerable)

A comprehensive vulnerable application designed for OWASP Top 10 training and educational purposes.

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

## API Examples

### SQL Injection
```bash
curl "http://localhost:8000/api/vulnerabilities/sql-injection/basic/user-search?username=admin' OR '1'='1"
```

### XSS
```bash
curl "http://localhost:8000/api/vulnerabilities/xss/basic/reflected?user_input=<script>alert('XSS')</script>"
```

### Command Injection
```bash
curl "http://localhost:8000/api/vulnerabilities/command-injection/basic/ping?host=127.0.0.1; ls -la"
```

## License

This project is licensed under the MIT License.

## Acknowledgments

- OWASP Foundation for the OWASP Top 10
- FastAPI for the excellent web framework
- All contributors and security researchers

---

**Remember: This application is intentionally vulnerable. Use responsibly and only in controlled environments.**
