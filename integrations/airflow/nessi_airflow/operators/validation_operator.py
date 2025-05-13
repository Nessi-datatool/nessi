"""
Nessi.dev Validation Operator for Apache Airflow.

This module provides an operator for running data validation with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiValidationOperator(BaseOperator):
    """
    Operator for running data validation with Nessi.dev.

    This operator runs validation rules on a table and waits for the results.
    It can be configured to fail the task if the validation fails.

    :param table_name: Name of the table to validate
    :type table_name: str
    :param rules: List of validation rules
    :type rules: List[Dict[str, Any]]
    :param fail_on_validation_failure: Whether to fail the task if validation fails
    :type fail_on_validation_failure: bool
    :param timeout: Timeout in seconds for the validation
    :type timeout: Optional[int]
    :param poll_interval: Interval in seconds between polling for results
    :type poll_interval: int
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_validation_results: Whether to push validation results to XCom
    :type xcom_push_validation_results: bool
    """

    template_fields: Sequence[str] = ('table_name', 'rules')
    ui_color = '#ff7f50'  # Coral color for validation

    @apply_defaults
    def __init__(
        self,
        *,
        table_name: str,
        rules: List[Dict[str, Any]],
        fail_on_validation_failure: bool = True,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_validation_results: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.table_name = table_name
        self.rules = rules
        self.fail_on_validation_failure = fail_on_validation_failure
        self.timeout = timeout
        self.poll_interval = poll_interval
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_validation_results = xcom_push_validation_results

    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the data validation.

        :param context: Airflow context
        :return: Validation results
        """
        self.log.info("Running Nessi data validation on table: %s", self.table_name)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Start validation
        response = hook.run_validation(
            table_name=self.table_name,
            rules=self.rules,
            timeout=self.timeout,
        )
        
        validation_id = response.get('validation_id')
        if not validation_id:
            raise AirflowException(f"Failed to start validation: {response}")
        
        self.log.info("Validation started with ID: %s", validation_id)
        
        # Poll for results
        start_time = time.time()
        max_time = None if self.timeout is None else start_time + self.timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise AirflowException(
                    f"Timeout reached while waiting for validation results: {validation_id}"
                )
            
            # Get results
            results = hook.get_validation_results(validation_id)
            status = results.get('status')
            
            if status == 'completed':
                self.log.info("Validation completed: %s", validation_id)
                break
            elif status == 'failed':
                raise AirflowException(f"Validation failed: {results.get('error')}")
            elif status == 'running':
                self.log.info("Validation still running, waiting...")
                time.sleep(self.poll_interval)
            else:
                raise AirflowException(f"Unknown validation status: {status}")
        
        # Check validation results
        if self.fail_on_validation_failure:
            rule_results = results.get('rule_results', [])
            failed_rules = [r for r in rule_results if r.get('status') == 'failed']
            if failed_rules:
                failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
                raise AirflowException(
                    f"Validation has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
                )
        
        # Log validation summary
        rule_count = len(results.get('rule_results', []))
        passed_count = len([r for r in results.get('rule_results', []) if r.get('status') == 'passed'])
        self.log.info(
            "Validation completed for table %s: %d/%d rules passed",
            self.table_name,
            passed_count,
            rule_count,
        )
        
        # Push results to XCom if requested
        if self.xcom_push_validation_results:
            return results
        
        return None
