"""
Nessi.dev quality tasks for Prefect.

This module provides tasks for running data quality checks with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union
import logging

from prefect import task

from nessi_prefect.client import NessiClient

logger = logging.getLogger(__name__)


@task(name="Run Nessi Quality Check")
def run_quality_check(
    table_name: str,
    rules: Optional[List[Dict[str, Any]]] = None,
    profile: bool = True,
    timeout: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Run a data quality check on a table.

    Args:
        table_name: Name of the table to check
        rules: List of quality rules to apply
        profile: Whether to generate a profile
        timeout: Timeout in seconds for the quality check
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Response containing the quality check ID
    """
    logger.info("Running Nessi data quality check on table: %s", table_name)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Run quality check
    response = client.run_quality_check(
        table_name=table_name,
        rules=rules,
        profile=profile,
        timeout=timeout,
    )
    
    check_id = response.get('check_id')
    if not check_id:
        raise RuntimeError(f"Failed to start quality check: {response}")
    
    logger.info("Quality check started with ID: %s", check_id)
    return response


@task(name="Wait for Nessi Quality Results")
def wait_for_quality_results(
    check_id: str,
    quality_threshold: Optional[float] = None,
    fail_on_rule_failure: bool = True,
    timeout: Optional[int] = None,
    poll_interval: int = 10,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Wait for quality check results.

    Args:
        check_id: ID of the quality check
        quality_threshold: Minimum quality score to consider the check successful
        fail_on_rule_failure: Whether to fail the task if any rule fails
        timeout: Maximum time to wait in seconds
        poll_interval: Time between polls in seconds
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Quality check results
    """
    logger.info("Waiting for Nessi data quality results: %s", check_id)
    
    # Initialize client
    if secret_block_name:
        client = NessiClient.from_secret_block(secret_block_name)
    else:
        if not api_host or not api_key:
            raise ValueError("Either secret_block_name or api_host and api_key must be provided")
        client = NessiClient(api_host=api_host, api_key=api_key, api_secret=api_secret)
    
    # Wait for results
    results = client.wait_for_quality_results(
        check_id=check_id,
        timeout=timeout,
        poll_interval=poll_interval,
    )
    
    # Check quality score
    quality_score = results.get('quality_score')
    if quality_score is not None and quality_threshold is not None:
        if quality_score < quality_threshold:
            raise RuntimeError(
                f"Quality score {quality_score} is below threshold {quality_threshold}"
            )
        logger.info(
            "Quality score %s meets threshold %s",
            quality_score,
            quality_threshold,
        )
    
    # Check rule failures
    if fail_on_rule_failure:
        rule_results = results.get('rule_results', [])
        failed_rules = [r for r in rule_results if r.get('status') == 'failed']
        if failed_rules:
            failed_rule_names = [r.get('name', 'unnamed') for r in failed_rules]
            raise RuntimeError(
                f"Quality check has {len(failed_rules)} failed rules: {', '.join(failed_rule_names)}"
            )
    
    return results
