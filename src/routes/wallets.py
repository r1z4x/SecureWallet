from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
from typing import List
from src.config.database import get_db
from src.services.auth import get_current_user_dependency
from src.services.wallet import (
    create_wallet, get_user_wallets, get_wallet_by_id, update_wallet,
    delete_wallet, transfer_funds, deposit_funds, withdraw_funds,
    get_wallet_transactions
)
from src.schemas.wallet import (
    WalletCreate, WalletUpdate, WalletResponse, WalletTransfer,
    WalletDeposit, WalletWithdrawal
)
from src.schemas.transaction import TransactionResponse
from src.models.user import User
from src.models.wallet import Wallet

router = APIRouter()
security = HTTPBearer()


@router.post("/", response_model=WalletResponse)
async def create_user_wallet(
    wallet_data: WalletCreate,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Create a new wallet for the current user"""
    wallet = create_wallet(db, current_user.id, wallet_data)
    return wallet


@router.get("/", response_model=List[WalletResponse])
async def get_wallets(
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get all wallets for the current user"""
    wallets = get_user_wallets(db, current_user.id)
    return wallets


@router.get("/balance")
async def get_wallet_balance(
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get current user's wallet balance and transaction count"""
    wallet = db.query(Wallet).filter(
        Wallet.user_id == current_user.id,
        Wallet.is_active == True
    ).first()
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Get transaction count
    from src.models.transaction import Transaction
    transaction_count = db.query(Transaction).filter(
        (Transaction.from_wallet_id == wallet.id) | (Transaction.to_wallet_id == wallet.id)
    ).count()
    
    return {
        "balance": float(wallet.balance),
        "transaction_count": transaction_count,
        "currency": wallet.currency
    }


@router.get("/{wallet_id}", response_model=WalletResponse)
async def get_wallet(
    wallet_id: int,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get a specific wallet by ID"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to access this wallet"
        )
    
    return wallet


@router.put("/{wallet_id}", response_model=WalletResponse)
async def update_user_wallet(
    wallet_id: int,
    wallet_data: WalletUpdate,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Update a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to modify this wallet"
        )
    
    updated_wallet = update_wallet(db, wallet_id, wallet_data)
    return updated_wallet


@router.delete("/{wallet_id}")
async def delete_user_wallet(
    wallet_id: int,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Delete a wallet (soft delete)"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to delete this wallet"
        )
    
    success = delete_wallet(db, wallet_id)
    if not success:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to delete wallet"
        )
    
    return {"message": "Wallet deleted successfully"}


@router.post("/{wallet_id}/transfer", response_model=TransactionResponse)
async def transfer_from_wallet(
    wallet_id: int,
    transfer_data: WalletTransfer,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Transfer funds from a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Source wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to transfer from this wallet"
        )
    
    transaction = transfer_funds(db, wallet_id, transfer_data)
    
    if not transaction:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Transfer failed. Check wallet balance and destination wallet."
        )
    
    return transaction


@router.post("/{wallet_id}/deposit", response_model=TransactionResponse)
async def deposit_to_wallet(
    wallet_id: int,
    deposit_data: WalletDeposit,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Deposit funds to a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to deposit to this wallet"
        )
    
    transaction = deposit_funds(db, wallet_id, deposit_data)
    
    if not transaction:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Deposit failed"
        )
    
    return transaction


@router.post("/{wallet_id}/withdraw", response_model=TransactionResponse)
async def withdraw_from_wallet(
    wallet_id: int,
    withdrawal_data: WalletWithdrawal,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Withdraw funds from a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to withdraw from this wallet"
        )
    
    transaction = withdraw_funds(db, wallet_id, withdrawal_data)
    
    if not transaction:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Withdrawal failed. Check wallet balance."
        )
    
    return transaction


@router.get("/{wallet_id}/transactions", response_model=List[TransactionResponse])
async def get_wallet_transaction_history(
    wallet_id: int,
    limit: int = 50,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get transaction history for a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user owns the wallet
    if wallet.user_id != current_user.id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to view this wallet's transactions"
        )
    
    transactions = get_wallet_transactions(db, wallet_id, limit)
    return transactions


