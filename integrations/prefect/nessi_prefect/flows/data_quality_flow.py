"""
Nessi.dev data quality flow for Prefect.

This module provides a flow for running data quality checks with Nessi.dev.
"""

from typing import Dict, Any, Optional, List, Union
import logging

from prefect import flow, get_run_logger

from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results
from nessi_prefect.tasks.profile import run_profile, wait_for_profile


@flow(name="Nessi Data Quality Flow")
def data_quality_flow(
    table_name: str,
    rules: Optional[List[Dict[str, Any]]] = None,
    profile: bool = True,
    quality_threshold: Optional[float] = None,
    fail_on_rule_failure: bool = True,
    timeout: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Flow for running data quality checks with Nessi.dev.

    This flow runs a data quality check on a table and waits for the results.
    It can be configured to fail if the quality check fails.

    Args:
        table_name: Name of the table to check
        rules: List of quality rules to apply
        profile: Whether to generate a profile
        quality_threshold: Minimum quality score to consider the check successful
        fail_on_rule_failure: Whether to fail the flow if any rule fails
        timeout: Timeout in seconds for the quality check
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Quality check results
    """
    logger = get_run_logger()
    logger.info("Starting Nessi data quality flow for table: %s", table_name)
    
    # Run quality check
    quality_response = run_quality_check(
        table_name=table_name,
        rules=rules,
        profile=profile,
        timeout=timeout,
        api_host=api_host,
        api_key=api_key,
        api_secret=api_secret,
        secret_block_name=secret_block_name,
    )
    
    # Wait for quality results
    quality_results = wait_for_quality_results(
        check_id=quality_response["check_id"],
        quality_threshold=quality_threshold,
        fail_on_rule_failure=fail_on_rule_failure,
        timeout=timeout,
        api_host=api_host,
        api_key=api_key,
        api_secret=api_secret,
        secret_block_name=secret_block_name,
    )
    
    # If profiling is enabled and not included in quality check
    if profile and not quality_results.get("profile"):
        logger.info("Running separate profile for table: %s", table_name)
        
        # Run profile
        profile_response = run_profile(
            table_name=table_name,
            timeout=timeout,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
        
        # Wait for profile results
        profile_results = wait_for_profile(
            profile_id=profile_response["profile_id"],
            timeout=timeout,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
        
        # Add profile to quality results
        quality_results["profile"] = profile_results
    
    # Log summary
    quality_score = quality_results.get("quality_score", 0)
    rule_count = len(quality_results.get("rule_results", []))
    passed_count = len([r for r in quality_results.get("rule_results", []) if r.get("status") == "passed"])
    
    logger.info(
        "Data quality flow completed for table %s: Quality score: %.2f, %d/%d rules passed",
        table_name,
        quality_score,
        passed_count,
        rule_count,
    )
    
    return quality_results
