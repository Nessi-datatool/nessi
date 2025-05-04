from fastapi import APIRouter, HTTPException, Depends
from typing import List, Dict, Any
from .pipeline import (
    PipelineManager,
    Pipeline,
    FilterTransformation,
    RenameTransformation,
    TypeCastTransformation,
    FillNullTransformation
)
from ..security.rbac import Permission, get_api_key
from ..security.user_manager import user_manager
from ..security.rbac import rbac

router = APIRouter(prefix="/pipelines", tags=["pipelines"])
pipeline_manager = PipelineManager()

@router.post("/", response_model=Dict)
async def create_pipeline(
    name: str,
    api_key: str = Depends(get_api_key)
):
    """Create a new transformation pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    try:
        pipeline = pipeline_manager.create_pipeline(name)
        return pipeline.to_dict()
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.get("/", response_model=List[Dict])
async def list_pipelines(api_key: str = Depends(get_api_key)):
    """List all transformation pipelines."""
    # Check if requester has permission to read data
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.READ_DATA):
        raise HTTPException(
            status_code=403,
            detail="Permission to read data required"
        )
    
    return pipeline_manager.list_pipelines()

@router.post("/{pipeline_name}/transformations/filter", response_model=Dict)
async def add_filter_transformation(
    pipeline_name: str,
    condition: str,
    api_key: str = Depends(get_api_key)
):
    """Add a filter transformation to a pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    pipeline = pipeline_manager.get_pipeline(pipeline_name)
    if not pipeline:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_name} not found")
    
    transformation = FilterTransformation(condition)
    pipeline.add_transformation(transformation)
    return pipeline.to_dict()

@router.post("/{pipeline_name}/transformations/rename", response_model=Dict)
async def add_rename_transformation(
    pipeline_name: str,
    column_mapping: Dict[str, str],
    api_key: str = Depends(get_api_key)
):
    """Add a rename transformation to a pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    pipeline = pipeline_manager.get_pipeline(pipeline_name)
    if not pipeline:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_name} not found")
    
    transformation = RenameTransformation(column_mapping)
    pipeline.add_transformation(transformation)
    return pipeline.to_dict()

@router.post("/{pipeline_name}/transformations/type_cast", response_model=Dict)
async def add_type_cast_transformation(
    pipeline_name: str,
    type_mapping: Dict[str, str],
    api_key: str = Depends(get_api_key)
):
    """Add a type cast transformation to a pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    pipeline = pipeline_manager.get_pipeline(pipeline_name)
    if not pipeline:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_name} not found")
    
    transformation = TypeCastTransformation(type_mapping)
    pipeline.add_transformation(transformation)
    return pipeline.to_dict()

@router.post("/{pipeline_name}/transformations/fill_null", response_model=Dict)
async def add_fill_null_transformation(
    pipeline_name: str,
    fill_values: Dict[str, Any],
    api_key: str = Depends(get_api_key)
):
    """Add a fill null transformation to a pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    pipeline = pipeline_manager.get_pipeline(pipeline_name)
    if not pipeline:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_name} not found")
    
    transformation = FillNullTransformation(fill_values)
    pipeline.add_transformation(transformation)
    return pipeline.to_dict()

@router.delete("/{pipeline_name}", status_code=204)
async def delete_pipeline(
    pipeline_name: str,
    api_key: str = Depends(get_api_key)
):
    """Delete a transformation pipeline."""
    # Check if requester has permission to configure system
    requester = user_manager.get_user_by_api_key(api_key)
    if not requester or not rbac.has_permission(requester.role, Permission.CONFIGURE_SYSTEM):
        raise HTTPException(
            status_code=403,
            detail="Permission to configure system required"
        )
    
    pipeline_manager.delete_pipeline(pipeline_name) 