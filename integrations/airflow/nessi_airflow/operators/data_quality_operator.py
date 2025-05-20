"""
Nessi.dev Data Quality Operator for Apache Airflow.

This module provides an operator for running data quality checks with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiDataQualityOperator(BaseOperator):
    """
    Operator for running data quality checks with Nessi.dev.

    This operator runs a data quality check on a table and waits for the results.
    It can be configured to fail the task if the quality check fails.

    :param table_name: Name of the table to check
    :type table_name: str
    :param rules: List of quality rules to apply
    :type rules: Optional[List[Dict[str, Any]]]
    :param profile: Whether to generate a profile
    :type profile: bool
    :param quality_threshold: Minimum quality score to consider the check successful
    :type quality_threshold: Optional[float]
    :param fail_on_rule_failure: Whether to fail the task if any rule fails
    :type fail_on_rule_failure: bool
    :param timeout: Timeout in seconds for the quality check
    :type timeout: Optional[int]
    :param poll_interval: Interval in seconds between polling for results
    :type poll_interval: int
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_quality_results: Whether to push quality results to XCom
    :type xcom_push_quality_results: bool
    """

    template_fields: Sequence[str] = ('table_name', 'rules')
    ui_color = '#f0e68c'  # Khaki color for data quality

    @apply_defaults
    def __init__(
        self,
        *,
        table_name: str,
        rules: Optional[List[Dict[str, Any]]] = None,
        profile: bool = True,
        quality_threshold: Optional[float] = None,
        fail_on_rule_failure: bool = True,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_quality_results: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.table_name = table_name
        self.rules = rules
        self.profile = profile
        self.quality_threshold = quality_threshold
        self.fail_on_rule_failure = fail_on_rule_failure
        self.timeout = timeout
        self.poll_interval = poll_interval
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_quality_results = xcom_push_quality_results

    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the data quality check.

        :param context: Airflow context
        :return: Quality check results
        """
        self.log.info("Running Nessi data quality check on table: %s", self.table_name)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Start quality check
        response = hook.run_quality_check(
            table_name=self.table_name,
            rules=self.rules,
            profile=self.profile,
            timeout=self.timeout,
        )
        
        check_id = response.get('check_id')
        if not check_id:
            raise AirflowException(f"Failed to start quality check: {response}")
        
        self.log.info("Quality check started with ID: %s", check_id)
        
        # Poll for results
        start_time = time.time()
        max_time = None if self.timeout is None else start_time + self.timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise AirflowException(
                    f"Timeout reached while waiting for quality check results: {check_id}"
                )
            
            # Get results
            results = hook.get_quality_results(check_id)
            status = results.get('status')
            
            if status == 'completed':
                self.log.info("Quality check completed: %s", check_id)
                break
            elif status == 'failed':
                raise AirflowException(f"Quality check failed: {results.get('error')}")
            elif status == 'running':
                self.log.info("Quality check still running, waiting...")
                time.sleep(self.poll_interval)
            else:
                raise AirflowException(f"Unknown quality check status: {status}")
        
        # Check quality score
        quality_score = results.get('quality_score')
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
            rule_results = results.get('rule_results', [])
            failed_rules = [r for r in rule_results if r.get('status') == 'failed']
            if failed_rules:
                failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
                raise AirflowException(
                    f"Quality check has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
                )
        
        # Push results to XCom if requested
        if self.xcom_push_quality_results:
            return results
        
        return None
