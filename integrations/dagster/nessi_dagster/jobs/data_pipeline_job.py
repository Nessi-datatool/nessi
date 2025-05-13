"""
Nessi.dev data pipeline job for Dagster.

This module provides a job for running a complete data pipeline with Nessi.dev
integration for data quality, validation, and lineage tracking.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import job, op, In, Out, Nothing

from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
from nessi_dagster.ops.validation_ops import run_validation, wait_for_validation_results
from nessi_dagster.ops.lineage_ops import get_lineage


@op(
    name="process_data",
    description="Process data based on quality results",
    ins={
        "input_table": In(str, description="Name of the input table"),
        "output_table": In(str, description="Name of the output table"),
        "quality_results": In(Dict[str, Any], description="Quality check results"),
    },
    out=Out(Dict[str, Any], description="Processing results"),
)
def process_data(
    context,
    input_table: str,
    output_table: str,
    quality_results: Dict[str, Any],
) -> Dict[str, Any]:
    """
    Process data based on quality results.

    Args:
        context: Dagster execution context
        input_table: Name of the input table
        output_table: Name of the output table
        quality_results: Quality check results

    Returns:
        Processing results
    """
    context.log.info(f"Processing data from {input_table} to {output_table}")
    
    # Get quality score
    quality_score = quality_results.get("quality_score", 0)
    context.log.info(f"Input table quality score: {quality_score}")
    
    # In a real implementation, this would process the data
    # For now, just return a dummy result
    return {
        "input_table": input_table,
        "output_table": output_table,
        "processed_at": context.get_run_time_str(),
        "quality_score": quality_score,
        "status": "completed",
    }


@job(
    name="nessi_data_pipeline_job",
    description="Run a complete data pipeline with Nessi.dev integration",
    resource_defs={
        "nessi": {
            "config": {
                "api_host": {"env": "NESSI_API_HOST"},
                "api_key": {"env": "NESSI_API_KEY"},
                "api_secret": {"env": "NESSI_API_SECRET"},
                "timeout": 300,
            }
        }
    },
)
def data_pipeline_job():
    """
    Job for running a complete data pipeline with Nessi.dev integration.

    This job runs a data pipeline with quality checks, validation, and lineage tracking.
    It can be configured to fail if the quality check or validation fails.
    """
    # Get input lineage
    input_lineage = get_lineage()
    
    # Run quality check on input table
    quality_response = run_quality_check()
    quality_results = wait_for_quality_results(check_id=quality_response["check_id"])
    
    # Process data
    input_table = quality_response["table_name"]
    output_table = f"{input_table}_processed"
    
    processing_results = process_data(
        input_table=input_table,
        output_table=output_table,
        quality_results=quality_results,
    )
    
    # Run validation on output table
    validation_response = run_validation(table_name=output_table)
    validation_results = wait_for_validation_results(
        validation_id=validation_response["validation_id"]
    )
    
    # Get output lineage
    output_lineage = get_lineage(table_name=output_table)
    
    return {
        "input_table": input_table,
        "output_table": output_table,
        "quality_results": quality_results,
        "processing_results": processing_results,
        "validation_results": validation_results,
        "input_lineage": input_lineage,
        "output_lineage": output_lineage,
    }
