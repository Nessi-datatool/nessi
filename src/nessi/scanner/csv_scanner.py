from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.functions import col, count, mean, stddev, min, max, countDistinct
from typing import Dict, List, Optional, Any, Union
import logging
from pathlib import Path
import pandas as pd
from datetime import datetime
from nessi.scanner.base_scanner import BaseScanner
from nessi.scanner.models import ScanResult

class CSVScanner(BaseScanner):
    def __init__(self, spark: Optional[SparkSession] = None, config: Optional[Dict[str, Any]] = None):
        super().__init__(spark, config)
        self.logger = logging.getLogger(__name__)

    def scan(self, path: str, **kwargs) -> Dict[str, Any]:
        """Scan a CSV file and return analysis results."""
        try:
            df = self._load_data(path, **kwargs)
            analysis = self.analyzer.analyze_dataframe(df)
            scan_results = {
                "source": path,
                "format": "csv",
                "row_count": df.count() if isinstance(df, DataFrame) else len(df),
                "column_count": len(df.columns),
                "column_stats": self._calculate_column_stats(df),
                "quality_score": analysis.get("quality_score", 0.0),
                "completeness_score": analysis.get("completeness_score", 0.0),
                "consistency_score": analysis.get("consistency_score", 0.0),
                "anomalies": analysis.get("anomalies", []),
                "patterns": analysis.get("patterns", []),
                "validation_results": {},
                "metadata": self._get_csv_metadata(path)
            }
            return self._generate_report(scan_results)
        except Exception as e:
            self.logger.error(f"Error scanning CSV file {path}: {str(e)}")
            raise

    def validate(self, path: str, rules: List[Dict[str, Any]], **kwargs) -> Dict[str, Any]:
        """Validate a CSV file against rules."""
        try:
            df = self._load_data(path, **kwargs)
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
            self.logger.error(f"Error validating CSV file {path}: {str(e)}")
            raise

    def scan_table(self, path: str, format: str = "csv", **kwargs) -> Dict[str, Any]:
        """Scan a CSV table and return analysis results.
        
        Args:
            path: Path to the CSV file
            format: Expected format (must be 'csv')
            **kwargs: Additional options for reading CSV
            
        Returns:
            Dictionary containing scan results
        """
        if format.lower() != "csv":
            raise ValueError("CSVScanner only supports 'csv' format")
        return self.scan(path, **kwargs)

    def read_table(self, path: str, **kwargs) -> Union[pd.DataFrame, DataFrame]:
        """Read a CSV table into a DataFrame.
        
        Args:
            path: Path to the CSV file
            **kwargs: Additional options for reading CSV
            
        Returns:
            DataFrame containing the data (pandas or Spark)
        """
        return self._load_data(path, **kwargs)

    def _read_data(self, path: Path, **kwargs) -> Union[pd.DataFrame, DataFrame]:
        """Read CSV data from the specified path."""
        try:
            if self.spark:
                return self.spark.read.csv(str(path), header=True, inferSchema=True, **kwargs)
            return pd.read_csv(path, **kwargs)
        except Exception as e:
            self.logger.error(f"Error reading CSV file {path}: {str(e)}")
            raise

    def _get_csv_metadata(self, path: str) -> Dict[str, Any]:
        """Get CSV file metadata."""
        try:
            file_path = Path(path)
            return {
                "file_size": file_path.stat().st_size,
                "last_modified": datetime.fromtimestamp(file_path.stat().st_mtime).isoformat(),
                "encoding": "utf-8"  # Default encoding, can be enhanced
            }
        except Exception as e:
            self.logger.error(f"Error getting CSV metadata: {str(e)}")
            raise

    def _calculate_column_stats(self, df: Union[pd.DataFrame, DataFrame]) -> Dict[str, Any]:
        """Calculate statistics for each column."""
        stats = {}
        
        if isinstance(df, pd.DataFrame):
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
        else:
            # PySpark DataFrame
            for column in df.columns:
                stats_df = df.agg(
                    count(col(column)).alias("count"),
                    countDistinct(col(column)).alias("unique_count"),
                    mean(col(column)).alias("mean"),
                    stddev(col(column)).alias("std"),
                    min(col(column)).alias("min"),
                    max(col(column)).alias("max")
                ).collect()[0]
                
                col_stats = {
                    "type": str(df.schema[column].dataType),
                    "count": stats_df["count"],
                    "null_count": df.filter(col(column).isNull()).count(),
                    "unique_count": stats_df["unique_count"]
                }
                
                if df.schema[column].dataType.typeName() in ["integer", "long", "double", "float"]:
                    col_stats.update({
                        "mean": float(stats_df["mean"]),
                        "std": float(stats_df["std"]),
                        "min": float(stats_df["min"]),
                        "max": float(stats_df["max"])
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