#!/usr/bin/env python3
"""
Example script for running Nessi data quality checks in Kubernetes.

This script demonstrates how to use the Nessi Kubernetes integration
to run data quality checks on a Delta Lake table.
"""

import os
import json
import argparse
import logging
from typing import Dict, Any, Optional, List

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
        description="Run Nessi data quality checks in Kubernetes"
    )
    parser.add_argument(
        "--table",
        required=True,
        help="Name or path of the table to check",
    )
    parser.add_argument(
        "--rules",
        help="Path to JSON file containing quality rules",
    )
    parser.add_argument(
        "--profile",
        action="store_true",
        default=True,
        help="Generate a profile (default: True)",
    )
    parser.add_argument(
        "--namespace",
        default="default",
        help="Kubernetes namespace to run the job in (default: default)",
    )
    parser.add_argument(
        "--image",
        default="nessi/nessi:latest",
        help="Docker image to use for the job (default: nessi/nessi:latest)",
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
        "--output-format",
        choices=["json", "yaml", "text"],
        default="json",
        help="Output format (default: json)",
    )
    parser.add_argument(
        "--output-path",
        help="Path to write output to",
    )
    parser.add_argument(
        "--wait",
        action="store_true",
        default=True,
        help="Wait for job completion (default: True)",
    )
    parser.add_argument(
        "--timeout",
        type=int,
        default=600,
        help="Maximum time to wait in seconds (default: 600)",
    )
    parser.add_argument(
        "--poll-interval",
        type=int,
        default=10,
        help="Time between polls in seconds (default: 10)",
    )
    parser.add_argument(
        "--job-name",
        help="Name for the job (generated if not provided)",
    )
    parser.add_argument(
        "--config",
        help="Path to configuration file",
    )
    parser.add_argument(
        "--create-config",
        help="Create a default configuration file and exit",
    )
    parser.add_argument(
        "--config-format",
        choices=["yaml", "json"],
        default="yaml",
        help="Configuration file format (default: yaml)",
    )
    
    return parser.parse_args()


def load_rules(rules_file: str) -> List[Dict[str, Any]]:
    """
    Load quality rules from a JSON file.
    
    Args:
        rules_file: Path to JSON file containing quality rules
        
    Returns:
        List of quality rules
    """
    try:
        with open(rules_file, "r") as f:
            return json.load(f)
    except Exception as e:
        logger.error(f"Error loading rules from {rules_file}: {e}")
        return []


def main():
    """Run the example script."""
    args = parse_args()
    
    # Create default configuration file if requested
    if args.create_config:
        create_default_config(args.create_config, args.config_format)
        logger.info(f"Created default configuration file: {args.create_config}")
        return
    
    # Load rules if provided
    rules = None
    if args.rules:
        rules = load_rules(args.rules)
        logger.info(f"Loaded {len(rules)} rules from {args.rules}")
    
    # Create client
    client = NessiK8sClient(
        config_file=args.config,
        namespace=args.namespace,
        image=args.image,
        service_account_name=args.service_account,
        api_host=args.api_host,
        api_key=args.api_key,
        api_secret=args.api_secret,
    )
    
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
        poll_interval=args.poll_interval,
        job_name=args.job_name,
    )
    
    # Print result
    if args.wait:
        logger.info("Quality check completed")
        logger.info(f"Job status: {result}")
    else:
        logger.info(f"Quality check job created: {result['name']}")
        logger.info("Use the following command to check the job status:")
        logger.info(f"kubectl get job {result['name']} -n {args.namespace}")


if __name__ == "__main__":
    main()
