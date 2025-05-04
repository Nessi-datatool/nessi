"""Configuration for quality checks."""

from dataclasses import dataclass
from typing import Dict, List, Optional, Union


@dataclass
class QualityConfig:
    """Configuration for quality checks."""

    # Thresholds for data quality checks
    min_completeness: float = 0.95  # Minimum completeness ratio (1 - null ratio)
    min_uniqueness: float = 0.0  # Minimum uniqueness ratio for key columns
    max_duplicates: float = 0.05  # Maximum ratio of duplicate rows
    min_data_points: int = 100  # Minimum number of data points required
    max_missing_values: float = 0.05  # Maximum ratio of missing values allowed

    # Column-specific configurations
    column_types: Dict[str, str] = None  # Expected data types for columns
    required_columns: List[str] = None  # List of required columns
    unique_columns: List[str] = None  # List of columns that should be unique
    value_ranges: Dict[str, Dict[str, Union[int, float]]] = None  # Valid ranges for numeric columns
    allowed_values: Dict[str, List[str]] = None  # Valid values for categorical columns

    # Time-based configurations
    max_data_age_days: Optional[int] = None  # Maximum age of data in days
    time_column: Optional[str] = None  # Column containing timestamps
    time_format: Optional[str] = None  # Format of timestamp column

    def __post_init__(self):
        """Initialize default values for None fields."""
        if self.column_types is None:
            self.column_types = {}
        if self.required_columns is None:
            self.required_columns = []
        if self.unique_columns is None:
            self.unique_columns = []
        if self.value_ranges is None:
            self.value_ranges = {}
        if self.allowed_values is None:
            self.allowed_values = {} 