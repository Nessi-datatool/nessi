from pyspark.sql import SparkSession, DataFrame
from typing import Dict, Any, Union, List, Optional
import logging
from pathlib import Path
import pandas as pd
from pyspark.sql.functions import col, count, countDistinct, when, mean, stddev, min, max

from nessi.scanner.base_scanner import BaseScanner
from nessi.scanner.models import ScanResult
from .data_quality_analyzer import DataQualityAnalyzer

class DeltaTableScanner(BaseScanner):
    def __init__(self, spark: SparkSession):
        super().__init__(spark)
        self.logger = logging.getLogger(__name__)
        self.quality_analyzer = DataQualityAnalyzer(spark)

    def _read_data(self, path: str) -> DataFrame:
        """Read data from a Delta table.
        
        Args:
            path: Path to the Delta table
            
        Returns:
            DataFrame containing the data
        """
        try:
            return self.spark.read.format("delta").load(path)
        except Exception as e:
            self.logger.error(f"Error reading Delta table {path}: {str(e)}")
            raise

    def scan(self, path: str) -> ScanResult:
        """
        Scan a Delta table and return scan results.
        
        Args:
            path (str): Path to the Delta table
            
        Returns:
            ScanResult: Results of the scan
        """
        # Read the Delta table
        df = self._read_data(path)
        
        # Get basic statistics
        stats = self._get_basic_stats(df)
        
        # Validate schema
        schema_validation = self._validate_schema(df)
        
        # Run quality checks
        quality_results = self.quality_analyzer.analyze(df)
        
        # Create scan result
        result = ScanResult(
            row_count=stats["row_count"],
            column_count=stats["column_count"],
            columns=stats["columns"],
            schema=schema_validation.fields,
            quality_results=quality_results
        )
        
        return result

    def validate(self, path: str, rules: List[Dict[str, Any]], **kwargs) -> Dict[str, Any]:
        """Validate a Delta table against rules."""
        try:
            df = self._read_data(path)
            validation_results = {
                "passed": True,
                "checks": []
            }
            
            for rule in rules:
                column = rule.get("column")
                if not column or column not in df.columns:
                    validation_results["checks"].append({
                        "column": column,
                        "status": "failed",
                        "error": "Column not found"
                    })
                    validation_results["passed"] = False
                    continue
                
                check_result = self._apply_rule(df, column, rule)
                validation_results["checks"].append(check_result)
                if not check_result["passed"]:
                    validation_results["passed"] = False
            
            return validation_results
        except Exception as e:
            self.logger.error(f"Error validating Delta table {path}: {str(e)}")
            raise

    def _get_delta_metadata(self, path: str) -> Dict[str, Any]:
        """Get Delta table metadata."""
        try:
            # Get Delta table history
            history = self.spark.sql(f"DESCRIBE HISTORY delta.`{path}`").collect()
            latest_version = history[0] if history else None
            
            return {
                "version": latest_version.version if latest_version else None,
                "timestamp": latest_version.timestamp.isoformat() if latest_version else None,
                "operation": latest_version.operation if latest_version else None,
                "operation_parameters": latest_version.operationParameters if latest_version else None,
                "is_time_travel_enabled": True  # Delta tables always support time travel
            }
        except Exception as e:
            self.logger.error(f"Error getting Delta metadata: {str(e)}")
            raise

    def _calculate_column_stats(self, df: DataFrame) -> Dict[str, Any]:
        """Calculate statistics for each column."""
        stats = {}
        for column in df.columns:
            col_stats = {
                "type": str(df.schema[column].dataType),
                "count": df.count(),
                "null_count": df.filter(df[column].isNull()).count(),
                "unique_count": df.select(column).distinct().count()
            }
            
            if df.schema[column].dataType.typeName() in ["integer", "long", "double", "float"]:
                numeric_stats = df.select(
                    mean(column).alias("mean"),
                    stddev(column).alias("std"),
                    min(column).alias("min"),
                    max(column).alias("max")
                ).collect()[0]
                
                col_stats.update({
                    "mean": float(numeric_stats["mean"]) if numeric_stats["mean"] is not None else 0.0,
                    "std": float(numeric_stats["std"]) if numeric_stats["std"] is not None else 0.0,
                    "min": float(numeric_stats["min"]) if numeric_stats["min"] is not None else 0.0,
                    "max": float(numeric_stats["max"]) if numeric_stats["max"] is not None else 0.0
                })
            
            stats[column] = col_stats
        
        return stats

    def _apply_rule(self, df: DataFrame, column: str, rule: Dict[str, Any]) -> Dict[str, Any]:
        """Apply a single data quality rule."""
        rule_type = rule.get("type")
        result = {
            "column": column,
            "rule": rule_type,
            "passed": True
        }
        
        if rule_type == "not_null":
            null_count = df.filter(df[column].isNull()).count()
            result["passed"] = null_count == 0
            result["null_count"] = null_count
            
        elif rule_type == "unique":
            total_count = df.count()
            unique_count = df.select(column).distinct().count()
            result["passed"] = total_count == unique_count
            result["duplicate_count"] = total_count - unique_count
            
        elif rule_type == "range":
            min_val = rule.get("min")
            max_val = rule.get("max")
            if min_val is not None:
                below_min = df.filter(df[column] < min_val).count()
                result["passed"] = result["passed"] and below_min == 0
                result["below_min"] = below_min
            if max_val is not None:
                above_max = df.filter(df[column] > max_val).count()
                result["passed"] = result["passed"] and above_max == 0
                result["above_max"] = above_max
                
        elif rule_type == "pattern":
            pattern = rule.get("pattern")
            if pattern:
                invalid_count = df.filter(~df[column].rlike(pattern)).count()
                result["passed"] = invalid_count == 0
                result["invalid_count"] = invalid_count
                
        return result
