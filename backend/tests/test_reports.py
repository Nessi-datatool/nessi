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

import pytest
from pathlib import Path
from src.reports.report_generator import ReportGenerator
from src.reports import (
    data_quality_report,
    performance_report,
    table_scan_report,
    custom_report
)

class TestReportGenerator(ReportGenerator):
    """Test implementation of ReportGenerator"""
    def generate(self, data, **kwargs):
        return self.save_report({"test": "data"}, "test_report", "json")

def test_data_quality_report_generation():
    """Test data quality report generation"""
    test_data = {
        "metrics": {"completeness": 0.95},
        "rows": [
            {"id": 1, "name": "test1", "value": None},
            {"id": 2, "name": "test2", "value": "valid"}
        ]
    }
    report = data_quality_report.generate_report(test_data)
    assert report is not None
    assert "data_quality" in report
    assert "metrics" in report["data_quality"]

def test_performance_report_generation():
    """Test performance report generation"""
    test_data = {
        "total_time": 100,
        "average_time": 50,
        "max_time": 150,
        "peak_memory": 1024,
        "average_memory": 512,
        "memory_growth": 256,
        "rows_per_second": 1000,
        "bytes_per_second": 10240
    }
    report = performance_report.generate_report(test_data)
    assert report is not None
    assert "performance" in report
    assert "metrics" in report["performance"]

def test_table_scan_report_generation():
    """Test table scan report generation"""
    test_data = {
        "columns": ["id", "name", "value"],
        "data_types": {"id": "int", "name": "string", "value": "string"},
        "nullable": {"id": False, "name": True, "value": True},
        "row_count": 100,
        "column_stats": {"id": {"min": 1, "max": 100}},
        "size_bytes": 1024,
        "missing_values": {"name": 5, "value": 10},
        "unique_values": {"id": 100, "name": 95}
    }
    report = table_scan_report.generate_report(test_data)
    assert report is not None
    assert "table_scan" in report
    assert "metrics" in report["table_scan"]

def test_custom_report_generation():
    """Test custom report generation"""
    test_data = {
        "metrics": {"custom_metric": 0.95},
        "thresholds": {"custom_metric": 0.9},
        "alerts": ["Warning: custom_metric near threshold"],
        "derived_metrics": {"ratio": 1.5},
        "aggregated_metrics": {"total": 100},
        "trend_metrics": {"slope": 0.1},
        "status": "healthy"
    }
    report = custom_report.generate_report(test_data)
    assert report is not None
    assert "custom" in report
    assert "metrics" in report["custom"]

def test_report_generator():
    """Test report generator functionality"""
    # Create test generator
    generator = TestReportGenerator()
    
    # Generate test report
    report_path = generator.generate({"test": "data"})
    
    # Verify report was created
    assert Path(report_path).exists()
    
    # Clean up
    Path(report_path).unlink()

def test_report_generator_validation():
    """Test report generator validation"""
    generator = TestReportGenerator()
    
    # Test validation success
    assert generator.validate_data({"field1": "value1"}, ["field1"]) is True
    
    # Test validation failure
    with pytest.raises(ValueError):
        generator.validate_data({"field1": "value1"}, ["field2"])

def test_report_generator_timestamp():
    """Test report generator timestamp"""
    generator = TestReportGenerator()
    timestamp = generator.get_timestamp()
    assert isinstance(timestamp, str)
    assert len(timestamp) == 15  # YYYYMMDD_HHMMSS