"""
Nessi.dev validation tasks for Prefect.

This module provides tasks for running data validation with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union
import logging

from prefect import task

from nessi_prefect.client import NessiClient

logger = logging.getLogger(__name__)


@task(name="Run Nessi Validation")
def run_validation(
    table_name: str,
    rules: List[Dict[str, Any]],
    timeout: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Run a validation on a table.

    Args:
        table_name: Name of the table to validate
        rules: List of validation rules
        timeout: Timeout in seconds for the validation
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Response containing the validation ID
    """
    logger.info("Running Nessi data validation on table: %s", table_name)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Run validation
    response = client.run_validation(
        table_name=table_name,
        rules=rules,
        timeout=timeout,
    )
    
    validation_id = response.get('validation_id')
    if not validation_id:
        raise RuntimeError(f"Failed to start validation: {response}")
    
    logger.info("Validation started with ID: %s", validation_id)
    return response


@task(name="Wait for Nessi Validation Results")
def wait_for_validation_results(
    validation_id: str,
    fail_on_validation_failure: bool = True,
    timeout: Optional[int] = None,
    poll_interval: int = 10,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Wait for validation results.

    Args:
        validation_id: ID of the validation
        fail_on_validation_failure: Whether to fail the task if validation fails
        timeout: Maximum time to wait in seconds
        poll_interval: Time between polls in seconds
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Validation results
    """
    logger.info("Waiting for Nessi data validation results: %s", validation_id)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Wait for results
    results = client.wait_for_validation_results(
        validation_id=validation_id,
        timeout=timeout,
        poll_interval=poll_interval,
    )
    
    # Check validation results
    if fail_on_validation_failure:
        rule_results = results.get('rule_results', [])
        failed_rules = [r for r in rule_results if r.get('status') == 'failed']
        if failed_rules:
            failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
            raise RuntimeError(
                f"Validation has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
            )
    
    # Log validation summary
    rule_count = len(results.get('rule_results', []))
    passed_count = len([r for r in results.get('rule_results', []) if r.get('status') == 'passed'])
    logger.info(
        "Validation completed for table: %d/%d rules passed",
        passed_count,
        rule_count,
    )
    
    return results
