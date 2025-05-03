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

import os
import shutil
from datetime import date
from typing import Optional, Union, List
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, DoubleType, DateType
from pyspark.sql.functions import lit
from .spark_config import get_spark_session

def create_test_delta_table(
    table_path: str,
    schema: Optional[StructType] = None,
    num_rows: int = 4,
    partition_by: Optional[Union[str, List[str]]] = None
) -> None:
    """Create a test Delta table with sample data."""
    if not table_path:
        raise ValueError("Table path cannot be empty")
    try:
        spark = get_spark_session()
        if not schema:
            schema = StructType([
                StructField("id", IntegerType(), False),
                StructField("name", StringType(), True),
                StructField("age", IntegerType(), True),
                StructField("salary", DoubleType(), True),
                StructField("department", StringType(), True),
                StructField("join_date", StringType(), True)
            ])
        data = []
        for i in range(1, num_rows + 1):
            data.append((
                i,
                f"User_{i}",
                25 + i,
                50000.0 + i * 5000.0,
                f"Department_{i}",
                f"2020-{i if i <= 12 else 12}-01"
            ))
        df = spark.createDataFrame(data, schema)
        if partition_by:
            df.write.format("delta").partitionBy(partition_by).mode("overwrite").save(table_path)
        else:
            df.write.format("delta").mode("overwrite").save(table_path)
    except Exception as e:
        if os.path.exists(table_path):
            shutil.rmtree(table_path, ignore_errors=True)
        raise RuntimeError(f"Failed to create test table: {str(e)}")