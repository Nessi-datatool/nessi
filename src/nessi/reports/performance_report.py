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

"""Performance report generator for nessi.dev."""

import logging
import json
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List, Optional

from .report_generator import ReportGenerator

# Configure logger
logger = logging.getLogger(__name__)

def validate_data(data: Dict[str, Any], required_fields: List[str]) -> None:
    """Validate that data contains all required fields.
    
    Args:
        data: Data to validate.
        required_fields: List of required field names.
        
    Raises:
        RuntimeError: If any required field is missing.
    """
    missing_fields = [field for field in required_fields if field not in data]
    if missing_fields:
        raise RuntimeError(f"Missing required fields: {', '.join(missing_fields)}")

def get_timestamp() -> str:
    """Get current timestamp in YYYYMMDD_HHMMSS format."""
    return datetime.now().strftime("%Y%m%d_%H%M%S")

def save_report(content: Any, filename: str, format: str) -> str:
    """Save report content to a file.
    
    Args:
        content: Report content to save.
        filename: Base filename without extension.
        format: File format (md, json, html).
        
    Returns:
        str: Path to the saved file.
    """
    # Create reports directory if it doesn't exist
    reports_dir = Path("reports")
    reports_dir.mkdir(exist_ok=True)
    
    # Determine file extension and content type
    if format == "json":
        extension = "json"
        if isinstance(content, dict):
            content = json.dumps(content, indent=2)
    elif format == "html":
        extension = "html"
    else:  # markdown
        extension = "md"
    
    # Save file
    filepath = reports_dir / f"{filename}.{extension}"
    with open(filepath, "w") as f:
        f.write(content)
    
    return str(filepath)

