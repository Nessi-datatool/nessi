import argparse
import json
import sys
import pandas as pd
import pyarrow.parquet as pq
from ydata_profiling import ProfileReport


def generate_profile(parquet_path, output_json=None, minimal=True, explorative=True, html_report=None):
    # Read Parquet as DataFrame
    table = pq.read_table(parquet_path)
    df = table.to_pandas()
    
    # Generate profile
    profile = ProfileReport(
        df,
        title=f"Data Profile: {parquet_path}",
        minimal=minimal,
        explorative=explorative
    )
    
    # Output summary as JSON
    summary = profile.get_description()
    if output_json:
        with open(output_json, 'w') as f:
            json.dump(summary, f, indent=2)
    else:
        json.dump(summary, sys.stdout, indent=2)
    
    # Optionally save HTML report
    if html_report:
        profile.to_file(html_report)


def main():
    parser = argparse.ArgumentParser(description='Generate a data profile from a Parquet file using ydata-profiling.')
    parser.add_argument('--input', required=True, help='Input Parquet file path')
    parser.add_argument('--output-json', help='Output JSON file for summary stats')
    parser.add_argument('--html', help='Output HTML report path (optional)')
    parser.add_argument('--full', action='store_true', help='Generate full (not minimal) profile')
    args = parser.parse_args()

    generate_profile(
        parquet_path=args.input,
        output_json=args.output_json,
        minimal=not args.full,
        explorative=True,
        html_report=args.html
    )

if __name__ == '__main__':
    main()
