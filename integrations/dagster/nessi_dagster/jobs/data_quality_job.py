"""
Nessi.dev data quality job for Dagster.

This module provides a job for running data quality checks with Nessi.dev.
"""

from typing import Dict, Any, Optional, List
import logging

from dagster import job, config_mapping, In, Out, Nothing

from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
from nessi_dagster.ops.profile_ops import run_profile, wait_for_profile


@job(
    name="nessi_data_quality_job",
    description="Run a data quality check on a table",
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
def data_quality_job():
    """
    Job for running data quality checks with Nessi.dev.

    This job runs a data quality check on a table and waits for the results.
    It can be configured to fail if the quality check fails.
    """
    # Run quality check
    quality_response = run_quality_check()
    
    # Wait for quality results
    quality_results = wait_for_quality_results(check_id=quality_response["check_id"])
    
    return quality_results