@router.post("/transfer")
async def transfer_funds_simple(
    transfer_data: dict,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Simple transfer endpoint for frontend"""
    recipient = transfer_data.get("recipient")
    amount = transfer_data.get("amount")
    description = transfer_data.get("description", "Transfer")
    
    if not recipient or not amount:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Recipient and amount are required"
        )
    
    # Get sender's wallet
    sender_wallet = db.query(Wallet).filter(
        Wallet.user_id == current_user.id,
        Wallet.is_active == True
    ).first()
    
    if not sender_wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if sender has sufficient balance
    if sender_wallet.balance < amount:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Insufficient balance"
        )
    
    # Find recipient by email or user ID
    recipient_user = None
    if "@" in recipient:  # Email
        recipient_user = db.query(User).filter(User.email == recipient).first()
    else:  # User ID
        try:
            user_id = int(recipient)
            recipient_user = db.query(User).filter(User.id == user_id).first()
        except ValueError:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Invalid recipient format"
            )
    
    if not recipient_user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Recipient not found"
        )
    
    # Get recipient's wallet
    recipient_wallet = db.query(Wallet).filter(
        Wallet.user_id == recipient_user.id,
        Wallet.is_active == True
    ).first()
    
    if not recipient_wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Recipient wallet not found"
        )
    
    # Create transaction
    from src.models.transaction import Transaction, TransactionType, TransactionStatus
    transaction = Transaction(
        from_wallet_id=sender_wallet.id,
        to_wallet_id=recipient_wallet.id,
        amount=amount,
        transaction_type=TransactionType.TRANSFER,
        status=TransactionStatus.COMPLETED,
        description=description
    )
    
    # Update balances
    sender_wallet.balance -= amount
    recipient_wallet.balance += amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    
    return {"message": "Transfer completed successfully", "transaction_id": transaction.id}


@router.post("/deposit")
async def deposit_funds_simple(
    deposit_data: dict,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Simple deposit endpoint for frontend"""
    amount = deposit_data.get("amount")
    description = deposit_data.get("description", "Deposit")
    
    if not amount or amount <= 0:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Invalid amount"
        )
    
    # Get user's wallet
    wallet = db.query(Wallet).filter(
        Wallet.user_id == current_user.id,
        Wallet.is_active == True
    ).first()
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Create deposit transaction
    from src.models.transaction import Transaction, TransactionType, TransactionStatus
    transaction = Transaction(
        from_wallet_id=None,  # External deposit
        to_wallet_id=wallet.id,
        amount=amount,
        transaction_type=TransactionType.DEPOSIT,
        status=TransactionStatus.COMPLETED,
        description=description
    )
    
    # Update balance
    wallet.balance += amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    
    return {"message": "Deposit completed successfully", "transaction_id": transaction.id}


@router.post("/withdraw")
async def withdraw_funds_simple(
    withdrawal_data: dict,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Simple withdrawal endpoint for frontend"""
    amount = withdrawal_data.get("amount")
    description = withdrawal_data.get("description", "Withdrawal")
    
    if not amount or amount <= 0:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Invalid amount"
        )
    
    # Get user's wallet
    wallet = db.query(Wallet).filter(
        Wallet.user_id == current_user.id,
        Wallet.is_active == True
    ).first()
    
    if not wallet:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Wallet not found"
        )
    
    # Check if user has sufficient balance
    if wallet.balance < amount:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Insufficient balance"
        )
    
    # Create withdrawal transaction
    from src.models.transaction import Transaction, TransactionType, TransactionStatus
    transaction = Transaction(
        from_wallet_id=wallet.id,
        to_wallet_id=None,  # External withdrawal
        amount=amount,
        transaction_type=TransactionType.WITHDRAWAL,
        status=TransactionStatus.COMPLETED,
        description=description
    )
    
    # Update balance
    wallet.balance -= amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    
    return {"message": "Withdrawal completed successfully", "transaction_id": transaction.id}
