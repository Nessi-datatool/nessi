import pytest
import os
import tempfile
from pathlib import Path
import pandas as pd
from backend.src.scanner.spark_config import get_spark_session

@pytest.fixture(scope="session")
def spark_session():
    spark = get_spark_session("TestSession")
    yield spark
    spark.stop()

@pytest.fixture
def temp_dir():
    with tempfile.TemporaryDirectory() as tmpdirname:
        yield tmpdirname 