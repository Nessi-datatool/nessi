"""Data quality report generator for Nessi."""

from typing import Dict, Any, List
from .report_generator import ReportGenerator

class DataQualityReport(ReportGenerator):
    """Generates data quality reports."""
    
    def __init__(self, output_dir: str = "reports/data_quality"):
        """Initialize the data quality report generator.
        
        Args:
            output_dir (str): Directory where reports will be saved.
        """
        super().__init__(output_dir)
        
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a data quality report.
        
        Args:
            data (Dict[str, Any]): Data quality analysis results.
            **kwargs: Additional arguments:
                - format (str): Output format (md, json, html)
                - include_details (bool): Include detailed analysis
                
        Returns:
            str: Path to the generated report.
        """
        try:
            # Validate required fields
            required_fields = ['completeness', 'consistency', 'uniqueness']
            self.validate_data(data, required_fields)
            
            # Get format and other options
            format = kwargs.get('format', 'md')
            include_details = kwargs.get('include_details', True)
            
            # Generate report content
            if format == 'json':
                content = self._generate_json_report(data, include_details)
            elif format == 'html':
                content = self._generate_html_report(data, include_details)
            else:  # default to markdown
                content = self._generate_markdown_report(data, include_details)
            
            # Save report
            timestamp = self.get_timestamp()
            filename = f"data_quality_{timestamp}"
            return self.save_report(content, filename, format)
            
        except Exception as e:
            logger.error(f"Failed to generate data quality report: {str(e)}")
            raise RuntimeError(f"Failed to generate data quality report: {str(e)}")
    
    def _generate_markdown_report(self, data: Dict[str, Any], include_details: bool) -> str:
        """Generate markdown report content.
        
        Args:
            data: Data quality analysis results.
            include_details: Whether to include detailed analysis.
            
        Returns:
            str: Markdown report content.
        """
        report = []
        
        # Title
        report.append("# Data Quality Report")
        report.append("")
        
        # Completeness Analysis
        report.append("## Completeness Analysis")
        completeness = data['completeness']
        report.append(f"Overall Completeness Score: {completeness['score']}%")
        report.append("")
        
        if include_details:
            report.append("### Missing Values by Column")
            report.append("| Column | Missing Count | Missing Percentage |")
            report.append("|--------|---------------|-------------------|")
            for col, stats in completeness['column_stats'].items():
                report.append(f"| {col} | {stats['missing_count']} | {stats['missing_percentage']}% |")
            report.append("")
        
        # Consistency Analysis
        report.append("## Consistency Analysis")
        consistency = data['consistency']
        report.append(f"Overall Consistency Score: {consistency['score']}%")
        report.append("")
        
        if include_details:
            report.append("### Data Type Validation")
            report.append("| Column | Expected Type | Actual Type | Status |")
            report.append("|--------|---------------|-------------|--------|")
            for col, stats in consistency['type_validation'].items():
                report.append(f"| {col} | {stats['expected']} | {stats['actual']} | {'✓' if stats['valid'] else '✗'} |")
            report.append("")
            
            report.append("### Value Range Validation")
            report.append("| Column | Min Value | Max Value | Out of Range |")
            report.append("|--------|-----------|-----------|--------------|")
            for col, stats in consistency['range_validation'].items():
                report.append(f"| {col} | {stats['min']} | {stats['max']} | {stats['out_of_range']} |")
            report.append("")
        
        # Uniqueness Analysis
        report.append("## Uniqueness Analysis")
        uniqueness = data['uniqueness']
        report.append(f"Overall Uniqueness Score: {uniqueness['score']}%")
        report.append("")
        
        if include_details:
            report.append("### Duplicate Analysis")
            report.append("| Column | Total Duplicates | Duplicate Percentage |")
            report.append("|--------|-----------------|---------------------|")
            for col, stats in uniqueness['duplicate_stats'].items():
                report.append(f"| {col} | {stats['duplicate_count']} | {stats['duplicate_percentage']}% |")
            report.append("")
            
            report.append("### Primary Key Validation")
            report.append("| Column | Status |")
            report.append("|--------|--------|")
            for col, status in uniqueness['primary_key_validation'].items():
                report.append(f"| {col} | {'✓' if status else '✗'} |")
            report.append("")
        
        return '\n'.join(report)
    
    def _generate_json_report(self, data: Dict[str, Any], include_details: bool) -> Dict[str, Any]:
        """Generate JSON report content.
        
        Args:
            data: Data quality analysis results.
            include_details: Whether to include detailed analysis.
            
        Returns:
            Dict[str, Any]: JSON report content.
        """
        report = {
            "completeness": {
                "score": data['completeness']['score']
            },
            "consistency": {
                "score": data['consistency']['score']
            },
            "uniqueness": {
                "score": data['uniqueness']['score']
            }
        }
        
        if include_details:
            report["completeness"]["column_stats"] = data['completeness']['column_stats']
            report["consistency"]["type_validation"] = data['consistency']['type_validation']
            report["consistency"]["range_validation"] = data['consistency']['range_validation']
            report["uniqueness"]["duplicate_stats"] = data['uniqueness']['duplicate_stats']
            report["uniqueness"]["primary_key_validation"] = data['uniqueness']['primary_key_validation']
            
        return report
    
    def _generate_html_report(self, data: Dict[str, Any], include_details: bool) -> str:
        """Generate HTML report content.
        
        Args:
            data: Data quality analysis results.
            include_details: Whether to include detailed analysis.
            
        Returns:
            str: HTML report content.
        """
        html = []
        html.append("<!DOCTYPE html>")
        html.append("<html>")
        html.append("<head>")
        html.append("<title>Data Quality Report</title>")
        html.append("<style>")
        html.append("body { font-family: Arial, sans-serif; margin: 20px; }")
        html.append("table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }")
        html.append("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
        html.append("th { background-color: #f2f2f2; }")
        html.append(".score { font-size: 1.2em; font-weight: bold; }")
        html.append(".valid { color: green; }")
        html.append(".invalid { color: red; }")
        html.append("</style>")
        html.append("</head>")
        html.append("<body>")
        
        # Title
        html.append("<h1>Data Quality Report</h1>")
        
        # Completeness Analysis
        html.append("<h2>Completeness Analysis</h2>")
        completeness = data['completeness']
        html.append(f"<p class='score'>Overall Completeness Score: {completeness['score']}%</p>")
        
        if include_details:
            html.append("<h3>Missing Values by Column</h3>")
            html.append("<table>")
            html.append("<tr><th>Column</th><th>Missing Count</th><th>Missing Percentage</th></tr>")
            for col, stats in completeness['column_stats'].items():
                html.append(f"<tr><td>{col}</td><td>{stats['missing_count']}</td><td>{stats['missing_percentage']}%</td></tr>")
            html.append("</table>")
        
        # Consistency Analysis
        html.append("<h2>Consistency Analysis</h2>")
        consistency = data['consistency']
        html.append(f"<p class='score'>Overall Consistency Score: {consistency['score']}%</p>")
        
        if include_details:
            html.append("<h3>Data Type Validation</h3>")
            html.append("<table>")
            html.append("<tr><th>Column</th><th>Expected Type</th><th>Actual Type</th><th>Status</th></tr>")
            for col, stats in consistency['type_validation'].items():
                status_class = "valid" if stats['valid'] else "invalid"
                html.append(f"<tr><td>{col}</td><td>{stats['expected']}</td><td>{stats['actual']}</td><td class='{status_class}'>{'✓' if stats['valid'] else '✗'}</td></tr>")
            html.append("</table>")
            
            html.append("<h3>Value Range Validation</h3>")
            html.append("<table>")
            html.append("<tr><th>Column</th><th>Min Value</th><th>Max Value</th><th>Out of Range</th></tr>")
            for col, stats in consistency['range_validation'].items():
                html.append(f"<tr><td>{col}</td><td>{stats['min']}</td><td>{stats['max']}</td><td>{stats['out_of_range']}</td></tr>")
            html.append("</table>")
        
        # Uniqueness Analysis
        html.append("<h2>Uniqueness Analysis</h2>")
        uniqueness = data['uniqueness']
        html.append(f"<p class='score'>Overall Uniqueness Score: {uniqueness['score']}%</p>")
        
        if include_details:
            html.append("<h3>Duplicate Analysis</h3>")
            html.append("<table>")
            html.append("<tr><th>Column</th><th>Total Duplicates</th><th>Duplicate Percentage</th></tr>")
            for col, stats in uniqueness['duplicate_stats'].items():
                html.append(f"<tr><td>{col}</td><td>{stats['duplicate_count']}</td><td>{stats['duplicate_percentage']}%</td></tr>")
            html.append("</table>")
            
            html.append("<h3>Primary Key Validation</h3>")
            html.append("<table>")
            html.append("<tr><th>Column</th><th>Status</th></tr>")
            for col, status in uniqueness['primary_key_validation'].items():
                status_class = "valid" if status else "invalid"
                html.append(f"<tr><td>{col}</td><td class='{status_class}'>{'✓' if status else '✗'}</td></tr>")
            html.append("</table>")
        
        html.append("</body>")
        html.append("</html>")
        
        return '\n'.join(html) 