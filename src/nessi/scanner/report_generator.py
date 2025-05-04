"""Report generator for scanner results."""
from typing import Dict, Any, Optional
import json
import os
from datetime import datetime
from pathlib import Path
import logging
from jinja2 import Environment, FileSystemLoader

class ReportGenerator:
    def __init__(self, template_dir: Optional[str] = None):
        self.logger = logging.getLogger(__name__)
        self.template_dir = template_dir or os.path.join(os.path.dirname(__file__), "templates")
        self.env = Environment(loader=FileSystemLoader(self.template_dir))

    def generate_report(self, scan_results: Dict[str, Any], output_format: str = "html") -> str:
        """Generate a report from scan results."""
        try:
            if output_format == "html":
                return self._generate_html_report(scan_results)
            elif output_format == "json":
                return self._generate_json_report(scan_results)
            else:
                raise ValueError(f"Unsupported output format: {output_format}")
        except Exception as e:
            self.logger.error(f"Error generating report: {str(e)}")
            raise

    def _generate_html_report(self, scan_results: Dict[str, Any]) -> str:
        """Generate an HTML report from scan results."""
        try:
            template = self.env.get_template("report_template.html")
            report = template.render(
                results=scan_results,
                timestamp=datetime.now().isoformat(),
                title="Scanner Report"
            )
            return report
        except Exception as e:
            self.logger.error(f"Error generating HTML report: {str(e)}")
            raise

    def _generate_json_report(self, scan_results: Dict[str, Any]) -> str:
        """Generate a JSON report from scan results."""
        try:
            report = {
                "timestamp": datetime.now().isoformat(),
                "results": scan_results
            }
            return json.dumps(report, indent=2)
        except Exception as e:
            self.logger.error(f"Error generating JSON report: {str(e)}")
            raise

    def save_report(self, report: str, output_path: str, output_format: str = "html") -> None:
        """Save the report to a file."""
        try:
            output_file = Path(output_path)
            output_file.parent.mkdir(parents=True, exist_ok=True)
            
            if output_format == "html":
                output_file = output_file.with_suffix(".html")
            elif output_format == "json":
                output_file = output_file.with_suffix(".json")
            
            output_file.write_text(report)
            self.logger.info(f"Report saved to {output_file}")
        except Exception as e:
            self.logger.error(f"Error saving report: {str(e)}")
            raise 