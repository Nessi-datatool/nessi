"""
Nessi.dev profile tasks for Prefect.

This module provides tasks for running data profiling with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union
import logging

from prefect import task

from nessi_prefect.client import NessiClient

logger = logging.getLogger(__name__)


@task(name="Run Nessi Profile")
def run_profile(
    table_name: str,
    columns: Optional[List[str]] = None,
    sample_size: Optional[int] = None,
    timeout: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Run a data profile on a table.

    Args:
        table_name: Name of the table to profile
        columns: List of columns to profile (all if not specified)
        sample_size: Number of rows to sample
        timeout: Timeout in seconds for the profiling
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Response containing the profile ID
    """
    logger.info("Running Nessi data profiling on table: %s", table_name)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Run profiling
    response = client.run_profile(
        table_name=table_name,
        columns=columns,
        sample_size=sample_size,
        timeout=timeout,
    )
    
    profile_id = response.get('profile_id')
    if not profile_id:
        raise RuntimeError(f"Failed to start profiling: {response}")
    
    logger.info("Profiling started with ID: %s", profile_id)
    return response


@task(name="Wait for Nessi Profile")
def wait_for_profile(
    profile_id: str,
    timeout: Optional[int] = None,
    poll_interval: int = 10,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Wait for profile results.

    Args:
        profile_id: ID of the profile
        timeout: Maximum time to wait in seconds
        poll_interval: Time between polls in seconds
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Profile results
    """
    logger.info("Waiting for Nessi data profile results: %s", profile_id)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Wait for profile
    profile = client.wait_for_profile(
        profile_id=profile_id,
        timeout=timeout,
        poll_interval=poll_interval,
    )
    
    # Log profile summary
    column_count = len(profile.get('columns', []))
    row_count = profile.get('row_count', 0)
    logger.info(
        "Profile completed for table: %d columns, %d rows",
        column_count,
        row_count,
    )
    
    return profile
