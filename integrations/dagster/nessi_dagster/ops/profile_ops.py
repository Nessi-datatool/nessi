"""
Nessi.dev profile operations for Dagster.

This module provides operations for running data profiling with Nessi.dev.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import op, Out, In, Nothing


@op(
    name="run_profile",
    description="Run a data profile on a table",
    required_resource_keys={"nessi"},
    ins={
        "table_name": In(str, description="Name of the table to profile"),
        "columns": In(
            Optional[List[str]],
            description="List of columns to profile (all if not specified)",
            default_value=None,
        ),
        "sample_size": In(
            Optional[int],
            description="Number of rows to sample",
            default_value=None,
        ),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for the profiling",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Profile response"),
)
def run_profile(
    context,
    table_name: str,
    columns: Optional[List[str]] = None,
    sample_size: Optional[int] = None,
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Run a data profile on a table.

    Args:
        context: Dagster execution context
        table_name: Name of the table to profile
        columns: List of columns to profile (all if not specified)
        sample_size: Number of rows to sample
        timeout: Timeout in seconds for the profiling

    Returns:
        Profile response
    """
    context.log.info(f"Running profile on table: {table_name}")
    
    response = context.resources.nessi.client.run_profile(
        table_name=table_name,
        columns=columns,
        sample_size=sample_size,
        timeout=timeout,
    )
    
    context.log.info(f"Profile started with ID: {response.get('profile_id')}")
    return response


@op(
    name="wait_for_profile",
    description="Wait for profile results",
    required_resource_keys={"nessi"},
    ins={
        "profile_id": In(str, description="ID of the profile"),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for waiting",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Profile results"),
)
def wait_for_profile(
    context,
    profile_id: str,
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Wait for profile results.

    Args:
        context: Dagster execution context
        profile_id: ID of the profile
        timeout: Timeout in seconds for waiting

    Returns:
        Profile results
    """
    context.log.info(f"Waiting for profile results: {profile_id}")
    
    results = context.resources.nessi.client.wait_for_profile(
        profile_id=profile_id,
        timeout=timeout,
    )
    
    # Log profile summary
    column_count = len(results.get("columns", []))
    row_count = results.get("row_count", 0)
    
    context.log.info(
        f"Profile completed for {column_count} columns and {row_count} rows"
    )
    
    return results
