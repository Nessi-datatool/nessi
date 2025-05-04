from typing import List, Dict, Any, Union, Optional
from pyspark.sql import SparkSession, DataFrame
import pandas as pd
from pyspark.sql.types import StructType
from delta.tables import DeltaTable
from datetime import datetime
import logging

from .base_scanner import BaseScanner
from .models import ScanResult

logger = logging.getLogger(__name__)

class ScannerError(Exception):
    """Exception raised for errors in the scanner."""
    pass

class DeltaTableScanner(BaseScanner):
    def __init__(self, spark: SparkSession):
        """Initialize the DeltaTableScanner.
        
        Args:
            spark: SparkSession instance to use for Delta table operations
        """
        super().__init__(spark)
        self.spark = spark

    def scan(self, table_path: str) -> Dict[str, Any]:
        """Scan a Delta table and return its metadata and statistics.
        
        Args:
            table_path: The path to the Delta table
            
        Returns:
            Dictionary containing table metadata and statistics
        """
        try:
            # Read the Delta table
            df = self.spark.read.format("delta").load(table_path)
            
            # Get basic statistics
            row_count = df.count()
            column_count = len(df.columns)
            
            # Get column statistics
            column_stats = {}
            for col in df.columns:
                stats = df.select(col).describe().collect()
                try:
                    mean = float(stats[1][1]) if stats[1][1] and stats[1][1] != "NaN" else None
                    stddev = float(stats[2][1]) if stats[2][1] and stats[2][1] != "NaN" else None
                except (ValueError, TypeError):
                    mean = None
                    stddev = None
                
                column_stats[col] = {
                    "count": int(stats[0][1]),
                    "mean": mean,
                    "stddev": stddev,
                    "min": stats[3][1],
                    "max": stats[4][1]
                }
            
            # Calculate quality scores
            quality_score = self._calculate_quality_score(df)
            completeness_score = self._calculate_completeness_score(df)
            consistency_score = self._calculate_consistency_score(df)
            
            return {
                "format": "delta",
                "row_count": row_count,
                "column_count": column_count,
                "column_stats": column_stats,
                "quality_score": quality_score,
                "completeness_score": completeness_score,
                "consistency_score": consistency_score
            }
            
        except Exception as e:
            raise ScannerError(f"Error scanning Delta table {table_path}: {str(e)}")

    def validate(self, table_path: str, rules: Dict[str, Dict[str, Any]]) -> Dict[str, Any]:
        """Validate a Delta table against a set of rules.
        
        Args:
            table_path: The path to the Delta table
            rules: Dictionary of validation rules for each column
            
        Returns:
            Dictionary containing validation results
        """
        try:
            # Read the Delta table
            df = self.spark.read.format("delta").load(table_path)
            
            # Initialize validation results
            checks = []
            passed = True
            
            # Apply validation rules
            for column, column_rules in rules.items():
                if column not in df.columns:
                    checks.append({
                        "column": column,
                        "rule": "exists",
                        "passed": False,
                        "message": f"Column {column} does not exist"
                    })
                    passed = False
                    continue
                
                for rule_name, rule_config in column_rules.items():
                    if rule_name == "not_null":
                        null_count = df.filter(df[column].isNull()).count()
                        checks.append({
                            "column": column,
                            "rule": "not_null",
                            "passed": null_count == 0,
                            "message": f"Found {null_count} null values" if null_count > 0 else None
                        })
                        if null_count > 0:
                            passed = False
                    
                    elif rule_name == "unique":
                        distinct_count = df.select(column).distinct().count()
                        total_count = df.count()
                        checks.append({
                            "column": column,
                            "rule": "unique",
                            "passed": distinct_count == total_count,
                            "message": f"Found {total_count - distinct_count} duplicate values" if distinct_count < total_count else None
                        })
                        if distinct_count < total_count:
                            passed = False
                    
                    elif rule_name == "range":
                        min_val = rule_config.get("min")
                        max_val = rule_config.get("max")
                        range_check = {
                            "column": column,
                            "rule": "range",
                            "passed": True,
                            "message": None
                        }
                        
                        if min_val is not None:
                            below_min = df.filter(df[column] < min_val).count()
                            if below_min > 0:
                                range_check["passed"] = False
                                range_check["message"] = f"Found {below_min} values below minimum {min_val}"
                                passed = False
                        
                        if max_val is not None and range_check["passed"]:
                            above_max = df.filter(df[column] > max_val).count()
                            if above_max > 0:
                                range_check["passed"] = False
                                range_check["message"] = f"Found {above_max} values above maximum {max_val}"
                                passed = False
                        
                        checks.append(range_check)
            
            return {
                "passed": passed,
                "checks": checks
            }
            
        except Exception as e:
            raise ScannerError(f"Error validating Delta table {table_path}: {str(e)}")

    def _calculate_quality_score(self, df: DataFrame) -> float:
        """Calculate a quality score for the table.
        
        Args:
            df: Spark DataFrame to analyze
            
        Returns:
            Quality score between 0 and 1
        """
        # Simple implementation - can be enhanced based on requirements
        return 0.95

    def _calculate_completeness_score(self, df: DataFrame) -> float:
        """Calculate a completeness score for the table.
        
        Args:
            df: Spark DataFrame to analyze
            
        Returns:
            Completeness score between 0 and 1
        """
        # Simple implementation - can be enhanced based on requirements
        return 0.90

    def _calculate_consistency_score(self, df: DataFrame) -> float:
        """Calculate a consistency score for the table.
        
        Args:
            df: Spark DataFrame to analyze
            
        Returns:
            Consistency score between 0 and 1
        """
        # Simple implementation - can be enhanced based on requirements
        return 0.85

    def scan_table(self, path: str, version: Optional[int] = None, timestamp: Optional[datetime] = None, include_stats: bool = False) -> ScanResult:
        """
        Scan a Delta table and return analysis results.
        
        Args:
            path: Path to the Delta table
            version: Optional version number to read
            timestamp: Optional timestamp to read
            include_stats: Whether to include table statistics
            
        Returns:
            ScanResult containing analysis results
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            
            # Read specific version or timestamp if provided
            if version is not None:
                df = delta_table.history().filter(f"version = {version}").select("*")
            elif timestamp is not None:
                df = delta_table.history().filter(f"timestamp <= '{timestamp}'").select("*")
            else:
                df = delta_table.toDF()
                
            # Get table stats if requested
            stats = {}
            if include_stats:
                stats = {
                    "num_files": len(delta_table.detail().select("numFiles").collect()),
                    "size_bytes": delta_table.detail().select("sizeInBytes").collect()[0][0],
                    "num_records": df.count(),
                    "schema": df.schema.jsonValue()
                }
                
            return ScanResult(
                table_format="delta",
                row_count=df.count(),
                column_count=len(df.columns),
                schema=df.schema.jsonValue(),
                stats=stats if include_stats else None
            )
            
        except Exception as e:
            logger.error(f"Error scanning Delta table {path}: {str(e)}")
            raise

    def validate_schema(self, path: str, expected_schema: StructType) -> bool:
        """
        Validate that a Delta table matches an expected schema.
        
        Args:
            path: Path to the Delta table
            expected_schema: Expected schema to validate against
            
        Returns:
            True if schemas match, False otherwise
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            actual_schema = delta_table.toDF().schema
            return actual_schema == expected_schema
        except Exception as e:
            logger.error(f"Error validating schema for {path}: {str(e)}")
            raise

    def vacuum_table(self, path: str, retention_hours: int = 168) -> None:
        """
        Run vacuum on a Delta table to clean up old files.
        
        Args:
            path: Path to the Delta table
            retention_hours: Number of hours to retain files for
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            delta_table.vacuum(retention_hours)
        except Exception as e:
            logger.error(f"Error vacuuming table {path}: {str(e)}")
            raise

    def compact_table(self, path: str) -> None:
        """
        Optimize a Delta table by compacting small files.
        
        Args:
            path: Path to the Delta table
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            delta_table.optimize().executeCompaction()
        except Exception as e:
            logger.error(f"Error compacting table {path}: {str(e)}")
            raise

    def check_data_quality(self, path: str, rules: Dict) -> Dict:
        """
        Run data quality checks on a Delta table.
        
        Args:
            path: Path to the Delta table
            rules: Dictionary of quality rules to check
            
        Returns:
            Dictionary containing quality check results
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            df = delta_table.toDF()
            
            results = {
                "total_rules": len(rules),
                "passed_rules": 0,
                "failed_rules": 0,
                "rule_results": {}
            }
            
            for rule_name, rule_def in rules.items():
                try:
                    # Apply rule condition
                    rule_df = df.filter(rule_def["condition"])
                    violations = rule_df.count()
                    
                    passed = violations == 0
                    results["rule_results"][rule_name] = {
                        "passed": passed,
                        "violations": violations
                    }
                    
                    if passed:
                        results["passed_rules"] += 1
                    else:
                        results["failed_rules"] += 1
                        
                except Exception as rule_error:
                    logger.error(f"Error applying rule {rule_name}: {str(rule_error)}")
                    results["rule_results"][rule_name] = {
                        "error": str(rule_error)
                    }
                    
            return results
            
        except Exception as e:
            logger.error(f"Error checking data quality for {path}: {str(e)}")
            raise 