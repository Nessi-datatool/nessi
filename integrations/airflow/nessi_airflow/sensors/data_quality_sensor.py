"""
Nessi.dev Data Quality Sensor for Apache Airflow.

This module provides a sensor for waiting for data quality checks to complete in Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.sensors.base import BaseSensorOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiDataQualitySensor(BaseSensorOperator):
    """
    Sensor for waiting for data quality checks to complete in Nessi.dev.

    This sensor waits for a data quality check to complete and can be configured
    to fail if the quality check fails.

    :param check_id: ID of the quality check to wait for
    :type check_id: str
    :param quality_threshold: Minimum quality score to consider the check successful
    :type quality_threshold: Optional[float]
    :param fail_on_rule_failure: Whether to fail the task if any rule fails
    :type fail_on_rule_failure: bool
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_quality_results: Whether to push quality results to XCom
    :type xcom_push_quality_results: bool
    """

    template_fields: Sequence[str] = ('check_id',)
    ui_color = '#f0e68c'  # Khaki color for data quality

    @apply_defaults
    def __init__(
        self,
        *,
        check_id: str,
        quality_threshold: Optional[float] = None,
        fail_on_rule_failure: bool = True,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_quality_results: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.check_id = check_id
        self.quality_threshold = quality_threshold
        self.fail_on_rule_failure = fail_on_rule_failure
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_quality_results = xcom_push_quality_results
        self.results = None

    def poke(self, context: Dict[str, Any]) -> bool:
        """
        Check if the data quality check has completed.

        :param context: Airflow context
        :return: True if the check has completed, False otherwise
        """
        self.log.info("Checking status of Nessi data quality check: %s", self.check_id)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Get results
        results = hook.get_quality_results(self.check_id)
        status = results.get('status')
        
        if status == 'completed':
            self.log.info("Quality check completed: %s", self.check_id)
            self.results = results
            return True
        elif status == 'failed':
            raise AirflowException(f"Quality check failed: {results.get('error')}")
        elif status == 'running':
            self.log.info("Quality check still running, waiting...")
            return False
        else:
            raise AirflowException(f"Unknown quality check status: {status}")
    
    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the sensor after the poke method returns True.

        :param context: Airflow context
        :return: Quality check results
        """
        # Check quality score
        quality_score = self.results.get('quality_score')
        if quality_score is not None and self.quality_threshold is not None:
            if quality_score < self.quality_threshold:
                raise AirflowException(
                    f"Quality score {quality_score} is below threshold {self.quality_threshold}"
                )
            self.log.info(
                "Quality score %s meets threshold %s",
                quality_score,
                self.quality_threshold,
            )
        
        # Check rule failures
        if self.fail_on_rule_failure:
            rule_results = self.results.get('rule_results', [])
            failed_rules = [r for r in rule_results if r.get('status') == 'failed']
            if failed_rules:
                failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
                raise AirflowException(
                    f"Quality check has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
                )
        
        # Push results to XCom if requested
        if self.xcom_push_quality_results:
            return self.results
        
        return None
