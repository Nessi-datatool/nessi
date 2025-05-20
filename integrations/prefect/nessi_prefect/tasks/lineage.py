"""
Nessi.dev lineage tasks for Prefect.

This module provides tasks for working with data lineage in Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union
import logging

from prefect import task

from nessi_prefect.client import NessiClient

logger = logging.getLogger(__name__)


@task(name="Get Nessi Lineage")
def get_lineage(
    table_name: str,
    max_depth: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Get lineage for a table.

    Args:
        table_name: Name of the table
        max_depth: Maximum depth of lineage graph
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Lineage graph
    """
    logger.info("Getting Nessi lineage for table: %s", table_name)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Get lineage
    lineage = client.get_lineage(
        table_name=table_name,
        max_depth=max_depth,
    )
    
    logger.info(
        "Retrieved lineage for table %s: %d nodes, %d edges",
        table_name,
        len(lineage.get('nodes', [])),
        len(lineage.get('edges', [])),
    )
    
    return lineage


@task(name="Visualize Nessi Lineage")
def visualize_lineage(
    table_name: str,
    format: str = "html",
    max_depth: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Get a visualization of lineage for a table.

    Args:
        table_name: Name of the table
        format: Visualization format (html, json, svg, dot)
        max_depth: Maximum depth of lineage graph
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Lineage visualization
    """
    logger.info("Visualizing Nessi lineage for table: %s", table_name)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Get visualization
    visualization = client.visualize_lineage(
        table_name=table_name,
        format=format,
        max_depth=max_depth,
    )
    
    logger.info(
        "Retrieved %s lineage visualization for table %s",
        format,
        table_name,
    )
    
    return visualization
