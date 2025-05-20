"""
Nessi.dev Profile Operator for Apache Airflow.

This module provides an operator for running data profiling with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiProfileOperator(BaseOperator):
    """
    Operator for running data profiling with Nessi.dev.

    This operator runs a data profile on a table and waits for the results.

    :param table_name: Name of the table to profile
    :type table_name: str
    :param columns: List of columns to profile (all if not specified)
    :type columns: Optional[List[str]]
    :param sample_size: Number of rows to sample
    :type sample_size: Optional[int]
    :param timeout: Timeout in seconds for the profiling
    :type timeout: Optional[int]
    :param poll_interval: Interval in seconds between polling for results
    :type poll_interval: int
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_profile: Whether to push profile results to XCom
    :type xcom_push_profile: bool
    """

    template_fields: Sequence[str] = ('table_name', 'columns')
    ui_color = '#87cefa'  # Light sky blue for profiling

    @apply_defaults
    def __init__(
        self,
        *,
        table_name: str,
        columns: Optional[List[str]] = None,
        sample_size: Optional[int] = None,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_profile: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.table_name = table_name
        self.columns = columns
        self.sample_size = sample_size
        self.timeout = timeout
        self.poll_interval = poll_interval
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_profile = xcom_push_profile

    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the data profiling.

        :param context: Airflow context
        :return: Profile results
        """
        self.log.info("Running Nessi data profiling on table: %s", self.table_name)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Start profiling
        response = hook.run_profile(
            table_name=self.table_name,
            columns=self.columns,
            sample_size=self.sample_size,
            timeout=self.timeout,
        )
        
        profile_id = response.get('profile_id')
        if not profile_id:
            raise AirflowException(f"Failed to start profiling: {response}")
        
        self.log.info("Profiling started with ID: %s", profile_id)
        
        # Poll for results
        start_time = time.time()
        max_time = None if self.timeout is None else start_time + self.timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise AirflowException(
                    f"Timeout reached while waiting for profile results: {profile_id}"
                )
            
            # Get profile
            profile = hook.get_profile(profile_id)
            status = profile.get('status')
            
            if status == 'completed':
                self.log.info("Profiling completed: %s", profile_id)
                break
            elif status == 'failed':
                raise AirflowException(f"Profiling failed: {profile.get('error')}")
            elif status == 'running':
                self.log.info("Profiling still running, waiting...")
                time.sleep(self.poll_interval)
            else:
                raise AirflowException(f"Unknown profiling status: {status}")
        
        # Log profile summary
        column_count = len(profile.get('columns', []))
        row_count = profile.get('row_count', 0)
        self.log.info(
            "Profile completed for table %s: %d columns, %d rows",
            self.table_name,
            column_count,
            row_count,
        )
        
        # Push profile to XCom if requested
        if self.xcom_push_profile:
            return profile
        
        return None
