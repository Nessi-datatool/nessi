"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

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
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Script to generate example reports and metrics for demonstration."""

import time
import random
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, Any, List
from prometheus_client import start_http_server, Gauge, Counter, Histogram
import json

# Initialize Prometheus metrics
table_scan_rows = Gauge('nessi_table_scan_total_rows', 'Total rows scanned in table')
data_quality_score = Gauge('nessi_data_quality_score', 'Data quality score (0-100)')
processing_time = Histogram('nessi_processing_time_seconds', 'Processing time in seconds')
custom_report_count = Counter('nessi_custom_report_count', 'Number of custom reports generated')

class ExampleGenerator:
    """Generates example reports and metrics."""
    
    def __init__(self):
        self.base_dir = Path(__file__).parent
        self.reports_dir = self.base_dir.parent / "reports"
        self.reports_dir.mkdir(exist_ok=True)
        self.scenarios = [
            "normal_operation",
            "high_load",
            "data_quality_issues",
            "performance_degradation",
            "custom_analysis"
        ]
        
    def _get_current_scenario(self) -> str:
        """Get current scenario based on time."""
        current_minute = datetime.now().minute
        return self.scenarios[current_minute % len(self.scenarios)]
        
    def generate_table_scan_report(self) -> Dict[str, Any]:
        """Generate example table scan report with different scenarios."""
        scenario = self._get_current_scenario()
        
        if scenario == "normal_operation":
            rows = random.randint(1000000, 2000000)
            null_counts = [0, 2, 1, 0]
            analysis = {
                "scan_efficiency": "Optimal",
                "partition_utilization": "Balanced",
                "index_usage": "Efficient",
                "recommendations": [
                    "Current scan patterns are optimal",
                    "Consider adding indexes for frequently queried columns",
                    "Monitor growth patterns for capacity planning"
                ]
            }
        elif scenario == "high_load":
            rows = random.randint(4000000, 5000000)
            null_counts = [0, 10, 5, 0]
            analysis = {
                "scan_efficiency": "High",
                "partition_utilization": "High",
                "index_usage": "Moderate",
                "recommendations": [
                    "Consider partitioning large tables",
                    "Review index strategy for high-volume queries",
                    "Implement query optimization for frequent patterns"
                ]
            }
        elif scenario == "data_quality_issues":
            rows = random.randint(1000000, 2000000)
            null_counts = [0, 50, 20, 0]
            analysis = {
                "scan_efficiency": "Moderate",
                "partition_utilization": "Unbalanced",
                "index_usage": "Inefficient",
                "recommendations": [
                    "Implement data validation rules",
                    "Add constraints to prevent null values",
                    "Review data ingestion process"
                ]
            }
        elif scenario == "performance_degradation":
            rows = random.randint(3000000, 4000000)
            null_counts = [0, 5, 3, 0]
            analysis = {
                "scan_efficiency": "Low",
                "partition_utilization": "Overloaded",
                "index_usage": "Poor",
                "recommendations": [
                    "Optimize table partitioning",
                    "Rebuild fragmented indexes",
                    "Consider table maintenance operations"
                ]
            }
        else:  # custom_analysis
            rows = random.randint(2000000, 3000000)
            null_counts = [0, 3, 2, 0]
            analysis = {
                "scan_efficiency": "Very High",
                "partition_utilization": "Optimal",
                "index_usage": "Excellent",
                "recommendations": [
                    "Current configuration is optimal",
                    "Consider implementing advanced indexing strategies",
                    "Explore materialized views for complex queries"
                ]
            }
            
        return {
            "timestamp": datetime.now().isoformat(),
            "scenario": scenario,
            "table_name": "example_table",
            "total_rows": rows,
            "columns": [
                {"name": "id", "type": "integer", "null_count": null_counts[0]},
                {"name": "name", "type": "string", "null_count": null_counts[1]},
                {"name": "value", "type": "float", "null_count": null_counts[2]},
                {"name": "timestamp", "type": "timestamp", "null_count": null_counts[3]}
            ],
            "statistics": {
                "min_value": 0.0,
                "max_value": 100.0,
                "avg_value": 50.0,
                "std_dev": 25.0
            },
            "detailed_analysis": analysis
        }
        
    def generate_data_quality_report(self) -> Dict[str, Any]:
        """Generate example data quality report with different scenarios."""
        scenario = self._get_current_scenario()
        
        if scenario == "normal_operation":
            overall_score = random.randint(90, 100)
            missing_values = random.randint(0, 10)
            inconsistent_values = random.randint(0, 5)
            duplicate_values = random.randint(0, 3)
            analysis = {
                "data_health": "Excellent",
                "trends": {
                    "completeness": "Stable",
                    "consistency": "Improving",
                    "uniqueness": "Maintained"
                },
                "anomalies": {
                    "detected": False,
                    "description": "No significant anomalies detected",
                    "impact": "None"
                },
                "recommendations": [
                    "Continue current data quality practices",
                    "Implement automated validation checks",
                    "Schedule regular data quality audits"
                ]
            }
        elif scenario == "high_load":
            overall_score = random.randint(85, 95)
            missing_values = random.randint(10, 30)
            inconsistent_values = random.randint(5, 15)
            duplicate_values = random.randint(3, 8)
            analysis = {
                "data_health": "Good",
                "trends": {
                    "completeness": "Slightly Declining",
                    "consistency": "Fluctuating",
                    "uniqueness": "Stable"
                },
                "anomalies": {
                    "detected": True,
                    "description": "Increased missing values during peak hours",
                    "impact": "Moderate - Affects reporting accuracy"
                },
                "recommendations": [
                    "Implement load-based validation rules",
                    "Add data quality checks in ingestion pipeline",
                    "Monitor data patterns during peak hours"
                ]
            }
        elif scenario == "data_quality_issues":
            overall_score = random.randint(70, 80)
            missing_values = random.randint(50, 100)
            inconsistent_values = random.randint(20, 50)
            duplicate_values = random.randint(10, 30)
            analysis = {
                "data_health": "Poor",
                "trends": {
                    "completeness": "Declining",
                    "consistency": "Degrading",
                    "uniqueness": "Problematic"
                },
                "anomalies": {
                    "detected": True,
                    "description": "Multiple data quality issues detected",
                    "impact": "High - Affects business decisions"
                },
                "recommendations": [
                    "Implement strict data validation rules",
                    "Review and fix data ingestion process",
                    "Schedule immediate data cleanup"
                ]
            }
        elif scenario == "performance_degradation":
            overall_score = random.randint(80, 90)
            missing_values = random.randint(20, 40)
            inconsistent_values = random.randint(10, 20)
            duplicate_values = random.randint(5, 10)
            analysis = {
                "data_health": "Fair",
                "trends": {
                    "completeness": "Declining",
                    "consistency": "Stable",
                    "uniqueness": "Improving"
                },
                "anomalies": {
                    "detected": True,
                    "description": "Performance issues affecting data quality",
                    "impact": "Moderate - Affects system reliability"
                },
                "recommendations": [
                    "Optimize data processing pipeline",
                    "Implement performance monitoring",
                    "Review resource allocation"
                ]
            }
        else:  # custom_analysis
            overall_score = random.randint(95, 100)
            missing_values = random.randint(0, 5)
            inconsistent_values = random.randint(0, 3)
            duplicate_values = random.randint(0, 2)
            analysis = {
                "data_health": "Excellent",
                "trends": {
                    "completeness": "Improving",
                    "consistency": "Excellent",
                    "uniqueness": "Optimal"
                },
                "anomalies": {
                    "detected": False,
                    "description": "No anomalies detected",
                    "impact": "None"
                },
                "recommendations": [
                    "Implement advanced data quality monitoring",
                    "Explore predictive analytics for data quality",
                    "Consider automated data quality improvements"
                ]
            }
            
        return {
            "timestamp": datetime.now().isoformat(),
            "scenario": scenario,
            "overall_score": overall_score,
            "completeness": {
                "score": random.randint(overall_score - 5, overall_score),
                "missing_values": missing_values,
                "total_fields": 1000
            },
            "consistency": {
                "score": random.randint(overall_score - 5, overall_score),
                "inconsistent_values": inconsistent_values,
                "total_checks": 500
            },
            "uniqueness": {
                "score": random.randint(overall_score - 5, overall_score),
                "duplicate_values": duplicate_values,
                "total_records": 10000
            },
            "detailed_analysis": analysis
        }
        
    def generate_performance_report(self) -> Dict[str, Any]:
        """Generate example performance report with different scenarios."""
        scenario = self._get_current_scenario()
        
        if scenario == "normal_operation":
            total_time = random.uniform(1.0, 2.0)
            memory_usage = random.uniform(0.5, 1.0)
            cpu_usage = random.uniform(10.0, 30.0)
            analysis = {
                "performance_status": "Optimal",
                "bottlenecks": [],
                "resource_utilization": {
                    "memory": "Efficient",
                    "cpu": "Balanced",
                    "disk": "Normal",
                    "network": "Stable"
                },
                "recommendations": [
                    "Current performance is optimal",
                    "Consider implementing caching for frequent queries",
                    "Monitor for any performance degradation"
                ]
            }
        elif scenario == "high_load":
            total_time = random.uniform(3.0, 5.0)
            memory_usage = random.uniform(1.5, 2.0)
            cpu_usage = random.uniform(40.0, 70.0)
            analysis = {
                "performance_status": "High Load",
                "bottlenecks": [
                    "CPU utilization approaching limits",
                    "Memory usage increased",
                    "Query response time affected"
                ],
                "resource_utilization": {
                    "memory": "High",
                    "cpu": "High",
                    "disk": "Moderate",
                    "network": "Stable"
                },
                "recommendations": [
                    "Consider scaling resources",
                    "Implement query optimization",
                    "Review and optimize heavy queries"
                ]
            }
        elif scenario == "data_quality_issues":
            total_time = random.uniform(2.0, 3.0)
            memory_usage = random.uniform(1.0, 1.5)
            cpu_usage = random.uniform(30.0, 50.0)
            analysis = {
                "performance_status": "Affected",
                "bottlenecks": [
                    "Data validation overhead",
                    "Increased processing time",
                    "Resource contention"
                ],
                "resource_utilization": {
                    "memory": "Moderate",
                    "cpu": "Moderate",
                    "disk": "High",
                    "network": "Normal"
                },
                "recommendations": [
                    "Optimize data validation process",
                    "Implement parallel processing",
                    "Review data quality checks"
                ]
            }
        elif scenario == "performance_degradation":
            total_time = random.uniform(4.0, 5.0)
            memory_usage = random.uniform(2.0, 3.0)
            cpu_usage = random.uniform(60.0, 80.0)
            analysis = {
                "performance_status": "Degraded",
                "bottlenecks": [
                    "Severe CPU contention",
                    "Memory pressure",
                    "I/O bottlenecks",
                    "Query performance issues"
                ],
                "resource_utilization": {
                    "memory": "Critical",
                    "cpu": "Critical",
                    "disk": "High",
                    "network": "Moderate"
                },
                "recommendations": [
                    "Immediate resource scaling required",
                    "Review and optimize all queries",
                    "Implement performance monitoring",
                    "Consider database maintenance"
                ]
            }
        else:  # custom_analysis
            total_time = random.uniform(0.5, 1.5)
            memory_usage = random.uniform(0.3, 0.8)
            cpu_usage = random.uniform(5.0, 20.0)
            analysis = {
                "performance_status": "Excellent",
                "bottlenecks": [],
                "resource_utilization": {
                    "memory": "Optimal",
                    "cpu": "Efficient",
                    "disk": "Normal",
                    "network": "Stable"
                },
                "recommendations": [
                    "Performance is optimal",
                    "Consider implementing advanced monitoring",
                    "Explore performance optimization opportunities"
                ]
            }
            
        return {
            "timestamp": datetime.now().isoformat(),
            "scenario": scenario,
            "processing_metrics": {
                "total_time": total_time,
                "rows_per_second": random.randint(10000, 50000),
                "memory_usage": memory_usage,
                "cpu_usage": cpu_usage
            },
            "resource_utilization": {
                "memory_peak": memory_usage * 1.5,
                "cpu_peak": cpu_usage * 1.2,
                "disk_io": random.randint(100, 500),
                "network_io": random.randint(50, 200)
            },
            "detailed_analysis": analysis
        }
        
    def generate_custom_report(self) -> Dict[str, Any]:
        """Generate example custom report with different scenarios."""
        scenario = self._get_current_scenario()
        
        if scenario == "normal_operation":
            success_rate = random.uniform(0.98, 1.0)
            error_count = random.randint(0, 2)
            anomalies = random.randint(0, 1)
            analysis = {
                "status": "Healthy",
                "trends": {
                    "success_rate": "Stable",
                    "error_rate": "Low",
                    "processing_time": "Consistent"
                },
                "anomalies": {
                    "detected": anomalies > 0,
                    "count": anomalies,
                    "severity": "Low" if anomalies > 0 else "None",
                    "description": "Minor processing variations" if anomalies > 0 else "No anomalies detected"
                },
                "recommendations": [
                    "Continue current monitoring practices",
                    "Implement automated alerting for anomalies",
                    "Review processing patterns for optimization"
                ]
            }
        elif scenario == "high_load":
            success_rate = random.uniform(0.95, 0.98)
            error_count = random.randint(2, 5)
            anomalies = random.randint(1, 2)
            analysis = {
                "status": "Under Pressure",
                "trends": {
                    "success_rate": "Declining",
                    "error_rate": "Increasing",
                    "processing_time": "Variable"
                },
                "anomalies": {
                    "detected": True,
                    "count": anomalies,
                    "severity": "Moderate",
                    "description": "Increased error rate during peak load"
                },
                "recommendations": [
                    "Implement load balancing",
                    "Review error handling procedures",
                    "Optimize processing during peak hours"
                ]
            }
        elif scenario == "data_quality_issues":
            success_rate = random.uniform(0.90, 0.95)
            error_count = random.randint(5, 10)
            anomalies = random.randint(2, 4)
            analysis = {
                "status": "Degraded",
                "trends": {
                    "success_rate": "Declining",
                    "error_rate": "High",
                    "processing_time": "Inconsistent"
                },
                "anomalies": {
                    "detected": True,
                    "count": anomalies,
                    "severity": "High",
                    "description": "Multiple processing issues detected"
                },
                "recommendations": [
                    "Implement strict data validation",
                    "Review error handling procedures",
                    "Schedule system maintenance"
                ]
            }
        elif scenario == "performance_degradation":
            success_rate = random.uniform(0.95, 0.98)
            error_count = random.randint(3, 7)
            anomalies = random.randint(1, 3)
            analysis = {
                "status": "Problematic",
                "trends": {
                    "success_rate": "Declining",
                    "error_rate": "Increasing",
                    "processing_time": "Degrading"
                },
                "anomalies": {
                    "detected": True,
                    "count": anomalies,
                    "severity": "Moderate",
                    "description": "Performance issues affecting processing"
                },
                "recommendations": [
                    "Review system resources",
                    "Optimize processing pipeline",
                    "Implement performance monitoring"
                ]
            }
        else:  # custom_analysis
            success_rate = random.uniform(0.99, 1.0)
            error_count = random.randint(0, 1)
            anomalies = random.randint(0, 1)
            analysis = {
                "status": "Excellent",
                "trends": {
                    "success_rate": "Stable",
                    "error_rate": "Minimal",
                    "processing_time": "Optimal"
                },
                "anomalies": {
                    "detected": anomalies > 0,
                    "count": anomalies,
                    "severity": "Low" if anomalies > 0 else "None",
                    "description": "Minor processing variations" if anomalies > 0 else "No anomalies detected"
                },
                "recommendations": [
                    "Implement advanced monitoring",
                    "Explore optimization opportunities",
                    "Consider predictive analytics"
                ]
            }
            
        return {
            "timestamp": datetime.now().isoformat(),
            "scenario": scenario,
            "report_type": "custom_analysis",
            "metrics": {
                "total_processed": random.randint(1000, 5000),
                "success_rate": success_rate,
                "error_count": error_count,
                "processing_duration": random.uniform(0.5, 2.0)
            },
            "analysis": {
                "trends": ["increasing", "stable", "decreasing"][random.randint(0, 2)],
                "anomalies": anomalies,
                "recommendations": [
                    "Optimize query performance",
                    "Increase resource allocation",
                    "Review data quality",
                    "Check network connectivity",
                    "Monitor system resources"
                ][random.randint(0, 4)]
            },
            "detailed_analysis": analysis
        }
        
    def update_metrics(self) -> None:
        """Update Prometheus metrics with new values based on scenario."""
        scenario = self._get_current_scenario()
        
        # Update table scan metrics
        if scenario == "normal_operation":
            rows = random.randint(1000000, 2000000)
        elif scenario == "high_load":
            rows = random.randint(4000000, 5000000)
        elif scenario == "data_quality_issues":
            rows = random.randint(1000000, 2000000)
        elif scenario == "performance_degradation":
            rows = random.randint(3000000, 4000000)
        else:  # custom_analysis
            rows = random.randint(2000000, 3000000)
        table_scan_rows.set(rows)
        
        # Update data quality metrics
        if scenario == "normal_operation":
            score = random.randint(90, 100)
        elif scenario == "high_load":
            score = random.randint(85, 95)
        elif scenario == "data_quality_issues":
            score = random.randint(70, 80)
        elif scenario == "performance_degradation":
            score = random.randint(80, 90)
        else:  # custom_analysis
            score = random.randint(95, 100)
        data_quality_score.set(score)
        
        # Update processing time metrics
        if scenario == "normal_operation":
            time = random.uniform(1.0, 2.0)
        elif scenario == "high_load":
            time = random.uniform(3.0, 5.0)
        elif scenario == "data_quality_issues":
            time = random.uniform(2.0, 3.0)
        elif scenario == "performance_degradation":
            time = random.uniform(4.0, 5.0)
        else:  # custom_analysis
            time = random.uniform(0.5, 1.5)
        processing_time.observe(time)
        
        # Update custom report metrics
        custom_report_count.inc()
        
    def run(self, duration: int = 300) -> None:
        """Run the example generator for specified duration."""
        print("Starting example report and metric generation...")
        print(f"Reports will be saved to: {self.reports_dir}")
        print(f"Metrics will be available at: http://localhost:8000/metrics")
        print(f"Running for {duration} seconds...")
        print("\nScenarios will cycle every minute:")
        for i, scenario in enumerate(self.scenarios):
            print(f"{i+1}. {scenario}")
        
        start_time = time.time()
        while time.time() - start_time < duration:
            # Generate reports
            table_scan = self.generate_table_scan_report()
            data_quality = self.generate_data_quality_report()
            performance = self.generate_performance_report()
            custom = self.generate_custom_report()
            
            # Save reports
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            with open(self.reports_dir / f"table_scan_{timestamp}.json", "w") as f:
                json.dump(table_scan, f, indent=2)
            with open(self.reports_dir / f"data_quality_{timestamp}.json", "w") as f:
                json.dump(data_quality, f, indent=2)
            with open(self.reports_dir / f"performance_{timestamp}.json", "w") as f:
                json.dump(performance, f, indent=2)
            with open(self.reports_dir / f"custom_{timestamp}.json", "w") as f:
                json.dump(custom, f, indent=2)
                
            # Update metrics
            self.update_metrics()
            
            # Wait for next interval
            time.sleep(5)
            
        print("Example generation completed!")
        
def main():
    """Start the example generator."""
    # Start Prometheus metrics server
    start_http_server(8000)
    
    # Create and run generator
    generator = ExampleGenerator()
    generator.run()
    
if __name__ == "__main__":
    main()