from pydantic_settings import BaseSettings
from typing import List, Optional
import os


class Settings(BaseSettings):
    # Database Configuration
    database_url: str = "mysql+pymysql://wallet_user:wallet_pass@localhost:3306/wallet_app"
    mongodb_url: str = "mongodb://admin:adminpass@localhost:27017/wallet_app?authSource=admin"
    redis_url: str = "redis://localhost:6379/0"
    
    # JWT Configuration
    jwt_secret_key: str = "your-super-secret-jwt-key-change-in-production"
    jwt_algorithm: str = "HS256"
    access_token_expire_minutes: int = 30
    
    # Application Configuration
    vulnerability_level: str = "basic"  # basic, medium, hard, expert
    debug: bool = True
    log_level: str = "INFO"
    
    # Security Configuration
    cors_origins: List[str] = ["http://localhost:3000", "http://localhost:8000"]
    allowed_hosts: List[str] = ["localhost", "127.0.0.1"]
    
    # Admin Configuration
    admin_email: str = "admin@vulnerable-app.com"
    admin_password: str = "admin123"
    
    # Logging Configuration
    log_file: str = "logs/app.log"
    log_format: str = "json"
    
    class Config:
        env_file = ".env"
        case_sensitive = False


# Global settings instance
settings = Settings()
