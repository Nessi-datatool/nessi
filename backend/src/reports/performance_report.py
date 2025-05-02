"""Performance report generator for Nessi."""

from typing import Dict, Any
from .report_generator import ReportGenerator

class PerformanceReport(ReportGenerator):
    """Generate performance reports."""
    
    def generate(self, data, **kwargs):
        """Generate a performance report."""
        report = {
            "performance": {
                "metrics": {
                    "execution_time": self._calculate_execution_time(data),
                    "memory_usage": self._calculate_memory_usage(data),
                    "throughput": self._calculate_throughput(data)
                }
            }
        }
        return report
    
    def _calculate_execution_time(self, data):
        """Calculate execution time metrics."""
        return {
            "total_time": data.get("total_time", 0),
            "average_time": data.get("average_time", 0),
            "max_time": data.get("max_time", 0)
        }
    
    def _calculate_memory_usage(self, data):
        """Calculate memory usage metrics."""
        return {
            "peak_memory": data.get("peak_memory", 0),
            "average_memory": data.get("average_memory", 0),
            "memory_growth": data.get("memory_growth", 0)
        }
    
    def _calculate_throughput(self, data):
        """Calculate throughput metrics."""
        return {
            "rows_per_second": data.get("rows_per_second", 0),
            "bytes_per_second": data.get("bytes_per_second", 0)
        }

def generate_report(data):
    """Generate a performance report."""
    report_generator = PerformanceReport()
    return report_generator.generate(data)

def generate_report_from_data(data: Dict[str, Any], **kwargs) -> str:
    """Generate a performance report from data.
    
    Args:
        data (Dict[str, Any]): Performance metrics.
        **kwargs: Additional arguments:
            - format (str): Output format (md, json, html)
            - include_details (bool): Include detailed metrics
            
    Returns:
        str: Path to the generated report.
    """
    try:
        # Validate required fields
        required_fields = ['processing_metrics', 'resource_utilization']
        validate_data(data, required_fields)
        
        # Get format and other options
        format = kwargs.get('format', 'md')
        include_details = kwargs.get('include_details', True)
        
        # Generate report content
        if format == 'json':
            content = _generate_json_report(data, include_details)
        elif format == 'html':
            content = _generate_html_report(data, include_details)
        else:  # default to markdown
            content = _generate_markdown_report(data, include_details)
        
        # Save report
        timestamp = get_timestamp()
        filename = f"performance_{timestamp}"
        return save_report(content, filename, format)
        
    except Exception as e:
        logger.error(f"Failed to generate performance report: {str(e)}")
        raise RuntimeError(f"Failed to generate performance report: {str(e)}")

def _generate_markdown_report(data: Dict[str, Any], include_details: bool) -> str:
    """Generate markdown report content.
    
    Args:
        data: Performance metrics.
        include_details: Whether to include detailed metrics.
        
    Returns:
        str: Markdown report content.
    """
    report = []
    
    # Title
    report.append("# Performance Report")
    report.append("")
    
    # Processing Metrics
    report.append("## Processing Metrics")
    processing = data['processing_metrics']
    report.append(f"Overall Performance Score: {processing['score']}%")
    report.append("")
    
    if include_details:
        report.append("### Execution Time")
        report.append("| Operation | Duration (ms) |")
        report.append("|-----------|---------------|")
        for op, duration in processing['execution_time'].items():
            report.append(f"| {op} | {duration} |")
        report.append("")
        
        report.append("### Memory Usage")
        report.append("| Component | Usage (MB) | Peak Usage (MB) |")
        report.append("|-----------|------------|----------------|")
        for comp, stats in processing['memory_usage'].items():
            report.append(f"| {comp} | {stats['current']} | {stats['peak']} |")
        report.append("")
        
        report.append("### CPU Utilization")
        report.append("| Component | Usage (%) |")
        report.append("|-----------|-----------|")
        for comp, usage in processing['cpu_usage'].items():
            report.append(f"| {comp} | {usage} |")
        report.append("")
    
    # Resource Utilization
    report.append("## Resource Utilization")
    resources = data['resource_utilization']
    report.append(f"Overall Resource Efficiency: {resources['efficiency_score']}%")
    report.append("")
    
    if include_details:
        report.append("### Disk I/O")
        report.append("| Operation | Read (MB/s) | Write (MB/s) |")
        report.append("|-----------|-------------|--------------|")
        for op, stats in resources['disk_io'].items():
            report.append(f"| {op} | {stats['read']} | {stats['write']} |")
        report.append("")
        
        report.append("### Network Traffic")
        report.append("| Component | Inbound (MB/s) | Outbound (MB/s) |")
        report.append("|-----------|----------------|-----------------|")
        for comp, stats in resources['network_traffic'].items():
            report.append(f"| {comp} | {stats['inbound']} | {stats['outbound']} |")
        report.append("")
        
        report.append("### Cache Performance")
        report.append("| Cache | Hit Rate (%) | Miss Rate (%) |")
        report.append("|-------|--------------|---------------|")
        for cache, stats in resources['cache_performance'].items():
            report.append(f"| {cache} | {stats['hit_rate']} | {stats['miss_rate']} |")
        report.append("")
    
    return '\n'.join(report)

