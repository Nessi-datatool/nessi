#!/usr/bin/env python3
"""
Nessi Custom Script Integration Example

This example demonstrates how to integrate Nessi with custom Python scripts
to run data quality checks and process the results programmatically.
"""

import argparse
import json
import os
import subprocess
import sys
from datetime import datetime
from typing import Dict, List, Any, Optional

# Define colors for terminal output
class Colors:
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BLUE = '\033[94m'
    BOLD = '\033[1m'
    END = '\033[0m'

class NessiIntegration:
    def __init__(self, nessi_path: Optional[str] = None):
        """
        Initialize the Nessi integration.
        
        Args:
            nessi_path: Path to the Nessi executable. If None, assumes 'nessi' is in PATH.
        """
        self.nessi_path = nessi_path or 'nessi'
        self._check_nessi_installed()
    
    def _check_nessi_installed(self) -> None:
        """Check if Nessi is installed and accessible."""
        try:
            result = subprocess.run(
                [self.nessi_path, 'version'], 
                capture_output=True, 
                text=True
            )
            if result.returncode != 0:
                print(f"{Colors.RED}Error: Nessi is not installed or not accessible.{Colors.END}")
                print(f"Error message: {result.stderr}")
                sys.exit(1)
            
            print(f"{Colors.BLUE}Using Nessi: {result.stdout.strip()}{Colors.END}")
        except FileNotFoundError:
            print(f"{Colors.RED}Error: Nessi executable not found at '{self.nessi_path}'.{Colors.END}")
            print("Please install Nessi or provide the correct path.")
            sys.exit(1)
    
    def validate_table(self, table_path: str, config_path: Optional[str] = None) -> Dict[str, Any]:
        """
        Validate a Delta table using Nessi.
        
        Args:
            table_path: Path to the Delta table.
            config_path: Path to the Nessi configuration file.
            
        Returns:
            Dict containing the validation results.
        """
        cmd = [self.nessi_path, 'validate', table_path, '--format', 'json']
        
        if config_path:
            cmd.extend(['--config', config_path])
        
        print(f"{Colors.BLUE}Running: {' '.join(cmd)}{Colors.END}")
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, check=True)
            return json.loads(result.stdout)
        except subprocess.CalledProcessError as e:
            print(f"{Colors.RED}Error running Nessi validation: {e}{Colors.END}")
            print(f"Stderr: {e.stderr}")
            return {"error": e.stderr, "passed": False}
        except json.JSONDecodeError as e:
            print(f"{Colors.RED}Error parsing Nessi output: {e}{Colors.END}")
            print(f"Output: {result.stdout}")
            return {"error": "Failed to parse JSON output", "passed": False}
    
    def batch_validate(self, tables_dir: str, config_path: Optional[str] = None) -> Dict[str, Dict[str, Any]]:
        """
        Validate all Delta tables in a directory.
        
        Args:
            tables_dir: Directory containing Delta tables.
            config_path: Path to the Nessi configuration file.
            
        Returns:
            Dict mapping table names to validation results.
        """
        results = {}
        
        # Get all subdirectories in the tables directory
        tables = [os.path.join(tables_dir, d) for d in os.listdir(tables_dir) 
                 if os.path.isdir(os.path.join(tables_dir, d))]
        
        for table_path in tables:
            table_name = os.path.basename(table_path)
            print(f"{Colors.BOLD}Validating table: {table_name}{Colors.END}")
            results[table_name] = self.validate_table(table_path, config_path)
        
        return results
    
    def generate_report(self, results: Dict[str, Dict[str, Any]], output_path: str, format: str = 'html') -> None:
        """
        Generate a report from validation results.
        
        Args:
            results: Dict mapping table names to validation results.
            output_path: Path to save the report.
            format: Report format (html, json, csv, pdf).
        """
        # Save results to temporary JSON files
        temp_dir = os.path.join(os.getcwd(), 'temp_results')
        os.makedirs(temp_dir, exist_ok=True)
        
        json_files = []
        for table_name, result in results.items():
            json_path = os.path.join(temp_dir, f"{table_name}_result.json")
            with open(json_path, 'w') as f:
                json.dump(result, f)
            json_files.append(json_path)
        
        # Generate report using Nessi
        cmd = [
            self.nessi_path, 'report', 'generate',
            '--input', ','.join(json_files),
            '--output', output_path,
            '--format', format
        ]
        
        print(f"{Colors.BLUE}Running: {' '.join(cmd)}{Colors.END}")
        
        try:
            subprocess.run(cmd, check=True)
            print(f"{Colors.GREEN}Report generated successfully: {output_path}{Colors.END}")
        except subprocess.CalledProcessError as e:
            print(f"{Colors.RED}Error generating report: {e}{Colors.END}")
        
        # Clean up temporary files
        for file in json_files:
            os.remove(file)
        os.rmdir(temp_dir)
    
    def trigger_webhook(self, event: str, payload: Dict[str, Any]) -> bool:
        """
        Trigger a Nessi webhook.
        
        Args:
            event: Event type.
            payload: Event payload.
            
        Returns:
            True if webhook was triggered successfully, False otherwise.
        """
        cmd = [
            self.nessi_path, 'webhook', 'trigger',
            '--event', event,
            '--payload', json.dumps(payload)
        ]
        
        print(f"{Colors.BLUE}Running: {' '.join(cmd)}{Colors.END}")
        
        try:
            subprocess.run(cmd, check=True)
            print(f"{Colors.GREEN}Webhook triggered successfully for event: {event}{Colors.END}")
            return True
        except subprocess.CalledProcessError as e:
            print(f"{Colors.RED}Error triggering webhook: {e}{Colors.END}")
            return False

