"""Tests for the report generator module."""

import pytest
from pathlib import Path
from typing import Dict, Any
from datetime import datetime
from nessi.scanner.models import (
    ScanResult,
    SchemaValidationResult,
    QualityMetrics,
    PerformanceMetrics
)
from nessi.reports.report_generator import ReportGenerator

@pytest.fixture
def sample_scan_result():
    """Create a sample scan result."""
    return ScanResult(
        table_name="test_table",
        row_count=1000,
        column_count=5,
        scan_timestamp=datetime.now().isoformat(),
        scan_duration=1.5,
        file_size=1024 * 1024,
        file_format="CSV",
        location="/path/to/data.csv",
        schema_validation=SchemaValidationResult(
            schema_name="test_schema",
            schema_definition={"fields": [
                {"name": "id", "type": "string", "nullable": False},
                {"name": "name", "type": "string", "nullable": True}
            ]},
            is_valid=True,
            errors=[]
        ),
        quality_metrics=QualityMetrics(
            completeness={"overall": 0.95, "id": 1.0, "name": 0.9},
            accuracy={"overall": 0.98, "id": 0.99, "name": 0.97},
            consistency={"overall": 0.90, "id": 0.95, "name": 0.85},
            uniqueness={"overall": 0.92, "id": 1.0, "name": 0.84},
            timeliness={"overall": 0.94, "id": 0.95, "name": 0.93},
            completeness_score=0.95,
            accuracy_score=0.98,
            consistency_score=0.90,
            uniqueness_score=0.92,
            timeliness_score=0.94,
            overall_score=0.94
        ),
        performance_metrics=PerformanceMetrics(
            scan_duration_ms=5000.0,
            memory_usage_mb=1024.0,
            cpu_usage_percent=50.0,
            io_operations=100,
            rows_processed=1000,
            bytes_processed=1024000,
            start_time=datetime.now(),
            end_time=datetime.now()
        )
    )

@pytest.fixture
def temp_report_dir(tmp_path):
    """Create a temporary directory for reports."""
    return tmp_path / "reports"

def test_report_generator_initialization():
    """Test report generator initialization."""
    generator = ReportGenerator()
    assert generator.title == "Data Quality Report"
    assert generator.author == "Nessi"
    assert generator.subject == "Data Quality Analysis"
    
    generator = ReportGenerator(
        title="Custom Report",
        author="Custom Author",
        subject="Custom Subject"
    )
    assert generator.title == "Custom Report"
    assert generator.author == "Custom Author"
    assert generator.subject == "Custom Subject"

def test_generate_report(sample_scan_result, temp_report_dir):
    """Test report generation."""
    generator = ReportGenerator()
    output_file = temp_report_dir / "report.html"
    
    generator.generate(sample_scan_result, output_file)
    assert output_file.exists()
    assert output_file.suffix == ".html"

def test_report_content(sample_scan_result, temp_report_dir):
    """Test report content."""
    generator = ReportGenerator()
    output_file = temp_report_dir / "report.html"
    
    generator.generate(sample_scan_result, output_file)
    
    # Verify report content
    content = output_file.read_text()
    assert "Data Quality Report" in content
    assert "test_table" in content
    assert "1000 rows" in content
    assert "5 columns" in content
    assert "95%" in content  # completeness
    assert "98%" in content  # accuracy
    assert "90%" in content  # consistency
    assert "5.0 seconds" in content  # scan time
    assert "1024 MB" in content  # memory usage
    assert "50.0%" in content  # cpu usage

def test_report_sections(sample_scan_result, temp_report_dir):
    """Test report sections."""
    generator = ReportGenerator()
    output_file = temp_report_dir / "report.html"
    
    generator.generate(sample_scan_result, output_file)
    
    # Verify report sections
    content = output_file.read_text()
    assert "Table Information" in content
    assert "Schema Validation" in content
    assert "Quality Metrics" in content
    assert "Performance Metrics" in content
    assert "Recommendations" in content

def test_report_formatting(sample_scan_result, temp_report_dir):
    """Test report formatting."""
    generator = ReportGenerator()
    output_file = temp_report_dir / "report.html"
    
    generator.generate(sample_scan_result, output_file)
    
    # Verify report formatting
    content = output_file.read_text()
    assert "<h1>" in content  # headings
    assert "<p>" in content  # paragraphs
    assert "<table>" in content  # tables
    assert "<tr>" in content  # table rows
    assert "<td>" in content  # table cells

def test_report_validation(sample_scan_result, temp_report_dir):
    """Test report validation."""
    generator = ReportGenerator()
    
    # Test with invalid scan result
    invalid_result = ScanResult(
        table_name="",
        row_count=-1,
        column_count=0,
        scan_timestamp=datetime.now().isoformat(),
        scan_duration=-1.0,
        file_size=-1,
        file_format="INVALID",
        location="/invalid/path",
        schema_validation=SchemaValidationResult(
            schema_name="invalid_schema",
            schema_definition={"fields": []},
            is_valid=False,
            errors=["Invalid schema"]
        ),
        quality_metrics=QualityMetrics(
            completeness={"overall": -0.1, "id": -0.1},
            accuracy={"overall": 1.5, "id": 1.5},
            consistency={"overall": 2.0, "id": 2.0},
            uniqueness={"overall": -0.1, "id": -0.1},
            timeliness={"overall": -0.1, "id": -0.1}
        ),
        performance_metrics=PerformanceMetrics(
            scan_duration_ms=-1000.0,
            memory_usage_mb=-1.0,
            cpu_usage_percent=-1.0,
            io_operations=-1,
            rows_processed=-1,
            bytes_processed=-1,
            start_time=datetime.now(),
            end_time=datetime.now()
        )
    )
    
    with pytest.raises(ValueError):
        generator.generate(invalid_result, temp_report_dir / "report.html")
    
    # Test with invalid output directory
    with pytest.raises(ValueError):
        generator.generate(sample_scan_result, Path("nonexistent/report.html"))
    
    # Test with invalid file extension
    with pytest.raises(ValueError):
        generator.generate(sample_scan_result, temp_report_dir / "report.txt")

def test_report_customization(sample_scan_result, temp_report_dir):
    """Test report customization."""
    # Test with custom title
    generator = ReportGenerator(title="Custom Title")
    output_file = temp_report_dir / "report.html"
    generator.generate(sample_scan_result, output_file)
    content = output_file.read_text()
    assert "Custom Title" in content
    
    # Test with custom author
    generator = ReportGenerator(author="Custom Author")
    output_file = temp_report_dir / "report.html"
    generator.generate(sample_scan_result, output_file)
    content = output_file.read_text()
    assert "Custom Author" in content
    
    # Test with custom subject
    generator = ReportGenerator(subject="Custom Subject")
    output_file = temp_report_dir / "report.html"
    generator.generate(sample_scan_result, output_file)
    content = output_file.read_text()
    assert "Custom Subject" in content 