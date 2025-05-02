from dataclasses import dataclass
from typing import Dict, List, Any, Optional

@dataclass
class ColumnStats:
    """Statistics for a column"""
    null_count: int
    distinct_count: int
    min_value: Any
    max_value: Any
    avg_value: Optional[float] = None
    std_value: Optional[float] = None

@dataclass
class ColumnMetadata:
    """Metadata for a column"""
    name: str
    type: str
    stats: Optional[ColumnStats] = None

@dataclass
class TableScanResult:
    """Result of scanning a table"""
    format: str
    columns: List[ColumnMetadata]
    row_count: int
    size: int
    compression: Optional[str] = None
    partition_columns: List[str] = None
    quality_checks: Dict[str, Any] = None 