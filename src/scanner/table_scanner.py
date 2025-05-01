"""Table scanning and profiling module.

This module provides functionality to scan and profile data tables in various formats
(Delta, Parquet, CSV) and generate detailed reports about the data.
"""

import os
import argparse
from typing import Dict, List, Any, Optional
from pyspark.sql import SparkSession, DataFrame
from jinja2 import Environment, FileSystemLoader, select_autoescape
import tempfile
import pandas as pd
import json
from delta.tables import DeltaTable
from pyspark.sql import functions as f
from datetime import datetime
import sys
from .models import TableScanResult, ColumnMetadata, ColumnStats
import logging
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, LongType, DoubleType, DateType, BooleanType, NumericType
from pyspark.sql.functions import col, count, mean, stddev, min, max, isnan, when

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class TableScanner:
    """Scanner for analyzing data tables in various formats.
    
    This class provides methods to scan and profile data tables, including:
    - Basic statistics (row count, column count)
    - Column-level statistics (min, max, mean, stddev, distinct count)
    - Data type information
    - Null value analysis
    """
    
    def __init__(self):
        """Initialize the TableScanner with a Spark session."""
        self.spark = self._get_or_create_spark_session()
    
    def _get_or_create_spark_session(self) -> SparkSession:
        """Get or create a Spark session with Delta Lake configuration.
        
        Returns:
            SparkSession: A configured Spark session with Delta Lake support.
            
        Raises:
            RuntimeError: If Spark session creation fails.
        """
        try:
            spark = SparkSession.builder \
                .appName("TableScanner") \
                .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
                .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
                .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0,io.delta:delta-storage:2.4.0") \
                .getOrCreate()
            return spark
        except Exception as e:
            logger.error(f"Failed to create Spark session: {str(e)}")
            raise RuntimeError(f"Failed to create Spark session: {str(e)}")
    
    def scan(self, table_path: str, auto_profile: bool = False) -> TableScanResult:
        """
        Scan a table and return its metadata
        
        Args:
            table_path: Path to the table
            auto_profile: Whether to automatically profile the data
            
        Returns:
            TableScanResult containing table metadata
        """
        # Detect table format
        format = self._get_table_format(table_path)
        if not format:
            raise ValueError(f"Could not detect table format at {table_path}")
        
        # Read the table
        df = self._read_table(table_path, format)
        
        # Get basic metadata
        columns = []
        stats = self._perform_quality_checks(df) if auto_profile else {}
        
        for field in df.schema.fields:
            col_stats = None
            if field.name in stats:
                col_stats = ColumnStats(
                    null_count=stats[field.name]['null_count'],
                    distinct_count=stats[field.name]['distinct_count'],
                    min_value=stats[field.name]['min_value'],
                    max_value=stats[field.name]['max_value'],
                    avg_value=stats[field.name]['avg_value'],
                    std_value=stats[field.name]['std_value']
                )
            
            columns.append(ColumnMetadata(
                name=field.name,
                type=field.dataType.simpleString(),
                stats=col_stats
            ))
        
        # Get table size
        size = 0
        try:
            table_path = str(table_path)
            for root, _, files in os.walk(table_path):
                for file in files:
                    if not file.startswith('.') and not file.startswith('_'):
                        size += os.path.getsize(os.path.join(root, file))
        except Exception as e:
            print(f"⚠️ Warning: Error getting table size: {str(e)}")
        
        # Get compression info
        compression = None
        if format == "CSV":
            for file in os.listdir(table_path):
                if file.endswith('.gz'):
                    compression = 'gzip'
                    break
                elif file.endswith('.bz2'):
                    compression = 'bzip2'
                    break
        
        return TableScanResult(
            format=format,
            columns=columns,
            row_count=df.count(),
            size=size,
            compression=compression,
            partition_columns=df.schema.names if format == "DELTA" else [],
            quality_checks=None
        )

    def _get_table_format(self, path):
        """Determine the format of the table at the given path.
        
        Args:
            path (str): Path to the table.
            
        Returns:
            str: The table format ('delta', 'parquet', or 'csv').
            
        Raises:
            ValueError: If the format cannot be determined.
        """
        try:
            # Convert Path to string if needed
            path = str(path)
            
            # Check if it's a Delta table
            if os.path.exists(os.path.join(path, '_delta_log')):
                return "DELTA"
            
            # Check if it's a single CSV file or Parquet file
            for file in os.listdir(path):
                if file.endswith('.csv') or file.endswith('.csv.gz') or file.endswith('.csv.bz2'):
                    return "CSV"
                elif file.endswith('.parquet'):
                    return "PARQUET"
                elif file == 'data.csv' and os.path.getsize(os.path.join(path, file)) == 0:
                    return "CSV"  # Empty CSV file
            
            return None
        except Exception as e:
            print(f"⚠️ Warning: Error detecting table format: {str(e)}")
            return None

    def _read_table(self, path, format):
        """
        Read a table using the appropriate format.
        
        Args:
            path: Path to the table
            format: Table format (DELTA, PARQUET, CSV)
            
        Returns:
            DataFrame containing the table data
        """
        try:
            # Convert Path to string if needed
            path = str(path)
            
            if format == "DELTA":
                return DeltaTable.forPath(self.spark, path).toDF()
            elif format == "PARQUET":
                return self.spark.read.parquet(path)
            elif format == "CSV":
                # Find the CSV file in the directory
                csv_file = None
                for file in os.listdir(path):
                    if file.endswith('.csv') or file.endswith('.csv.gz') or file.endswith('.csv.bz2'):
                        csv_file = os.path.join(path, file)
                        break
                
                if not csv_file:
                    raise ValueError(f"No CSV file found in {path}")
                
                return self.spark.read.option("header", "true").option("inferSchema", "true").csv(csv_file)
            else:
                raise ValueError(f"Unsupported format: {format}")
        except Exception as e:
            print(f"⚠️ Warning: Error reading table: {str(e)}")
            raise

    def _perform_quality_checks(self, df):
        """
        Perform quality checks on a DataFrame.
        
        Args:
            df: DataFrame to check
            
        Returns:
            Dictionary of column statistics
        """
        try:
            # Get basic statistics
            stats = {}
            for field in df.schema.fields:
                col = field.name
                col_stats = {
                    'null_count': df.filter(df[col].isNull()).count(),
                    'distinct_count': df.select(col).distinct().count()
                }
                
                # Get min/max/avg/std for numeric columns
                if field.dataType.simpleString() in ['int', 'bigint', 'double', 'float']:
                    agg_stats = df.agg(
                        f.min(col).alias('min'),
                        f.max(col).alias('max'),
                        f.avg(col).alias('avg'),
                        f.stddev(col).alias('std')
                    ).collect()[0]
                    
                    col_stats.update({
                        'min_value': agg_stats['min'],
                        'max_value': agg_stats['max'],
                        'avg_value': float(agg_stats['avg']) if agg_stats['avg'] is not None else None,
                        'std_value': float(agg_stats['std']) if agg_stats['std'] is not None else None
                    })
                else:
                    # For non-numeric columns, just get min/max
                    agg_stats = df.agg(
                        f.min(col).alias('min'),
                        f.max(col).alias('max')
                    ).collect()[0]
                    
                    col_stats.update({
                        'min_value': agg_stats['min'],
                        'max_value': agg_stats['max'],
                        'avg_value': None,
                        'std_value': None
                    })
                
                stats[col] = col_stats
            
            return stats
        except Exception as e:
            print(f"⚠️ Warning: Error performing quality checks: {str(e)}")
            return {}

    def _get_historical_stats(self, path, format):
        """
        Get historical statistics for a table.
        
        Args:
            path: Path to the table
            format: Format of the table (DELTA, PARQUET, CSV)
            
        Returns:
            Dictionary of historical statistics
        """
        try:
            # Read the table
            df = self._read_table(path, format)
            
            # Get basic statistics
            row_count = df.count()
            column_count = len(df.columns)
            
            # Get data size
            if format == "DELTA":
                size_bytes = sum(os.path.getsize(os.path.join(root, f))
                               for root, _, files in os.walk(path)
                               for f in files
                               if not f.startswith('.') and not f.endswith('.crc'))
            else:
                size_bytes = sum(os.path.getsize(os.path.join(path, f))
                               for f in os.listdir(path)
                               if os.path.isfile(os.path.join(path, f))
                               and not f.startswith('.')
                               and not f.endswith('.crc'))
            
            # Get schema
            schema = df.schema.json()
            
            # Get column statistics
            column_stats = {}
            for col in df.columns:
                stats = df.select(
                    f.count(col).alias("count"),
                    f.count(f.when(f.col(col).isNull(), True)).alias("null_count"),
                    f.countDistinct(col).alias("distinct_count")
                ).collect()[0]
                
                column_stats[col] = {
                    "count": stats["count"],
                    "null_count": stats["null_count"],
                    "null_percentage": round(stats["null_count"] * 100.0 / stats["count"], 2) if stats["count"] > 0 else 0,
                    "distinct_count": stats["distinct_count"],
                    "distinct_percentage": round(stats["distinct_count"] * 100.0 / stats["count"], 2) if stats["count"] > 0 else 0
                }
            
            return {
                "row_count": row_count,
                "column_count": column_count,
                "size_bytes": size_bytes,
                "schema": schema,
                "column_stats": column_stats,
                "last_modified": datetime.fromtimestamp(os.path.getmtime(path)).isoformat()
            }
            
        except Exception as e:
            print(f"⚠️ Warning: Error getting historical stats: {str(e)}")
            return {}

    def _get_partition_changes(self, path, format, num_versions=5):
        """Analyze partition-level changes across versions."""
        try:
            if format != 'DELTA':
                return None
            
            delta_table = DeltaTable.forPath(self.spark, path)
            history = delta_table.history().limit(num_versions).collect()
            
            if len(history) < 2:
                return None
            
            partition_changes = {}
            
            # Get partition columns
            partition_cols = self.spark.read.format("delta").load(path).schema.partitionColumns()
            
            if not partition_cols:
                return None
            
            # For each version pair, calculate changes
            for i in range(len(history) - 1):
                current_version = history[i].version
                previous_version = history[i - 1].version if i > 0 else None
                
                # Read current and previous versions
                current_df = self.spark.read.format("delta").option("versionAsOf", current_version).load(path)
                previous_df = self.spark.read.format("delta").option("versionAsOf", previous_version).load(path) if previous_version is not None else None
                
                # Group by partition and count changes
                if previous_df is not None:
                    # Get added/modified rows
                    added_rows = current_df.exceptAll(previous_df)
                    # Get deleted rows
                    deleted_rows = previous_df.exceptAll(current_df)
                    
                    # Combine changes
                    changes = added_rows.union(deleted_rows)
                    
                    # Count changes per partition
                    partition_counts = changes.groupBy(*partition_cols).count()
                    
                    # Convert to dictionary
                    for row in partition_counts.collect():
                        partition_key = "_".join(str(row[col]) for col in partition_cols)
                        partition_changes[f"{previous_version}_{current_version}"] = {
                            "partition": partition_key,
                            "change_count": row["count"],
                            "total_rows": current_df.count(),
                            "change_pct": (row["count"] / current_df.count()) * 100 if current_df.count() > 0 else 0
                        }
            
            return {
                "partition_columns": [col.name for col in partition_cols],
                "changes": partition_changes
            }
        except Exception as e:
            print(f"Warning: Could not analyze partition changes: {str(e)}")
            return None

    def _analyze_table(self, path, format, auto_profile=False):
        """
        Analyze a table and generate a report.
        
        Args:
            path: Path to the table
            format: Format of the table (DELTA, PARQUET, CSV)
            auto_profile: Whether to run auto-profiling
            
        Returns:
            Dictionary containing analysis results
        """
        try:
            # Get historical stats
            print("🔍 Analyzing historical statistics...")
            historical_stats = self._get_historical_stats(path, format)
            
            # Get data quality checks
            print("🔍 Performing quality checks...")
            quality_checks = self._perform_quality_checks(
                self._read_table(path, format)
            )
            
            # Get partition changes for Delta tables
            partition_changes = []
            if format == "DELTA":
                try:
                    print("🔍 Analyzing partition changes...")
                    partition_changes = self._get_partition_changes(path)
                except Exception as e:
                    print(f"⚠️ Warning: Could not get partition changes: {str(e)}")
            
            # Generate report
            print("📝 Generating report...")
            report_path = self._generate_report(
                path=path,
                format=format,
                historical_stats=historical_stats,
                quality_checks=quality_checks,
                partition_changes=partition_changes
            )
            
            return {
                'table_path': path,
                'format': format,
                'historical_stats': historical_stats,
                'quality_checks': quality_checks,
                'partition_changes': partition_changes,
                'report_path': report_path
            }
            
        except Exception as e:
            print(f"❌ Error analyzing table: {str(e)}")
            return None

    def _scan_project(self, project_path: str, auto_profile: bool = False) -> List[Dict]:
        """
        Scan all tables in a project directory.
        
        Args:
            project_path: Path to the project directory
            auto_profile: Whether to run auto-profiling
            
        Returns:
            List of scan results for each table
        """
        results = []
        
        # Skip these directories and files
        skip_patterns = {
            '_delta_log',  # Delta Lake transaction logs
            '_SUCCESS',    # Spark success marker
            '.DS_Store',   # macOS metadata
            '._',          # macOS resource forks
            '.git',        # Git metadata
            '.idea',       # IDE metadata
            '__pycache__', # Python cache
            'part-',       # Individual Parquet files
            '.crc',        # Checksum files
        }
        
        # Get all potential table directories
        table_dirs = set()
        for root, dirs, files in os.walk(project_path):
            # Skip hidden directories and metadata
            dirs[:] = [d for d in dirs if not any(d.startswith(p) for p in skip_patterns)]
            
            # Skip if this is a subdirectory of a table
            if any(root.startswith(os.path.join(d, '')) for d in table_dirs):
                continue
            
            # Check if this directory contains table files
            has_table_files = False
            for file in files:
                if not any(file.startswith(p) or file.endswith(p) for p in skip_patterns):
                    if file.endswith('.parquet') or file.endswith('.csv') or '_delta_log' in dirs:
                        has_table_files = True
                        break
            
            if has_table_files:
                table_dirs.add(root)
        
        # Process each table directory
        for table_path in table_dirs:
            try:
                # Try to detect table format
                table_format = self._get_table_format(table_path)
                if not table_format:
                    continue
                
                print(f"📊 Detected table format: {table_format}")
                
                # Analyze table
                result = self._analyze_table(table_path, table_format, auto_profile)
                if result:
                    results.append(result)
                
            except Exception as e:
                print(f"❌ Error scanning table: {str(e)}")
                continue
        
        return results

    def _render_report(self, context, output_html=None):
        """Render the report using the template and save as HTML."""
        # Get the absolute path to the logo
        logo_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', '..', 'Nessi_Logo.png'))
        
        # Add logo path to context
        context['logo_path'] = logo_path
        
        # Get the template directory
        template_dir = os.path.join(os.path.dirname(__file__), 'templates')
        
        # Create Jinja2 environment
        env = Environment(
            loader=FileSystemLoader(template_dir),
            autoescape=select_autoescape(['html', 'xml'])
        )
        
        # Load the template
        template = env.get_template('nessi_report_template.html')
        
        # Generate unique filename based on table name and timestamp if not provided
        if output_html is None:
            # Extract table name from path
            table_path = context.get('table_path', 'unknown_table')
            table_name = os.path.basename(table_path.rstrip('/'))
            if table_name == '':
                table_name = 'unknown_table'
            
            # Generate timestamp
            timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
            
            # Create reports directory
            reports_dir = "backend/data/reports"
            os.makedirs(reports_dir, exist_ok=True)
            
            # Generate unique filename
            output_html = os.path.join(reports_dir, f"nessi_report_{table_name}_{timestamp}.html")
        
        # Create reports directory if it doesn't exist
        os.makedirs(os.path.dirname(output_html), exist_ok=True)
        
        # Render HTML
        html_content = template.render(**context)
        with open(output_html, 'w') as f:
            f.write(html_content)
        
        print(f"✅ Report written to {output_html}")

    def _generate_report(self, path, format, historical_stats, quality_checks, partition_changes=None):
        """
        Generate an HTML report for the table scan results.
        
        Args:
            path: Path to the table
            format: Format of the table (DELTA, PARQUET, CSV)
            historical_stats: Historical statistics
            quality_checks: Data quality check results
            partition_changes: List of partition changes (for Delta tables)
            
        Returns:
            Path to the generated report
        """
        # Create reports directory if it doesn't exist
        reports_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), "reports")
        os.makedirs(reports_dir, exist_ok=True)
        
        # Generate report filename
        table_name = os.path.basename(path)
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        report_path = os.path.join(reports_dir, f"{table_name}_{timestamp}.html")
        
        # Generate HTML report
        html = f"""
        <html>
        <head>
            <title>Table Scan Report - {table_name}</title>
            <style>
                body {{ font-family: Arial, sans-serif; margin: 20px; }}
                h1, h2 {{ color: #333; }}
                table {{ border-collapse: collapse; width: 100%; margin: 10px 0; }}
                th, td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
                th {{ background-color: #f5f5f5; }}
                .metric {{ font-weight: bold; }}
                .warning {{ color: orange; }}
                .error {{ color: red; }}
            </style>
        </head>
        <body>
            <h1>Table Scan Report</h1>
            <p><b>Table:</b> {path}</p>
            <p><b>Format:</b> {format}</p>
            <p><b>Scan Date:</b> {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}</p>
            
            <h2>Historical Statistics</h2>
            <table>
                <tr><th>Metric</th><th>Value</th></tr>
        """
        
        # Add historical stats
        for metric, value in historical_stats.items():
            html += f"<tr><td class='metric'>{metric}</td><td>{value}</td></tr>"
        
        html += """
            </table>
            
            <h2>Data Quality Checks</h2>
            <table>
                <tr><th>Check</th><th>Result</th><th>Details</th></tr>
        """
        
        # Add quality checks
        for check, result in quality_checks.items():
            status = result.get("status", "")
            details = result.get("details", "")
            status_class = "error" if status == "FAILED" else "warning" if status == "WARNING" else ""
            html += f"<tr><td class='metric'>{check}</td><td class='{status_class}'>{status}</td><td>{details}</td></tr>"
        
        html += "</table>"
        
        # Add partition changes for Delta tables
        if partition_changes:
            html += """
                <h2>Partition Changes</h2>
                <table>
                    <tr><th>Version</th><th>Change Type</th><th>Partitions</th></tr>
            """
            for change in partition_changes:
                html += f"<tr><td>{change['version']}</td><td>{change['type']}</td><td>{change['partitions']}</td></tr>"
            html += "</table>"
        
        html += """
        </body>
        </html>
        """
        
        # Write report to file
        with open(report_path, "w") as f:
            f.write(html)
        
        return report_path

    def scan_table(self, table_path: str) -> Dict[str, Any]:
        """Scan a table and return its metadata and statistics.

        Args:
            table_path (str): Path to the table.

        Returns:
            Dict[str, Any]: Table information including:
                - schema: Table schema
                - row_count: Number of rows
                - column_stats: Statistics for each column
                - file_info: File format and path

        Raises:
            ValueError: If table_path is invalid.
            RuntimeError: If scanning fails.
        """
        if not table_path or not os.path.exists(table_path):
            raise ValueError(f"Invalid table path: {table_path}")

        try:
            spark = self._get_or_create_spark_session()

            # Determine file format and calculate size
            file_format = "unknown"
            file_size = 0
            if os.path.exists(os.path.join(table_path, "_delta_log")):
                file_format = "delta"
                df = spark.read.format("delta").load(table_path)
                file_size = sum(os.path.getsize(os.path.join(dirpath, filename))
                              for dirpath, _, filenames in os.walk(table_path)
                              for filename in filenames
                              if not filename.startswith('.') and not filename.endswith('.crc'))
            elif os.path.isfile(table_path) and table_path.endswith(".parquet"):
                file_format = "parquet"
                df = spark.read.parquet(table_path)
                file_size = os.path.getsize(table_path)
            elif os.path.isfile(table_path) and table_path.endswith(".csv"):
                file_format = "csv"
                df = spark.read.csv(table_path, header=True, inferSchema=True)
                file_size = os.path.getsize(table_path)
            elif os.path.isdir(table_path) and any(f.endswith(".parquet") for f in os.listdir(table_path)):
                file_format = "parquet"
                df = spark.read.parquet(table_path)
                file_size = sum(os.path.getsize(os.path.join(dirpath, filename))
                              for dirpath, _, filenames in os.walk(table_path)
                              for filename in filenames
                              if not filename.startswith('.') and not filename.endswith('.crc'))
            else:
                raise ValueError(f"Unsupported file format at {table_path}")

            # Get schema
            schema_fields = []
            for field in df.schema.fields:
                schema_fields.append({
                    'name': field.name,
                    'type': str(field.dataType),
                    'nullable': field.nullable
                })

            # Get row count
            row_count = df.count()

            # Get column statistics
            column_stats = {}
            for field in df.schema.fields:
                col_name = field.name
                stats = df.select(
                    count(col_name).alias("count"),
                    count(when(col(col_name).isNull(), True)).alias("null_count"),
                    min(col_name).alias("min"),
                    max(col_name).alias("max")
                ).collect()[0]

                # Add mean and stddev for numeric columns
                if isinstance(field.dataType, (IntegerType, LongType, DoubleType)):
                    agg_stats = df.select(
                        mean(col_name).alias("mean"),
                        stddev(col_name).alias("stddev")
                    ).collect()[0]
                    column_stats[col_name] = {
                        'count': stats['count'],
                        'null_count': stats['null_count'],
                        'min': stats['min'],
                        'max': stats['max'],
                        'mean': float(agg_stats['mean']) if agg_stats['mean'] is not None else None,
                        'stddev': float(agg_stats['stddev']) if agg_stats['stddev'] is not None else None
                    }
                else:
                    column_stats[col_name] = {
                        'count': stats['count'],
                        'null_count': stats['null_count'],
                        'min': stats['min'],
                        'max': stats['max']
                    }

            # Get file information
            file_info = {
                'format': file_format,
                'path': table_path,
                'size': file_size
            }

            result = {
                'schema': schema_fields,
                'row_count': row_count,
                'column_stats': column_stats,
                'file_info': file_info
            }

            logger.info(f"Successfully scanned table at {table_path}")
            return result
        except Exception as e:
            logger.error(f"Failed to scan table: {str(e)}")
            raise RuntimeError(f"Failed to scan table: {str(e)}")
    
    def generate_report(self, scan_result: Dict[str, Any], output_path: Optional[str] = None) -> str:
        """Generate a report from scan results.
        
        Args:
            scan_result (Dict[str, Any]): Results from scan_table.
            output_path (Optional[str]): Path to save the report. If None, returns the report as a string.
            
        Returns:
            str: The generated report.
            
        Raises:
            ValueError: If scan_result is empty or invalid.
        """
        if not scan_result:
            raise ValueError("Scan result cannot be empty")
        
        report = []
        report.append("# Table Analysis Report\n")
        
        # File Information
        report.append("## File Information\n")
        file_info = scan_result["file_info"]
        report.append(f"- Format: {file_info['format']}")
        report.append(f"- Path: {file_info['path']}")
        report.append(f"- Size: {file_info['size']} bytes\n")
        
        # Schema Information
        report.append("## Schema Information\n")
        for field in scan_result["schema"]:
            report.append(f"- {field['name']}: {field['type']}")
        report.append("")
        
        # Row Count
        report.append("## Row Count\n")
        report.append(f"Total rows: {scan_result['row_count']}\n")
        
        # Column Statistics
        report.append("## Column Statistics\n")
        for col, stats in scan_result["column_stats"].items():
            report.append(f"### {col}\n")
            report.append(f"- Count: {stats['count']}")
            report.append(f"- Null count: {stats['null_count']}")
            report.append(f"- Min: {stats['min']}")
            report.append(f"- Max: {stats['max']}")
            if 'mean' in stats and stats['mean'] is not None:
                report.append(f"- Mean: {stats['mean']:.2f}")
            if 'stddev' in stats and stats['stddev'] is not None:
                report.append(f"- Standard deviation: {stats['stddev']:.2f}")
            report.append("")
        
        # Join report lines
        report_text = "\n".join(report)
        
        # Save to file if path provided
        if output_path:
            try:
                os.makedirs(os.path.dirname(output_path), exist_ok=True)
                with open(output_path, "w") as f:
                    f.write(report_text)
                logger.info(f"Successfully saved report to {output_path}")
            except Exception as e:
                logger.error(f"Failed to save report: {str(e)}")
                raise RuntimeError(f"Failed to save report: {str(e)}")
        
        return report_text

