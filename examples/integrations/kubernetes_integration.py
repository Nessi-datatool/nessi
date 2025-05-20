#!/usr/bin/env python3
"""
Kubernetes Integration Example for Nessi.dev

This example demonstrates how to use the Nessi Kubernetes integration
to run data quality checks, profiling, and validation as Kubernetes jobs.

Prerequisites:
- Kubernetes cluster with access configured via kubeconfig
- Nessi.dev instance accessible from the Kubernetes cluster
- nessi-k8s package installed: pip install nessi-k8s

Usage:
python kubernetes_integration.py --table my_table --namespace nessi
"""

import argparse
import json
import logging
import os
import sys
import time
from typing import Dict, List, Any, Optional

from nessi_k8s.client import NessiK8sClient
from nessi_k8s.config import NessiK8sConfig

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger("nessi-k8s-example")


def parse_args() -> argparse.Namespace:
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(description="Run Nessi operations in Kubernetes")
    parser.add_argument("--table", required=True, help="Table name to check")
    parser.add_argument("--namespace", default="default", help="Kubernetes namespace")
    parser.add_argument("--config", help="Path to configuration file")
    parser.add_argument("--rules", help="Path to rules JSON file")
    parser.add_argument(
        "--operation",
        choices=["quality", "profile", "validation"],
        default="quality",
        help="Operation to perform",
    )
    parser.add_argument(
        "--wait", action="store_true", help="Wait for job completion"
    )
    return parser.parse_args()


def load_rules(rules_file: Optional[str]) -> List[Dict[str, Any]]:
    """Load rules from a JSON file or use default rules."""
    if rules_file and os.path.exists(rules_file):
        with open(rules_file, "r") as f:
            return json.load(f)
    
    # Default rules if no file is provided
    return [
        {
            "name": "not_null_check",
            "description": "Check that id column is not null",
            "rule_type": "not_null",
            "column": "id",
        },
        {
            "name": "unique_check",
            "description": "Check that id column is unique",
            "rule_type": "unique",
            "column": "id",
        },
    ]


def run_quality_check(
    client: NessiK8sClient, table_name: str, rules: List[Dict[str, Any]], wait: bool
) -> Dict[str, Any]:
    """Run a data quality check as a Kubernetes job."""
    logger.info(f"Running quality check on table {table_name}")
    
    result = client.run_quality_check(
        table_name=table_name,
        rules=rules,
        profile=True,
        wait_for_completion=wait,
    )
    
    if wait:
        logger.info(f"Quality check completed with status: {result.get('status', 'unknown')}")
        if result.get("passed", False):
            logger.info("All quality checks passed!")
        else:
            logger.warning("Some quality checks failed!")
    else:
        logger.info(f"Quality check job submitted: {result.get('job_name', 'unknown')}")
    
    return result


def run_profile(
    client: NessiK8sClient, table_name: str, wait: bool
) -> Dict[str, Any]:
    """Run a data profile as a Kubernetes job."""
    logger.info(f"Running profile on table {table_name}")
    
    result = client.run_profile(
        table_name=table_name,
        columns=None,  # Profile all columns
        wait_for_completion=wait,
    )
    
    if wait:
        logger.info(f"Profile completed with status: {result.get('status', 'unknown')}")
    else:
        logger.info(f"Profile job submitted: {result.get('job_name', 'unknown')}")
    
    return result


def run_validation(
    client: NessiK8sClient, table_name: str, rules: List[Dict[str, Any]], wait: bool
) -> Dict[str, Any]:
    """Run a validation as a Kubernetes job."""
    logger.info(f"Running validation on table {table_name}")
    
    result = client.run_validation(
        table_name=table_name,
        rules=rules,
        wait_for_completion=wait,
    )
    
    if wait:
        logger.info(f"Validation completed with status: {result.get('status', 'unknown')}")
        if result.get("passed", False):
            logger.info("All validation rules passed!")
        else:
            logger.warning("Some validation rules failed!")
    else:
        logger.info(f"Validation job submitted: {result.get('job_name', 'unknown')}")
    
    return result


def main() -> int:
    """Main entry point."""
    args = parse_args()
    
    try:
        # Load configuration
        config = NessiK8sConfig(config_file=args.config)
        
        # Override namespace if provided
        if args.namespace:
            config.namespace = args.namespace
        
        # Create client
        client = NessiK8sClient(config=config)
        
        # Load rules
        rules = load_rules(args.rules)
        
        # Run the requested operation
        if args.operation == "quality":
            result = run_quality_check(client, args.table, rules, args.wait)
        elif args.operation == "profile":
            result = run_profile(client, args.table, args.wait)
        elif args.operation == "validation":
            result = run_validation(client, args.table, rules, args.wait)
        
        # Output the result as JSON
        print(json.dumps(result, indent=2))
        
        return 0
    
    except Exception as e:
        logger.error(f"Error: {e}", exc_info=True)
        return 1


if __name__ == "__main__":
    sys.exit(main())
