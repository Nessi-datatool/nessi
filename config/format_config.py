from typing import Dict, Any
from pydantic import BaseSettings

class FormatConfig(BaseSettings):
    """Configuration settings for data formats."""
    
    # Delta Lake settings
    delta_auto_optimize: bool = True
    delta_optimize_interval: int = 3600  # seconds
    delta_auto_compact: bool = True
    delta_compact_interval: int = 86400  # seconds
    
    # Parquet settings
    parquet_compression: str = "snappy"
    parquet_row_group_size: int = 128 * 1024 * 1024  # 128MB
    parquet_page_size: int = 1 * 1024 * 1024  # 1MB
    
    # CSV settings
    csv_encoding: str = "utf-8"
    csv_delimiter: str = ","
    csv_quote_char: str = '"'
    csv_escape_char: str = "\\"
    
    # Schema inference settings
    schema_sample_size: int = 1000
    schema_inference_timeout: int = 30  # seconds
    
    # Format-specific options
    format_options: Dict[str, Dict[str, Any]] = {
        "delta": {
            "mergeSchema": True,
            "overwriteSchema": False
        },
        "parquet": {
            "coerce_timestamps": "ms",
            "allow_truncated_timestamps": False
        },
        "csv": {
            "infer_datetime_format": True,
            "keep_default_na": True
        }
    }
    
    class Config:
        env_prefix = "NESSI_FORMAT_" 