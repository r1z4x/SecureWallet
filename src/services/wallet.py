from decimal import Decimal
from typing import List, Optional
from sqlalchemy.orm import Session
from src.models.wallet import Wallet
from src.models.transaction import Transaction, TransactionType, TransactionStatus
from src.schemas.wallet import WalletCreate, WalletUpdate, WalletTransfer, WalletDeposit, WalletWithdrawal


def create_wallet(db: Session, user_id: int, wallet_data: WalletCreate) -> Wallet:
    """Create a new wallet for a user"""
    wallet = Wallet(
        user_id=user_id,
        wallet_name=wallet_data.wallet_name,
        currency=wallet_data.currency
    )
    db.add(wallet)
    db.commit()
    db.refresh(wallet)
    return wallet


def get_user_wallets(db: Session, user_id: int) -> List[Wallet]:
    """Get all wallets for a user"""
    return db.query(Wallet).filter(Wallet.user_id == user_id, Wallet.is_active == True).all()


def get_wallet_by_id(db: Session, wallet_id: int) -> Optional[Wallet]:
    """Get a wallet by ID"""
    return db.query(Wallet).filter(Wallet.id == wallet_id).first()


def update_wallet(db: Session, wallet_id: int, wallet_data: WalletUpdate) -> Optional[Wallet]:
    """Update a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    if not wallet:
        return None
    
    update_data = wallet_data.dict(exclude_unset=True)
    for field, value in update_data.items():
        setattr(wallet, field, value)
    
    db.commit()
    db.refresh(wallet)
    return wallet


def delete_wallet(db: Session, wallet_id: int) -> bool:
    """Delete a wallet (soft delete)"""
    wallet = get_wallet_by_id(db, wallet_id)
    if not wallet:
        return False
    
    wallet.is_active = False
    db.commit()
    return True


def transfer_funds(db: Session, from_wallet_id: int, transfer_data: WalletTransfer) -> Optional[Transaction]:
    """Transfer funds between wallets"""
    from_wallet = get_wallet_by_id(db, from_wallet_id)
    to_wallet = get_wallet_by_id(db, transfer_data.to_wallet_id)
    
    if not from_wallet or not to_wallet:
        return None
    
    if from_wallet.balance < transfer_data.amount:
        return None
    
    # Create transaction
    transaction = Transaction(
        from_wallet_id=from_wallet_id,
        to_wallet_id=transfer_data.to_wallet_id,
        amount=transfer_data.amount,
        transaction_type=TransactionType.TRANSFER,
        status=TransactionStatus.COMPLETED,
        description=transfer_data.description
    )
    
    # Update balances
    from_wallet.balance -= transfer_data.amount
    to_wallet.balance += transfer_data.amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    return transaction


def deposit_funds(db: Session, wallet_id: int, deposit_data: WalletDeposit) -> Optional[Transaction]:
    """Deposit funds to a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    if not wallet:
        return None
    
    # Create transaction
    transaction = Transaction(
        to_wallet_id=wallet_id,
        amount=deposit_data.amount,
        transaction_type=TransactionType.DEPOSIT,
        status=TransactionStatus.COMPLETED,
        description=deposit_data.description
    )
    
    # Update balance
    wallet.balance += deposit_data.amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    return transaction


def withdraw_funds(db: Session, wallet_id: int, withdrawal_data: WalletWithdrawal) -> Optional[Transaction]:
    """Withdraw funds from a wallet"""
    wallet = get_wallet_by_id(db, wallet_id)
    if not wallet or wallet.balance < withdrawal_data.amount:
        return None
    
    # Create transaction
    transaction = Transaction(
        from_wallet_id=wallet_id,
        to_wallet_id=wallet_id,  # Same wallet for withdrawal
        amount=withdrawal_data.amount,
        transaction_type=TransactionType.WITHDRAWAL,
        status=TransactionStatus.COMPLETED,
        description=withdrawal_data.description
    )
    
    # Update balance
    wallet.balance -= withdrawal_data.amount
    
    db.add(transaction)
    db.commit()
    db.refresh(transaction)
    return transaction


def get_wallet_transactions(db: Session, wallet_id: int, limit: int = 50) -> List[Transaction]:
    """Get transaction history for a wallet"""
    return db.query(Transaction).filter(
        (Transaction.from_wallet_id == wallet_id) | (Transaction.to_wallet_id == wallet_id)
    ).order_by(Transaction.created_at.desc()).limit(limit).all()
