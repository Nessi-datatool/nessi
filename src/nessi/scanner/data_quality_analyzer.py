from typing import Dict, Any, List, Optional, Union
import pandas as pd
import numpy as np
from datetime import datetime
import logging
import re
from scipy import stats
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler
from pyspark.sql import DataFrame as SparkDataFrame
from pyspark.sql.functions import col, count, countDistinct, when, isnull, mean, stddev, min, max

logger = logging.getLogger(__name__)

class DataQualityAnalyzer:
    """Analyzes data quality and provides comprehensive profiling."""
    
    def __init__(self, spark):
        self.spark = spark
        self._setup_logging()
    
    def _setup_logging(self):
        """Setup logging configuration."""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    
    def analyze_data_quality(self, df: Union[pd.DataFrame, SparkDataFrame]) -> Dict[str, Any]:
        """
        Analyze data quality metrics for a DataFrame.
        
        Args:
            df: Input DataFrame to analyze (pandas or Spark)
            
        Returns:
            Dictionary containing data quality metrics
        """
        if isinstance(df, SparkDataFrame):
            if df.isEmpty():
                return {}
            return self._analyze_spark_dataframe(df)
        else:
            if df.empty:
                return {}
            return self._analyze_pandas_dataframe(df)

    def _analyze_spark_dataframe(self, df: SparkDataFrame) -> Dict[str, Any]:
        """Analyze a Spark DataFrame."""
        total_rows = df.count()
        metrics = {
            "completeness": self._calculate_completeness(df, total_rows),
            "uniqueness": self._calculate_uniqueness(df, total_rows),
            "consistency": self._calculate_consistency(df),
            "anomalies": self._detect_spark_anomalies(df),
            "patterns": self._detect_spark_patterns(df)
        }
        return metrics

    def _analyze_pandas_dataframe(self, df: pd.DataFrame) -> Dict[str, Any]:
        """Analyze a pandas DataFrame."""
        analysis = {
            "timestamp": datetime.now().isoformat(),
            "row_count": len(df),
            "column_count": len(df.columns),
            "columns": {},
            "quality_score": 0.0,
            "completeness_score": 0.0,
            "consistency_score": 0.0,
            "anomalies": [],
            "patterns": []
        }
        
        # Analyze each column
        for column in df.columns:
            analysis["columns"][column] = self._analyze_column(df[column])
        
        # Calculate overall scores
        analysis["completeness_score"] = self._calculate_completeness_score(analysis["columns"])
        analysis["consistency_score"] = self._calculate_consistency_score(analysis["columns"])
        analysis["quality_score"] = self._calculate_quality_score(analysis)
        
        # Detect anomalies and patterns
        analysis["anomalies"] = self._detect_anomalies(df)
        analysis["patterns"] = self._detect_patterns(df)
        
        return analysis

    def _calculate_completeness(self, df: SparkDataFrame, total_rows: int) -> Dict[str, float]:
        """Calculate completeness metrics for each column in a Spark DataFrame."""
        completeness = {}
        for column in df.columns:
            non_null_count = df.select(count(when(col(column).isNotNull(), 1))).collect()[0][0]
            completeness[column] = non_null_count / total_rows if total_rows > 0 else 0.0
        return completeness

    def _calculate_uniqueness(self, df: SparkDataFrame, total_rows: int) -> Dict[str, float]:
        """Calculate uniqueness metrics for each column in a Spark DataFrame."""
        uniqueness = {}
        for column in df.columns:
            distinct_count = df.select(countDistinct(col(column))).collect()[0][0]
            uniqueness[column] = distinct_count / total_rows if total_rows > 0 else 0.0
        return uniqueness

    def _calculate_consistency(self, df: SparkDataFrame) -> Dict[str, Any]:
        """Calculate consistency metrics for a Spark DataFrame."""
        consistency = {
            "schema_consistency": self._check_schema_consistency(df),
            "data_type_consistency": self._check_data_type_consistency(df)
        }
        return consistency

    def _check_schema_consistency(self, df: SparkDataFrame) -> bool:
        """Check if the schema is consistent across the Spark DataFrame."""
        try:
            df.limit(1).collect()
            return True
        except Exception:
            return False

    def _check_data_type_consistency(self, df: SparkDataFrame) -> Dict[str, bool]:
        """Check if data types are consistent for each column in a Spark DataFrame."""
        type_consistency = {}
        for column in df.columns:
            try:
                if df.schema[column].dataType.typeName() == "integer":
                    df.select(col(column).cast("integer")).limit(1).collect()
                elif df.schema[column].dataType.typeName() == "string":
                    df.select(col(column).cast("string")).limit(1).collect()
                type_consistency[column] = True
            except Exception:
                type_consistency[column] = False
        return type_consistency

    def _detect_spark_anomalies(self, df: SparkDataFrame) -> List[Dict[str, Any]]:
        """Detect anomalies in a Spark DataFrame."""
        anomalies = []
        for column in df.columns:
            if df.schema[column].dataType.typeName() in ["integer", "long", "double", "float"]:
                stats = df.select(
                    mean(col(column)).alias("mean"),
                    stddev(col(column)).alias("std"),
                    min(col(column)).alias("min"),
                    max(col(column)).alias("max")
                ).collect()[0]
                
                if stats["std"] is not None:
                    threshold = 3.0  # 3 standard deviations
                    anomalies.extend(
                        df.filter(
                            (col(column) < stats["mean"] - threshold * stats["std"]) |
                            (col(column) > stats["mean"] + threshold * stats["std"])
                        ).select(column).collect()
                    )
        return anomalies

    def _detect_spark_patterns(self, df: SparkDataFrame) -> List[Dict[str, Any]]:
        """Detect patterns in string columns of a Spark DataFrame."""
        patterns = []
        pattern_checks = {
            "email": r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
            "phone": r'^\+?[\d\s-]{10,}$',
            "url": r'^https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+',
            "date": r'^\d{4}-\d{2}-\d{2}$',
            "time": r'^\d{2}:\d{2}(:\d{2})?$',
            "ip": r'^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$'
        }
        
        for column in df.columns:
            if df.schema[column].dataType.typeName() == "string":
                for pattern_name, regex in pattern_checks.items():
                    matches = df.filter(col(column).rlike(regex)).count()
                    if matches > 0:
                        patterns.append({
                            "column": column,
                            "pattern": pattern_name,
                            "count": matches,
                            "percentage": matches / df.count() * 100
                        })
        return patterns

    def _analyze_column(self, series: pd.Series) -> Dict[str, Any]:
        """Analyze a single column of a pandas DataFrame."""
        analysis = {
            "type": str(series.dtype),
            "null_count": series.isnull().sum(),
            "unique_count": series.nunique(),
            "value_counts": {},
            "statistics": {},
            "patterns": []
        }
        
        if pd.api.types.is_numeric_dtype(series):
            analysis["statistics"] = {
                "min": series.min(),
                "max": series.max(),
                "mean": series.mean(),
                "median": series.median(),
                "std": series.std(),
                "skew": series.skew(),
                "kurtosis": series.kurtosis()
            }
        
        value_counts = series.value_counts().head(10)
        analysis["value_counts"] = value_counts.to_dict()
        
        if pd.api.types.is_string_dtype(series):
            analysis["patterns"] = self._detect_string_patterns(series)
        
        return analysis

    def _calculate_completeness_score(self, columns: Dict[str, Dict[str, Any]]) -> float:
        """Calculate completeness score based on null values."""
        total_cells = sum(col["null_count"] for col in columns.values())
        total_values = sum(len(col["value_counts"]) for col in columns.values())
        return 1.0 - (total_cells / (total_cells + total_values))

    def _calculate_consistency_score(self, columns: Dict[str, Dict[str, Any]]) -> float:
        """Calculate consistency score based on data type consistency."""
        type_consistency = sum(
            1 for col in columns.values()
            if col["type"] not in ["object", "string"]
        )
        return type_consistency / len(columns)

    def _calculate_quality_score(self, analysis: Dict[str, Any]) -> float:
        """Calculate overall quality score."""
        weights = {
            "completeness": 0.4,
            "consistency": 0.3,
            "anomaly": 0.2,
            "pattern": 0.1
        }
        
        scores = {
            "completeness": analysis["completeness_score"],
            "consistency": analysis["consistency_score"],
            "anomaly": 1.0 - (len(analysis["anomalies"]) / analysis["row_count"]),
            "pattern": 1.0 - (len(analysis["patterns"]) / len(analysis["columns"]))
        }
        
        return sum(score * weights[metric] for metric, score in scores.items())

    def _detect_anomalies(self, df: pd.DataFrame) -> List[Dict[str, Any]]:
        """Detect anomalies in a pandas DataFrame using Isolation Forest."""
        anomalies = []
        
        for column in df.columns:
            if pd.api.types.is_numeric_dtype(df[column]):
                data = df[column].dropna().values.reshape(-1, 1)
                if len(data) < 2:
                    continue
                
                scaler = StandardScaler()
                scaled_data = scaler.fit_transform(data)
                
                clf = IsolationForest(contamination=0.1, random_state=42)
                clf.fit(scaled_data)
                
                scores = clf.decision_function(scaled_data)
                threshold = np.percentile(scores, 10)
                
                anomaly_indices = np.where(scores < threshold)[0]
                for idx in anomaly_indices:
                    anomalies.append({
                        "column": column,
                        "row": int(df.index[idx]),
                        "value": float(data[idx][0]),
                        "score": float(scores[idx])
                    })
        
        return anomalies

    def _detect_string_patterns(self, series: pd.Series) -> List[Dict[str, Any]]:
        """Detect patterns in string columns of a pandas DataFrame."""
        patterns = []
        pattern_checks = {
            "email": r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
            "phone": r'^\+?[\d\s-]{10,}$',
            "url": r'^https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+',
            "date": r'^\d{4}-\d{2}-\d{2}$',
            "time": r'^\d{2}:\d{2}(:\d{2})?$',
            "ip": r'^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$'
        }
        
        for pattern_name, regex in pattern_checks.items():
            matches = series.str.match(regex, na=False)
            if matches.any():
                patterns.append({
                    "name": pattern_name,
                    "count": int(matches.sum()),
                    "percentage": float(matches.mean() * 100)
                })
        
        return patterns

    def _detect_patterns(self, df: pd.DataFrame) -> List[Dict[str, Any]]:
        """Detect patterns in string columns of a pandas DataFrame."""
        patterns = []
        for column in df.columns:
            if pd.api.types.is_string_dtype(df[column]):
                column_patterns = self._detect_string_patterns(df[column])
                for pattern in column_patterns:
                    pattern["column"] = column
                    patterns.append(pattern)
        return patterns 