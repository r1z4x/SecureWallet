from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from src.config.database import get_db, get_mongodb, get_redis
from src.config.settings import settings
from src.routes import auth, users, wallets, transactions, admin, vulnerabilities, support, data_management
from src.models import user, wallet, transaction, session, audit_log
from src.models.user import User
from src.services.auth import get_current_user_dependency
import os

# Create FastAPI app
app = FastAPI(
    title="SecureWallet - Digital Banking Platform (Vulnerable)",
    description="A comprehensive vulnerable application for OWASP Top 10",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

# Security scheme
security = HTTPBearer()

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://127.0.0.1:3000"],  # Vue.js dev server
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(auth.router, prefix="/api/auth", tags=["Authentication"])
app.include_router(users.router, prefix="/api/users", tags=["Users"])
app.include_router(wallets.router, prefix="/api/wallets", tags=["Wallets"])
app.include_router(transactions.router, prefix="/api/transactions", tags=["Transactions"])
app.include_router(admin.router, prefix="/api/admin", tags=["Admin"])
app.include_router(vulnerabilities.router, prefix="/api/vulnerabilities", tags=["Vulnerabilities"])
app.include_router(support.router, prefix="/api/support", tags=["Support"])
app.include_router(data_management.router, prefix="/api/data", tags=["Data Management"])


@app.get("/")
async def root():
    """Root endpoint - API information"""
    return {
        "message": "SecureWallet - Digital Banking Platform (Vulnerable) API",
        "version": "1.0.0",
        "docs": "/docs",
        "health": "/health"
    }


@app.get("/health")
async def health_check(current_user: User = Depends(get_current_user_dependency)):
    """Health check endpoint"""
    return {
        "status": "healthy",
        "app_name": "SecureWallet - Digital Banking Platform (Vulnerable)",
        "user": current_user.username,
        "timestamp": "2024-01-01T00:00:00Z"
    }


@app.get("/api/info")
async def get_api_info():
    """Get API information"""
    return {
        "message": "SecureWallet - Digital Banking Platform (Vulnerable) API",
        "version": "1.0.0",
        "status": "running",
        "environment": os.getenv("ENVIRONMENT", "development"),
        "vulnerability_level": os.getenv("VULNERABILITY_LEVEL", "basic")
    }


@app.get("/api/auth/me")
async def get_current_user_info(
    current_user: User = Depends(get_current_user_dependency)
):
    """Get current user information"""
    return current_user


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
