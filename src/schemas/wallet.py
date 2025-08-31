from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime
from decimal import Decimal


class WalletBase(BaseModel):
    wallet_name: str
    currency: str = "USD"


class WalletCreate(WalletBase):
    pass


class WalletUpdate(BaseModel):
    wallet_name: Optional[str] = None
    currency: Optional[str] = None
    is_active: Optional[bool] = None


class WalletResponse(WalletBase):
    id: int
    user_id: int
    balance: Decimal
    is_active: bool
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class WalletTransfer(BaseModel):
    to_wallet_id: int
    amount: Decimal = Field(gt=0)
    description: Optional[str] = None


class WalletDeposit(BaseModel):
    amount: Decimal = Field(gt=0)
    description: Optional[str] = None


class WalletWithdrawal(BaseModel):
    amount: Decimal = Field(gt=0)
    description: Optional[str] = None
