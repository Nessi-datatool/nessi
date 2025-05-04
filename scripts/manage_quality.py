#!/usr/bin/env python3
import argparse
import sys
from pathlib import Path
import json
import pandas as pd
from typing import Dict, Any

from src.quality.analyzer import DataQualityAnalyzer
from src.config.quality_config import QualityConfig

def load_data(file_path: str) -> pd.DataFrame:
    """Load data from various formats."""
    path = Path(file_path)
    if path.suffix == '.csv':
        return pd.read_csv(path)
    elif path.suffix == '.parquet':
        return pd.read_parquet(path)
    else:
        raise ValueError(f"Unsupported file format: {path.suffix}")

def save_results(results: Dict[str, Any], output_path: str):
    """Save analysis results to a file."""
    with open(output_path, 'w') as f:
        json.dump(results, f, indent=2)

def main():
    parser = argparse.ArgumentParser(description="Data Quality Management Tool")
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # Analyze command
    analyze_parser = subparsers.add_parser("analyze", help="Analyze data quality")
    analyze_parser.add_argument("input", help="Input data file path")
    analyze_parser.add_argument("--output", help="Output results file path")
    analyze_parser.add_argument("--config", help="Configuration file path")
    
    # Validate command
    validate_parser = subparsers.add_parser("validate", help="Validate data against rules")
    validate_parser.add_argument("input", help="Input data file path")
    validate_parser.add_argument("--rules", help="Rules file path")
    validate_parser.add_argument("--output", help="Output results file path")
    
    args = parser.parse_args()
    
    try:
        if args.command == "analyze":
            # Load data
            df = load_data(args.input)
            
            # Load configuration
            config = QualityConfig()
            if args.config:
                with open(args.config) as f:
                    config_dict = json.load(f)
                    config = QualityConfig.from_dict(config_dict)
            
            # Analyze data
            analyzer = DataQualityAnalyzer(config.to_dict())
            results = analyzer.analyze_dataframe(df)
            
            # Save results
            if args.output:
                save_results(results, args.output)
            else:
                print(json.dumps(results, indent=2))
        
        elif args.command == "validate":
            # Load data
            df = load_data(args.input)
            
            # Load rules
            if not args.rules:
                print("Error: Rules file path is required for validation")
                sys.exit(1)
            
            with open(args.rules) as f:
                rules = json.load(f)
            
            # Validate data
            analyzer = DataQualityAnalyzer()
            results = analyzer.validate_rules(df, rules)
            
            # Save results
            if args.output:
                save_results(results, args.output)
            else:
                print(json.dumps(results, indent=2))
        
        else:
            parser.print_help()
            sys.exit(1)
    
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main() 