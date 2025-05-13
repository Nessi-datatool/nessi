#!/usr/bin/env python3
"""
Example script demonstrating the multi-format support features in the Nessi client.
"""

import os
import sys
import tempfile
import pandas as pd
from datetime import datetime

# Add the parent directory to the path so we can import the nessi_client package
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from nessi_client import NessiClient
from nessi_client.format_models import FormatConfig


def create_sample_files():
    """Create sample files in different formats for testing."""
    # Create a temporary directory
    temp_dir = tempfile.mkdtemp()
    print(f"Created temporary directory: {temp_dir}")
    
    # Sample data
    data = {
        'id': [1, 2, 3, 4, 5],
        'name': ['Alice', 'Bob', 'Charlie', 'David', 'Eve'],
        'age': [25, 30, 35, 40, 45],
        'email': ['alice@example.com', 'bob@example.com', 'charlie@example.com', 'david@example.com', 'eve@example.com'],
        'active': [True, False, True, True, False],
        'created_at': [datetime.now() for _ in range(5)]
    }
    
    df = pd.DataFrame(data)
    
    # Create CSV file
    csv_path = os.path.join(temp_dir, 'sample.csv')
    df.to_csv(csv_path, index=False)
    print(f"Created CSV file: {csv_path}")
    
    # Create Parquet file
    parquet_path = os.path.join(temp_dir, 'sample.parquet')
    df.to_parquet(parquet_path, index=False)
    print(f"Created Parquet file: {parquet_path}")
    
    return temp_dir, csv_path, parquet_path


def main():
    # Create sample files
    temp_dir, csv_path, parquet_path = create_sample_files()
    
    try:
        # Initialize client with custom format configuration
        format_config = FormatConfig(
            date_formats=["%Y-%m-%d", "%Y/%m/%d", "%d-%m-%Y"],
            csv_delimiter=",",
            csv_has_header=True,
            max_rows_for_inference=500,
            min_confidence_threshold=0.8
        )
        
        client = NessiClient(format_config=format_config)
        
        # Detect formats
        print("\n=== Format Detection ===")
        
        csv_format = client.detect_format(csv_path)
        print(f"CSV Format: {csv_format.format} (confidence: {csv_format.confidence})")
        print(f"CSV Schema: {len(csv_format.schema.fields)} fields")
        
        parquet_format = client.detect_format(parquet_path)
        print(f"Parquet Format: {parquet_format.format} (confidence: {parquet_format.confidence})")
        print(f"Parquet Schema: {len(parquet_format.schema.fields)} fields")
        
        # Read data
        print("\n=== Reading Data ===")
        
        csv_data = client.read_data(csv_path)
        print(f"CSV Data: {len(csv_data.data)} rows, {len(csv_data.schema.fields)} columns")
        print(f"First row: {csv_data.data[0]}")
        
        parquet_data = client.read_data(parquet_path)
        print(f"Parquet Data: {len(parquet_data.data)} rows, {len(parquet_data.schema.fields)} columns")
        print(f"First row: {parquet_data.data[0]}")
        
        # Schema inference from data
        print("\n=== Schema Inference ===")
        
        # Infer schema from list of records
        records = [
            {"id": 1, "name": "Alice", "score": 95.5},
            {"id": 2, "name": "Bob", "score": 87.0},
            {"id": 3, "name": "Charlie", "score": 92.5}
        ]
        
        inferred_schema = client.infer_schema(records)
        print("Inferred Schema from Records:")
        for field in inferred_schema.fields:
            print(f"  - {field.name}: {field.data_type} (nullable: {field.nullable})")
        
        # Validate data against rules
        print("\n=== Data Validation ===")
        
        # This would normally validate against rules defined in the Nessi server
        # For demonstration, we'll just show the API call
        print("To validate data against rules:")
        print("results = client.validate_file(csv_path)")
        print("# or")
        print("results = client.validate_data(csv_data.data)")
        
    finally:
        # Clean up
        print(f"\nCleaning up temporary directory: {temp_dir}")
        import shutil
        shutil.rmtree(temp_dir)


if __name__ == "__main__":
    main()
