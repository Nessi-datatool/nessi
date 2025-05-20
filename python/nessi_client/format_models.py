"""
Format handling models for the Nessi monitoring client.
"""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Dict, List, Optional, Any, Union


@dataclass
class SchemaField:
    """Represents a field in a data schema."""
    name: str
    data_type: str
    nullable: bool = True
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class Schema:
    """Represents a data schema."""
    fields: List[SchemaField]
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class FormatConfig:
    """Configuration for format handling."""
    date_formats: List[str] = field(default_factory=lambda: ["%Y-%m-%d", "%Y/%m/%d", "%d-%m-%Y", "%d/%m/%Y"])
    csv_delimiter: str = ","
    csv_has_header: bool = True
    max_rows_for_inference: int = 1000
    min_confidence_threshold: float = 0.7


@dataclass
class FormatDetectionResult:
    """Result of format detection."""
    format: str
    confidence: float
    schema: Optional[Schema] = None
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class DataBatch:
    """Represents a batch of data with schema information."""
    data: List[Dict[str, Any]]
    schema: Schema
    format: str
    row_count: int
    metadata: Dict[str, Any] = field(default_factory=dict)
