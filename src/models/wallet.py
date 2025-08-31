from sqlalchemy import Column, Integer, String, Boolean, DateTime, Numeric, ForeignKey
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from src.config.database import Base


class Wallet(Base):
    __tablename__ = "wallets"
    
    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    wallet_name = Column(String(100), nullable=False)
    balance = Column(Numeric(15, 2), default=0.00)
    currency = Column(String(3), default="USD")
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    
    # Relationships
    user = relationship("User", back_populates="wallets")
    from_transactions = relationship("Transaction", foreign_keys="Transaction.from_wallet_id", back_populates="from_wallet")
    to_transactions = relationship("Transaction", foreign_keys="Transaction.to_wallet_id", back_populates="to_wallet")
    
    def __repr__(self):
        return f"<Wallet(id={self.id}, user_id={self.user_id}, balance={self.balance})>"
