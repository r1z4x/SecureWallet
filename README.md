# SecureWallet - Digital Banking Platform

A modern digital banking and wallet application that provides secure financial services including money transfers, transaction management, and user account administration. This application demonstrates real-world financial technology capabilities while maintaining high security standards.

## 🚨 **IMPORTANT DISCLAIMER**

**This application is designed for educational and testing purposes in controlled environments. It demonstrates various security concepts and should not be used in production environments without proper security hardening.**

## Features

### Core Application Features
- **User Authentication System** - Secure login, logout with JWT tokens
- **Digital Wallet Management** - Create wallets, view balances, transfer funds
- **Transaction Processing** - Deposit, withdrawal, and transfer operations
- **Admin Panel** - User management and system administration
- **RESTful API** - Complete API with Swagger/OpenAPI documentation
- **Modern Web Interface** - Beautiful frontend with TailwindCSS and Vue.js
- **Real-time Banking** - Live transaction updates and account management

### Security Framework
- **Multi-layered Security**: Authentication, authorization, and data protection
- **Comprehensive Testing**: Various security scenarios and edge cases
- **Organized by Category**: Different security aspects and implementations
- **Progressive Complexity**: From basic to advanced security features

## Security Categories

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
- Node.js 18+ (for frontend development)
- Python 3.12+ (for backend development)

### Using Docker (Recommended)

#### Production Setup
1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/OWASP-WSTG-Vulnerable-App.git
   cd OWASP-WSTG-Vulnerable-App
   ```

2. **Setup environment variables**
   ```bash
   cp env.example .env
   # Edit .env file for production settings
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - **Frontend**: http://localhost:3000
   - **Backend API**: http://localhost:8080
   - **API Documentation**: http://localhost:8080/docs
   - **Health Check**: http://localhost:8080/health

#### Development Setup

##### Option 1: Full Docker Development (Recommended)
```bash
# 1. Setup environment variables
cp env.example .env

# 2. Start all services (frontend + backend)
./start-dev.sh

# 3. Access the application
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# API Docs: http://localhost:8080/docs
```

##### Option 2: Hybrid Development (Backend Docker + Frontend Local)
```bash
# 1. Setup environment variables
cp env.example .env

# 2. Start backend services only
docker-compose -f docker-compose.dev.yml up -d mysql mongodb redis backend

# 3. Start frontend development server
cd frontend
npm install
npm run dev

# 4. Access the application
# Frontend (Local): http://localhost:3000
# Backend (Docker): http://localhost:8080
# API Docs: http://localhost:8080/docs
```

##### Option 3: Hybrid Development with Script
```bash
# 1. Setup environment variables
cp env.example .env

# 2. Start hybrid development
./start-dev-hybrid.sh

# 3. Access the application
# Frontend (Local): http://localhost:3000
# Backend (Docker): http://localhost:8080
# API Docs: http://localhost:8080/docs
```

### Local Development

1. **Install dependencies**
   ```bash
   pip install -r requirements.txt
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Start the application**
   ```bash
   uvicorn src.app.main:app --reload --host 0.0.0.0 --port 8000
   ```

## Configuration

### Environment Variables

Create a `.env` file based on `env.example`:

```env
# Database Configuration
DATABASE_URL=mysql+pymysql://wallet_user:wallet_pass@localhost:3306/wallet_app
MONGODB_URL=mongodb://admin:adminpass@localhost:27017/wallet_app?authSource=admin
REDIS_URL=redis://localhost:6379/0

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
JWT_ALGORITHM=HS256
ACCESS_TOKEN_EXPIRE_MINUTES=30

# Application Configuration
VULNERABILITY_LEVEL=basic  # basic, medium, hard, expert
DEBUG=True
LOG_LEVEL=INFO
```

### Vulnerability Levels

Set `VULNERABILITY_LEVEL` in your environment:

- **basic**: Simple vulnerabilities for beginners
- **medium**: Moderate complexity with some protection
- **hard**: Advanced scenarios requiring bypass techniques
- **expert**: Complex chained attacks and advanced techniques

## Web Interface

The application includes a modern, responsive web interface built with TailwindCSS and Vue.js that provides:

- **Interactive Dashboard** - Overview of vulnerability categories and difficulty levels
- **Real-time Testing** - Test vulnerabilities directly through the web interface
- **User Management** - Login, logout, and user profile management
- **Wallet Interface** - Manage funds and view transaction history
- **Admin Panel** - User management and system administration (admin users only)
- **Vulnerability Categories** - Organized sections for each vulnerability type
- **Difficulty Levels** - Color-coded interface for different complexity levels

### Using the Web Interface

1. **Access the interface**: Navigate to http://localhost:8000
2. **Login**: Use the login button to authenticate with your credentials
3. **Navigate**: Use the sidebar to access different vulnerability categories
4. **Test vulnerabilities**: Enter payloads and test different attack scenarios
5. **View results**: See real-time responses and vulnerability information

## API Usage

### Authentication

1. **Register a new user**
   ```bash
   curl -X POST "http://localhost:8000/api/auth/register" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}'
   ```

2. **Login to get access token**
   ```bash
   curl -X POST "http://localhost:8000/api/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"username": "testuser", "password": "password123"}'
   ```

3. **Use the token for authenticated requests**
   ```bash
   curl -X GET "http://localhost:8000/api/wallets/" \
        -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
   ```

### Testing Vulnerabilities

#### SQL Injection Example
```bash
# Basic SQL Injection
curl "http://localhost:8000/api/vulnerabilities/sql-injection/basic/user-search?username=admin' OR '1'='1"

