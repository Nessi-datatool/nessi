"""Base report generator class for Nessi."""

import os
import json
import logging
from abc import ABC, abstractmethod
from datetime import datetime
from typing import Dict, Any, Optional, Union
from pathlib import Path

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class ReportGenerator(ABC):
    """Base class for all report generators."""
    
    def __init__(self, output_dir: str = "reports"):
        """Initialize the report generator.
        
        Args:
            output_dir (str): Directory where reports will be saved.
        """
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
    @abstractmethod
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a report from the provided data.
        
        Args:
            data (Dict[str, Any]): Data to generate the report from.
            **kwargs: Additional arguments for report generation.
            
        Returns:
            str: Path to the generated report.
        """
        pass
    
    def save_report(self, content: Union[str, Dict[str, Any]], 
                   filename: str, format: str = "md") -> str:
        """Save the report content to a file.
        
        Args:
            content: Report content to save.
            filename: Name of the file without extension.
            format: Format of the report (md, json, html, csv).
            
        Returns:
            str: Path to the saved report.
        """
        try:
            # Create output path
            output_path = self.output_dir / f"{filename}.{format}"
            
            # Save based on format
            if format == "json":
                with open(output_path, 'w') as f:
                    json.dump(content, f, indent=2)
            else:
                with open(output_path, 'w') as f:
                    f.write(content)
                    
            logger.info(f"Report saved to {output_path}")
            return str(output_path)
            
        except Exception as e:
            logger.error(f"Failed to save report: {str(e)}")
            raise RuntimeError(f"Failed to save report: {str(e)}")
    
    def get_timestamp(self) -> str:
        """Get current timestamp for report naming.
        
        Returns:
            str: Current timestamp in YYYYMMDD_HHMMSS format.
        """
        return datetime.now().strftime("%Y%m%d_%H%M%S")
    
    def validate_data(self, data: Dict[str, Any], required_fields: list) -> bool:
        """Validate that required fields are present in the data.
        
        Args:
            data: Data to validate.
            required_fields: List of required field names.
            
        Returns:
            bool: True if all required fields are present.
        """
        missing_fields = [field for field in required_fields if field not in data]
        if missing_fields:
            raise ValueError(f"Missing required fields: {', '.join(missing_fields)}")
        return True 