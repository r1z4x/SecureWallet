from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List
from src.config.database import get_db
from src.services.auth import get_current_user_dependency
from src.models.user import User
from src.models.support_ticket import SupportTicket
from src.schemas.support import SupportTicketCreate, SupportTicketResponse

router = APIRouter()


@router.post("/ticket", response_model=SupportTicketResponse)
async def create_support_ticket(
    ticket_data: SupportTicketCreate,
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Create a new support ticket"""
    ticket = SupportTicket(
        user_id=current_user.id,
        subject=ticket_data.subject,
        message=ticket_data.message,
        category=ticket_data.category,
        priority=ticket_data.priority,
        status="open"
    )
    
    db.add(ticket)
    db.commit()
    db.refresh(ticket)
    
    return ticket


@router.get("/tickets", response_model=List[SupportTicketResponse])
async def get_user_tickets(
    current_user: User = Depends(get_current_user_dependency),
    db: Session = Depends(get_db)
):
    """Get all tickets for the current user"""
    tickets = db.query(SupportTicket).filter(SupportTicket.user_id == current_user.id).all()
    return tickets
