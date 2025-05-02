"""Table scan report generator for Nessi."""

from typing import Dict, Any
from .report_generator import ReportGenerator

class TableScanReport(ReportGenerator):
    """Generates reports from table scan results."""
    
    def __init__(self, output_dir: str = "reports/table_scan"):
        """Initialize the table scan report generator.
        
        Args:
            output_dir (str): Directory where reports will be saved.
        """
        super().__init__(output_dir)
        
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a table scan report.
        
        Args:
            data (Dict[str, Any]): Table scan results.
            **kwargs: Additional arguments:
                - format (str): Output format (md, json, html)
                - include_stats (bool): Include detailed statistics
                
        Returns:
            str: Path to the generated report.
        """
        try:
            # Validate required fields
            required_fields = ['file_info', 'schema', 'row_count', 'column_stats']
            self.validate_data(data, required_fields)
            
            # Get format and other options
            format = kwargs.get('format', 'md')
            include_stats = kwargs.get('include_stats', True)
            
            # Generate report content
            if format == 'json':
                content = self._generate_json_report(data, include_stats)
            elif format == 'html':
                content = self._generate_html_report(data, include_stats)
            else:  # default to markdown
                content = self._generate_markdown_report(data, include_stats)
            
            # Save report
            timestamp = self.get_timestamp()
            filename = f"table_scan_{timestamp}"
            return self.save_report(content, filename, format)
            
        except Exception as e:
            logger.error(f"Failed to generate table scan report: {str(e)}")
            raise RuntimeError(f"Failed to generate table scan report: {str(e)}")
    
    def _generate_markdown_report(self, data: Dict[str, Any], include_stats: bool) -> str:
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
    
    def _generate_json_report(self, data: Dict[str, Any], include_stats: bool) -> Dict[str, Any]:
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
    
    def _generate_html_report(self, data: Dict[str, Any], include_stats: bool) -> str:
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