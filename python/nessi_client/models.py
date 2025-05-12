"""
Data models for the Nessi monitoring client.
"""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Dict, List, Optional, Any, Union


@dataclass
class Metric:
    """Represents a monitoring metric."""
    name: str
    value: Union[float, int, str]
    timestamp: datetime = field(default_factory=datetime.now)
    tags: Dict[str, str] = field(default_factory=dict)
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class Alert:
    """Represents a monitoring alert."""
    id: str
    name: str
    description: str
    severity: str
    triggered: bool
    timestamp: datetime = field(default_factory=datetime.now)
    metric_name: Optional[str] = None
    metric_value: Optional[Any] = None
    threshold: Optional[Any] = None
    tags: Dict[str, str] = field(default_factory=dict)
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class Rule:
    """Represents a data quality rule."""
    id: str
    name: str
    description: str
    severity: str
    field: Optional[str] = None
    rule_type: str = "custom"
    config: Dict[str, Any] = field(default_factory=dict)
    tags: List[str] = field(default_factory=list)
    enabled: bool = True


@dataclass
class ValidationResult:
    """Represents the result of a data quality validation."""
    rule_id: str
    field: Optional[str]
    value: Optional[Any]
    message: str
    row_index: Optional[int] = None
    passed: bool = False
    timestamp: datetime = field(default_factory=datetime.now)


@dataclass
class Profile:
    """Represents a data profile."""
    id: str
    name: str
    dataset_name: str
    timestamp: datetime = field(default_factory=datetime.now)
    row_count: int = 0
    column_count: int = 0
    column_stats: Dict[str, Dict[str, Any]] = field(default_factory=dict)
    outliers: Dict[str, List[Any]] = field(default_factory=dict)
    metadata: Dict[str, Any] = field(default_factory=dict)
    detection_method: str = "zscore"
