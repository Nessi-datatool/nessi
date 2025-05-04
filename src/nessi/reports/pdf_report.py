"""
PDF report generator for nessi.dev.
"""

from reportlab.lib import colors
from reportlab.lib.pagesizes import letter
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import inch
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List, Optional
import logging
from .report_generator import ReportGenerator

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class PDFReport(ReportGenerator):
    """Generate PDF reports."""
    
    def __init__(self, template_dir: str = "templates", output_dir: str = "reports"):
        """Initialize the PDF report generator.
        
        Args:
            template_dir (str): Directory containing report templates
            output_dir (str): Directory where reports will be saved
        """
        super().__init__(template_dir, output_dir)
        self.styles = getSampleStyleSheet()
        self._add_custom_styles()
    
    def _add_custom_styles(self):
        """Add custom styles for PDF reports."""
        self.styles.add(ParagraphStyle(
            name='CustomTitle',
            parent=self.styles['Title'],
            fontSize=24,
            spaceAfter=30
        ))
        self.styles.add(ParagraphStyle(
            name='CustomHeading1',
            parent=self.styles['Heading1'],
            fontSize=18,
            spaceAfter=20
        ))
        self.styles.add(ParagraphStyle(
            name='CustomHeading2',
            parent=self.styles['Heading2'],
            fontSize=14,
            spaceAfter=15
        ))
        self.styles.add(ParagraphStyle(
            name='CustomNormal',
            parent=self.styles['Normal'],
            fontSize=12,
            spaceAfter=10
        ))
    
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a PDF report.
        
        Args:
            data: Data to include in the report
            **kwargs: Additional options for report generation
            
        Returns:
            str: Path to the generated PDF file
        """
        # Process data for PDF generation
        processed_data = self._process_data(data)
        
        # Generate PDF
        output_path = kwargs.get('output_path')
        if not output_path:
            output_path = self.output_dir / f"report_{self._get_timestamp()}.pdf"
        else:
            output_path = Path(output_path)
            
        return self._generate_pdf(processed_data, output_path)
        
    def _process_data(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Process data for PDF generation.
        
        Args:
            data: Raw report data
            
        Returns:
            Dict[str, Any]: Processed data ready for PDF generation
        """
        processed = {
            "title": data.get("title", "Report"),
            "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            "summary": data.get("summary", {}),
            "metrics": data.get("metrics", {}),
            "quality_checks": data.get("quality_checks", []),
            "performance": data.get("performance", {}),
            "recommendations": data.get("recommendations", [])
        }
        
        # Add visualizations if available
        if "visualizations" in data:
            processed["visualizations"] = data["visualizations"]
            
        return processed
        
    def _generate_pdf(self, data: Dict[str, Any], output_path: Path) -> str:
        """Generate PDF report using ReportLab.
        
        Args:
            data: Processed report data
            output_path: Path to save the PDF file
            
        Returns:
            str: Path to the generated PDF file
        """
        doc = SimpleDocTemplate(str(output_path), pagesize=letter)
        story = []
        
        # Add title
        story.append(Paragraph(data["title"], self.styles['CustomTitle']))
        story.append(Spacer(1, 12))
        
        # Add timestamp
        story.append(Paragraph(f"Generated on {data['timestamp']}", self.styles['CustomNormal']))
        story.append(Spacer(1, 20))
        
        # Add summary if present
        if "summary" in data and data["summary"]:
            story.append(Paragraph("Summary", self.styles['CustomHeading1']))
            summary_data = [[key, str(value)] for key, value in data["summary"].items()]
            summary_table = Table(summary_data, colWidths=[2*inch, 4*inch])
            summary_table.setStyle(TableStyle([
                ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                ('GRID', (0, 0), (-1, -1), 1, colors.black)
            ]))
            story.append(summary_table)
            story.append(Spacer(1, 20))
            
        # Add metrics if present
        if "metrics" in data and data["metrics"]:
            story.append(Paragraph("Metrics", self.styles['CustomHeading1']))
            for category, metrics in data["metrics"].items():
                story.append(Paragraph(category, self.styles['CustomHeading2']))
                metrics_data = [[key, str(value)] for key, value in metrics.items()]
                metrics_table = Table(metrics_data, colWidths=[2*inch, 4*inch])
                metrics_table.setStyle(TableStyle([
                    ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                    ('GRID', (0, 0), (-1, -1), 1, colors.black)
                ]))
                story.append(metrics_table)
                story.append(Spacer(1, 12))
                
        # Add quality checks if present
        if "quality_checks" in data and data["quality_checks"]:
            story.append(Paragraph("Quality Checks", self.styles['CustomHeading1']))
            for check in data["quality_checks"]:
                status = "PASSED" if check.get("passed", False) else "FAILED"
                status_color = colors.green if check.get("passed", False) else colors.red
                
                story.append(Paragraph(
                    f"{check.get('column', '')} - {check.get('rule', '')}",
                    self.styles['CustomHeading2']
                ))
                
                status_style = ParagraphStyle(
                    'Status',
                    parent=self.styles['CustomNormal'],
                    textColor=status_color
                )
                story.append(Paragraph(status, status_style))
                
                if "details" in check:
                    details_data = [[key, str(value)] for key, value in check["details"].items()]
                    details_table = Table(details_data, colWidths=[2*inch, 4*inch])
                    details_table.setStyle(TableStyle([
                        ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                        ('GRID', (0, 0), (-1, -1), 1, colors.black)
                    ]))
                    story.append(details_table)
                    
                story.append(Spacer(1, 12))
                
        # Add performance metrics if present
        if "performance" in data and data["performance"]:
            story.append(Paragraph("Performance", self.styles['CustomHeading1']))
            performance_data = [[key, str(value)] for key, value in data["performance"].items()]
            performance_table = Table(performance_data, colWidths=[2*inch, 4*inch])
            performance_table.setStyle(TableStyle([
                ('BACKGROUND', (0, 0), (-1, -1), colors.beige),
                ('GRID', (0, 0), (-1, -1), 1, colors.black)
            ]))
            story.append(performance_table)
            story.append(Spacer(1, 20))
            
        # Add recommendations if present
        if "recommendations" in data and data["recommendations"]:
            story.append(Paragraph("Recommendations", self.styles['CustomHeading1']))
            for rec in data["recommendations"]:
                story.append(Paragraph(f"• {rec}", self.styles['CustomNormal']))
            story.append(Spacer(1, 20))
            
        # Build the PDF
        doc.build(story)
        
        return str(output_path) 