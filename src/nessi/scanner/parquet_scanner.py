from pyspark.sql import SparkSession
from typing import Dict, List, Optional, Any, Union, Tuple
import logging
from pathlib import Path
import pandas as pd
from datetime import datetime
import os
import shutil
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, TimestampType
from pyspark.sql import DataFrame
from pyspark.sql.window import Window
from pyspark.sql import functions as F

from nessi.scanner.base_scanner import BaseScanner
from nessi.scanner.models import ScanResult

class ParquetScanner(BaseScanner):
    def __init__(self, spark: Optional[SparkSession] = None, config: Optional[Dict[str, Any]] = None):
        super().__init__(spark, config)
        self.spark = spark
        self.logger = logging.getLogger(__name__)
        self.chunk_size = config.get("chunk_size", 100000) if config else 100000
    def scan(self, path: str) -> Dict:
        """Scan a Parquet file and return analysis results."""
        try:
            # Check for corrupted data first
            if self._is_corrupted(path):
                raise RuntimeError("Parquet file appears to be corrupted")

            df = self._load_data(path)
            scan_results = {
                "source": path,
                "format": "parquet",
                "row_count": df.count(),
                "column_count": len(df.columns),
                "column_stats": self._calculate_column_stats(df),
                "quality_score": self.analyzer.calculate_quality_score(df),
                "completeness_score": self.analyzer.calculate_completeness_score(df),
                "consistency_score": self.analyzer.calculate_consistency_score(df),
                "anomalies": self.analyzer.detect_anomalies(df),
                "patterns": self.analyzer.detect_patterns(df),
                "validation_results": {},
                "metadata": self._get_parquet_metadata(path),
                "schema_evolution": self._detect_schema_evolution(path)
            }
            return self._generate_report(scan_results)
        except Exception as e:
            self.logger.error(f"Error scanning Parquet file {path}: {str(e)}")
            raise
    def scan_large_table(self, path: str) -> Dict:
        """Scan a large Parquet table using chunked processing."""
        try:
            if self._is_corrupted(path):
                raise RuntimeError("Parquet file appears to be corrupted")

            results = {
                "source": path,
                "format": "parquet",
                "row_count": 0,
                "column_count": 0,
                "column_stats": {},
                "quality_score": 0.0,
                "completeness_score": 0.0,
                "consistency_score": 0.0,
                "anomalies": [],
                "patterns": [],
                "validation_results": {},
                "metadata": self._get_parquet_metadata(path),
                "schema_evolution": self._detect_schema_evolution(path)
            }

            # Process in chunks
            for chunk_df in self._read_data_chunks(path):
                results["row_count"] += chunk_df.count()
                results["column_count"] = len(chunk_df.columns)
                
                # Update statistics
                chunk_stats = self._calculate_column_stats(chunk_df)
                for col_name, stats in chunk_stats.items():
                    if col_name not in results["column_stats"]:
                        results["column_stats"][col_name] = stats
                    else:
                        # Merge statistics
                        results["column_stats"][col_name].update({
                            "count": results["column_stats"][col_name]["count"] + stats["count"],
                            "null_count": results["column_stats"][col_name]["null_count"] + stats["null_count"],
                            "unique_count": max(results["column_stats"][col_name]["unique_count"], stats["unique_count"])
                        })

                # Update quality metrics
                results["quality_score"] = max(results["quality_score"], self.analyzer.calculate_quality_score(chunk_df))
                results["completeness_score"] = max(results["completeness_score"], self.analyzer.calculate_completeness_score(chunk_df))
                results["consistency_score"] = max(results["consistency_score"], self.analyzer.calculate_consistency_score(chunk_df))

                # Collect anomalies and patterns
                results["anomalies"].extend(self.analyzer.detect_anomalies(chunk_df))
                results["patterns"].extend(self.analyzer.detect_patterns(chunk_df))

            return self._generate_report(results)
        except Exception as e:
            self.logger.error(f"Error scanning large Parquet file {path}: {str(e)}")
            raise
    def _is_corrupted(self, path: str) -> bool:
        """Check if a Parquet file is corrupted."""
        try:
            if self.spark:
                # Try to read the schema
                self.spark.read.parquet(path).schema
            else:
                # Try to read metadata
                pd.read_parquet(path, engine='pyarrow')
            return False
        except Exception as e:
            self.logger.warning(f"Potential corruption detected in {path}: {str(e)}")
            return True

    def _read_data_chunks(self, path: str) -> Union[pd.DataFrame, SparkSession]:
        """Read Parquet data in chunks."""
        if self.spark:
            df = self.spark.read.parquet(path)
            total_rows = df.count()
            for i in range(0, total_rows, self.chunk_size):
                yield df.limit(self.chunk_size).offset(i)
        else:
            # For pandas, we need to read the entire file and split it
            df = pd.read_parquet(path)
            for i in range(0, len(df), self.chunk_size):
                yield df.iloc[i:i + self.chunk_size]

    def _detect_schema_evolution(self, path: str) -> Dict:
        """Detect schema evolution in Parquet files."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                schema = df.schema
            else:
                df = pd.read_parquet(path)
                schema = df.dtypes

            evolution_info = {
                "current_schema": str(schema),
                "schema_changes": [],
                "compatibility_level": "FULL"  # FULL, BACKWARD, FORWARD, NONE
            }

            # Check for schema changes
            if isinstance(schema, StructType):
                for field in schema.fields:
                    if field.nullable:
                        evolution_info["schema_changes"].append({
                            "field": field.name,
                            "change": "nullable",
                            "type": str(field.dataType)
                        })

            return evolution_info
        except Exception as e:
            self.logger.error(f"Error detecting schema evolution: {str(e)}")
            return {
                "current_schema": "unknown",
                "schema_changes": [],
                "compatibility_level": "UNKNOWN"
            }

    def _get_parquet_metadata(self, path: str) -> Dict:
        """Get enhanced Parquet file metadata."""
        try:
            file_path = Path(path)
            metadata = {
                "file_size": file_path.stat().st_size,
                "last_modified": datetime.fromtimestamp(file_path.stat().st_mtime).isoformat(),
                "compression": self._detect_compression(path),
                "row_groups": self._get_row_group_info(path),
                "encoding": self._detect_encoding(path),
                "created_by": self._get_created_by(path)
            }
            return metadata
        except Exception as e:
            self.logger.error(f"Error getting Parquet metadata: {str(e)}")
            raise

    def _detect_compression(self, path: str) -> str:
        """Detect the compression used in the Parquet file."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                return df.inputFiles()[0].split('.')[-1]  # Simple detection based on file extension
            return "unknown"
        except Exception:
            return "unknown"

    def _get_row_group_info(self, path: str) -> Dict:
        """Get information about row groups in the Parquet file."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                return {
                    "count": len(df.inputFiles()),
                    "size_bytes": sum(os.path.getsize(f) for f in df.inputFiles())
                }
            return {"count": 1, "size_bytes": os.path.getsize(path)}
        except Exception:
            return {"count": 0, "size_bytes": 0}

    def _detect_encoding(self, path: str) -> str:
        """Detect the encoding used in the Parquet file."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                return "PLAIN"  # Default encoding
            return "unknown"
        except Exception:
            return "unknown"

    def _get_created_by(self, path: str) -> str:
        """Get information about what created the Parquet file."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                return "Spark"  # Default creator
            return "unknown"
        except Exception:
            return "unknown"
    def _read_data(self, path: Path) -> pd.DataFrame:
        """Read Parquet data from the specified path."""
        try:
            if self.spark:
                return self.spark.read.parquet(str(path)).toPandas()
            else:
                return pd.read_parquet(path)
        except Exception as e:
            self.logger.error(f"Error reading Parquet file {path}: {str(e)}")
            raise

    def _calculate_column_stats(self, df: pd.DataFrame) -> Dict:
        """Calculate statistics for each column."""
        stats = {}
        for column in df.columns:
            col_stats = {
                "type": str(df[column].dtype),
                "count": len(df[column]),
                "null_count": df[column].isnull().sum(),
                "unique_count": df[column].nunique()
            }
            
            if df[column].dtype in ["int64", "float64"]:
                col_stats.update({
                    "mean": float(df[column].mean()),
                    "std": float(df[column].std()),
                    "min": float(df[column].min()),
                    "max": float(df[column].max())
                })
            
            stats[column] = col_stats
        
        return stats

    def _apply_rule(self, df: pd.DataFrame, column: str, rule: Dict) -> Dict:
        """Apply a single data quality rule."""
        rule_type = rule.get("type")
        result = {
            "column": column,
            "rule": rule_type,
            "passed": True
        }
        
        if rule_type == "not_null":
            null_count = df[column].isnull().sum()
            result["passed"] = null_count == 0
            result["null_count"] = null_count
            
        elif rule_type == "unique":
            total_count = len(df)
            unique_count = df[column].nunique()
            result["passed"] = total_count == unique_count
            result["duplicate_count"] = total_count - unique_count
            
        elif rule_type == "range":
            min_val = rule.get("min")
            max_val = rule.get("max")
            if min_val is not None:
                below_min = (df[column] < min_val).sum()
                result["passed"] = result["passed"] and below_min == 0
                result["below_min"] = below_min
            if max_val is not None:
                above_max = (df[column] > max_val).sum()
                result["passed"] = result["passed"] and above_max == 0
                result["above_max"] = above_max
                
        elif rule_type == "pattern":
            pattern = rule.get("pattern")
            if pattern:
                invalid_count = df[~df[column].str.match(pattern)].shape[0]
                result["passed"] = invalid_count == 0
                result["invalid_count"] = invalid_count
                
        return result
    def scan_partitioned_table(self, path: str, partition_cols: List[str] = None) -> Dict:
        """Scan a partitioned Parquet table with partition information."""
        try:
            if self._is_corrupted(path):
                raise RuntimeError("Parquet file appears to be corrupted")

            df = self._load_data(path)
            partition_info = self._get_partition_info(path, partition_cols)
            
            scan_results = {
                "source": path,
                "format": "parquet",
                "row_count": df.count(),
                "column_count": len(df.columns),
                "partition_info": partition_info,
                "column_stats": self._calculate_column_stats(df),
                "quality_score": self.analyzer.calculate_quality_score(df),
                "completeness_score": self.analyzer.calculate_completeness_score(df),
                "consistency_score": self.analyzer.calculate_consistency_score(df),
                "anomalies": self.analyzer.detect_anomalies(df),
                "patterns": self.analyzer.detect_patterns(df),
                "validation_results": {},
                "metadata": self._get_parquet_metadata(path),
                "schema_evolution": self._detect_schema_evolution(path)
            }
            return self._generate_report(scan_results)
        except Exception as e:
            self.logger.error(f"Error scanning partitioned Parquet table {path}: {str(e)}")
            raise

    def _get_partition_info(self, path: str, partition_cols: List[str] = None) -> Dict:
        """Get information about table partitions."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                if partition_cols:
                    partition_values = df.select(partition_cols).distinct().collect()
                    return {
                        "partition_columns": partition_cols,
                        "partition_count": len(partition_values),
                        "partition_values": [dict(zip(partition_cols, row)) for row in partition_values]
                    }
                else:
                    # Try to infer partition columns
                    partition_cols = self._infer_partition_columns(path)
                    if partition_cols:
                        return self._get_partition_info(path, partition_cols)
            return {"partition_columns": [], "partition_count": 0, "partition_values": []}
        except Exception as e:
            self.logger.error(f"Error getting partition info: {str(e)}")
            return {"partition_columns": [], "partition_count": 0, "partition_values": []}

    def _infer_partition_columns(self, path: str) -> List[str]:
        """Infer partition columns from the directory structure."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                # Look for columns that match directory structure patterns
                partition_cols = []
                for col in df.columns:
                    if df.select(col).distinct().count() < 100:  # Heuristic for partition columns
                        partition_cols.append(col)
                return partition_cols
            return []
        except Exception:
            return []

    def optimize_table(self, path: str, options: Dict[str, Any] = None) -> Dict:
        """Optimize Parquet table for better performance."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                
                # Default optimization options
                default_options = {
                    "compression": "snappy",
                    "row_group_size": 128 * 1024 * 1024,  # 128MB
                    "page_size": 1 * 1024 * 1024,  # 1MB
                    "enable_statistics": True,
                    "max_statistics_size": 4096
                }
                
                # Merge with user options
                options = {**default_options, **(options or {})}
                
                # Write optimized table
                optimized_path = f"{path}_optimized"
                df.write.parquet(
                    optimized_path,
                    compression=options["compression"],
                    rowGroupSize=options["row_group_size"],
                    pageSize=options["page_size"],
                    enableStatistics=options["enable_statistics"],
                    maxStatisticsSize=options["max_statistics_size"]
                )
                
                # Compare sizes and performance
                original_size = sum(os.path.getsize(f) for f in self.spark.read.parquet(path).inputFiles())
                optimized_size = sum(os.path.getsize(f) for f in self.spark.read.parquet(optimized_path).inputFiles())
                
                return {
                    "original_size": original_size,
                    "optimized_size": optimized_size,
                    "compression_ratio": original_size / optimized_size if optimized_size > 0 else 0,
                    "options_used": options
                }
            return {"error": "Optimization requires Spark"}
        except Exception as e:
            self.logger.error(f"Error optimizing Parquet table: {str(e)}")
            raise

    def validate_schema(self, path: str, expected_schema: Dict) -> Dict:
        """Validate the schema of a Parquet table against an expected schema."""
        try:
            if self.spark:
                df = self.spark.read.parquet(path)
                actual_schema = df.schema
                
                validation_results = {
                    "passed": True,
                    "missing_columns": [],
                    "extra_columns": [],
                    "type_mismatches": [],
                    "nullable_mismatches": []
                }
                
                # Check for missing columns
                expected_columns = set(expected_schema.keys())
                actual_columns = set(df.columns)
                validation_results["missing_columns"] = list(expected_columns - actual_columns)
                validation_results["extra_columns"] = list(actual_columns - expected_columns)
                
                # Check column types and nullability
                for col_name, col_info in expected_schema.items():
                    if col_name in df.columns:
                        actual_type = str(df.schema[col_name].dataType)
                        actual_nullable = df.schema[col_name].nullable
                        
                        if col_info["type"] != actual_type:
                            validation_results["type_mismatches"].append({
                                "column": col_name,
                                "expected": col_info["type"],
                                "actual": actual_type
                            })
                        
                        if col_info.get("nullable") != actual_nullable:
                            validation_results["nullable_mismatches"].append({
                                "column": col_name,
                                "expected": col_info.get("nullable"),
                                "actual": actual_nullable
                            })
                
                validation_results["passed"] = all(
                    not validation_results[key] for key in [
                        "missing_columns", "extra_columns",
                        "type_mismatches", "nullable_mismatches"
                    ]
                )
                
                return validation_results
            return {"error": "Schema validation requires Spark"}
        except Exception as e:
            self.logger.error(f"Error validating schema: {str(e)}")
            raise
    def scan_time_series_table(self, path: str, time_column: str, time_range: Tuple[str, str] = None, 
                             aggregation: str = "1h", metrics: List[str] = None) -> Dict:
        """Scan a time-series Parquet table with time-based filtering and aggregation."""
        try:
            if self._is_corrupted(path):
                raise RuntimeError("Parquet file appears to be corrupted")

            df = self._load_data(path)
            
            if time_range and time_column in df.columns:
                start_time, end_time = time_range
                df = df.filter(
                    (df[time_column] >= start_time) & 
                    (df[time_column] <= end_time)
                )
            
            time_stats = self._calculate_time_stats(df, time_column)
            
            scan_results = {
                "source": path,
                "format": "parquet",
                "row_count": df.count(),
                "column_count": len(df.columns),
                "time_column": time_column,
                "time_range": time_range,
                "aggregation": aggregation,
                "time_stats": time_stats,
                "column_stats": self._calculate_column_stats(df),
                "quality_score": self.analyzer.calculate_quality_score(df),
                "completeness_score": self.analyzer.calculate_completeness_score(df),
                "consistency_score": self.analyzer.calculate_consistency_score(df),
                "anomalies": self.analyzer.detect_anomalies(df),
                "patterns": self.analyzer.detect_patterns(df),
                "seasonality": self._detect_seasonality(df, time_column, metrics) if metrics else None,
                "trends": self._detect_trends(df, time_column, metrics) if metrics else None,
                "validation_results": {},
                "metadata": self._get_parquet_metadata(path),
                "schema_evolution": self._detect_schema_evolution(path)
            }
            return self._generate_report(scan_results)
        except Exception as e:
            self.logger.error(f"Error scanning time series Parquet table {path}: {str(e)}")
            raise

    def _calculate_time_stats(self, df: DataFrame, time_column: str, 
                            aggregation: str = "1h", metrics: List[str] = None) -> Dict:
        """Calculate comprehensive time-based statistics for a DataFrame."""
        try:
            if self.spark:
                # Basic time statistics
                time_stats = {
                    "min_time": df.select(F.min(time_column)).collect()[0][0],
                    "max_time": df.select(F.max(time_column)).collect()[0][0],
                    "time_range": df.select(
                        F.max(time_column) - F.min(time_column)
                    ).collect()[0][0],
                    "record_count_by_period": {},
                    "metrics_by_period": {}
                }
                
                # Calculate record counts and metrics by different time periods
                for period in ["hour", "day", "month", "year"]:
                    period_col = f"{time_column}_{period}"
                    df_with_period = df.withColumn(
                        period_col,
                        F.date_trunc(period, df[time_column])
                    )
                    
                    # Record counts
                    counts = df_with_period.groupBy(period_col).count().collect()
                    time_stats["record_count_by_period"][period] = {
                        str(row[period_col]): row["count"] for row in counts
                    }
                    
                    # Calculate metrics if specified
                    if metrics:
                        for metric in metrics:
                            if metric in df.columns:
                                aggs = df_with_period.groupBy(period_col).agg(
                                    F.avg(metric).alias("avg"),
                                    F.min(metric).alias("min"),
                                    F.max(metric).alias("max"),
                                    F.stddev(metric).alias("std")
                                ).collect()
                                
                                if period not in time_stats["metrics_by_period"]:
                                    time_stats["metrics_by_period"][period] = {}
                                
                                time_stats["metrics_by_period"][period][metric] = {
                                    str(row[period_col]): {
                                        "avg": row["avg"],
                                        "min": row["min"],
                                        "max": row["max"],
                                        "std": row["std"]
                                    } for row in aggs
                                }
                
                return time_stats
            return {}
        except Exception as e:
            self.logger.error(f"Error calculating time stats: {str(e)}")
            return {}

    def _detect_seasonality(self, df: DataFrame, time_column: str, metrics: List[str]) -> Dict:
        """Detect seasonality patterns in time series data."""
        try:
            seasonality = {}
            for metric in metrics:
                if metric in df.columns:
                    daily = df.withColumn(
                        "hour", F.hour(df[time_column])
                    ).groupBy("hour").agg(
                        F.avg(metric).alias("avg_value")
                    ).collect()
                    
                    weekly = df.withColumn(
                        "day_of_week", F.dayofweek(df[time_column])
                    ).groupBy("day_of_week").agg(
                        F.avg(metric).alias("avg_value")
                    ).collect()
                    
                    seasonality[metric] = {
                        "daily": {row["hour"]: row["avg_value"] for row in daily},
                        "weekly": {row["day_of_week"]: row["avg_value"] for row in weekly}
                    }
            return seasonality
        except Exception as e:
            self.logger.error(f"Error detecting seasonality: {str(e)}")
            return {}

    def _detect_trends(self, df: DataFrame, time_column: str, metrics: List[str]) -> Dict:
        """Detect trends in time series data."""
        try:
            trends = {}
            for metric in metrics:
                if metric in df.columns:
                    window_spec = Window.orderBy(time_column).rowsBetween(-7, 7)
                    df_with_ma = df.withColumn(
                        f"{metric}_ma", F.avg(metric).over(window_spec)
                    )
                    
                    start_value = df_with_ma.select(F.first(f"{metric}_ma")).collect()[0][0]
                    end_value = df_with_ma.select(F.last(f"{metric}_ma")).collect()[0][0]
                    
                    trends[metric] = {
                        "direction": "up" if end_value > start_value else "down",
                        "change_percent": ((end_value - start_value) / start_value) * 100 if start_value != 0 else 0,
                        "start_value": start_value,
                        "end_value": end_value
                    }
            return trends
        except Exception as e:
            self.logger.error(f"Error detecting trends: {str(e)}")
            return {}

    def validate(self, path: str, rules: List[Dict[str, Any]], **kwargs) -> Dict[str, Any]:
        """Validate a Parquet table against rules.
        
        Args:
            path: Path to the Parquet table
            rules: List of validation rules
            **kwargs: Additional validation options
            
        Returns:
            Dictionary containing validation results
        """
        try:
            df = self._read_data(path)
            validation_results = {
                "passed": True,
                "checks": []
            }
            
            for rule in rules:
                column = rule.get("column")
                rule_type = rule.get("type")
                rule_params = rule.get("params", {})
                
                if column not in df.columns:
                    validation_results["checks"].append({
                        "column": column,
                        "status": "failed",
                        "error": "Column not found"
                    })
                    validation_results["passed"] = False
                    continue
                
                check_result = self._apply_rule(df, column, rule_type, rule_params)
                validation_results["checks"].append(check_result)
                if not check_result["passed"]:
                    validation_results["passed"] = False
            
            return validation_results
        except Exception as e:
            self.logger.error(f"Error validating Parquet table {path}: {str(e)}")
            raise

    def _read_data(self, path: str) -> Union[pd.DataFrame, DataFrame]:
        """Read data from a Parquet file.
        
        Args:
            path: Path to the Parquet file
            
        Returns:
            DataFrame containing the data
        """
        try:
            if self.spark:
                return self.spark.read.format("parquet").load(path)
            else:
                return pd.read_parquet(path)
        except Exception as e:
            self.logger.error(f"Error reading Parquet file {path}: {str(e)}")
            raise

    def scan_table(self, path: str, format: str = "parquet", **kwargs) -> Dict[str, Any]:
        """Scan a Parquet table and return analysis results.
        
        Args:
            path: Path to the Parquet table
            format: Expected format (must be 'parquet')
            **kwargs: Additional scan options
            
        Returns:
            Dictionary containing scan results
        """
        try:
            if format != "parquet":
                raise ValueError("ParquetScanner only supports 'parquet' format")

            df = self._read_data(path)
            
            scan_results = {
                "source": path,
                "format": "parquet",
                "row_count": len(df) if isinstance(df, pd.DataFrame) else df.count(),
                "column_count": len(df.columns),
                "column_stats": self._calculate_column_stats(df),
                "quality_score": self.analyzer.calculate_quality_score(df),
                "completeness_score": self.analyzer.calculate_completeness_score(df),
                "consistency_score": self.analyzer.calculate_consistency_score(df),
                "anomalies": self.analyzer.detect_anomalies(df),
                "patterns": self.analyzer.detect_patterns(df),
                "validation_results": {},
                "metadata": self._get_metadata(path)
            }

            return scan_results
        except Exception as e:
            self.logger.error(f"Error scanning Parquet table {path}: {str(e)}")
            raise

    def read_table(self, path: str, **kwargs) -> Union[pd.DataFrame, DataFrame]:
        """Read a Parquet table.
        
        Args:
            path: Path to the Parquet table
            **kwargs: Additional read options
            
        Returns:
            DataFrame containing the data
        """
        return self._read_data(path)