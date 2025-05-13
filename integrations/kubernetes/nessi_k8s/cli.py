#!/usr/bin/env python3
"""
Command-line interface for Nessi Kubernetes integration.

This module provides a command-line interface for running Nessi operations
in Kubernetes.
"""

import os
import sys
import json
import argparse
import logging
from typing import Dict, Any, Optional, List, Union

from nessi_k8s.client import NessiK8sClient
from nessi_k8s.config import create_default_config

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


def parse_args():
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(
        description="Nessi Kubernetes CLI",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    
    # Common arguments
    parser.add_argument(
        "--config",
        help="Path to configuration file",
    )
    parser.add_argument(
        "--namespace",
        help="Kubernetes namespace to run the job in",
    )
    parser.add_argument(
        "--image",
        help="Docker image to use for the job",
    )
    parser.add_argument(
        "--service-account",
        help="Service account to use for the job",
    )
    parser.add_argument(
        "--api-host",
        help="Nessi API host URL",
    )
    parser.add_argument(
        "--api-key",
        help="Nessi API key",
    )
    parser.add_argument(
        "--api-secret",
        help="Nessi API secret",
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Enable verbose logging",
    )
    
    # Subparsers for commands
    subparsers = parser.add_subparsers(dest="command", help="Command to run")
    
    # Create config command
    create_config_parser = subparsers.add_parser(
        "create-config",
        help="Create a default configuration file",
    )
    create_config_parser.add_argument(
        "--output",
        required=True,
        help="Path to output configuration file",
    )
    create_config_parser.add_argument(
        "--format",
        choices=["yaml", "json"],
        default="yaml",
        help="Configuration file format",
    )
    
    # Quality check command
    quality_check_parser = subparsers.add_parser(
        "quality-check",
        help="Run a data quality check",
    )
    quality_check_parser.add_argument(
        "--table",
        required=True,
        help="Name or path of the table to check",
    )
    quality_check_parser.add_argument(
        "--rules",
        help="Path to JSON file containing quality rules",
    )
    quality_check_parser.add_argument(
        "--profile",
        action="store_true",
        default=True,
        help="Generate a profile",
    )
    quality_check_parser.add_argument(
        "--output-format",
        choices=["json", "yaml", "text"],
        default="json",
        help="Output format",
    )
    quality_check_parser.add_argument(
        "--output-path",
        help="Path to write output to",
    )
    quality_check_parser.add_argument(
        "--wait",
        action="store_true",
        default=True,
        help="Wait for job completion",
    )
    quality_check_parser.add_argument(
        "--timeout",
        type=int,
        default=600,
        help="Maximum time to wait in seconds",
    )
    quality_check_parser.add_argument(
        "--job-name",
        help="Name for the job (generated if not provided)",
    )
    
    # Profile command
    profile_parser = subparsers.add_parser(
        "profile",
        help="Run a data profile",
    )
    profile_parser.add_argument(
        "--table",
        required=True,
        help="Name or path of the table to profile",
    )
    profile_parser.add_argument(
        "--columns",
        help="Comma-separated list of columns to profile",
    )
    profile_parser.add_argument(
        "--sample-size",
        type=int,
        help="Number of rows to sample",
    )
    profile_parser.add_argument(
        "--output-format",
        choices=["json", "yaml", "text"],
        default="json",
        help="Output format",
    )
    profile_parser.add_argument(
        "--output-path",
        help="Path to write output to",
    )
    profile_parser.add_argument(
        "--wait",
        action="store_true",
        default=True,
        help="Wait for job completion",
    )
    profile_parser.add_argument(
        "--timeout",
        type=int,
        default=600,
        help="Maximum time to wait in seconds",
    )
    profile_parser.add_argument(
        "--job-name",
        help="Name for the job (generated if not provided)",
    )
    
    # Validation command
    validation_parser = subparsers.add_parser(
        "validate",
        help="Run a data validation",
    )
    validation_parser.add_argument(
        "--table",
        required=True,
        help="Name or path of the table to validate",
    )
    validation_parser.add_argument(
        "--rules",
        required=True,
        help="Path to JSON file containing validation rules",
    )
    validation_parser.add_argument(
        "--output-format",
        choices=["json", "yaml", "text"],
        default="json",
        help="Output format",
    )
    validation_parser.add_argument(
        "--output-path",
        help="Path to write output to",
    )
    validation_parser.add_argument(
        "--wait",
        action="store_true",
        default=True,
        help="Wait for job completion",
    )
    validation_parser.add_argument(
        "--timeout",
        type=int,
        default=600,
        help="Maximum time to wait in seconds",
    )
    validation_parser.add_argument(
        "--job-name",
        help="Name for the job (generated if not provided)",
    )
    
    # Get results command
    get_results_parser = subparsers.add_parser(
        "get-results",
        help="Get results of a job",
    )
    get_results_parser.add_argument(
        "--job-name",
        required=True,
        help="Name of the job",
    )
    get_results_parser.add_argument(
        "--job-type",
        choices=["quality", "profile", "validation"],
        required=True,
        help="Type of job",
    )
    get_results_parser.add_argument(
        "--wait",
        action="store_true",
        default=True,
        help="Wait for job completion",
    )
    get_results_parser.add_argument(
        "--timeout",
        type=int,
        default=600,
        help="Maximum time to wait in seconds",
    )
    
    return parser.parse_args()


def load_rules(rules_file: str) -> List[Dict[str, Any]]:
    """
    Load rules from a JSON file.
    
    Args:
        rules_file: Path to JSON file containing rules
        
    Returns:
        List of rules
    """
    try:
        with open(rules_file, "r") as f:
            return json.load(f)
    except Exception as e:
        logger.error(f"Error loading rules from {rules_file}: {e}")
        return []


def main():
    """Run the CLI."""
    args = parse_args()
    
    # Set log level
    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)
    
    # Handle create-config command
    if args.command == "create-config":
        create_default_config(args.output, args.format)
        logger.info(f"Created default configuration file: {args.output}")
        return
    
    # Create client
    client_kwargs = {}
    if args.config:
        client_kwargs["config_file"] = args.config
    if args.namespace:
        client_kwargs["namespace"] = args.namespace
    if args.image:
        client_kwargs["image"] = args.image
    if args.service_account:
        client_kwargs["service_account_name"] = args.service_account
    if args.api_host:
        client_kwargs["api_host"] = args.api_host
    if args.api_key:
        client_kwargs["api_key"] = args.api_key
    if args.api_secret:
        client_kwargs["api_secret"] = args.api_secret
    
    client = NessiK8sClient(**client_kwargs)
    
    # Handle commands
    if args.command == "quality-check":
        # Load rules if provided
        rules = None
        if args.rules:
            rules = load_rules(args.rules)
            logger.info(f"Loaded {len(rules)} rules from {args.rules}")
        
        # Run quality check
        logger.info(f"Running quality check on table: {args.table}")
        result = client.run_quality_check(
            table_name=args.table,
            rules=rules,
            profile=args.profile,
            output_format=args.output_format,
            output_path=args.output_path,
            wait_for_completion=args.wait,
            timeout=args.timeout,
            job_name=args.job_name,
        )
        
        # Print result
        if args.wait:
            logger.info("Quality check completed")
            logger.info(f"Job status: {result}")
        else:
            logger.info(f"Quality check job created: {result['name']}")
            logger.info("Use the following command to check the job status:")
            logger.info(f"kubectl get job {result['name']} -n {client.config['namespace']}")
    
    elif args.command == "profile":
        # Parse columns
        columns = None
        if args.columns:
            columns = [col.strip() for col in args.columns.split(",")]
        
        # Run profile
        logger.info(f"Running profile on table: {args.table}")
        result = client.run_profile(
            table_name=args.table,
            columns=columns,
            sample_size=args.sample_size,
            output_format=args.output_format,
            output_path=args.output_path,
            wait_for_completion=args.wait,
            timeout=args.timeout,
            job_name=args.job_name,
        )
        
        # Print result
        if args.wait:
            logger.info("Profile completed")
            logger.info(f"Job status: {result}")
        else:
            logger.info(f"Profile job created: {result['name']}")
            logger.info("Use the following command to check the job status:")
            logger.info(f"kubectl get job {result['name']} -n {client.config['namespace']}")
    
    elif args.command == "validate":
        # Load rules
        rules = load_rules(args.rules)
        logger.info(f"Loaded {len(rules)} rules from {args.rules}")
        
        # Run validation
        logger.info(f"Running validation on table: {args.table}")
        result = client.run_validation(
            table_name=args.table,
            rules=rules,
            output_format=args.output_format,
            output_path=args.output_path,
            wait_for_completion=args.wait,
            timeout=args.timeout,
            job_name=args.job_name,
        )
        
        # Print result
        if args.wait:
            logger.info("Validation completed")
            logger.info(f"Job status: {result}")
        else:
            logger.info(f"Validation job created: {result['name']}")
            logger.info("Use the following command to check the job status:")
            logger.info(f"kubectl get job {result['name']} -n {client.config['namespace']}")
    
    elif args.command == "get-results":
        # Get results
        logger.info(f"Getting results for job: {args.job_name}")
        
        if args.job_type == "quality":
            result = client.get_quality_results(
                job_name=args.job_name,
                wait_for_completion=args.wait,
                timeout=args.timeout,
            )
        elif args.job_type == "profile":
            result = client.get_profile_results(
                job_name=args.job_name,
                wait_for_completion=args.wait,
                timeout=args.timeout,
            )
        elif args.job_type == "validation":
            result = client.get_validation_results(
                job_name=args.job_name,
                wait_for_completion=args.wait,
                timeout=args.timeout,
            )
        
        # Print result
        logger.info(f"Job status: {result}")
    
    else:
        logger.error("No command specified")
        return 1
    
    return 0


if __name__ == "__main__":
    sys.exit(main())