def _generate_json_report(data: Dict[str, Any], include_details: bool) -> Dict[str, Any]:
    """Generate JSON report content.
    
    Args:
        data: Performance metrics.
        include_details: Whether to include detailed metrics.
        
    Returns:
        Dict[str, Any]: JSON report content.
    """
    report = {
        "processing_metrics": {
            "score": data['processing_metrics']['score']
        },
        "resource_utilization": {
            "efficiency_score": data['resource_utilization']['efficiency_score']
        }
    }
    
    if include_details:
        report["processing_metrics"].update({
            "execution_time": data['processing_metrics']['execution_time'],
            "memory_usage": data['processing_metrics']['memory_usage'],
            "cpu_usage": data['processing_metrics']['cpu_usage']
        })
        report["resource_utilization"].update({
            "disk_io": data['resource_utilization']['disk_io'],
            "network_traffic": data['resource_utilization']['network_traffic'],
            "cache_performance": data['resource_utilization']['cache_performance']
        })
        
    return report

def _generate_html_report(data: Dict[str, Any], include_details: bool) -> str:
    """Generate HTML report content.
    
    Args:
        data: Performance metrics.
        include_details: Whether to include detailed metrics.
        
    Returns:
        str: HTML report content.
    """
    html = []
    html.append("<!DOCTYPE html>")
    html.append("<html>")
    html.append("<head>")
    html.append("<title>Performance Report</title>")
    html.append("<style>")
    html.append("body { font-family: Arial, sans-serif; margin: 20px; }")
    html.append("table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }")
    html.append("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
    html.append("th { background-color: #f2f2f2; }")
    html.append(".score { font-size: 1.2em; font-weight: bold; }")
    html.append(".good { color: green; }")
    html.append(".warning { color: orange; }")
    html.append(".critical { color: red; }")
    html.append("</style>")
    html.append("</head>")
    html.append("<body>")
    
    # Title
    html.append("<h1>Performance Report</h1>")
    
    # Processing Metrics
    html.append("<h2>Processing Metrics</h2>")
    processing = data['processing_metrics']
    html.append(f"<p class='score'>Overall Performance Score: {processing['score']}%</p>")
    
    if include_details:
        html.append("<h3>Execution Time</h3>")
        html.append("<table>")
        html.append("<tr><th>Operation</th><th>Duration (ms)</th></tr>")
        for op, duration in processing['execution_time'].items():
            html.append(f"<tr><td>{op}</td><td>{duration}</td></tr>")
        html.append("</table>")
        
        html.append("<h3>Memory Usage</h3>")
        html.append("<table>")
        html.append("<tr><th>Component</th><th>Usage (MB)</th><th>Peak Usage (MB)</th></tr>")
        for comp, stats in processing['memory_usage'].items():
            html.append(f"<tr><td>{comp}</td><td>{stats['current']}</td><td>{stats['peak']}</td></tr>")
        html.append("</table>")
        
        html.append("<h3>CPU Utilization</h3>")
        html.append("<table>")
        html.append("<tr><th>Component</th><th>Usage (%)</th></tr>")
        for comp, usage in processing['cpu_usage'].items():
            html.append(f"<tr><td>{comp}</td><td>{usage}</td></tr>")
        html.append("</table>")
    
    # Resource Utilization
    html.append("<h2>Resource Utilization</h2>")
    resources = data['resource_utilization']
    html.append(f"<p class='score'>Overall Resource Efficiency: {resources['efficiency_score']}%</p>")
    
    if include_details:
        html.append("<h3>Disk I/O</h3>")
        html.append("<table>")
        html.append("<tr><th>Operation</th><th>Read (MB/s)</th><th>Write (MB/s)</th></tr>")
        for op, stats in resources['disk_io'].items():
            html.append(f"<tr><td>{op}</td><td>{stats['read']}</td><td>{stats['write']}</td></tr>")
        html.append("</table>")
        
        html.append("<h3>Network Traffic</h3>")
        html.append("<table>")
        html.append("<tr><th>Component</th><th>Inbound (MB/s)</th><th>Outbound (MB/s)</th></tr>")
        for comp, stats in resources['network_traffic'].items():
            html.append(f"<tr><td>{comp}</td><td>{stats['inbound']}</td><td>{stats['outbound']}</td></tr>")
        html.append("</table>")
        
        html.append("<h3>Cache Performance</h3>")
        html.append("<table>")
        html.append("<tr><th>Cache</th><th>Hit Rate (%)</th><th>Miss Rate (%)</th></tr>")
        for cache, stats in resources['cache_performance'].items():
            html.append(f"<tr><td>{cache}</td><td>{stats['hit_rate']}</td><td>{stats['miss_rate']}</td></tr>")
        html.append("</table>")
    
    html.append("</body>")
    html.append("</html>")
    
    return '\n'.join(html) 