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

"""Table scan report generator for nessi.dev."""

from typing import Dict, Any
from .report_generator import ReportGenerator

class TableScanReport(ReportGenerator):
    """Generate table scan reports."""
    
    def generate(self, data, **kwargs):
        """Generate a table scan report."""
        report = {
            "table_scan": {
                "metrics": {
                    "schema": self._analyze_schema(data),
                    "statistics": self._calculate_statistics(data),
                    "quality": self._assess_quality(data)
                }
            }
        }
        return report
    
    def _analyze_schema(self, data):
        """Analyze table schema."""
        return {
            "columns": data.get("columns", []),
            "data_types": data.get("data_types", {}),
            "nullable": data.get("nullable", {})
        }
    
    def _calculate_statistics(self, data):
        """Calculate table statistics."""
        return {
            "row_count": data.get("row_count", 0),
            "column_stats": data.get("column_stats", {}),
            "size_bytes": data.get("size_bytes", 0)
        }
    
    def _assess_quality(self, data):
        """Assess table quality."""
        return {
            "missing_values": data.get("missing_values", {}),
            "unique_values": data.get("unique_values", {}),
            "data_types": data.get("data_types", {})
        }

def generate_report(data):
    """Generate a table scan report."""
    report_generator = TableScanReport()
    return report_generator.generate(data)

def generate_markdown_report(data: Dict[str, Any], include_stats: bool) -> str:
    """Generate markdown report content.
    
    Args:
        data: Table scan results.
        include_stats: Whether to include detailed statistics.
        
    Returns:
        str: Markdown report content.
    """
    report = []
    
    # Title
    report.append("# Table Scan Report")
    report.append("")
    
    # File Information
    report.append("## File Information")
    file_info = data['file_info']
    report.append(f"- Format: {file_info['format']}")
    report.append(f"- Path: {file_info['path']}")
    report.append(f"- Size: {file_info['size']} bytes")
    report.append(f"- Last Modified: {file_info['last_modified']}")
    report.append("")
    
    # Schema
    report.append("## Schema")
    report.append("| Column | Type | Nullable |")
    report.append("|--------|------|----------|")
    for field in data['schema']:
        report.append(f"| {field['name']} | {field['type']} | {field['nullable']} |")
    report.append("")
    
    # Row Count
    report.append("## Row Count")
    report.append(f"Total Rows: {data['row_count']}")
    report.append("")
    
    # Column Statistics
    if include_stats:
        report.append("## Column Statistics")
        for col, stats in data['column_stats'].items():
            report.append(f"### {col}")
            report.append(f"- Count: {stats['count']}")
            report.append(f"- Null Count: {stats['null_count']}")
            if stats.get('min') is not None:
                report.append(f"- Min: {stats['min']}")
            if stats.get('max') is not None:
                report.append(f"- Max: {stats['max']}")
            if stats.get('mean') is not None:
                report.append(f"- Mean: {stats['mean']}")
            if stats.get('stddev') is not None:
                report.append(f"- Standard Deviation: {stats['stddev']}")
            report.append("")
    
    return '\n'.join(report)

def generate_json_report(data: Dict[str, Any], include_stats: bool) -> Dict[str, Any]:
    """Generate JSON report content.
    
    Args:
        data: Table scan results.
        include_stats: Whether to include detailed statistics.
        
    Returns:
        Dict[str, Any]: JSON report content.
    """
    report = {
        "file_info": data['file_info'],
        "schema": data['schema'],
        "row_count": data['row_count']
    }
    
    if include_stats:
        report["column_stats"] = data['column_stats']
        
    return report

def generate_html_report(data: Dict[str, Any], include_stats: bool) -> str:
    """Generate HTML report content.
    
    Args:
        data: Table scan results.
        include_stats: Whether to include detailed statistics.
        
    Returns:
        str: HTML report content.
    """
    html = []
    html.append("<!DOCTYPE html>")
    html.append("<html>")
    html.append("<head>")
    html.append("<title>Table Scan Report</title>")
    html.append("<style>")
    html.append("body { font-family: Arial, sans-serif; margin: 20px; }")
    html.append("table { border-collapse: collapse; width: 100%; }")
    html.append("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
    html.append("th { background-color: #f2f2f2; }")
    html.append("</style>")
    html.append("</head>")
    html.append("<body>")
    
    # Title
    html.append("<h1>Table Scan Report</h1>")
    
    # File Information
    html.append("<h2>File Information</h2>")
    file_info = data['file_info']
    html.append("<ul>")
    html.append(f"<li>Format: {file_info['format']}</li>")
    html.append(f"<li>Path: {file_info['path']}</li>")
    html.append(f"<li>Size: {file_info['size']} bytes</li>")
    html.append(f"<li>Last Modified: {file_info['last_modified']}</li>")
    html.append("</ul>")
    
    # Schema
    html.append("<h2>Schema</h2>")
    html.append("<table>")
    html.append("<tr><th>Column</th><th>Type</th><th>Nullable</th></tr>")
    for field in data['schema']:
        html.append(f"<tr><td>{field['name']}</td><td>{field['type']}</td><td>{field['nullable']}</td></tr>")
    html.append("</table>")
    
    # Row Count
    html.append("<h2>Row Count</h2>")
    html.append(f"<p>Total Rows: {data['row_count']}</p>")
    
    # Column Statistics
    if include_stats:
        html.append("<h2>Column Statistics</h2>")
        for col, stats in data['column_stats'].items():
            html.append(f"<h3>{col}</h3>")
            html.append("<ul>")
            html.append(f"<li>Count: {stats['count']}</li>")
            html.append(f"<li>Null Count: {stats['null_count']}</li>")
            if stats.get('min') is not None:
                html.append(f"<li>Min: {stats['min']}</li>")
            if stats.get('max') is not None:
                html.append(f"<li>Max: {stats['max']}</li>")
            if stats.get('mean') is not None:
                html.append(f"<li>Mean: {stats['mean']}</li>")
            if stats.get('stddev') is not None:
                html.append(f"<li>Standard Deviation: {stats['stddev']}</li>")
            html.append("</ul>")
    
    html.append("</body>")
    html.append("</html>")
    
    return '\n'.join(html)