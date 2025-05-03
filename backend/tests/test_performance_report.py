"""Tests for performance report generation."""

import pytest
from datetime import datetime
from src.reports.performance_report import (
    PerformanceReport,
    generate_report,
    generate_report_from_data
)

@pytest.fixture
def sample_performance_data():
    """Sample performance data for testing."""
    return {
        "processing_metrics": {
            "score": 85.5,
            "execution_time": {
                "query_execution": 150,
                "data_loading": 200,
                "data_processing": 300
            },
            "memory_usage": {
                "spark": {"current": 1024, "peak": 2048},
                "system": {"current": 512, "peak": 1024}
            },
            "cpu_usage": {
                "spark": 75.5,
                "system": 60.0
            }
        },
        "resource_utilization": {
            "efficiency_score": 90.0,
            "disk_io": {
                "read": {"read": 100, "write": 50},
                "write": {"read": 0, "write": 200}
            },
            "network_traffic": {
                "spark": {"inbound": 10, "outbound": 20},
                "system": {"inbound": 5, "outbound": 10}
            },
            "cache_performance": {
                "memory": {"hit_rate": 80.0, "miss_rate": 20.0},
                "disk": {"hit_rate": 70.0, "miss_rate": 30.0}
            }
        }
    }

def test_performance_report_generation(sample_performance_data):
    """Test basic performance report generation."""
    report = PerformanceReport()
    result = report.generate(sample_performance_data)
    
    assert "performance" in result
    assert "metrics" in result["performance"]
    assert "execution_time" in result["performance"]["metrics"]
    assert "memory_usage" in result["performance"]["metrics"]
    assert "throughput" in result["performance"]["metrics"]

def test_calculate_execution_time(sample_performance_data):
    """Test execution time calculation."""
    report = PerformanceReport()
    result = report._calculate_execution_time(sample_performance_data)
    
    assert "total_time" in result
    assert "average_time" in result
    assert "max_time" in result

def test_calculate_memory_usage(sample_performance_data):
    """Test memory usage calculation."""
    report = PerformanceReport()
    result = report._calculate_memory_usage(sample_performance_data)
    
    assert "peak_memory" in result
    assert "average_memory" in result
    assert "memory_growth" in result

def test_calculate_throughput(sample_performance_data):
    """Test throughput calculation."""
    report = PerformanceReport()
    result = report._calculate_throughput(sample_performance_data)
    
    assert "rows_per_second" in result
    assert "bytes_per_second" in result

def test_generate_report_function(sample_performance_data):
    """Test the generate_report function."""
    result = generate_report(sample_performance_data)
    assert "performance" in result
    assert "metrics" in result["performance"]

def test_generate_report_from_data_markdown(sample_performance_data, tmp_path):
    """Test markdown report generation."""
    result = generate_report_from_data(
        sample_performance_data,
        format="md",
        include_details=True
    )
    assert result.endswith(".md")
    assert "Performance Report" in open(result).read()

def test_generate_report_from_data_json(sample_performance_data, tmp_path):
    """Test JSON report generation."""
    result = generate_report_from_data(
        sample_performance_data,
        format="json",
        include_details=True
    )
    assert result.endswith(".json")
    content = open(result).read()
    assert '"performance"' in content
    assert '"metrics"' in content

def test_generate_report_from_data_html(sample_performance_data, tmp_path):
    """Test HTML report generation."""
    result = generate_report_from_data(
        sample_performance_data,
        format="html",
        include_details=True
    )
    assert result.endswith(".html")
    content = open(result).read()
    assert "<html" in content
    assert "Performance Report" in content

def test_invalid_data_format():
    """Test handling of invalid data format."""
    with pytest.raises(RuntimeError):
        generate_report_from_data({}, format="md")

def test_missing_required_fields():
    """Test handling of missing required fields."""
    with pytest.raises(RuntimeError):
        generate_report_from_data({"processing_metrics": {}}, format="md")

def test_invalid_format():
    """Test handling of invalid format."""
    with pytest.raises(RuntimeError):
        generate_report_from_data(
            {"processing_metrics": {}, "resource_utilization": {}},
            format="invalid"
        ) 