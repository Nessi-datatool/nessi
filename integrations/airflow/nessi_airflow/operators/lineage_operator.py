"""
Nessi.dev Lineage Operator for Apache Airflow.

This module provides an operator for working with data lineage in Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time
import json

from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException
from airflow.lineage import apply_lineage, get_lineage

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiLineageOperator(BaseOperator):
    """
    Operator for working with data lineage in Nessi.dev.

    This operator can fetch lineage information from Nessi.dev and also
    push Airflow lineage information to Nessi.dev.

    :param table_name: Name of the table to get lineage for
    :type table_name: str
    :param max_depth: Maximum depth of lineage graph
    :type max_depth: Optional[int]
    :param format: Visualization format (html, json, svg, dot)
    :type format: str
    :param push_airflow_lineage: Whether to push Airflow lineage to Nessi.dev
    :type push_airflow_lineage: bool
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_lineage: Whether to push lineage results to XCom
    :type xcom_push_lineage: bool
    """

    template_fields: Sequence[str] = ('table_name',)
    ui_color = '#9370db'  # Medium purple for lineage

    @apply_defaults
    def __init__(
        self,
        *,
        table_name: str,
        max_depth: Optional[int] = None,
        format: str = 'json',
        push_airflow_lineage: bool = False,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_lineage: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.table_name = table_name
        self.max_depth = max_depth
        self.format = format
        self.push_airflow_lineage = push_airflow_lineage
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_lineage = xcom_push_lineage

    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the lineage operation.

        :param context: Airflow context
        :return: Lineage results
        """
        self.log.info("Working with Nessi lineage for table: %s", self.table_name)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Get lineage from Nessi.dev
        if self.format == 'json':
            lineage = hook.get_lineage(
                table_name=self.table_name,
                max_depth=self.max_depth,
            )
            self.log.info(
                "Retrieved lineage for table %s: %d nodes, %d edges",
                self.table_name,
                len(lineage.get('nodes', [])),
                len(lineage.get('edges', [])),
            )
        else:
            lineage = hook.visualize_lineage(
                table_name=self.table_name,
                format=self.format,
                max_depth=self.max_depth,
            )
            self.log.info(
                "Retrieved %s lineage visualization for table %s",
                self.format,
                self.table_name,
            )
        
        # Push Airflow lineage to Nessi.dev if requested
        if self.push_airflow_lineage:
            self._push_airflow_lineage(context, hook)
        
        # Push lineage to XCom if requested
        if self.xcom_push_lineage:
            return lineage
        
        return None
    
    def _push_airflow_lineage(self, context: Dict[str, Any], hook: NessiHook) -> None:
        """
        Push Airflow lineage information to Nessi.dev.

        :param context: Airflow context
        :param hook: Nessi hook
        """
        try:
            # Get Airflow lineage
            lineage_data = get_lineage(self, context)
            
            if not lineage_data:
                self.log.warning("No Airflow lineage data available to push to Nessi.dev")
                return
            
            # Convert Airflow lineage to Nessi.dev format
            nessi_lineage = self._convert_airflow_lineage(lineage_data)
            
            # Push to Nessi.dev
            self.log.info("Pushing Airflow lineage to Nessi.dev")
            
            # This would call a Nessi.dev API endpoint to push lineage
            # For now, we'll just log it as this endpoint might not exist yet
            self.log.info(
                "Would push lineage to Nessi.dev: %s",
                json.dumps(nessi_lineage, indent=2),
            )
            
        except Exception as e:
            self.log.error("Failed to push Airflow lineage to Nessi.dev: %s", e)
    
    def _convert_airflow_lineage(self, lineage_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        Convert Airflow lineage to Nessi.dev format.

        :param lineage_data: Airflow lineage data
        :return: Nessi.dev lineage data
        """
        # This is a placeholder implementation
        # The actual conversion would depend on the Nessi.dev lineage format
        
        nodes = []
        edges = []
        
        # Process inlets (upstream)
        for inlet in lineage_data.get('inlets', []):
            node_id = f"inlet_{inlet.get('name', 'unknown')}"
            nodes.append({
                "id": node_id,
                "name": inlet.get('name', 'unknown'),
                "type": "table",
                "description": inlet.get('description', ''),
            })
            
            edges.append({
                "source_id": node_id,
                "target_id": "task",
                "type": "read",
            })
        
        # Process outlets (downstream)
        for outlet in lineage_data.get('outlets', []):
            node_id = f"outlet_{outlet.get('name', 'unknown')}"
            nodes.append({
                "id": node_id,
                "name": outlet.get('name', 'unknown'),
                "type": "table",
                "description": outlet.get('description', ''),
            })
            
            edges.append({
                "source_id": "task",
                "target_id": node_id,
                "type": "write",
            })
        
        # Add task node
        nodes.append({
            "id": "task",
            "name": self.task_id,
            "type": "process",
            "description": f"Airflow task: {self.task_id}",
        })
        
        return {
            "name": f"Airflow Lineage - {self.task_id}",
            "description": f"Lineage from Airflow task {self.task_id}",
            "nodes": nodes,
            "edges": edges,
        }
