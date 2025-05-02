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