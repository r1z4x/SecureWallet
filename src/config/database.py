from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from pymongo import MongoClient
import redis
from .settings import settings

# SQLAlchemy setup
engine = create_engine(
    settings.database_url,
    pool_pre_ping=True,
    pool_recycle=300,
    echo=settings.debug
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()


def get_db():
    """Dependency to get database session"""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


# MongoDB setup
mongodb_client = MongoClient(settings.mongodb_url)
mongodb_db = mongodb_client.get_default_database()


def get_mongodb():
    """Dependency to get MongoDB database"""
    return mongodb_db


# Redis setup
redis_client = redis.from_url(settings.redis_url)


def get_redis():
    """Dependency to get Redis client"""
    return redis_client
