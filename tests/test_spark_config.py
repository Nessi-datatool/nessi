import pytest
import os
from backend.src.scanner.spark_config import get_spark_session

def test_spark_session_creation():
    """Test basic Spark session creation"""
    spark = get_spark_session("TestApp")
    assert spark is not None
    assert spark.conf.get("spark.app.name") == "TestApp"
    spark.stop()

def test_delta_lake_configuration():
    """Test Delta Lake related configurations"""
    spark = get_spark_session()
    
    # Check Delta Lake extensions
    assert spark.conf.get("spark.sql.extensions") == "io.delta.sql.DeltaSparkSessionExtension"
    assert spark.conf.get("spark.sql.catalog.spark_catalog") == "org.apache.spark.sql.delta.catalog.DeltaCatalog"
    
    # Check Delta Lake jars
    jars = spark.conf.get("spark.jars")
    assert "delta-core" in jars
    assert "delta-storage" in jars
    
    spark.stop()

def test_memory_settings():
    """Test memory configurations"""
    spark = get_spark_session()
    
    assert spark.conf.get("spark.driver.memory") == "2g"
    assert spark.conf.get("spark.executor.memory") == "2g"
    
    spark.stop()

def test_legacy_settings():
    """Test legacy mode settings"""
    spark = get_spark_session()
    
    assert spark.conf.get("spark.sql.legacy.timeParserPolicy") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.int96RebaseModeInRead") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.int96RebaseModeInWrite") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.datetimeRebaseModeInRead") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.datetimeRebaseModeInWrite") == "LEGACY"
    
    spark.stop()

def test_log_level():
    """Test log level configuration"""
    spark = get_spark_session()
    assert spark.sparkContext.getLogLevel() == "WARN"
    spark.stop()

def test_custom_jars_directory(monkeypatch):
    """Test custom jars directory configuration"""
    custom_jars_dir = "/custom/jars/path"
    monkeypatch.setenv("SPARK_JARS_DIR", custom_jars_dir)
    
    spark = get_spark_session()
    jars = spark.conf.get("spark.jars")
    
    assert custom_jars_dir in jars
    spark.stop() 