# Union-based SQL Injection
curl "http://localhost:8000/api/vulnerabilities/sql-injection/expert/user-search?username=admin' UNION SELECT 1,2,3,4,5--"
```

#### XSS Example
```bash
# Basic Reflected XSS
curl "http://localhost:8000/api/vulnerabilities/xss/basic/reflected?user_input=<script>alert('XSS')</script>"
```

#### Command Injection Example
```bash
# Basic Command Injection
curl "http://localhost:8000/api/vulnerabilities/command-injection/basic/ping?host=127.0.0.1; ls -la"
```

#### XXE Example
```bash
# Basic XXE - File read
curl -X POST "http://localhost:8000/api/vulnerabilities/xxe/basic/xml-upload" \
     -H "Content-Type: application/json" \
     -d '{"xml_content": "<?xml version=\"1.0\"?><!DOCTYPE test [<!ENTITY xxe SYSTEM \"file:///etc/passwd\">]><root><data>&xxe;</data></root>"}'
```

#### Pickle Injection Example
```bash
# Basic Pickle Injection
curl -X POST "http://localhost:8000/api/vulnerabilities/pickle-injection/basic/deserialize" \
     -H "Content-Type: application/json" \
     -d '{"data": "cos\\nsystem\\n(S\"ls -la\"\\ntR."}'
```

## Default Credentials

- **Admin User**: `admin@vulnerable-app.com` / `admin123`
- **Database**: `wallet_user` / `wallet_pass`
- **MongoDB**: `admin` / `adminpass`

## Project Structure

```
OWASP-WSTG-Vulnerable-App/
├── src/
│   ├── app/
│   │   └── main.py                 # FastAPI application
│   ├── config/
│   │   ├── settings.py             # Configuration settings
│   │   └── database.py             # Database configuration
│   ├── models/
│   │   ├── user.py                 # User model
│   │   ├── wallet.py               # Wallet model
│   │   ├── transaction.py          # Transaction model
│   │   ├── session.py              # Session model
│   │   └── audit_log.py            # Audit log model
│   ├── schemas/
│   │   ├── user.py                 # User schemas
│   │   ├── wallet.py               # Wallet schemas
│   │   └── transaction.py          # Transaction schemas
│   ├── services/
│   │   ├── auth.py                 # Authentication service
│   │   └── wallet.py               # Wallet service
│   ├── routes/
│   │   ├── auth.py                 # Authentication routes
│   │   ├── users.py                # User management routes
│   │   ├── wallets.py              # Wallet routes
│   │   ├── transactions.py         # Transaction routes
│   │   ├── admin.py                # Admin routes
│   │   └── vulnerabilities.py      # Vulnerability testing routes
│   └── vulnerabilities/
│       ├── injection/
│       │   ├── sql_injection.py    # SQL injection vulnerabilities
│       │   ├── nosql_injection.py  # NoSQL injection vulnerabilities
│       │   └── command_injection.py # Command injection vulnerabilities
│       ├── authentication/
│       │   └── weak_auth.py        # Weak authentication vulnerabilities
│       ├── xss/
│       │   └── xss_vulnerabilities.py # XSS vulnerabilities
│       ├── xxe/
│       │   └── xxe_vulnerabilities.py # XXE vulnerabilities
│       └── deserialization/
│           └── pickle_injection.py # Pickle injection vulnerabilities
├── db/
│   └── init.sql                    # Database initialization
├── docker-compose.yml              # Docker services
├── Dockerfile                      # Application container
├── requirements.txt                # Python dependencies
└── README.md                       # This file
```

## Contributing

### Adding New Vulnerabilities

1. **Create vulnerability class** in appropriate category
2. **Add difficulty levels** (basic, medium, hard, expert)
3. **Create route endpoints** in `src/routes/vulnerabilities.py`
4. **Update documentation** and examples

### Coding Standards

- Follow PEP 8 for Python code
- Use type hints
- Add comprehensive docstrings
- Include vulnerability descriptions
- Test all difficulty levels

## Security Considerations

### For Educational Use
- Use in isolated environments only
- Never expose to public networks
- Regularly update dependencies
- Monitor for unauthorized access

### For Testing
- Document all test scenarios
- Use appropriate tools (Burp Suite, OWASP ZAP, etc.)
- Follow responsible disclosure practices
- Report any unintended vulnerabilities

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- OWASP Foundation for the Web Security Testing Guide
- FastAPI for the excellent web framework
- SQLAlchemy for database ORM
- All contributors and security researchers

## Support

For questions, issues, or contributions:
- Create an issue on GitHub
- Follow security best practices
- Provide detailed reproduction steps

---

**Remember: This application is intentionally vulnerable. Use responsibly and only in controlled environments.**
