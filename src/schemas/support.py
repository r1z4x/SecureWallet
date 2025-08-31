from pydantic import BaseModel
from typing import Optional
from datetime import datetime


class SupportTicketBase(BaseModel):
    subject: str
    message: str
    category: Optional[str] = "general"
    priority: Optional[str] = "medium"


class SupportTicketCreate(SupportTicketBase):
    pass


class SupportTicketResponse(SupportTicketBase):
    id: int
    user_id: int
    status: str
    priority: str
    category: str
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    class Config:
        from_attributes = True


class SupportTicketUpdate(BaseModel):
    subject: Optional[str] = None
    message: Optional[str] = None
    category: Optional[str] = None
    status: Optional[str] = None
    priority: Optional[str] = None
