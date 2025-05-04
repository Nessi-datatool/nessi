from typing import List, Optional
from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel, EmailStr
from .user_manager import UserManager, User
from .rbac import Role, Permission, RBAC

router = APIRouter(prefix="/users", tags=["users"])
user_manager = UserManager()
rbac = RBAC()

class UserCreate(BaseModel):
    """User creation request model."""
    username: str
    email: EmailStr
    role: Role

class UserResponse(BaseModel):
    """User response model."""
    username: str
    email: EmailStr
    role: Role
    created_at: str
    last_login: Optional[str]
    is_active: bool

class UserUpdate(BaseModel):
    """User update request model."""
    role: Optional[Role] = None
    is_active: Optional[bool] = None

@router.post("/", response_model=UserResponse)
@rbac.require_permission(Permission.MANAGE_USERS)
async def create_user(user: UserCreate, api_key: str = Depends(rbac.get_api_key)):
    """Create a new user."""
    try:
        new_user = user_manager.create_user(
            username=user.username,
            email=user.email,
            role=user.role
        )
        return UserResponse(**new_user.to_dict())
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.get("/", response_model=List[UserResponse])
@rbac.require_permission(Permission.MANAGE_USERS)
async def list_users(api_key: str = Depends(rbac.get_api_key)):
    """List all users."""
    return [UserResponse(**user) for user in user_manager.list_users()]

@router.get("/{username}", response_model=UserResponse)
@rbac.require_permission(Permission.MANAGE_USERS)
async def get_user(username: str, api_key: str = Depends(rbac.get_api_key)):
    """Get user by username."""
    user = user_manager.get_user(username)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return UserResponse(**user.to_dict())

@router.patch("/{username}", response_model=UserResponse)
@rbac.require_permission(Permission.MANAGE_USERS)
async def update_user(
    username: str,
    user_update: UserUpdate,
    api_key: str = Depends(rbac.get_api_key)
):
    """Update user information."""
    user = user_manager.get_user(username)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    if user_update.role:
        user = user_manager.update_user_role(username, user_update.role)
    
    if user_update.is_active is not None:
        if not user_update.is_active:
            user_manager.deactivate_user(username)
    
    return UserResponse(**user.to_dict())

@router.get("/{username}/permissions", response_model=List[str])
@rbac.require_permission(Permission.MANAGE_USERS)
async def get_user_permissions(
    username: str,
    api_key: str = Depends(rbac.get_api_key)
):
    """Get user permissions."""
    user = user_manager.get_user(username)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user_manager.get_user_permissions(username) 