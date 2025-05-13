"""
Nessi.dev quality operations for Dagster.

This module provides operations for running data quality checks with Nessi.dev.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import op, Out, In, Nothing


@op(
    name="run_quality_check",
    description="Run a data quality check on a table",
    required_resource_keys={"nessi"},
    ins={
        "table_name": In(str, description="Name of the table to check"),
        "rules": In(
            Optional[List[Dict[str, Any]]],
            description="List of quality rules to apply",
            default_value=None,
        ),
        "profile": In(
            bool,
            description="Whether to generate a profile",
            default_value=True,
        ),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for the quality check",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Quality check response"),
)
def run_quality_check(
    context,
    table_name: str,
    rules: Optional[List[Dict[str, Any]]] = None,
    profile: bool = True,
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Run a data quality check on a table.

    Args:
        context: Dagster execution context
        table_name: Name of the table to check
        rules: List of quality rules to apply
        profile: Whether to generate a profile
        timeout: Timeout in seconds for the quality check

    Returns:
        Quality check response
    """
    context.log.info(f"Running quality check on table: {table_name}")
    
    response = context.resources.nessi.client.run_quality_check(
        table_name=table_name,
        rules=rules,
        profile=profile,
        timeout=timeout,
    )
    
    context.log.info(f"Quality check started with ID: {response.get('check_id')}")
    return response


@op(
    name="wait_for_quality_results",
    description="Wait for quality check results",
    required_resource_keys={"nessi"},
    ins={
        "check_id": In(str, description="ID of the quality check"),
        "quality_threshold": In(
            Optional[float],
            description="Minimum quality score to consider the check successful",
            default_value=None,
        ),
        "fail_on_rule_failure": In(
            bool,
            description="Whether to fail the op if any rule fails",
            default_value=True,
        ),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for waiting",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Quality check results"),
)
def wait_for_quality_results(
    context,
    check_id: str,
    quality_threshold: Optional[float] = None,
    fail_on_rule_failure: bool = True,
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Wait for quality check results.

    Args:
        context: Dagster execution context
        check_id: ID of the quality check
        quality_threshold: Minimum quality score to consider the check successful
        fail_on_rule_failure: Whether to fail the op if any rule fails
        timeout: Timeout in seconds for waiting

    Returns:
        Quality check results
    """
    context.log.info(f"Waiting for quality check results: {check_id}")
    
    results = context.resources.nessi.client.wait_for_quality_results(
        check_id=check_id,
        timeout=timeout,
    )
    
    # Check quality score
    quality_score = results.get("quality_score", 0)
    context.log.info(f"Quality score: {quality_score}")
    
    if quality_threshold is not None and quality_score < quality_threshold:
        raise RuntimeError(
            f"Quality score {quality_score} is below threshold {quality_threshold}"
        )
    
    # Check rule failures
    if fail_on_rule_failure:
        rule_results = results.get("rule_results", [])
        failed_rules = [
            r.get("name") for r in rule_results if r.get("status") == "failed"
        ]
        
        if failed_rules:
            raise RuntimeError(f"The following rules failed: {', '.join(failed_rules)}")
    
    return results