def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Scan tables in a project directory.")
    parser.add_argument("--project-path", required=True, help="Path to the project directory")
    parser.add_argument("--auto-profile", action="store_true", help="Enable auto-profiling")
    args = parser.parse_args()
    
    try:
        # Scan project
        print(f"🔍 Scanning project at {args.project_path}")
        scanner = TableScanner()
        results = scanner._scan_project(args.project_path, args.auto_profile)
        
        # Print summary
        print("\n📊 Scan Summary:")
        print(f"Total tables scanned: {len(results)}")
        for result in results:
            print(f"\nTable: {result['table_path']}")
            print(f"Format: {result['format']}")
            print(f"Row count: {result['historical_stats'].get('row_count', 'N/A')}")
            print(f"Column count: {result['historical_stats'].get('column_count', 'N/A')}")
            print(f"Report: {result['report_path']}")
            
            # Print quality check summary
            failed_checks = sum(1 for check in result['quality_checks'].values() if check['status'] == 'FAILED')
            warning_checks = sum(1 for check in result['quality_checks'].values() if check['status'] == 'WARNING')
            if failed_checks > 0:
                print(f"❌ Failed checks: {failed_checks}")
            if warning_checks > 0:
                print(f"⚠️ Warning checks: {warning_checks}")
        
    except Exception as e:
        print(f"❌ Error: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    main()
