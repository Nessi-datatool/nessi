#!/usr/bin/env python3

"""Table scanner module.

This module provides functionality to scan and analyze tables in various formats
(Parquet, CSV, Delta) and generate profile information.
"""

import os
import logging
from datetime import datetime
from typing import Dict, Any, List, Optional
import pandas as pd
import numpy as np
from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, DoubleType, DateType
from ..spark_config import get_spark_session

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class TableScanner:
    """Class for scanning and profiling tables."""

    def __init__(self):
        """Initialize the table scanner."""
        self.spark = get_spark_session()
    
    def _get_table_format(self, table_path: str) -> str:
        """Determine the format of the table.
        
        Args:
            table_path (str): Path to the table.
            
        Returns:
            str: Format of the table ('parquet', 'csv', or 'delta').
            
        Raises:
            ValueError: If the table format cannot be determined.
        """
        if not os.path.exists(table_path):
            raise ValueError(f"Table path does not exist: {table_path}")
        
        if os.path.isfile(table_path):
            if table_path.lower().endswith('.csv'):
                return 'csv'
            elif table_path.lower().endswith('.parquet'):
                return 'parquet'
        elif os.path.isdir(table_path):
            # Check for Delta table
            if os.path.exists(os.path.join(table_path, '_delta_log')):
                return 'delta'
            # Check for Parquet files
            parquet_files = [f for f in os.listdir(table_path) if f.lower().endswith('.parquet')]
            if parquet_files:
                return 'parquet'
            # Check for CSV files
            csv_files = [f for f in os.listdir(table_path) if f.lower().endswith('.csv')]
            if csv_files:
                return 'csv'
        
        raise ValueError(f"Unsupported table format at path: {table_path}")
    
    def _read_table(self, table_path: str, table_format: str) -> pd.DataFrame:
        """Read a table into a pandas DataFrame.
        
        Args:
            table_path (str): Path to the table.
            table_format (str): Format of the table.
            
        Returns:
            pd.DataFrame: Table data.
            
        Raises:
            ValueError: If the table cannot be read.
        """
        try:
            if table_format == 'parquet':
                if os.path.isfile(table_path):
                    return pd.read_parquet(table_path)
                else:
                    # Read all Parquet files in directory
                    return pd.concat([
                        pd.read_parquet(os.path.join(table_path, f))
                        for f in os.listdir(table_path)
                        if f.lower().endswith('.parquet')
                    ])
            elif table_format == 'csv':
                if os.path.isfile(table_path):
                    return pd.read_csv(table_path)
                else:
                    # Read all CSV files in directory
                    return pd.concat([
                        pd.read_csv(os.path.join(table_path, f))
                        for f in os.listdir(table_path)
                        if f.lower().endswith('.csv')
                    ])
            elif table_format == 'delta':
                return self.spark.read.format("delta").load(table_path).toPandas()
            else:
                raise ValueError(f"Unsupported table format: {table_format}")
        except Exception as e:
            logger.error(f"Failed to read table: {str(e)}")
            raise ValueError(f"Failed to read table: {str(e)}")
    
    def scan_table(self, table_path: str) -> Dict[str, Any]:
        """Scan a table and generate profile information.
        
        Args:
            table_path (str): Path to the table.
            
        Returns:
            Dict[str, Any]: Profile information including:
                - file_info: Information about the table files
                - schema: Table schema information
                - row_count: Number of rows
                - column_stats: Statistics for each column
                
        Raises:
            ValueError: If the table cannot be scanned.
        """
        try:
            # Get table format
            table_format = self._get_table_format(table_path)
            
            # Read table
            df = self._read_table(table_path, table_format)
            
            # Get file information
            file_info = {
                'format': table_format,
                'path': table_path,
                'size': os.path.getsize(table_path) if os.path.isfile(table_path) else sum(
                    os.path.getsize(os.path.join(dirpath, filename))
                    for dirpath, _, filenames in os.walk(table_path)
                    for filename in filenames
                ),
                'last_modified': datetime.fromtimestamp(os.path.getmtime(table_path)).isoformat()
            }
            
            # Get schema information
            schema = [
                {
                    'name': col,
                    'type': str(df[col].dtype)
                }
                for col in df.columns
            ]
            
            # Get row count
            row_count = len(df)
            
            # Get column statistics
            column_stats = {}
            for col in df.columns:
                stats = {
                    'count': len(df[col]),
                    'null_count': df[col].isna().sum(),
                    'min': df[col].min() if df[col].dtype in ['int64', 'float64'] else None,
                    'max': df[col].max() if df[col].dtype in ['int64', 'float64'] else None,
                }
                
                if df[col].dtype in ['int64', 'float64']:
                    stats.update({
                        'mean': df[col].mean(),
                        'stddev': df[col].std()
                    })
                
                column_stats[col] = stats
            
            return {
                'file_info': file_info,
                'schema': schema,
                'row_count': row_count,
                'column_stats': column_stats
            }
        except Exception as e:
            logger.error(f"Failed to scan table: {str(e)}")
            raise ValueError(f"Failed to scan table: {str(e)}")
    
    def generate_report(self, scan_result: Dict[str, Any], output_path: str) -> str:
        """Generate a human-readable report from scan results.
        
        Args:
            scan_result (Dict[str, Any]): Results from scan_table().
            output_path (str): Path where to save the report.
            
        Returns:
            str: Path to the generated report.
        """
        try:
            report = []
            
            # Title
            report.append("# Table Scan Report")
            report.append("")
            
            # File information
            report.append("## File Information")
            file_info = scan_result['file_info']
            report.append(f"- Format: {file_info['format']}")
            report.append(f"- Path: {file_info['path']}")
            report.append(f"- Size: {file_info['size']} bytes")
            report.append(f"- Last Modified: {file_info['last_modified']}")
            report.append("")
            
            # Schema information
            report.append("## Schema")
            report.append("| Column | Type |")
            report.append("|--------|------|")
            for field in scan_result['schema']:
                report.append(f"| {field['name']} | {field['type']} |")
            report.append("")
            
            # Row count
            report.append("## Row Count")
            report.append(f"Total Rows: {scan_result['row_count']}")
            report.append("")
            
            # Column statistics
            report.append("## Column Statistics")
            for col, stats in scan_result['column_stats'].items():
                report.append(f"### {col}")
                report.append(f"- Count: {stats['count']}")
                report.append(f"- Null Count: {stats['null_count']}")
                if stats.get('min') is not None:
                    report.append(f"- Min: {stats['min']}")
                if stats.get('max') is not None:
                    report.append(f"- Max: {stats['max']}")
                if stats.get('mean') is not None:
                    report.append(f"- Mean: {stats['mean']}")
                if stats.get('stddev') is not None:
                    report.append(f"- Standard Deviation: {stats['stddev']}")
                report.append("")
            
            # Write report to file
            os.makedirs(os.path.dirname(output_path), exist_ok=True)
            with open(output_path, 'w') as f:
                f.write('\n'.join(report))
            
            return output_path
        except Exception as e:
            logger.error(f"Failed to generate report: {str(e)}")
            raise ValueError(f"Failed to generate report: {str(e)}")
