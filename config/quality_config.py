from typing import Dict, Any, Optional, List
from dataclasses import dataclass, field
from datetime import timedelta

@dataclass
class QualityConfig:
    """Configuration for data quality analysis."""
    
    # Analysis settings
    sample_size: int = 10000
    analysis_timeout: int = 300  # seconds
    
    # Scoring weights
    completeness_weight: float = 0.4
    consistency_weight: float = 0.3
    anomaly_weight: float = 0.2
    pattern_weight: float = 0.1
    
    # Anomaly detection
    anomaly_contamination: float = 0.1
    anomaly_threshold: float = 0.1  # bottom 10% considered anomalies
    
    # Pattern detection
    pattern_threshold: float = 0.8  # 80% match required for pattern detection
    custom_patterns: Dict[str, str] = field(default_factory=dict)
    
    # Validation rules
    default_rules: Dict[str, List[Dict[str, Any]]] = field(default_factory=dict)
    
    # Logging
    log_level: str = "INFO"
    log_format: str = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    
    @classmethod
    def from_dict(cls, config_dict: Optional[Dict[str, Any]] = None) -> "QualityConfig":
        """Create a QualityConfig instance from a dictionary."""
        if config_dict is None:
            return cls()
        
        return cls(**{
            k: v for k, v in config_dict.items()
            if k in cls.__dataclass_fields__
        })
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert the configuration to a dictionary."""
        return {
            k: v for k, v in self.__dict__.items()
            if not k.startswith('_')
        } 