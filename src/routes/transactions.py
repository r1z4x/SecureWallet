from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from typing import List
from src.config.database import get_db
from src.services.auth import get_current_user_dependency
from src.schemas.transaction import TransactionResponse
from src.models.user import User
from src.models.transaction import Transaction

router = APIRouter()
security = HTTPBearer()


@router.get("/", response_model=List[TransactionResponse])
async def get_transactions(
    limit: int = 50,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get all transactions for the current user"""
    # Get transactions where user is involved (either as sender or receiver)
    transactions = db.query(Transaction).join(
        Transaction.from_wallet
    ).filter(
        Transaction.from_wallet.has(user_id=current_user.id)
    ).union(
        db.query(Transaction).join(
            Transaction.to_wallet
        ).filter(
            Transaction.to_wallet.has(user_id=current_user.id)
        )
    ).order_by(Transaction.created_at.desc()).limit(limit).all()
    
    return transactions


@router.get("/{transaction_id}", response_model=TransactionResponse)
async def get_transaction(
    transaction_id: int,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get a specific transaction by ID"""
    transaction = db.query(Transaction).filter(Transaction.id == transaction_id).first()
    
    if not transaction:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Transaction not found"
        )
    
    # Check if user is involved in the transaction
    user_involved = False
    if transaction.from_wallet and transaction.from_wallet.user_id == current_user.id:
        user_involved = True
    if transaction.to_wallet and transaction.to_wallet.user_id == current_user.id:
        user_involved = True
    
    if not user_involved and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to view this transaction"
        )
    
    return transaction
