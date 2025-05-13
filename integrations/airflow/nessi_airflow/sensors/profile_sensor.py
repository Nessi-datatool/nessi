"""
Nessi.dev Profile Sensor for Apache Airflow.

This module provides a sensor for waiting for data profiling to complete in Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.sensors.base import BaseSensorOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiProfileSensor(BaseSensorOperator):
    """
    Sensor for waiting for data profiling to complete in Nessi.dev.

    This sensor waits for a data profile to complete.

    :param profile_id: ID of the profile to wait for
    :type profile_id: str
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_profile: Whether to push profile results to XCom
    :type xcom_push_profile: bool
    """

    template_fields: Sequence[str] = ('profile_id',)
    ui_color = '#87cefa'  # Light sky blue for profiling

    @apply_defaults
    def __init__(
        self,
        *,
        profile_id: str,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_profile: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.profile_id = profile_id
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_profile = xcom_push_profile
        self.profile = None

    def poke(self, context: Dict[str, Any]) -> bool:
        """
        Check if the data profiling has completed.

        :param context: Airflow context
        :return: True if the profiling has completed, False otherwise
        """
        self.log.info("Checking status of Nessi data profiling: %s", self.profile_id)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Get profile
        profile = hook.get_profile(self.profile_id)
        status = profile.get('status')
        
        if status == 'completed':
            self.log.info("Profiling completed: %s", self.profile_id)
            self.profile = profile
            return True
        elif status == 'failed':
            raise AirflowException(f"Profiling failed: {profile.get('error')}")
        elif status == 'running':
            self.log.info("Profiling still running, waiting...")
            return False
        else:
            raise AirflowException(f"Unknown profiling status: {status}")
    
    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the sensor after the poke method returns True.

        :param context: Airflow context
        :return: Profile results
        """
        # Log profile summary
        column_count = len(self.profile.get('columns', []))
        row_count = self.profile.get('row_count', 0)
        self.log.info(
            "Profile completed: %d columns, %d rows",
            column_count,
            row_count,
        )
        
        # Push profile to XCom if requested
        if self.xcom_push_profile:
            return self.profile
        
        return None
