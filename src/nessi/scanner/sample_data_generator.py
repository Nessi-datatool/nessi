"""Sample data generator for testing and demonstration."""

import logging
from typing import Dict, List, Optional
from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, TimestampType
from .spark_config import get_spark_session
import os
from pathlib import Path

// ... existing code ... 