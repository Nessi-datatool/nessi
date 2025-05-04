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
from statsmodels.tsa.seasonal import seasonal_decompose
from statsmodels.tsa.stattools import adfuller, acf, pacf

logger = logging.getLogger(__name__)

class DataQualityAnalyzer:
    """Analyzes data quality and provides comprehensive profiling."""
    
    def __init__(self, config: Optional[Dict[str, Any]] = None):
        self.config = config or {}
        self._setup_logging()
    
    def _setup_logging(self):
        """Setup logging configuration."""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    
    def analyze_dataframe(self, df: Union[pd.DataFrame, SparkDataFrame]) -> Dict[str, Any]:
        """Perform comprehensive data quality analysis on a DataFrame."""
        try:
            # Convert PySpark DataFrame to pandas if needed
            if isinstance(df, SparkDataFrame):
                df = df.toPandas()
            
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
            
            # Detect anomalies
            analysis["anomalies"] = self._detect_anomalies(df)
            
            # Detect patterns
            analysis["patterns"] = self._detect_patterns(df)
            
            return analysis
        
        except Exception as e:
            logger.error(f"Error analyzing DataFrame: {str(e)}")
            raise
    
    def _analyze_column(self, series: pd.Series) -> Dict[str, Any]:
        """Analyze a single column of data."""
        analysis = {
            "type": str(series.dtype),
            "null_count": series.isnull().sum(),
            "unique_count": series.nunique(),
            "value_counts": {},
            "statistics": {},
            "patterns": []
        }
        
        # Calculate basic statistics
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
        
        # Calculate value distribution
        value_counts = series.value_counts().head(10)
        analysis["value_counts"] = value_counts.to_dict()
        
        # Detect patterns
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
        """Detect anomalies using multiple ML techniques."""
        anomalies = []
        
        # Isolation Forest for point anomalies
        from sklearn.ensemble import IsolationForest
        clf = IsolationForest(contamination=0.1, random_state=42)
        
        for column in df.columns:
            if pd.api.types.is_numeric_dtype(df[column]):
                data = df[column].dropna().values.reshape(-1, 1)
                if len(data) < 2:
                    continue
                
                clf.fit(data)
                scores = clf.decision_function(data)
                threshold = np.percentile(scores, 10)
                
                anomaly_indices = np.where(scores < threshold)[0]
                for idx in anomaly_indices:
                    anomalies.append({
                        'type': 'point_anomaly',
                        'column': column,
                        'row': int(df.index[idx]),
                        'value': float(data[idx][0]),
                        'score': float(scores[idx])
                    })
        
        # LSTM for sequence anomalies
        if 'timestamp' in df.columns:
            anomalies.extend(self._detect_sequence_anomalies(df))
        
        return anomalies
    
    def _detect_string_patterns(self, series: pd.Series) -> List[Dict[str, Any]]:
        """Detect patterns in string data."""
        patterns = []
        
        # Detect common formats
        formats = {
            'email': r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
            'phone': r'^\+?1?\d{9,15}$',
            'url': r'^https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+',
            'date': r'^\d{4}-\d{2}-\d{2}$'
        }
        
        for format_name, pattern in formats.items():
            matches = series.str.match(pattern, na=False)
            if matches.any():
                patterns.append({
                    'type': 'format',
                    'format': format_name,
                    'count': matches.sum(),
                    'percentage': matches.sum() / len(series) * 100
                })
        
        return patterns
    
    def _detect_patterns(self, df: pd.DataFrame) -> List[Dict[str, Any]]:
        """Detect patterns in the data using advanced techniques."""
        patterns = []
        
        # Time series patterns
        if 'timestamp' in df.columns:
            patterns.extend(self._detect_time_series_patterns(df))
        
        # Value patterns
        for column in df.columns:
            if pd.api.types.is_numeric_dtype(df[column]):
                patterns.extend(self._detect_numeric_patterns(df[column]))
            elif pd.api.types.is_string_dtype(df[column]):
                patterns.extend(self._detect_string_patterns(df[column]))
        
        return patterns

    def _detect_time_series_patterns(self, df: pd.DataFrame) -> List[Dict[str, Any]]:
        """Detect patterns in time series data."""
        patterns = []
        
        # Resample to consistent frequency
        df_resampled = df.set_index('timestamp').resample('1H').mean()
        
        # Detect seasonality
        decomposition = seasonal_decompose(df_resampled, model='additive', period=24)
        
        # Calculate seasonality strength
        seasonal_strength = 1 - np.var(decomposition.resid) / np.var(decomposition.seasonal + decomposition.resid)
        
        if seasonal_strength > 0.5:
            patterns.append({
                'type': 'seasonality',
                'column': 'timestamp',
                'period': 24,
                'strength': seasonal_strength
            })
        
        # Detect trends
        result = adfuller(df_resampled)
        if result[1] > 0.05:
            patterns.append({
                'type': 'trend',
                'column': 'timestamp',
                'direction': 'increasing' if df_resampled.diff().mean() > 0 else 'decreasing'
            })
        
        return patterns

    def _detect_numeric_patterns(self, series: pd.Series) -> List[Dict[str, Any]]:
        """Detect patterns in numeric data."""
        patterns = []
        
        # Detect normal distribution
        k2, p = stats.normaltest(series.dropna())
        if p > 0.05:
            patterns.append({
                'type': 'distribution',
                'distribution': 'normal',
                'p_value': p
            })
        
        # Detect outliers using IQR
        Q1 = series.quantile(0.25)
        Q3 = series.quantile(0.75)
        IQR = Q3 - Q1
        outliers = series[(series < (Q1 - 1.5 * IQR)) | (series > (Q3 + 1.5 * IQR))]
        if len(outliers) > 0:
            patterns.append({
                'type': 'outliers',
                'count': len(outliers),
                'percentage': len(outliers) / len(series) * 100
            })
        
        return patterns

    def _detect_sequence_anomalies(self, df: pd.DataFrame) -> List[Dict[str, Any]]:
        """Detect anomalies in time series sequences using LSTM."""
        anomalies = []
        
        try:
            from tensorflow.keras.models import Sequential
            from tensorflow.keras.layers import LSTM, Dense
            
            # Prepare data
            sequence_length = 24
            data = df.set_index('timestamp').resample('1H').mean()
            data = data.fillna(method='ffill')
            
            # Create sequences
            X = []
            for i in range(len(data) - sequence_length):
                X.append(data[i:(i + sequence_length)])
            X = np.array(X)
            
            # Build and train LSTM model
            model = Sequential([
                LSTM(50, return_sequences=True, input_shape=(sequence_length, 1)),
                LSTM(50),
                Dense(1)
            ])
            model.compile(optimizer='adam', loss='mse')
            model.fit(X, data[sequence_length:], epochs=10, verbose=0)
            
            # Detect anomalies
            predictions = model.predict(X)
            mse = np.mean(np.power(data[sequence_length:] - predictions, 2), axis=1)
            threshold = np.percentile(mse, 95)
            
            anomaly_indices = np.where(mse > threshold)[0]
            for idx in anomaly_indices:
                anomalies.append({
                    'type': 'sequence_anomaly',
                    'timestamp': data.index[idx + sequence_length],
                    'score': float(mse[idx])
                })
        
        except Exception as e:
            logger.warning(f"Error in sequence anomaly detection: {str(e)}")
        
        return anomalies

    def validate_rules(self, df: pd.DataFrame, rules: Dict[str, List[Dict[str, Any]]]) -> Dict[str, Any]:
        """Validate data against custom rules."""
        validation_results = {
            "timestamp": datetime.now().isoformat(),
            "total_rules": len(rules),
            "passed_rules": 0,
            "failed_rules": [],
            "validation_score": 0.0
        }
        
        for column, column_rules in rules.items():
            if column not in df.columns:
                continue
            
            for rule in column_rules:
                rule_type = rule.get("type")
                rule_value = rule.get("value")
                rule_description = rule.get("description", "")
                
                try:
                    if rule_type == "not_null":
                        passed = df[column].notnull().all()
                    elif rule_type == "unique":
                        passed = df[column].nunique() == len(df)
                    elif rule_type == "min":
                        passed = df[column].min() >= rule_value
                    elif rule_type == "max":
                        passed = df[column].max() <= rule_value
                    elif rule_type == "regex":
                        passed = df[column].str.match(rule_value, na=False).all()
                    elif rule_type == "custom":
                        passed = eval(rule_value)(df[column])
                    else:
                        continue
                    
                    if passed:
                        validation_results["passed_rules"] += 1
                    else:
                        validation_results["failed_rules"].append({
                            "column": column,
                            "rule": rule_type,
                            "description": rule_description
                        })
                
                except Exception as e:
                    logger.error(f"Error validating rule: {str(e)}")
                    validation_results["failed_rules"].append({
                        "column": column,
                        "rule": rule_type,
                        "error": str(e)
                    })
        
        validation_results["validation_score"] = (
            validation_results["passed_rules"] / validation_results["total_rules"]
            if validation_results["total_rules"] > 0 else 1.0
        )
        
        return validation_results

    def detect_data_drift(self, df1: pd.DataFrame, df2: pd.DataFrame) -> Dict[str, Any]:
        """Detect data drift between two datasets."""
        drift_results = {
            'columns': {},
            'overall_drift': 0.0
        }
        
        # Calculate drift for each column
        for column in df1.columns:
            if column in df2.columns:
                if pd.api.types.is_numeric_dtype(df1[column]):
                    drift = self._calculate_numeric_drift(df1[column], df2[column])
                else:
                    drift = self._calculate_categorical_drift(df1[column], df2[column])
                
                drift_results['columns'][column] = drift
        
        # Calculate overall drift
        if drift_results['columns']:
            drift_results['overall_drift'] = np.mean(list(drift_results['columns'].values()))
        
        return drift_results

    def _calculate_numeric_drift(self, series1: pd.Series, series2: pd.Series) -> float:
        """Calculate drift between two numeric series."""
        from scipy import stats
        
        # Kolmogorov-Smirnov test
        _, p_value = stats.ks_2samp(series1.dropna(), series2.dropna())
        
        # Wasserstein distance
        from scipy.stats import wasserstein_distance
        w_distance = wasserstein_distance(series1.dropna(), series2.dropna())
        
        return 1 - min(p_value, 1 - w_distance)

    def _calculate_categorical_drift(self, series1: pd.Series, series2: pd.Series) -> float:
        """Calculate drift between two categorical series."""
        # Chi-square test
        from scipy.stats import chi2_contingency
        
        contingency = pd.crosstab(series1, series2)
        _, p_value, _, _ = chi2_contingency(contingency)
        
        # Jensen-Shannon divergence
        from scipy.spatial.distance import jensenshannon
        p = series1.value_counts(normalize=True)
        q = series2.value_counts(normalize=True)
        js_divergence = jensenshannon(p, q)
        
        return 1 - min(p_value, 1 - js_divergence)

    def _detect_anomalous_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect anomalous patterns in the data.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected anomalous patterns
        """
        patterns = []
        for column in df.columns:
            if pd.api.types.is_numeric_dtype(df[column]):
                # Use IQR method for anomaly detection
                Q1 = df[column].quantile(0.25)
                Q3 = df[column].quantile(0.75)
                IQR = Q3 - Q1
                outliers = df[column][(df[column] < (Q1 - 1.5 * IQR)) | (df[column] > (Q3 + 1.5 * IQR))]
                if len(outliers) > 0:
                    patterns.append({
                        "type": "anomaly",
                        "column": column,
                        "count": len(outliers),
                        "percentage": (len(outliers) / len(df)) * 100
                    })
        return patterns

    def _detect_correlation_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect correlation patterns between numeric columns.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected correlation patterns
        """
        patterns = []
        numeric_cols = df.select_dtypes(include=[np.number]).columns
        if len(numeric_cols) > 1:
            corr_matrix = df[numeric_cols].corr()
            for i in range(len(numeric_cols)):
                for j in range(i+1, len(numeric_cols)):
                    corr = corr_matrix.iloc[i,j]
                    if abs(corr) > 0.7:  # Strong correlation threshold
                        patterns.append({
                            "type": "correlation",
                            "columns": [numeric_cols[i], numeric_cols[j]],
                            "correlation": corr
                        })
        return patterns

    def _detect_duplicate_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect duplicate patterns in the data.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected duplicate patterns
        """
        duplicates = df.duplicated()
        if duplicates.any():
            return [{
                "type": "duplicates",
                "count": duplicates.sum(),
                "percentage": (duplicates.sum() / len(df)) * 100
            }]
        return []

    def _detect_missing_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect missing value patterns in the data.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected missing patterns
        """
        patterns = []
        for column in df.columns:
            missing = df[column].isnull()
            if missing.any():
                patterns.append({
                    "type": "missing",
                    "column": column,
                    "count": missing.sum(),
                    "percentage": (missing.sum() / len(df)) * 100
                })
        return patterns

    def _detect_inconsistent_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect inconsistent patterns in the data.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected inconsistent patterns
        """
        patterns = []
        for column in df.columns:
            if pd.api.types.is_string_dtype(df[column]):
                # Check for mixed case
                mixed_case = df[column].str.islower() != df[column].str.isupper()
                if mixed_case.any():
                    patterns.append({
                        "type": "inconsistent_case",
                        "column": column,
                        "count": mixed_case.sum(),
                        "percentage": (mixed_case.sum() / len(df)) * 100
                    })
        return patterns

    def _detect_invalid_patterns(self, df: pd.DataFrame) -> List[Dict]:
        """
        Detect invalid patterns in the data.
        
        Args:
            df: Input DataFrame
            
        Returns:
            List of detected invalid patterns
        """
        patterns = []
        for column in df.columns:
            if pd.api.types.is_numeric_dtype(df[column]):
                # Check for negative values in positive-only fields
                if df[column].min() < 0:
                    patterns.append({
                        "type": "invalid_negative",
                        "column": column,
                        "count": (df[column] < 0).sum(),
                        "percentage": ((df[column] < 0).sum() / len(df)) * 100
                    })
        return patterns

    def _detect_seasonal_patterns(self, series: pd.Series) -> Dict:
        """
        Detect seasonal patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing seasonal pattern information
        """
        try:
            # Decompose the time series
            decomposition = seasonal_decompose(series, period=12)
            return {
                "type": "seasonal",
                "seasonal_strength": np.std(decomposition.seasonal),
                "trend_strength": np.std(decomposition.trend),
                "residual_strength": np.std(decomposition.resid)
            }
        except:
            return {}

    def _detect_trend_patterns(self, series: pd.Series) -> Dict:
        """
        Detect trend patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing trend pattern information
        """
        try:
            # Use Mann-Kendall test for trend detection
            trend = stats.kendalltau(range(len(series)), series)[0]
            return {
                "type": "trend",
                "trend_strength": trend,
                "direction": "increasing" if trend > 0 else "decreasing"
            }
        except:
            return {}

    def _detect_cyclic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect cyclic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing cyclic pattern information
        """
        try:
            # Use autocorrelation to detect cycles
            acf_values = acf(series, nlags=50)
            peaks = np.where((acf_values[1:-1] > acf_values[:-2]) & 
                           (acf_values[1:-1] > acf_values[2:]))[0] + 1
            if len(peaks) > 0:
                return {
                    "type": "cyclic",
                    "cycle_length": peaks[0],
                    "cycle_strength": acf_values[peaks[0]]
                }
        except:
            pass
        return {}

    def _detect_stationary_patterns(self, series: pd.Series) -> Dict:
        """
        Detect stationarity in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing stationarity information
        """
        try:
            # Use Augmented Dickey-Fuller test
            result = adfuller(series)
            return {
                "type": "stationarity",
                "is_stationary": result[1] < 0.05,
                "p_value": result[1]
            }
        except:
            return {}

    def _detect_non_stationary_patterns(self, series: pd.Series) -> Dict:
        """
        Detect non-stationarity in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing non-stationarity information
        """
        try:
            # Use rolling statistics
            rolling_mean = series.rolling(window=12).mean()
            rolling_std = series.rolling(window=12).std()
            return {
                "type": "non_stationarity",
                "mean_variation": rolling_mean.std(),
                "std_variation": rolling_std.std()
            }
        except:
            return {}

    def _detect_auto_correlation_patterns(self, series: pd.Series) -> Dict:
        """
        Detect autocorrelation patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing autocorrelation information
        """
        try:
            # Calculate autocorrelation
            acf_values = acf(series, nlags=10)
            return {
                "type": "autocorrelation",
                "lag_1": acf_values[1],
                "significant_lags": sum(np.abs(acf_values[1:]) > 1.96/np.sqrt(len(series)))
            }
        except:
            return {}

    def _detect_partial_auto_correlation_patterns(self, series: pd.Series) -> Dict:
        """
        Detect partial autocorrelation patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing partial autocorrelation information
        """
        try:
            # Calculate partial autocorrelation
            pacf_values = pacf(series, nlags=10)
            return {
                "type": "partial_autocorrelation",
                "lag_1": pacf_values[1],
                "significant_lags": sum(np.abs(pacf_values[1:]) > 1.96/np.sqrt(len(series)))
            }
        except:
            return {}

    def _detect_spectral_patterns(self, series: pd.Series) -> Dict:
        """
        Detect spectral patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing spectral pattern information
        """
        try:
            # Use FFT for spectral analysis
            fft = np.fft.fft(series)
            freqs = np.fft.fftfreq(len(series))
            dominant_freq = freqs[np.argmax(np.abs(fft[1:]))+1]
            return {
                "type": "spectral",
                "dominant_frequency": abs(dominant_freq),
                "amplitude": np.max(np.abs(fft[1:]))
            }
        except:
            return {}

    def _detect_wavelet_patterns(self, series: pd.Series) -> Dict:
        """
        Detect wavelet patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing wavelet pattern information
        """
        # Placeholder for wavelet analysis
        return {}

    def _detect_fractal_patterns(self, series: pd.Series) -> Dict:
        """
        Detect fractal patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing fractal pattern information
        """
        # Placeholder for fractal analysis
        return {}

    def _detect_chaotic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect chaotic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing chaotic pattern information
        """
        # Placeholder for chaos analysis
        return {}

    def _detect_deterministic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect deterministic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing deterministic pattern information
        """
        # Placeholder for deterministic pattern analysis
        return {}

    def _detect_periodic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect periodic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing periodic pattern information
        """
        try:
            # Use FFT to detect periodicity
            fft = np.fft.fft(series)
            freqs = np.fft.fftfreq(len(series))
            peaks = np.where((np.abs(fft[1:-1]) > np.abs(fft[:-2])) & 
                           (np.abs(fft[1:-1]) > np.abs(fft[2:])))[0] + 1
            if len(peaks) > 0:
                return {
                    "type": "periodic",
                    "period": 1/abs(freqs[peaks[0]]),
                    "strength": np.abs(fft[peaks[0]])
                }
        except:
            pass
        return {}

    def _detect_aperiodic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect aperiodic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing aperiodic pattern information
        """
        # Placeholder for aperiodic pattern analysis
        return {}

    def _detect_stochastic_patterns(self, series: pd.Series) -> Dict:
        """
        Detect stochastic patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing stochastic pattern information
        """
        # Placeholder for stochastic pattern analysis
        return {}

    def _detect_random_patterns(self, series: pd.Series) -> Dict:
        """
        Detect random patterns in time series data.
        
        Args:
            series: Input time series
            
        Returns:
            Dictionary containing random pattern information
        """
        try:
            # Use runs test for randomness
            median = np.median(series)
            runs = np.sum(np.diff(series > median) != 0) + 1
            n1 = np.sum(series > median)
            n2 = len(series) - n1
            expected_runs = ((2 * n1 * n2) / len(series)) + 1
            return {
                "type": "random",
                "is_random": abs(runs - expected_runs) < 1.96,
                "runs": runs,
                "expected_runs": expected_runs
            }
        except:
            return {}

    def _detect_text_patterns(self, series: pd.Series) -> List[Dict]:
        """
        Detect patterns in text data.
        
        Args:
            series: Input series
            
        Returns:
            List of detected text patterns
        """
        patterns = []
        if pd.api.types.is_string_dtype(series):
            # Check for email pattern
            email_pattern = series.str.match(r'^[\w\.-]+@[\w\.-]+\.\w+$')
            if email_pattern.any():
                patterns.append({
                    "type": "format",
                    "format": "email",
                    "count": email_pattern.sum(),
                    "percentage": (email_pattern.sum() / len(series)) * 100
                })
        return patterns

    def _detect_date_patterns(self, series: pd.Series) -> List[Dict]:
        """
        Detect patterns in date data.
        
        Args:
            series: Input series
            
        Returns:
            List of detected date patterns
        """
        patterns = []
        if pd.api.types.is_datetime64_any_dtype(series):
            # Check for weekday pattern
            weekday_counts = series.dt.dayofweek.value_counts()
            if len(weekday_counts) > 0:
                patterns.append({
                    "type": "weekday_distribution",
                    "most_common": weekday_counts.index[0],
                    "count": weekday_counts.iloc[0]
                })
        return patterns

    def _detect_categorical_patterns(self, series: pd.Series) -> List[Dict]:
        """
        Detect patterns in categorical data.
        
        Args:
            series: Input series
            
        Returns:
            List of detected categorical patterns
        """
        patterns = []
        if pd.api.types.is_categorical_dtype(series) or pd.api.types.is_object_dtype(series):
            value_counts = series.value_counts()
            if len(value_counts) > 0:
                patterns.append({
                    "type": "category_distribution",
                    "most_common": value_counts.index[0],
                    "count": value_counts.iloc[0],
                    "unique_values": len(value_counts)
                })
        return patterns

    def _detect_outlier_patterns(self, series: pd.Series) -> List[Dict]:
        """
        Detect outlier patterns in numeric data.
        
        Args:
            series: Input series
            
        Returns:
            List of detected outlier patterns
        """
        patterns = []
        if pd.api.types.is_numeric_dtype(series):
            # Use Z-score method
            z_scores = np.abs(stats.zscore(series))
            outliers = z_scores > 3
            if outliers.any():
                patterns.append({
                    "type": "outliers",
                    "count": outliers.sum(),
                    "percentage": (outliers.sum() / len(series)) * 100
                })
        return patterns 