class PerformanceReport(ReportGenerator):
    """Generate performance reports."""
    
    def __init__(self, template_dir: str = "templates", output_dir: str = "reports"):
        """Initialize the performance report generator.
        
        Args:
            template_dir (str): Directory containing report templates
            output_dir (str): Directory where reports will be saved
        """
        super().__init__(template_dir, output_dir)
        
    def generate(self, data: Dict[str, Any], **kwargs) -> str:
        """Generate a performance report.
        
        Args:
            data: Performance data to include in the report
            **kwargs: Additional options for report generation
            
        Returns:
            str: Path to the generated report
        """
        # Calculate performance metrics
        processed_data = self._process_data(data)
        
        # Generate report in specified format
        format = kwargs.get('format', 'html')
        output_path = kwargs.get('output_path')
        
        return self.generate_report(processed_data, format=format, output_path=output_path)
        
    def _process_data(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Process raw performance data into report format.
        
        Args:
            data: Raw performance data
            
        Returns:
            Dict[str, Any]: Processed data ready for report generation
        """
        processed = {
            "title": "Performance Report",
            "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            "metrics": self._calculate_metrics(data),
            "performance": {
                "metrics": {
                    "execution_time": self._calculate_execution_time(data),
                    "memory_usage": self._calculate_memory_usage(data),
                    "throughput": self._calculate_throughput(data)
                }
            }
        }
        
        return processed

    def _calculate_execution_time(self, data: Dict[str, Any]) -> Dict[str, float]:
        """Calculate execution time metrics.
        
        Args:
            data: Raw performance data
            
        Returns:
            Dict[str, float]: Execution time metrics
        """
        execution_times = data["processing_metrics"]["execution_time"]
        total_time = sum(execution_times.values())
        return {
            "total_time": total_time,
            "average_time": total_time / len(execution_times),
            "max_time": max(execution_times.values())
        }

    def _calculate_memory_usage(self, data: Dict[str, Any]) -> Dict[str, float]:
        """Calculate memory usage metrics.
        
        Args:
            data: Raw performance data
            
        Returns:
            Dict[str, float]: Memory usage metrics
        """
        memory_usage = data["processing_metrics"]["memory_usage"]
        return {
            "peak_memory": max(memory_usage["spark"]["peak"], memory_usage["system"]["peak"]),
            "average_memory": (memory_usage["spark"]["current"] + memory_usage["system"]["current"]) / 2,
            "memory_growth": (memory_usage["spark"]["peak"] - memory_usage["spark"]["current"]) / memory_usage["spark"]["current"] * 100
        }

    def _calculate_throughput(self, data: Dict[str, Any]) -> float:
        """Calculate throughput metrics.
        
        Args:
            data: Raw performance data
            
        Returns:
            float: Throughput in requests per second
        """
        total_time = sum(data["processing_metrics"]["execution_time"].values())
        if total_time == 0:
            return 0
        return 1000 / total_time  # Convert to requests per second

    def _calculate_metrics(self, data: Dict[str, Any]) -> Dict[str, Dict[str, Any]]:
        """Calculate detailed performance metrics.
        
        Args:
            data: Raw performance data
            
        Returns:
            Dict[str, Dict[str, Any]]: Detailed metrics by category
        """
        metrics = {
            "execution": {
                "min_time": self._calculate_min_time(data),
                "max_time": self._calculate_max_time(data),
                "median_time": self._calculate_median_time(data),
                "std_dev_time": self._calculate_std_dev_time(data)
            },
            "memory": {
                "min_usage": self._calculate_min_memory(data),
                "max_usage": self._calculate_max_memory(data),
                "average_usage": self._calculate_average_memory(data)
            },
            "throughput": {
                "requests_per_second": self._calculate_throughput(data),
                "concurrent_requests": self._calculate_concurrent_requests(data)
            }
        }
        return metrics
        
    def _calculate_min_time(self, data: Dict[str, Any]) -> float:
        """Calculate minimum execution time."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        return min(e.get("time", 0) for e in executions)
        
    def _calculate_max_time(self, data: Dict[str, Any]) -> float:
        """Calculate maximum execution time."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        return max(e.get("time", 0) for e in executions)
        
    def _calculate_median_time(self, data: Dict[str, Any]) -> float:
        """Calculate median execution time."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        times = sorted(e.get("time", 0) for e in executions)
        n = len(times)
        if n % 2 == 0:
            return (times[n//2 - 1] + times[n//2]) / 2
        return times[n//2]
        
    def _calculate_std_dev_time(self, data: Dict[str, Any]) -> float:
        """Calculate standard deviation of execution time."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        times = [e.get("time", 0) for e in executions]
        mean = sum(times) / len(times)
        variance = sum((x - mean) ** 2 for x in times) / len(times)
        return variance ** 0.5
        
    def _calculate_min_memory(self, data: Dict[str, Any]) -> float:
        """Calculate minimum memory usage."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        return min(e.get("memory", 0) for e in executions)
        
    def _calculate_max_memory(self, data: Dict[str, Any]) -> float:
        """Calculate maximum memory usage."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        return max(e.get("memory", 0) for e in executions)
        
    def _calculate_average_memory(self, data: Dict[str, Any]) -> float:
        """Calculate average memory usage."""
        executions = data.get("executions", [])
        if not executions:
            return 0.0
        return sum(e.get("memory", 0) for e in executions) / len(executions)
        
    def _calculate_concurrent_requests(self, data: Dict[str, Any]) -> int:
        """Calculate maximum concurrent requests."""
        executions = data.get("executions", [])
        if not executions:
            return 0
        # This is a simplified calculation - in reality, you'd need timestamps
        return max(len([e for e in executions if e.get("concurrent", False)]), 1)

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
        "performance": {
            "processing_metrics": {
                "score": data['processing_metrics']['score']
            },
            "resource_utilization": {
                "efficiency_score": data['resource_utilization']['efficiency_score']
            }
        },
        "metrics": {
            "processing_metrics": {
                "score": data['processing_metrics']['score']
            },
            "resource_utilization": {
                "efficiency_score": data['resource_utilization']['efficiency_score']
            }
        }
    }
    
    if include_details:
        report["performance"]["processing_metrics"].update({
            "execution_time": data['processing_metrics']['execution_time'],
            "memory_usage": data['processing_metrics']['memory_usage'],
            "cpu_usage": data['processing_metrics']['cpu_usage']
        })
        report["performance"]["resource_utilization"].update({
            "disk_io": data['resource_utilization']['disk_io'],
            "network_traffic": data['resource_utilization']['network_traffic'],
            "cache_performance": data['resource_utilization']['cache_performance']
        })
        report["metrics"]["processing_metrics"].update({
            "execution_time": data['processing_metrics']['execution_time'],
            "memory_usage": data['processing_metrics']['memory_usage'],
            "cpu_usage": data['processing_metrics']['cpu_usage']
        })
        report["metrics"]["resource_utilization"].update({
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