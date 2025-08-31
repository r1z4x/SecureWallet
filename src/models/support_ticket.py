from sqlalchemy import Column, Integer, String, Text, DateTime, ForeignKey, Boolean
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from src.config.database import Base


class SupportTicket(Base):
    __tablename__ = "support_tickets"
    
    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    subject = Column(String(255), nullable=False)
    message = Column(Text, nullable=False)
    category = Column(String(50), default="general")  # account, transaction, security, technical, billing, general
    status = Column(String(50), default="open")  # open, in_progress, closed
    priority = Column(String(20), default="medium")  # low, medium, high, urgent
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    # Relationships
    user = relationship("User", back_populates="support_tickets")
    
    def __repr__(self):
        return f"<SupportTicket(id={self.id}, subject='{self.subject}', status='{self.status}')>"
