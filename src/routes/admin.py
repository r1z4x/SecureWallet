from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from typing import List, Dict, Any
from src.config.database import get_db
from src.services.auth import get_current_user_dependency
from src.models.user import User
from src.models.wallet import Wallet
from src.models.transaction import Transaction
from src.models.audit_log import AuditLog
from sqlalchemy import func

router = APIRouter()
security = HTTPBearer()


def require_admin(current_user: User = Depends(get_current_user_dependency)):
    """Dependency to require admin access"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Admin access required"
        )
    return current_user


@router.get("/dashboard")
async def admin_dashboard(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Admin dashboard with system statistics"""
    # Get system statistics
    total_users = db.query(User).count()
    active_users = db.query(User).filter(User.is_active == True).count()
    total_wallets = db.query(Wallet).count()
    active_wallets = db.query(Wallet).filter(Wallet.is_active == True).count()
    total_transactions = db.query(Transaction).count()
    
    # Get recent audit logs
    recent_logs = db.query(AuditLog).order_by(AuditLog.created_at.desc()).limit(10).all()
    
    return {
        "statistics": {
            "total_users": total_users,
            "active_users": active_users,
            "total_wallets": total_wallets,
            "active_wallets": active_wallets,
            "total_transactions": total_transactions
        },
        "recent_audit_logs": recent_logs
    }


@router.get("/users", response_model=List[Dict[str, Any]])
async def admin_get_users(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Get all users with detailed information"""
    users = db.query(User).all()
    return [
        {
            "id": user.id,
            "username": user.username,
            "email": user.email,
            "is_active": user.is_active,
            "is_admin": user.is_admin,
            "created_at": user.created_at,
            "wallet_count": len(user.wallets)
        }
        for user in users
    ]


@router.get("/wallets", response_model=List[Dict[str, Any]])
async def admin_get_wallets(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Get all wallets with detailed information"""
    wallets = db.query(Wallet).all()
    return [
        {
            "id": wallet.id,
            "user_id": wallet.user_id,
            "wallet_name": wallet.wallet_name,
            "balance": float(wallet.balance),
            "currency": wallet.currency,
            "is_active": wallet.is_active,
            "created_at": wallet.created_at
        }
        for wallet in wallets
    ]


@router.get("/transactions", response_model=List[Dict[str, Any]])
async def admin_get_transactions(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Get all transactions with detailed information"""
    transactions = db.query(Transaction).all()
    return [
        {
            "id": transaction.id,
            "from_wallet_id": transaction.from_wallet_id,
            "to_wallet_id": transaction.to_wallet_id,
            "amount": float(transaction.amount),
            "transaction_type": transaction.transaction_type.value,
            "status": transaction.status.value,
            "description": transaction.description,
            "created_at": transaction.created_at
        }
        for transaction in transactions
    ]


@router.get("/audit-logs", response_model=List[Dict[str, Any]])
async def admin_get_audit_logs(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Get all audit logs"""
    logs = db.query(AuditLog).order_by(AuditLog.created_at.desc()).all()
    return [
        {
            "id": log.id,
            "user_id": log.user_id,
            "action": log.action,
            "resource_type": log.resource_type,
            "resource_id": log.resource_id,
            "details": log.details,
            "ip_address": log.ip_address,
            "user_agent": log.user_agent,
            "created_at": log.created_at
        }
        for log in logs
    ]


@router.post("/users/{user_id}/toggle-status")
async def admin_toggle_user_status(
    user_id: int,
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Toggle user active status"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    user.is_active = not user.is_active
    db.commit()
    
    return {"message": f"User status toggled to {'active' if user.is_active else 'inactive'}"}


@router.post("/users/{user_id}/toggle-admin")
async def admin_toggle_user_admin(
    user_id: int,
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Toggle user admin status"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    user.is_admin = not user.is_admin
    db.commit()
    
    return {"message": f"User admin status toggled to {'admin' if user.is_admin else 'user'}"}


@router.get("/stats")
async def get_admin_stats(
    current_user: User = Depends(require_admin),
    db: Session = Depends(get_db)
):
    """Get admin statistics (admin only)"""
    # Get total users
    total_users = db.query(User).count()
    
    # Get total transactions
    total_transactions = db.query(Transaction).count()
    
    # Get total transaction volume
    total_volume = db.query(func.sum(Transaction.amount)).scalar() or 0
    
    return {
        "total_users": total_users,
        "total_transactions": total_transactions,
        "total_volume": float(total_volume)
    }


@router.get("/system-info")
async def admin_system_info(
    current_user: User = Depends(require_admin)
):
    """Get system information"""
    from src.config.settings import settings
    
    return {
        "vulnerability_level": settings.vulnerability_level,
        "debug_mode": settings.debug,
        "log_level": settings.log_level,
        "cors_origins": settings.cors_origins,
        "allowed_hosts": settings.allowed_hosts
    }
