"""
Nessi.dev data pipeline flow for Prefect.

This module provides a flow for running a complete data pipeline with Nessi.dev
integration for data quality, validation, and lineage tracking.
"""

from typing import Dict, Any, Optional, List, Union, Callable
import logging
from datetime import datetime

from prefect import flow, task, get_run_logger

from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results
from nessi_prefect.tasks.validation import run_validation, wait_for_validation_results
from nessi_prefect.tasks.lineage import get_lineage, visualize_lineage


@task(name="Process Data")
def process_data(
    input_table: str,
    output_table: str,
    quality_results: Dict[str, Any],
    process_fn: Optional[Callable] = None,
) -> Dict[str, Any]:
    """
    Process data based on quality results.

    Args:
        input_table: Name of the input table
        output_table: Name of the output table
        quality_results: Quality check results
        process_fn: Optional function to process the data

    Returns:
        Processing results
    """
    logger = get_run_logger()
    logger.info("Processing data from %s to %s", input_table, output_table)
    
    # Get quality score
    quality_score = quality_results.get("quality_score", 0)
    logger.info("Input table quality score: %.2f", quality_score)
    
    # Process data using custom function if provided
    if process_fn and callable(process_fn):
        logger.info("Running custom processing function")
        result = process_fn(input_table, output_table, quality_results)
        if result:
            return result
    
    # Default processing result
    return {
        "input_table": input_table,
        "output_table": output_table,
        "processed_at": datetime.now().isoformat(),
        "quality_score": quality_score,
        "status": "completed",
    }


@flow(name="Nessi Data Pipeline Flow")
def data_pipeline_flow(
    input_table: str,
    output_table: str,
    quality_rules: Optional[List[Dict[str, Any]]] = None,
    validation_rules: Optional[List[Dict[str, Any]]] = None,
    quality_threshold: Optional[float] = 0.8,
    fail_on_rule_failure: bool = True,
    capture_lineage: bool = True,
    process_fn: Optional[Callable] = None,
    timeout: Optional[int] = None,
    api_host: Optional[str] = None,
    api_key: Optional[str] = None,
    api_secret: Optional[str] = None,
    secret_block_name: Optional[str] = None,
) -> Dict[str, Any]:
    """
    Flow for running a complete data pipeline with Nessi.dev integration.

    This flow runs a data pipeline with quality checks, validation, and lineage tracking.
    It can be configured to fail if the quality check or validation fails.

    Args:
        input_table: Name of the input table
        output_table: Name of the output table
        quality_rules: List of quality rules to apply to the input table
        validation_rules: List of validation rules to apply to the output table
        quality_threshold: Minimum quality score to consider the check successful
        fail_on_rule_failure: Whether to fail the flow if any rule fails
        capture_lineage: Whether to capture lineage information
        process_fn: Optional function to process the data
        timeout: Timeout in seconds for operations
        api_host: Host URL for the Nessi.dev API
        api_key: API key for authentication
        api_secret: API secret for authentication (optional)
        secret_block_name: Name of the Secret block containing Nessi.dev credentials

    Returns:
        Pipeline results
    """
    logger = get_run_logger()
    logger.info("Starting Nessi data pipeline flow: %s -> %s", input_table, output_table)
    
    # Capture input lineage if enabled
    input_lineage = None
    if capture_lineage:
        logger.info("Capturing input lineage for table: %s", input_table)
        input_lineage = get_lineage(
            table_name=input_table,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
    
    # Run quality check on input table
    quality_response = None
    quality_results = None
    
    if quality_rules:
        logger.info("Running quality check on input table: %s", input_table)
        quality_response = run_quality_check(
            table_name=input_table,
            rules=quality_rules,
            profile=True,
            timeout=timeout,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
        
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
    
    # Process data
    processing_results = process_data(
        input_table=input_table,
        output_table=output_table,
        quality_results=quality_results or {},
        process_fn=process_fn,
    )
    
    # Run validation on output table
    validation_response = None
    validation_results = None
    
    if validation_rules:
        logger.info("Running validation on output table: %s", output_table)
        validation_response = run_validation(
            table_name=output_table,
            rules=validation_rules,
            timeout=timeout,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
        
        validation_results = wait_for_validation_results(
            validation_id=validation_response["validation_id"],
            fail_on_validation_failure=fail_on_rule_failure,
            timeout=timeout,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
    
    # Capture output lineage if enabled
    output_lineage = None
    if capture_lineage:
        logger.info("Capturing output lineage for table: %s", output_table)
        output_lineage = get_lineage(
            table_name=output_table,
            api_host=api_host,
            api_key=api_key,
            api_secret=api_secret,
            secret_block_name=secret_block_name,
        )
    
    # Compile pipeline results
    pipeline_results = {
        "input_table": input_table,
        "output_table": output_table,
        "completed_at": datetime.now().isoformat(),
        "processing_results": processing_results,
    }
    
    if quality_results:
        pipeline_results["quality_results"] = quality_results
    
    if validation_results:
        pipeline_results["validation_results"] = validation_results
    
    if input_lineage:
        pipeline_results["input_lineage"] = {
            "node_count": len(input_lineage.get("nodes", [])),
            "edge_count": len(input_lineage.get("edges", [])),
        }
    
    if output_lineage:
        pipeline_results["output_lineage"] = {
            "node_count": len(output_lineage.get("nodes", [])),
            "edge_count": len(output_lineage.get("edges", [])),
        }
    
    # Log summary
    logger.info(
        "Data pipeline flow completed: %s -> %s",
        input_table,
        output_table,
    )
    
    return pipeline_results
