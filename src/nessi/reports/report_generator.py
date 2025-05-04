"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Base report generator class for nessi.dev."""

import os
import json
import csv
import logging
from abc import ABC, abstractmethod
from datetime import datetime
from typing import Dict, Any, Optional, Union
from pathlib import Path
import pdfkit
from jinja2 import Environment, FileSystemLoader, select_autoescape
from reportlab.lib import colors
from reportlab.lib.pagesizes import letter
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class ReportGenerator(ABC):
    """Base class for all report generators."""
    
    def __init__(self, template_dir: str = "templates", output_dir: str = "reports"):
        """Initialize the report generator.
        
        Args:
            template_dir (str): Directory containing report templates
            output_dir (str): Directory where reports will be saved
        """
        self.template_dir = Path(template_dir)
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Set up Jinja2 environment
        self.env = Environment(
            loader=FileSystemLoader(str(self.template_dir)),
            autoescape=select_autoescape(['html', 'xml']),
            trim_blocks=True,
            lstrip_blocks=True
        )
        
        # Add custom filters and globals
        self.env.globals['now'] = datetime.now
        
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

    def generate_report(self, data: Dict[str, Any], format: str = "html", 
                       output_path: Optional[str] = None, **kwargs) -> str:
        """Generate a report in the specified format.
        
        Args:
            data: Data to include in the report
            format: Output format (html, pdf, json, csv)
            output_path: Optional path to save the report
            **kwargs: Additional format-specific options
            
        Returns:
            str: Generated report content or path to saved file
        """
        # Validate input data
        self._validate_data(data)
        
        # Generate based on format
        if format == "html":
            return self._generate_html(data, output_path)
        elif format == "pdf":
            return self._generate_pdf(data, output_path)
        elif format == "json":
            return self._generate_json(data, output_path)
        elif format == "csv":
            return self.export_to_csv(data, output_path)
        else:
            raise ValueError(f"Unsupported format: {format}")
            
    def _validate_data(self, data: Dict[str, Any]) -> None:
        """Validate the input data structure.
        
        Args:
            data: Data to validate
            
        Raises:
            ValueError: If data is invalid
        """
        required_fields = ["title"]
        missing_fields = [field for field in required_fields if field not in data]
        if missing_fields:
            raise ValueError(f"Missing required fields: {', '.join(missing_fields)}")
            
        # Validate numeric fields
        if "summary" in data:
            for key, value in data["summary"].items():
                if isinstance(value, (int, float)) and value < 0:
                    raise ValueError(f"Invalid negative value for {key}: {value}")
                    
    def _generate_html(self, data: Dict[str, Any], output_path: Optional[str] = None) -> str:
        """Generate HTML report.
        
        Args:
            data: Report data
            output_path: Optional path to save the HTML file
            
        Returns:
            str: HTML content or path to saved file
        """
        template = self.env.get_template("report_template.html")
        html_content = template.render(**data)
        
        if output_path:
            output_path = Path(output_path)
            output_path.write_text(html_content)
            return str(output_path)
        return html_content
        
    def _generate_pdf(self, data: Dict[str, Any], output_path: Optional[str] = None) -> str:
        """Generate PDF report using wkhtmltopdf or ReportLab as fallback.
        
        Args:
            data: Report data
            output_path: Path to save the PDF file
            
        Returns:
            str: Path to the generated PDF file
        """
        if not output_path:
            output_path = self.output_dir / f"report_{self._get_timestamp()}.pdf"
        else:
            output_path = Path(output_path)
            
        try:
            # Try wkhtmltopdf first
            html_content = self._generate_html(data)
            pdfkit.from_string(html_content, str(output_path))
        except Exception as e:
            logger.warning(f"Failed to generate PDF with wkhtmltopdf: {e}. Falling back to ReportLab.")
            self._generate_pdf_reportlab(data, output_path)
            
        return str(output_path)
        
    def _generate_pdf_reportlab(self, data: Dict[str, Any], output_path: Path) -> None:
        """Generate PDF report using ReportLab (fallback method).
        
        Args:
            data: Report data
            output_path: Path to save the PDF file
        """
        doc = SimpleDocTemplate(str(output_path), pagesize=letter)
        styles = getSampleStyleSheet()
        story = []
        
        # Add title
        title_style = ParagraphStyle(
            'CustomTitle',
            parent=styles['Heading1'],
            fontSize=24,
            spaceAfter=30
        )
        story.append(Paragraph(data["title"], title_style))
        story.append(Spacer(1, 12))
        
        # Add summary if present
        if "summary" in data:
            story.append(Paragraph("Summary", styles['Heading2']))
            summary_data = [[key, str(value)] for key, value in data["summary"].items()]
            summary_table = Table(summary_data, colWidths=[2*inch, 4*inch])
            summary_table.setStyle(TableStyle([
                ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                ('GRID', (0, 0), (-1, -1), 1, colors.black)
            ]))
            story.append(summary_table)
            story.append(Spacer(1, 12))
            
        # Add metrics if present
        if "metrics" in data:
            story.append(Paragraph("Metrics", styles['Heading2']))
            for category, metrics in data["metrics"].items():
                story.append(Paragraph(category, styles['Heading3']))
                metrics_data = [[key, str(value)] for key, value in metrics.items()]
                metrics_table = Table(metrics_data, colWidths=[2*inch, 4*inch])
                metrics_table.setStyle(TableStyle([
                    ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                    ('GRID', (0, 0), (-1, -1), 1, colors.black)
                ]))
                story.append(metrics_table)
                story.append(Spacer(1, 12))
                
        doc.build(story)
        
    def _generate_json(self, data: Dict[str, Any], output_path: Optional[str] = None) -> str:
        """Generate JSON report.
        
        Args:
            data: Report data
            output_path: Optional path to save the JSON file
            
        Returns:
            str: JSON string or path to saved file
        """
        json_content = json.dumps(data, indent=2)
        
        if output_path:
            output_path = Path(output_path)
            output_path.write_text(json_content)
            return str(output_path)
        return json_content
        
    def export_to_csv(self, data: Dict[str, Any], output_path: str) -> str:
        """Export data to CSV format.
        
        Args:
            data: Data to export (must be flat dictionary of lists)
            output_path: Path to save the CSV file
            
        Returns:
            str: Path to the saved CSV file
        """
        output_path = Path(output_path)
        
        # Validate data structure
        if not all(isinstance(v, list) for v in data.values()):
            raise ValueError("All values must be lists for CSV export")
            
        # Get lengths of all lists
        lengths = {k: len(v) for k, v in data.items()}
        if len(set(lengths.values())) > 1:
            raise ValueError("All lists must have the same length")
            
        # Write CSV
        with open(output_path, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(data.keys())  # Header
            writer.writerows(zip(*data.values()))  # Data
            
        return str(output_path)
        
    def _get_timestamp(self) -> str:
        """Get current timestamp for file naming.
        
        Returns:
            str: Timestamp in YYYYMMDD_HHMMSS format
        """
        return datetime.now().strftime("%Y%m%d_%H%M%S")