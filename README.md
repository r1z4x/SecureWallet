# OWASP WSTG Vulnerable Application

A comprehensive vulnerable web application designed for OWASP Web Security Testing Guide (WSTG) and Attack Simulator integration. This application simulates real-world vulnerabilities across multiple OWASP Top 10 categories with advanced attack techniques including second-order attacks, out-of-band (OOB) exfiltration, and complex bypass methods.

## 🎯 Purpose

This application serves as a testing platform for:
- **OWASPAttackSimulator** integration and testing
- **Security researchers** and penetration testers
- **Educational purposes** for learning web application security
- **Vulnerability assessment** training and certification

The application contains intentionally vulnerable code with vulnerabilities integrated into the existing application functionality, making it suitable for comprehensive security testing scenarios.

## 🏗️ Architecture

### Backend (Go Gin)
- **Framework**: Go Gin with GORM ORM
- **Database**: MySQL with Redis caching
- **Authentication**: JWT-based with multiple vulnerability patterns
- **Vulnerability Integration**: Payload-based aggregation system with webhook logging

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API
- **Styling**: Tailwind CSS
- **State Management**: Pinia
- **Build Tool**: Vite

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+
- Node.js 16+

### 1. Clone and Setup
```bash
git clone <repository-url>
cd OWASP-WSTG-Vulnerable-App
```

### 2. Environment Configuration
```bash
cp env.example .env
# Edit .env file with your configuration
```

### 3. Start the Application
```bash
# Start all services
docker-compose up -d

# Or start development environment
./start-go-dev.sh
```

### 4. Access the Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/docs

## 🔧 Vulnerability Integration Approach

The application uses a **payload-based aggregation system** that:

1. **Logs payloads** sent by the Attack Simulator
2. **Logs webhook requests** made by the Attack Simulator
3. **Aggregates vulnerability levels** per test case
4. **Calculates final difficulty** based on maximum level triggered
5. **Reports results** via webhook to the Attack Simulator with request logs

### Vulnerability Levels

The system supports four progressive difficulty levels:

- **basic**: Simple, direct vulnerabilities
- **medium**: Vulnerabilities with weak protection mechanisms and blind attacks
- **hard**: Complex vulnerabilities requiring advanced bypass techniques and time-based attacks
- **expert**: Sophisticated vulnerabilities with multiple attack vectors including OOB and second-order attacks

## 🎯 Integrated Vulnerability Types

### Authentication Routes (`/api/auth/*`)
- **Weak Password Storage**: Multiple formats (plain text, MD5, SHA1, Base64)
- **Weak Authentication Logic**: Multiple bypass techniques
- **Weak JWT Implementation**: Algorithm confusion, weak secrets, long expiration
- **Password Reset Vulnerabilities**: OOB attacks, second-order storage

### User Management Routes (`/api/users/*`)
- **SQL Injection**: Advanced techniques (OOB, second-order, union-based)
- **IDOR (Insecure Direct Object References)**: Complex bypass patterns
- **XSS (Cross-Site Scripting)**: Advanced sanitization bypass techniques

### Transaction Routes (`/api/transactions/*`)
- **SQL Injection**: Blind, time-based, union-based attacks
- **Command Injection**: DNS, HTTP, and advanced command chaining
- **IDOR**: Transaction access control bypass

### Support System Routes (`/api/support/*`)
- **XSS**: Reflected, stored, and DOM-based attacks
