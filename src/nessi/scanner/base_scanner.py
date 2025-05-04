from abc import ABC, abstractmethod
from typing import Dict, Any, Optional, List, Union
from pathlib import Path
import pandas as pd
from datetime import datetime
import logging
from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.functions import mean, stddev, min, max
from pyspark.sql.types import IntegerType, FloatType

from nessi.quality.analyzer import DataQualityAnalyzer
from nessi.config.quality_config import QualityConfig
from nessi.scanner.models import ScanResult, SchemaValidationResult

class BaseScanner(ABC):
    """Base class for all data scanners."""
    
    def __init__(self, spark: Optional[SparkSession] = None, config: Optional[Dict[str, Any]] = None):
        """Initialize the scanner with optional Spark session and configuration."""
        self.spark = spark
        self.config = config or {}
        self.quality_analyzer = DataQualityAnalyzer(spark)
        self.logger = logging.getLogger(self.__class__.__name__)
    
    @abstractmethod
    def _read_data(self, path: str) -> Union[pd.DataFrame, DataFrame]:
        """Read data from a data source.
        
        Args:
            path: Path to the data source
            
        Returns:
            DataFrame containing the data (pandas or Spark)
        """
        pass

    def scan(self, path: str, **kwargs) -> ScanResult:
        """Scan a data source and return analysis results.
        
        Args:
            path: Path to the data source
            **kwargs: Additional scan options
            
        Returns:
            ScanResult: Results of the scan
        """
        try:
            df = self._read_data(path)
            scan_results = {
                "source": path,
                "format": self._get_format(),
                "row_count": len(df) if isinstance(df, pd.DataFrame) else df.count(),
                "column_count": len(df.columns),
                "column_stats": self._calculate_column_stats(df),
                "quality_score": self.quality_analyzer.calculate_quality_score(df),
                "completeness_score": self.quality_analyzer.calculate_completeness_score(df),
                "consistency_score": self.quality_analyzer.calculate_consistency_score(df),
                "anomalies": self.quality_analyzer.detect_anomalies(df),
                "patterns": self.quality_analyzer.detect_patterns(df),
                "validation_results": {},
                "metadata": self._get_metadata(path)
            }
            return ScanResult(**scan_results)
        except Exception as e:
            self.logger.error(f"Error scanning data source {path}: {str(e)}")
            raise

    def validate(self, path: str, rules: List[Dict[str, Any]], **kwargs) -> Dict[str, Any]:
        """Validate a data source against rules.
        
        Args:
            path: Path to the data source
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
                
                if not column or column not in df.columns:
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
            self.logger.error(f"Error validating data source {path}: {str(e)}")
            raise

    def _get_format(self) -> str:
        """Get the format of the data source."""
        return self.__class__.__name__.replace("Scanner", "").lower()

    def _get_metadata(self, path: str) -> Dict[str, Any]:
        """Get metadata for the data source."""
        file_path = Path(path)
        return {
            "file_size": file_path.stat().st_size if file_path.exists() else 0,
            "last_modified": datetime.fromtimestamp(file_path.stat().st_mtime).isoformat() if file_path.exists() else None,
            "format": self._get_format()
        }

    def _calculate_column_stats(self, df: Union[pd.DataFrame, DataFrame]) -> Dict[str, Any]:
        """Calculate statistics for each column."""
        stats = {}
        for column in df.columns:
            if isinstance(df, pd.DataFrame):
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
            else:
                col_stats = {
                    "type": str(df.schema[column].dataType),
                    "count": df.select(column).count(),
                    "null_count": df.filter(col(column).isNull()).count(),
                    "unique_count": df.select(column).distinct().count()
                }
                
                if isinstance(df.schema[column].dataType, (IntegerType, FloatType)):
                    stats_df = df.select(
                        mean(column).alias("mean"),
                        stddev(column).alias("std"),
                        min(column).alias("min"),
                        max(column).alias("max")
                    ).collect()[0]
                    
                    col_stats.update({
                        "mean": float(stats_df["mean"]),
                        "std": float(stats_df["std"]),
                        "min": float(stats_df["min"]),
                        "max": float(stats_df["max"])
                    })
            
            stats[column] = col_stats
        
        return stats

    def _apply_rule(self, df: Union[pd.DataFrame, DataFrame], column: str, rule_type: str, rule_params: Dict[str, Any]) -> Dict[str, Any]:
        """Apply a validation rule to a column."""
        result = {
            "column": column,
            "rule": rule_type,
            "passed": True
        }
        
        if isinstance(df, pd.DataFrame):
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
                min_val = rule_params.get("min")
                max_val = rule_params.get("max")
                if min_val is not None:
                    below_min = (df[column] < min_val).sum()
                    result["passed"] = result["passed"] and below_min == 0
                    result["below_min"] = below_min
                if max_val is not None:
                    above_max = (df[column] > max_val).sum()
                    result["passed"] = result["passed"] and above_max == 0
                    result["above_max"] = above_max
        else:
            if rule_type == "not_null":
                null_count = df.filter(col(column).isNull()).count()
                result["passed"] = null_count == 0
                result["null_count"] = null_count
            elif rule_type == "unique":
                total_count = df.count()
                unique_count = df.select(column).distinct().count()
                result["passed"] = total_count == unique_count
                result["duplicate_count"] = total_count - unique_count
            elif rule_type == "range":
                min_val = rule_params.get("min")
                max_val = rule_params.get("max")
                if min_val is not None:
                    below_min = df.filter(col(column) < min_val).count()
                    result["passed"] = result["passed"] and below_min == 0
                    result["below_min"] = below_min
                if max_val is not None:
                    above_max = df.filter(col(column) > max_val).count()
                    result["passed"] = result["passed"] and above_max == 0
                    result["above_max"] = above_max
        
        return result

    def _generate_report(self, scan_results: Dict[str, Any]) -> Dict[str, Any]:
        """Generate a report from scan results."""
        report = {
            "timestamp": datetime.now().isoformat(),
            "source": scan_results.get("source", ""),
            "format": scan_results.get("format", ""),
            "row_count": scan_results.get("row_count", 0),
            "column_count": scan_results.get("column_count", 0),
            "column_stats": scan_results.get("column_stats", {}),
            "quality_score": scan_results.get("quality_score", 0.0),
            "completeness_score": scan_results.get("completeness_score", 0.0),
            "consistency_score": scan_results.get("consistency_score", 0.0),
            "anomalies": scan_results.get("anomalies", []),
            "patterns": scan_results.get("patterns", []),
            "validation_results": scan_results.get("validation_results", {}),
            "metadata": scan_results.get("metadata", {})
        }
        return report
    
    def _load_data(self, path: str, **kwargs) -> Union[pd.DataFrame, DataFrame]:
        """Load data from the specified path."""
        file_path = Path(path)
        if not file_path.exists():
            raise FileNotFoundError(f"Data source not found: {path}")
        
        return self._read_data(file_path, **kwargs)

    def _get_basic_stats(self, df: Union[pd.DataFrame, DataFrame]) -> Dict[str, Any]:
        """Get basic statistics for a DataFrame.
        
        Args:
            df: DataFrame to analyze
            
        Returns:
            Dictionary containing basic statistics
        """
        try:
            if isinstance(df, pd.DataFrame):
                return {
                    "row_count": len(df),
                    "column_count": len(df.columns),
                    "schema": [{"name": col, "type": str(df[col].dtype)} for col in df.columns]
                }
            else:
                return {
                    "row_count": df.count(),
                    "column_count": len(df.columns),
                    "schema": [{"name": field.name, "type": str(field.dataType)} for field in df.schema.fields]
                }
        except Exception as e:
            self.logger.error(f"Error getting basic stats: {str(e)}")
            raise

    def _validate_schema(self, df: DataFrame) -> SchemaValidationResult:
        """
        Validate the schema of a DataFrame.
        
        Args:
            df (DataFrame): DataFrame to validate
            
        Returns:
            SchemaValidationResult: Results of schema validation
        """
        schema = df.schema.jsonValue()
        fields = []
        for field in schema["fields"]:
            fields.append({
                "name": field["name"],
                "type": field["type"],
                "nullable": field["nullable"]
            })
        
        return SchemaValidationResult(
            fields=fields,
            valid=True,
            errors=[]
        ) 