# SecureWallet - Digital Banking Platform

A comprehensive digital banking platform built with Go (backend) and Vue.js (frontend), featuring secure wallet management, transaction processing, and user authentication.

## Features

- **User Management**: Secure user registration, authentication, and profile management
- **Wallet System**: Multi-currency wallet support with balance tracking
- **Transaction Processing**: Secure money transfers between users
- **Two-Factor Authentication**: Enhanced security with TOTP support
- **Admin Panel**: Comprehensive administrative tools and monitoring
- **Audit Logging**: Complete audit trail for all system activities
- **API Documentation**: Swagger/OpenAPI documentation

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### Environment Setup

1. Copy `env.example` to `.env` and configure your environment variables
2. Set up your MySQL database and Redis instance
3. Configure the database connection details in your `.env` file

### Database Schema Issues

If you encounter foreign key constraint errors like:
```
Error 3780 (HY000): Referencing column 'user_id' and referenced column 'id' in foreign key constraint 'fk_users_wallets' are incompatible.
```

This indicates a schema mismatch between your existing database and the new UUID-based schema. To resolve this:

#### Quick Fix Script (Recommended)
```bash
# Run the interactive database fix script
./fix-database.sh
```

This script will guide you through the available options and set the necessary environment variables.

#### Manual Options

**Option 1: Force Database Recreation (Recommended for Development)**
```bash
# Set environment variable
export FORCE_DATABASE_RECREATION=true
export RESET_DATABASE_ON_STARTUP=true

# Restart the application
./start-go-dev.sh
```

**Option 2: Use the Force Recreation API Endpoint**
```bash
# Call the force recreation endpoint
curl -X POST http://localhost:8080/api/data/force-recreate
```

**Option 3: Manual Database Reset**
```bash
# Set environment variable
export RESET_DATABASE_ON_STARTUP=true

# Restart the application
./start-go-dev.sh
```

### Running the Application

1. **Backend (Go)**
   ```bash
   # Development mode
   ./start-go-dev.sh
   
   # Production mode
   ./start-go.sh
   ```

2. **Frontend (Vue.js)**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

3. **Docker (Alternative)**
   ```bash
   # Development
   docker-compose -f docker-compose.dev.yml up
   
   # Production
   docker-compose up
   ```

## API Documentation

Once the application is running, you can access the Swagger documentation at:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Base URL**: `http://localhost:8080/api`

## Security Features

- JWT-based authentication
- Two-factor authentication (TOTP)
- Rate limiting and input validation
- Secure password hashing with bcrypt
- Comprehensive audit logging
- CORS protection

## Development

### Project Structure

```
SecureWallet/
├── internal/           # Go backend code
│   ├── config/        # Configuration management
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── routes/        # API routes
│   └── services/      # Business logic
├── frontend/          # Vue.js frontend
├── db/               # Database initialization scripts
├── docs/             # API documentation
└── docker/           # Docker configuration
```

### Database Models

The application uses UUIDs for all primary keys and foreign keys, ensuring:
- Global uniqueness
- No auto-increment conflicts
- Better security (no predictable IDs)
- Distributed system compatibility

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Check the API documentation at `/swagger/index.html`
- Review the logs for detailed error information
- Use the database reset endpoints if you encounter schema issues
