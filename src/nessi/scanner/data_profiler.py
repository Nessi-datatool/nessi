from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.functions import col, count, mean, stddev, min, max, approx_count_distinct, percentile_approx
from pyspark.sql.types import NumericType, StringType, DateType, TimestampType
from typing import Dict, List, Optional
import numpy as np
from scipy import stats
import logging

class DataProfiler:
    def __init__(self, spark: SparkSession):
        self.spark = spark
        self.logger = logging.getLogger(__name__)

    def profile_table(self, df: DataFrame) -> Dict:
        """Generate comprehensive data profile for a DataFrame."""
        try:
            profile = {
                "table_stats": self._get_table_stats(df),
                "column_profiles": self._get_column_profiles(df),
                "correlations": self._calculate_correlations(df),
                "anomalies": self._detect_anomalies(df)
            }
            return profile
        except Exception as e:
            self.logger.error(f"Error profiling table: {str(e)}")
            raise

    def _get_table_stats(self, df: DataFrame) -> Dict:
        """Calculate table-level statistics."""
        return {
            "row_count": df.count(),
            "column_count": len(df.columns),
            "duplicate_rows": self._count_duplicates(df),
            "completeness": self._calculate_completeness(df)
        }

    def _get_column_profiles(self, df: DataFrame) -> Dict:
        """Generate detailed profiles for each column."""
        profiles = {}
        for col_name in df.columns:
            col_type = df.schema[col_name].dataType
            profiles[col_name] = self._profile_column(df, col_name, col_type)
        return profiles

    def _profile_column(self, df: DataFrame, col_name: str, col_type: type) -> Dict:
        """Generate detailed profile for a single column."""
        profile = {
            "type": col_type.typeName(),
            "null_count": df.filter(col(col_name).isNull()).count(),
            "distinct_count": df.select(col_name).distinct().count()
        }

        if isinstance(col_type, NumericType):
            stats = df.select(
                mean(col(col_name)).alias("mean"),
                stddev(col(col_name)).alias("stddev"),
                min(col(col_name)).alias("min"),
                max(col(col_name)).alias("max"),
                percentile_approx(col(col_name), 0.25).alias("q1"),
                percentile_approx(col(col_name), 0.5).alias("median"),
                percentile_approx(col(col_name), 0.75).alias("q3")
            ).collect()[0]

            profile.update({
                "mean": float(stats["mean"]),
                "stddev": float(stats["stddev"]),
                "min": float(stats["min"]),
                "max": float(stats["max"]),
                "q1": float(stats["q1"]),
                "median": float(stats["median"]),
                "q3": float(stats["q3"]),
                "skewness": self._calculate_skewness(df, col_name),
                "kurtosis": self._calculate_kurtosis(df, col_name)
            })

        elif isinstance(col_type, StringType):
            profile.update({
                "min_length": df.select(col_name).rdd.map(lambda x: len(str(x[0]))).min(),
                "max_length": df.select(col_name).rdd.map(lambda x: len(str(x[0]))).max(),
                "avg_length": df.select(col_name).rdd.map(lambda x: len(str(x[0]))).mean(),
                "top_values": self._get_top_values(df, col_name)
            })

        elif isinstance(col_type, (DateType, TimestampType)):
            profile.update({
                "min_date": df.select(min(col(col_name))).collect()[0][0],
                "max_date": df.select(max(col(col_name))).collect()[0][0],
                "date_range": self._calculate_date_range(df, col_name)
            })

        return profile

    def _calculate_correlations(self, df: DataFrame) -> Dict:
        """Calculate correlations between numeric columns."""
        numeric_cols = [f.name for f in df.schema.fields if isinstance(f.dataType, NumericType)]
        if len(numeric_cols) < 2:
            return {}

        corr_matrix = []
        for i, col1 in enumerate(numeric_cols):
            row = []
            for j, col2 in enumerate(numeric_cols):
                if i <= j:
                    corr = df.stat.corr(col1, col2)
                    row.append(float(corr))
                else:
                    row.append(None)
            corr_matrix.append(row)

        return {
            "columns": numeric_cols,
            "matrix": corr_matrix
        }

    def _detect_anomalies(self, df: DataFrame) -> Dict:
        """Detect anomalies in numeric columns using statistical methods."""
        anomalies = {}
        for col_name in df.columns:
            if isinstance(df.schema[col_name].dataType, NumericType):
                col_anomalies = self._detect_column_anomalies(df, col_name)
                if col_anomalies:
                    anomalies[col_name] = col_anomalies
        return anomalies

    def _detect_column_anomalies(self, df: DataFrame, col_name: str) -> Dict:
        """Detect anomalies in a numeric column using z-score and IQR methods."""
        stats = df.select(
            mean(col(col_name)).alias("mean"),
            stddev(col(col_name)).alias("stddev"),
            percentile_approx(col(col_name), 0.25).alias("q1"),
            percentile_approx(col(col_name), 0.75).alias("q3")
        ).collect()[0]

        mean_val = float(stats["mean"])
        std_val = float(stats["stddev"])
        q1 = float(stats["q1"])
        q3 = float(stats["q3"])
        iqr = q3 - q1

        # Z-score anomalies
        z_scores = df.select(
            (col(col_name) - mean_val) / std_val
        ).rdd.map(lambda x: abs(x[0])).collect()

        z_threshold = 3.0
        z_anomalies = [i for i, z in enumerate(z_scores) if z > z_threshold]

        # IQR anomalies
        lower_bound = q1 - 1.5 * iqr
        upper_bound = q3 + 1.5 * iqr
        iqr_anomalies = df.filter(
            (col(col_name) < lower_bound) | (col(col_name) > upper_bound)
        ).count()

        return {
            "z_score_anomalies": len(z_anomalies),
            "iqr_anomalies": iqr_anomalies,
            "thresholds": {
                "z_score": z_threshold,
                "iqr_lower": lower_bound,
                "iqr_upper": upper_bound
            }
        }

    def _count_duplicates(self, df: DataFrame) -> int:
        """Count duplicate rows in the DataFrame."""
        return df.count() - df.distinct().count()

    def _calculate_completeness(self, df: DataFrame) -> float:
        """Calculate the completeness ratio of the DataFrame."""
        total_cells = df.count() * len(df.columns)
        null_cells = sum(df.select([count(col(c)) for c in df.columns]).collect()[0])
        return 1.0 - (null_cells / total_cells)

    def _calculate_skewness(self, df: DataFrame, col_name: str) -> float:
        """Calculate skewness of a numeric column."""
        values = df.select(col_name).rdd.map(lambda x: float(x[0])).collect()
        return float(stats.skew(values))

    def _calculate_kurtosis(self, df: DataFrame, col_name: str) -> float:
        """Calculate kurtosis of a numeric column."""
        values = df.select(col_name).rdd.map(lambda x: float(x[0])).collect()
        return float(stats.kurtosis(values))

    def _get_top_values(self, df: DataFrame, col_name: str, limit: int = 10) -> List[Dict]:
        """Get top N most frequent values in a column."""
        return df.groupBy(col_name).count().orderBy("count", ascending=False) \
            .limit(limit).collect()

    def _calculate_date_range(self, df: DataFrame, col_name: str) -> Dict:
        """Calculate date range statistics."""
        min_date = df.select(min(col(col_name))).collect()[0][0]
        max_date = df.select(max(col(col_name))).collect()[0][0]
        return {
            "min": min_date,
            "max": max_date,
            "range_days": (max_date - min_date).days
        } 