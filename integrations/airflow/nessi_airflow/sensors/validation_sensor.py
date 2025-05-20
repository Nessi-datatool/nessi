"""
Nessi.dev Validation Sensor for Apache Airflow.

This module provides a sensor for waiting for data validation to complete in Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union, Sequence
import time

from airflow.sensors.base import BaseSensorOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowException

from nessi_airflow.hooks.nessi_hook import NessiHook


class NessiValidationSensor(BaseSensorOperator):
    """
    Sensor for waiting for data validation to complete in Nessi.dev.

    This sensor waits for a data validation to complete and can be configured
    to fail if the validation fails.

    :param validation_id: ID of the validation to wait for
    :type validation_id: str
    :param fail_on_validation_failure: Whether to fail the task if validation fails
    :type fail_on_validation_failure: bool
    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
    :param xcom_push_validation_results: Whether to push validation results to XCom
    :type xcom_push_validation_results: bool
    """

    template_fields: Sequence[str] = ('validation_id',)
    ui_color = '#ff7f50'  # Coral color for validation

    @apply_defaults
    def __init__(
        self,
        *,
        validation_id: str,
        fail_on_validation_failure: bool = True,
        nessi_conn_id: str = 'nessi_default',
        xcom_push_validation_results: bool = True,
        **kwargs,
    ) -> None:
        super().__init__(**kwargs)
        self.validation_id = validation_id
        self.fail_on_validation_failure = fail_on_validation_failure
        self.nessi_conn_id = nessi_conn_id
        self.xcom_push_validation_results = xcom_push_validation_results
        self.results = None

    def poke(self, context: Dict[str, Any]) -> bool:
        """
        Check if the data validation has completed.

        :param context: Airflow context
        :return: True if the validation has completed, False otherwise
        """
        self.log.info("Checking status of Nessi data validation: %s", self.validation_id)
        
        # Initialize hook
        hook = NessiHook(nessi_conn_id=self.nessi_conn_id)
        
        # Get results
        results = hook.get_validation_results(self.validation_id)
        status = results.get('status')
        
        if status == 'completed':
            self.log.info("Validation completed: %s", self.validation_id)
            self.results = results
            return True
        elif status == 'failed':
            raise AirflowException(f"Validation failed: {results.get('error')}")
        elif status == 'running':
            self.log.info("Validation still running, waiting...")
            return False
        else:
            raise AirflowException(f"Unknown validation status: {status}")
    
    def execute(self, context: Dict[str, Any]) -> Dict[str, Any]:
        """
        Execute the sensor after the poke method returns True.

        :param context: Airflow context
        :return: Validation results
        """
        # Check validation results
        if self.fail_on_validation_failure:
            rule_results = self.results.get('rule_results', [])
            failed_rules = [r for r in rule_results if r.get('status') == 'failed']
            if failed_rules:
                failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
                raise AirflowException(
                    f"Validation has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
                )
        
        # Log validation summary
        rule_count = len(self.results.get('rule_results', []))
        passed_count = len([r for r in self.results.get('rule_results', []) if r.get('status') == 'passed'])
        self.log.info(
            "Validation completed: %d/%d rules passed",
            passed_count,
            rule_count,
        )
        
        # Push results to XCom if requested
        if self.xcom_push_validation_results:
            return self.results
        
        return None