def print_results_summary(results: Dict[str, Dict[str, Any]]) -> None:
    """Print a summary of validation results."""
    print(f"\n{Colors.BOLD}Validation Results Summary:{Colors.END}")
    print("-" * 80)
    
    all_passed = True
    for table_name, result in results.items():
        passed = result.get('passed', False)
        all_passed = all_passed and passed
        
        status = f"{Colors.GREEN}PASSED{Colors.END}" if passed else f"{Colors.RED}FAILED{Colors.END}"
        print(f"Table: {table_name} - Status: {status}")
        
        if not passed and 'rules' in result:
            failed_rules = [rule for rule in result['rules'] if not rule.get('passed', False)]
            for rule in failed_rules:
                print(f"  - {Colors.YELLOW}{rule.get('name')}: {rule.get('message')}{Colors.END}")
    
    print("-" * 80)
    overall = f"{Colors.GREEN}PASSED{Colors.END}" if all_passed else f"{Colors.RED}FAILED{Colors.END}"
    print(f"Overall Status: {overall}")

def main():
    parser = argparse.ArgumentParser(description='Nessi Integration Script')
    parser.add_argument('--tables-dir', required=True, help='Directory containing Delta tables')
    parser.add_argument('--config', help='Path to Nessi configuration file')
    parser.add_argument('--output', default='report.html', help='Path to save the report')
    parser.add_argument('--format', default='html', choices=['html', 'json', 'csv', 'pdf'], help='Report format')
    parser.add_argument('--webhook', action='store_true', help='Trigger webhook with results')
    parser.add_argument('--nessi-path', help='Path to Nessi executable')
    
    args = parser.parse_args()
    
    # Initialize Nessi integration
    nessi = NessiIntegration(args.nessi_path)
    
    # Validate tables
    print(f"{Colors.BOLD}Starting validation of tables in: {args.tables_dir}{Colors.END}")
    results = nessi.batch_validate(args.tables_dir, args.config)
    
    # Print results summary
    print_results_summary(results)
    
    # Generate report
    nessi.generate_report(results, args.output, args.format)
    
    # Trigger webhook if requested
    if args.webhook:
        webhook_payload = {
            'tables': list(results.keys()),
            'results': results,
            'timestamp': datetime.now().isoformat(),
            'summary': {
                'total_tables': len(results),
                'passed_tables': sum(1 for r in results.values() if r.get('passed', False)),
                'failed_tables': sum(1 for r in results.values() if not r.get('passed', False))
            }
        }
        nessi.trigger_webhook('validation.complete', webhook_payload)
    
    # Exit with appropriate status code
    all_passed = all(result.get('passed', False) for result in results.values())
    sys.exit(0 if all_passed else 1)

if __name__ == '__main__':
    main()
