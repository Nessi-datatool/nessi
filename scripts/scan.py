#!/usr/bin/env python3
import argparse
import sys
import json
from pathlib import Path
import logging
from typing import Dict, Any

from src.scanner import Scanner
from src.report import ReportGenerator
from src.config.quality_config import QualityConfig

logger = logging.getLogger(__name__)

def load_config(config_path: str) -> Dict[str, Any]:
    """Load configuration from a JSON file."""
    try:
        with open(config_path) as f:
            return json.load(f)
    except Exception as e:
        logger.error(f"Error loading config from {config_path}: {str(e)}")
        raise

def save_results(results: Dict[str, Any], output_path: str):
    """Save results to a JSON file."""
    try:
        with open(output_path, 'w') as f:
            json.dump(results, f, indent=2)
    except Exception as e:
        logger.error(f"Error saving results to {output_path}: {str(e)}")
        raise

def main():
    parser = argparse.ArgumentParser(description="Nessi Data Scanner")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Scan command
    scan_parser = subparsers.add_parser("scan", help="Scan data source")
    scan_parser.add_argument("path", help="Path to data source")
    scan_parser.add_argument("--config", help="Path to configuration file")
    scan_parser.add_argument("--output", help="Path to save results")
    scan_parser.add_argument("--format", choices=["html", "json", "text"], default="html", help="Report format")
    scan_parser.add_argument("--display", action="store_true", help="Display report in browser")
    
    # Validate command
    validate_parser = subparsers.add_parser("validate", help="Validate data source")
    validate_parser.add_argument("path", help="Path to data source")
    validate_parser.add_argument("--rules", required=True, help="Path to rules file")
    validate_parser.add_argument("--config", help="Path to configuration file")
    validate_parser.add_argument("--output", help="Path to save results")
    validate_parser.add_argument("--format", choices=["html", "json", "text"], default="html", help="Report format")
    validate_parser.add_argument("--display", action="store_true", help="Display report in browser")
    
    # Optimize command (Delta only)
    optimize_parser = subparsers.add_parser("optimize", help="Optimize Delta table")
    optimize_parser.add_argument("path", help="Path to Delta table")
    optimize_parser.add_argument("--output", help="Path to save results")
    optimize_parser.add_argument("--format", choices=["html", "json", "text"], default="html", help="Report format")
    optimize_parser.add_argument("--display", action="store_true", help="Display report in browser")
    
    # Vacuum command (Delta only)
    vacuum_parser = subparsers.add_parser("vacuum", help="Vacuum Delta table")
    vacuum_parser.add_argument("path", help="Path to Delta table")
    vacuum_parser.add_argument("--retention", type=int, default=168, help="Retention hours")
    vacuum_parser.add_argument("--output", help="Path to save results")
    vacuum_parser.add_argument("--format", choices=["html", "json", "text"], default="html", help="Report format")
    vacuum_parser.add_argument("--display", action="store_true", help="Display report in browser")
    
    args = parser.parse_args()
    
    try:
        # Load configuration
        config = None
        if args.config:
            config = load_config(args.config)
        
        # Initialize scanner and report generator
        scanner = Scanner(config)
        report_generator = ReportGenerator()
        
        if args.command == "scan":
            # Scan data source
            results = scanner.scan(args.path)
            
            # Generate report
            report = report_generator.generate_report(results, args.format)
            
            # Save or print results
            if args.output:
                report_generator.save_report(report, args.output)
            elif args.display:
                report_generator.display_report(report, args.format)
            else:
                print(report)
        
        elif args.command == "validate":
            # Load rules
            rules = load_config(args.rules)
            
            # Validate data source
            results = scanner.validate(args.path, rules)
            
            # Generate report
            report = report_generator.generate_report(results, args.format)
            
            # Save or print results
            if args.output:
                report_generator.save_report(report, args.output)
            elif args.display:
                report_generator.display_report(report, args.format)
            else:
                print(report)
        
        elif args.command == "optimize":
            # Optimize Delta table
            results = scanner.optimize_delta_table(args.path)
            
            # Generate report
            report = report_generator.generate_report(results, args.format)
            
            # Save or print results
            if args.output:
                report_generator.save_report(report, args.output)
            elif args.display:
                report_generator.display_report(report, args.format)
            else:
                print(report)
        
        elif args.command == "vacuum":
            # Vacuum Delta table
            results = scanner.vacuum_delta_table(args.path, args.retention)
            
            # Generate report
            report = report_generator.generate_report(results, args.format)
            
            # Save or print results
            if args.output:
                report_generator.save_report(report, args.output)
            elif args.display:
                report_generator.display_report(report, args.format)
            else:
                print(report)
        
        else:
            parser.print_help()
            sys.exit(1)
    
    except Exception as e:
        logger.error(str(e))
        sys.exit(1)

if __name__ == "__main__":
    main() 