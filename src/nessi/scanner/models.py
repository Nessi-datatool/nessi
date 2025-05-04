"""Models for the scanner module."""
from dataclasses import dataclass
from typing import Dict, List, Optional

from pyspark.sql.types import StructType


@dataclass
class SchemaValidationResult:
    """Result of schema validation."""
    is_valid: bool
    missing_fields: List[str]
    extra_fields: List[str]
    type_mismatches: Dict[str, str]
    schema: StructType


@dataclass
class ScanResult:
    """Result of scanning a table."""
    table_name: str
    row_count: int
    column_count: int
    schema: StructType
    column_stats: Dict[str, Dict[str, float]]
    data_quality_score: Optional[float] = None
    anomalies: Optional[Dict[str, List[int]]] = None
    schema_validation: Optional[SchemaValidationResult] = None


@dataclass
class QualityCheckResult:
    """Result of a quality check."""
    check_name: str
    passed: bool
    score: float
    details: Dict[str, any]
    column: Optional[str] = None
    threshold: Optional[float] = None
    actual_value: Optional[float] = None 