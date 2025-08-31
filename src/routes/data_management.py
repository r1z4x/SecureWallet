from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from typing import List, Dict, Any
from src.config.database import get_db
from src.services.data_manager import DataManager
from src.services.auth import get_current_user_dependency
from src.models.user import User

router = APIRouter()

@router.post("/demo-data")
async def create_demo_data(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Create demo data with sample users and transactions"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can create demo data"
        )
    
    data_manager = DataManager(db)
    result = data_manager.create_demo_data()
    
    return {
        "message": "Demo data created successfully",
        "data": result
    }

@router.post("/snapshot")
async def create_snapshot(
    version: str,
    description: str = "",
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Create a data snapshot"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can create snapshots"
        )
    
    data_manager = DataManager(db)
    snapshot = data_manager.create_snapshot(version, description)
    
    return {
        "message": "Snapshot created successfully",
        "snapshot": snapshot
    }

@router.get("/snapshots")
async def list_snapshots(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """List all available snapshots"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can list snapshots"
        )
    
    data_manager = DataManager(db)
    snapshots = data_manager.list_snapshots()
    
    return {
        "snapshots": snapshots
    }

@router.get("/snapshot/{version}")
async def get_snapshot(
    version: str,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Get a specific snapshot"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can view snapshots"
        )
    
    data_manager = DataManager(db)
    snapshot = data_manager.load_snapshot(version)
    
    if not snapshot:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Snapshot not found"
        )
    
    return {
        "snapshot": snapshot
    }

@router.post("/snapshot/{version}/restore")
async def restore_snapshot(
    version: str,
    clear_existing: bool = True,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Restore data from a snapshot"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can restore snapshots"
        )
    
    data_manager = DataManager(db)
    success = data_manager.restore_snapshot(version, clear_existing)
    
    if not success:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Failed to restore snapshot"
        )
    
    return {
        "message": f"Snapshot {version} restored successfully"
    }

@router.delete("/clear-data")
async def clear_all_data(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Clear all data from database"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can clear data"
        )
    
    data_manager = DataManager(db)
    data_manager._clear_all_data()
    
    return {
        "message": "All data cleared successfully"
    }

@router.get("/credentials")
async def get_demo_credentials(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Get demo user credentials"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can view credentials"
        )
    
    data_manager = DataManager(db)
    return data_manager.get_all_credentials()

@router.post("/setup-fresh")
async def setup_fresh_database(
    db: Session = Depends(get_db)
):
    """Setup fresh database with admin and demo data (no auth required)"""
    data_manager = DataManager(db)
    result = data_manager.setup_fresh_database()
    
    return result

@router.post("/create-admin")
async def create_admin_user(
    db: Session = Depends(get_db)
):
    """Create admin user (no auth required)"""
    data_manager = DataManager(db)
    result = data_manager.create_admin_only()
    
    return result

@router.post("/reset-demo")
async def reset_to_demo(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    """Reset database to demo state"""
    if not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Only admins can reset database"
        )
    
    data_manager = DataManager(db)
    result = data_manager.reset_to_demo()
    
    return result

@router.post("/make-admin/{username}")
async def make_user_admin(
    username: str,
    db: Session = Depends(get_db)
):
    """Make a user admin (no auth required)"""
    data_manager = DataManager(db)
    result = data_manager.make_user_admin(username)
    
    return result
