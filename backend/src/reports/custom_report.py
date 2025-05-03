"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

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
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Custom report generator for Nessi."""

import os
from typing import Dict, Any, Optional
from jinja2 import Environment, FileSystemLoader, select_autoescape
from .report_generator import ReportGenerator

class CustomReport(ReportGenerator):
    """Generate custom reports."""
    
    def generate(self, data, **kwargs):
        """Generate a custom report."""
        report = {
            "custom": {
                "metrics": {
                    "user_defined": self._process_user_metrics(data),
                    "calculated": self._calculate_metrics(data),
                    "summary": self._generate_summary(data)
                }
            }
        }
        return report
    
    def _process_user_metrics(self, data):
        """Process user-defined metrics."""
        return {
            "metrics": data.get("metrics", {}),
            "thresholds": data.get("thresholds", {}),
            "alerts": data.get("alerts", [])
        }
    
    def _calculate_metrics(self, data):
        """Calculate custom metrics."""
        return {
            "derived": data.get("derived_metrics", {}),
            "aggregated": data.get("aggregated_metrics", {}),
            "trends": data.get("trend_metrics", {})
        }
    
    def _generate_summary(self, data):
        """Generate summary metrics."""
        return {
            "total_metrics": len(data.get("metrics", {})),
            "alerts_triggered": len(data.get("alerts", [])),
            "status": data.get("status", "unknown")
        }

def generate_report(data):
    """Generate a custom report."""
    report_generator = CustomReport()
    return report_generator.generate(data)

def generate_report_from_template(data, template_name, format, template_vars):
    """Generate a custom report from a template."""
    report_generator = CustomReport()
    return report_generator.generate(data, template=template_name, format=format, template_vars=template_vars)

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