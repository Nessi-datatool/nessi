"""Custom report generator for Nessi."""

import os
from typing import Dict, Any, Optional
from jinja2 import Environment, FileSystemLoader, select_autoescape
from .report_generator import ReportGenerator

class CustomReport(ReportGenerator):
    """Generates custom reports using templates."""
    
    def __init__(self, output_dir: str = "reports/custom", template_dir: str = "templates"):
        """Initialize the custom report generator.
        
        Args:
            output_dir (str): Directory where reports will be saved.
            template_dir (str): Directory containing report templates.
        """
        super().__init__(output_dir)
        self.template_dir = Path(template_dir)
        self.template_dir.mkdir(parents=True, exist_ok=True)
        
        # Initialize Jinja2 environment
        self.env = Environment(
            loader=FileSystemLoader(str(self.template_dir)),
            autoescape=select_autoescape(['html', 'xml'])
        )
        
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a custom report using a template.
        
        Args:
            data (Dict[str, Any]): Data to include in the report.
            **kwargs: Additional arguments:
                - template (str): Template file name
                - format (str): Output format (md, json, html)
                - template_vars (Dict[str, Any]): Additional template variables
                
        Returns:
            str: Path to the generated report.
        """
        try:
            # Get template and format
            template_name = kwargs.get('template', 'default.html')
            format = kwargs.get('format', 'html')
            template_vars = kwargs.get('template_vars', {})
            
            # Load and render template
            template = self.env.get_template(template_name)
            content = template.render(data=data, **template_vars)
            
            # Save report
            timestamp = self.get_timestamp()
            filename = f"custom_report_{timestamp}"
            return self.save_report(content, filename, format)
            
        except Exception as e:
            logger.error(f"Failed to generate custom report: {str(e)}")
            raise RuntimeError(f"Failed to generate custom report: {str(e)}")
    
    def create_template(self, name: str, content: str) -> str:
        """Create a new report template.
        
        Args:
            name (str): Template file name.
            content (str): Template content.
            
        Returns:
            str: Path to the created template.
        """
        try:
            template_path = self.template_dir / name
            with open(template_path, 'w') as f:
                f.write(content)
            return str(template_path)
        except Exception as e:
            logger.error(f"Failed to create template: {str(e)}")
            raise RuntimeError(f"Failed to create template: {str(e)}")
    
    def list_templates(self) -> list:
        """List available report templates.
        
        Returns:
            list: List of template file names.
        """
        return [f.name for f in self.template_dir.glob('*') if f.is_file()]
    
    def get_template(self, name: str) -> Optional[str]:
        """Get the content of a template.
        
        Args:
            name (str): Template file name.
            
        Returns:
            Optional[str]: Template content if found, None otherwise.
        """
        try:
            template_path = self.template_dir / name
            if template_path.exists():
                with open(template_path, 'r') as f:
                    return f.read()
            return None
        except Exception as e:
            logger.error(f"Failed to get template: {str(e)}")
            raise RuntimeError(f"Failed to get template: {str(e)}")
    
    def delete_template(self, name: str) -> bool:
        """Delete a report template.
        
        Args:
            name (str): Template file name.
            
        Returns:
            bool: True if template was deleted, False otherwise.
        """
        try:
            template_path = self.template_dir / name
            if template_path.exists():
                template_path.unlink()
                return True
            return False
        except Exception as e:
            logger.error(f"Failed to delete template: {str(e)}")
            raise RuntimeError(f"Failed to delete template: {str(e)}") 