"""
Nessi.dev validation operations for Dagster.

This module provides operations for running data validation with Nessi.dev.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import op, Out, In, Nothing


@op(
    name="run_validation",
    description="Run a validation on a table",
    required_resource_keys={"nessi"},
    ins={
        "table_name": In(str, description="Name of the table to validate"),
        "rules": In(
            List[Dict[str, Any]],
            description="List of validation rules",
        ),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for the validation",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Validation response"),
)
def run_validation(
    context,
    table_name: str,
    rules: List[Dict[str, Any]],
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Run a validation on a table.

    Args:
        context: Dagster execution context
        table_name: Name of the table to validate
        rules: List of validation rules
        timeout: Timeout in seconds for the validation

    Returns:
        Validation response
    """
    context.log.info(f"Running validation on table: {table_name}")
    
    response = context.resources.nessi.client.run_validation(
        table_name=table_name,
        rules=rules,
        timeout=timeout,
    )
    
    context.log.info(f"Validation started with ID: {response.get('validation_id')}")
    return response


@op(
    name="wait_for_validation_results",
    description="Wait for validation results",
    required_resource_keys={"nessi"},
    ins={
        "validation_id": In(str, description="ID of the validation"),
        "fail_on_validation_failure": In(
            bool,
            description="Whether to fail the op if validation fails",
            default_value=True,
        ),
        "timeout": In(
            Optional[int],
            description="Timeout in seconds for waiting",
            default_value=None,
        ),
    },
    out=Out(Dict[str, Any], description="Validation results"),
)
def wait_for_validation_results(
    context,
    validation_id: str,
    fail_on_validation_failure: bool = True,
    timeout: Optional[int] = None,
) -> Dict[str, Any]:
    """
    Wait for validation results.

    Args:
        context: Dagster execution context
        validation_id: ID of the validation
        fail_on_validation_failure: Whether to fail the op if validation fails
        timeout: Timeout in seconds for waiting

    Returns:
        Validation results
    """
    context.log.info(f"Waiting for validation results: {validation_id}")
    
    results = context.resources.nessi.client.wait_for_validation_results(
        validation_id=validation_id,
        timeout=timeout,
    )
    
    # Check validation failures
    if fail_on_validation_failure:
        rule_results = results.get("rule_results", [])
        failed_rules = [
            r.get("name") for r in rule_results if r.get("status") == "failed"
        ]
        
        if failed_rules:
            raise RuntimeError(f"The following validation rules failed: {', '.join(failed_rules)}")
    
    # Log validation summary
    rule_count = len(results.get("rule_results", []))
    passed_count = len([r for r in results.get("rule_results", []) if r.get("status") == "passed"])
    
    context.log.info(
        f"Validation completed: {passed_count}/{rule_count} rules passed"
    )
    
    return results
