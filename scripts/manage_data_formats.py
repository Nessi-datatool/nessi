#!/usr/bin/env python3
import argparse
import sys
from pathlib import Path
from typing import Optional

from src.formats.data_handler import DataFormatManager
from src.config.format_config import FormatConfig
from pyspark.sql import SparkSession

def create_spark_session() -> SparkSession:
    """Create a Spark session for Delta operations."""
    return SparkSession.builder \
        .appName("NessiDataFormatManager") \
        .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
        .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
        .getOrCreate()

def main():
    parser = argparse.ArgumentParser(description="Manage data formats")
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")
    
    # Convert command
    convert_parser = subparsers.add_parser("convert", help="Convert between formats")
    convert_parser.add_argument("source_path", help="Source file/directory path")
    convert_parser.add_argument("source_format", help="Source format (delta/parquet/csv)")
    convert_parser.add_argument("target_path", help="Target file/directory path")
    convert_parser.add_argument("target_format", help="Target format (delta/parquet/csv)")
    
    # Schema command
    schema_parser = subparsers.add_parser("schema", help="Infer schema")
    schema_parser.add_argument("path", help="File/directory path")
    schema_parser.add_argument("format", help="Format (delta/parquet/csv)")
    
    # Optimize command
    optimize_parser = subparsers.add_parser("optimize", help="Optimize Delta table")
    optimize_parser.add_argument("path", help="Delta table path")
    
    args = parser.parse_args()
    
    # Initialize format manager
    config = FormatConfig()
    spark = create_spark_session() if args.command in ["convert", "optimize"] else None
    format_manager = DataFormatManager(spark)
    
    try:
        if args.command == "convert":
            # Read source data
            df = format_manager.read(args.source_path, args.source_format)
            
            # Write to target format
            format_manager.write(df, args.target_path, args.target_format)
            print(f"Successfully converted {args.source_path} to {args.target_format}")
        
        elif args.command == "schema":
            # Infer schema
            schema = format_manager.infer_schema(args.path, args.format)
            print(json.dumps(schema, indent=2))
        
        elif args.command == "optimize":
            if args.format != "delta":
                print("Optimization is only supported for Delta format", file=sys.stderr)
                sys.exit(1)
            
            # Optimize Delta table
            delta_table = DeltaTable.forPath(spark, args.path)
            delta_table.optimize().executeCompaction()
            print(f"Successfully optimized Delta table at {args.path}")
        
        else:
            parser.print_help()
            sys.exit(1)
    
    except Exception as e:
        print(f"Error: {str(e)}", file=sys.stderr)
        sys.exit(1)
    
    finally:
        if spark:
            spark.stop()

if __name__ == "__main__":
    main() 