# Testing Guide

## Overview

This guide covers all aspects of testing in Nessi, which is exclusively run in Docker containers. The testing environment includes all necessary services (Redis, MongoDB, Elasticsearch, etc.) and is isolated from the host system.

## Test Structure

### 1. Unit Tests

```python
# tests/test_delta_manager.py
from nessi import DeltaManager
import pytest

def test_create_table(spark):
    manager = DeltaManager(spark)
    manager.create_table(
        path="/app/data/test_table",
        schema="id INT, name STRING",
        partition_by=["id"]
    )
    assert manager.table_exists("/app/data/test_table")

def test_time_travel(spark):
    manager = DeltaManager(spark)
    df = manager.time_travel(
        path="/app/data/test_table",
        version=1
    )
    assert df.count() > 0

def test_schema_evolution():
    manager = DeltaManager(spark)
    manager.evolve_schema("/data/test_table", {
        "new_column": {
            "type": "STRING",
            "nullable": True
        }
    })
    assert "new_column" in manager.get_schema("/data/test_table")
```

### 2. Integration Tests

```python
# tests/test_quality_analyzer.py
from nessi import QualityAnalyzer, DeltaManager

def test_quality_analysis(spark):
    manager = DeltaManager(spark)
    analyzer = QualityAnalyzer()
    
    df = manager.read_table("/app/data/test_table")
    results = analyzer.analyze_dataframe(df)
    
    assert "overall_quality" in results
    assert "completeness" in results
    assert "consistency" in results

def test_anomaly_detection(spark):
    analyzer = QualityAnalyzer()
    df = spark.createDataFrame([
        (1, 100),
        (2, 200),
        (3, 300),
        (4, 1000)  # Anomaly
    ], ["id", "value"])
    
    anomalies = analyzer._detect_anomalies(df)
    assert len(anomalies) > 0
```

### 3. Performance Tests

```python
# tests/test_performance.py
import time
from nessi import DeltaManager

def test_optimization_performance():
    manager = DeltaManager(spark)
    start_time = time.time()
    
    manager.optimize_table(
        path="/data/large_table",
        z_order_by=["timestamp"]
    )
    
    duration = time.time() - start_time
    assert duration < 300  # Should complete within 5 minutes

def test_query_performance():
    manager = DeltaManager(spark)
    start_time = time.time()
    
    df = manager.read_table(
        path="/data/large_table",
        filters={"timestamp": "2024-01-01"}
    )
    df.count()
    
    duration = time.time() - start_time
    assert duration < 10  # Should complete within 10 seconds
```

## Running Tests

### 1. Using Docker

All tests are run in Docker containers. Use the provided script:

```bash
# Make the script executable
chmod +x docker/run-tests.sh

# Run all tests
./docker/run-tests.sh

# Run specific test file
./docker/run-tests.sh tests/test_delta_manager.py

# Run specific test function
./docker/run-tests.sh tests/test_delta_manager.py::test_create_table
```

### 2. Test Results

Test results are automatically copied to the host machine:
- Coverage reports: `test-results/coverage.html`
- JUnit XML: `test-results/junit.xml`
- HTML report: `test-results/report.html`
- Benchmarks: `test-results/benchmark.json`

### 3. Test Services

The test environment includes:
- Redis: For caching and session management
- MongoDB: For document storage
- Elasticsearch: For search functionality
- Grafana: For monitoring
- Prometheus: For metrics collection

## Test Configuration

### 1. Docker Configuration

```yaml
# docker/docker-compose.test.yml
services:
  test-runner:
    build:
      context: ..
      dockerfile: docker/Dockerfile.test
    volumes:
      - ../src:/app/src
      - ../tests:/app/tests
      - ../test-results:/app/test-results
    environment:
      - PYTHONPATH=/app:/app/src
      - ENVIRONMENT=test
      - SPARK_HOME=/opt/spark
```

### 2. Test Fixtures

```python
# tests/conftest.py
import pytest
from pyspark.sql import SparkSession

@pytest.fixture(scope="session")
def spark():
    return SparkSession.builder \
        .master("local[*]") \
        .appName("nessi_test") \
        .config("spark.sql.shuffle.partitions", "2") \
        .config("spark.default.parallelism", "2") \
        .getOrCreate()

@pytest.fixture(scope="function")
def test_table(spark):
    df = spark.createDataFrame([
        (1, "test1"),
        (2, "test2")
    ], ["id", "name"])
    df.write.format("delta").save("/app/data/test_table")
    yield "/app/data/test_table"
    spark.sql("DROP TABLE IF EXISTS delta.`/app/data/test_table`")
```

## Test Best Practices

### 1. Unit Testing

- Test one functionality at a time
- Use meaningful test names
- Include both success and failure cases
- Mock external dependencies
- Clean up test data

### 2. Integration Testing

- Test service interactions
- Verify data flow between components
- Check error handling
- Validate data consistency

### 3. Performance Testing

- Use benchmark markers
- Set realistic timeouts
- Monitor resource usage
- Clean up after tests

## Troubleshooting

Common issues and solutions:

1. **Test Container Issues**
   - Check logs: `docker-compose -f docker-compose.test.yml logs -f`
   - Rebuild containers: `docker-compose -f docker-compose.test.yml build --no-cache`
   - Clean up volumes: `docker-compose -f docker-compose.test.yml down -v`

2. **Test Failures**
   - Check test-results directory for detailed reports
   - Verify service connectivity in Docker network
   - Ensure sufficient resources for Spark

3. **Coverage Issues**
   - Check if all source files are mounted correctly
   - Verify PYTHONPATH settings
   - Ensure test files are properly discovered

## Test Maintenance

### 1. Test Documentation

```python
def test_create_table():
    """
    Test the creation of a Delta table.
    
    Steps:
    1. Initialize DeltaManager
    2. Create table with test schema
    3. Verify table exists
    4. Clean up test table
    
    Expected:
    - Table should be created successfully
    - Table should be accessible
    - Schema should match input
    """
    # Test implementation
```

### 2. Test Data Management

```python
# tests/data/test_data.py
def generate_test_data(spark, size=1000):
    """Generate test data for performance testing."""
    return spark.range(size).withColumn(
        "value",
        F.rand() * 1000
    )
```

### 3. Test Reporting

```bash
# Generate test report
python -m pytest tests/ --junitxml=test-results.xml

# Generate HTML report
python -m pytest tests/ --html=report.html
```

## Next Steps

- [Continuous Integration](ci/README.md)
- [Performance Benchmarking](benchmarks/README.md)
- [Test Automation](automation/README.md) 