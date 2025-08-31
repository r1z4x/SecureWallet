from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime
from decimal import Decimal
from src.models.transaction import TransactionType, TransactionStatus


class TransactionBase(BaseModel):
    amount: Decimal = Field(gt=0)
    description: Optional[str] = None


class TransactionResponse(TransactionBase):
    id: int
    from_wallet_id: Optional[int] = None
    to_wallet_id: Optional[int] = None
    transaction_type: TransactionType
    status: TransactionStatus
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class TransactionCreate(TransactionBase):
    from_wallet_id: Optional[int] = None
    to_wallet_id: Optional[int] = None
    transaction_type: TransactionType


class TransactionUpdate(BaseModel):
    status: Optional[TransactionStatus] = None
    description: Optional[str] = None
