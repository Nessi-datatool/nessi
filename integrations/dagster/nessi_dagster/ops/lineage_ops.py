"""
Nessi.dev lineage operations for Dagster.

This module provides operations for retrieving and visualizing data lineage with Nessi.dev.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import op, Out, In, Nothing


@op(
    name="get_lineage",
    description="Get lineage for a table",
    required_resource_keys={"nessi"},
    ins={
        "table_name": In(str, description="Name of the table"),
        "max_depth": In(
            Optional[int],
            description="Maximum depth of lineage graph",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Lineage graph"),
)
def get_lineage(
    context,
    table_name: str,
    max_depth: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Get lineage for a table.

    Args:
        context: Dagster execution context
        table_name: Name of the table
        max_depth: Maximum depth of lineage graph

    Returns:
        Lineage graph
    """
    context.log.info(f"Getting lineage for table: {table_name}")
    
    lineage = context.resources.nessi.client.get_lineage(
        table_name=table_name,
        max_depth=max_depth,
    )
    
    # Log lineage summary
    node_count = len(lineage.get("nodes", []))
    edge_count = len(lineage.get("edges", []))
    
    context.log.info(
        f"Lineage retrieved: {node_count} nodes and {edge_count} edges"
    )
    
    return lineage


@op(
    name="visualize_lineage",
    description="Get a visualization of lineage for a table",
    required_resource_keys={"nessi"},
    ins={
        "table_name": In(str, description="Name of the table"),
        "format": In(
            str,
            description="Visualization format (html, json, svg, dot)",
            default_value="html",
        ),
        "max_depth": In(
            Optional[int],
            description="Maximum depth of lineage graph",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Lineage visualization"),
)
def visualize_lineage(
    context,
    table_name: str,
    format: str = "html",
    max_depth: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Get a visualization of lineage for a table.

    Args:
        context: Dagster execution context
        table_name: Name of the table
        format: Visualization format (html, json, svg, dot)
        max_depth: Maximum depth of lineage graph

    Returns:
        Lineage visualization
    """
    context.log.info(f"Visualizing lineage for table: {table_name} in format: {format}")
    
    visualization = context.resources.nessi.client.visualize_lineage(
        table_name=table_name,
        format=format,
        max_depth=max_depth,
    )
    
    context.log.info(f"Lineage visualization generated in {format} format")
    
    return visualization
