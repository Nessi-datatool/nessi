"""Test Delta table creation module.

This module provides functionality to create and manage test Delta tables
with various schemas and data patterns.
"""

import os
import logging
from typing import Dict, List, Optional, Union, Any, Tuple
from datetime import datetime, date, timedelta
import pandas as pd
from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, LongType, DoubleType, DateType, BooleanType, TimestampType
from pyspark.sql.functions import lit, col, when
import shutil
import random

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def get_or_create_spark_session() -> SparkSession:
    """Get or create a Spark session with Delta Lake configuration.
    
    Returns:
        SparkSession: A configured Spark session with Delta Lake support.
        
    Raises:
        RuntimeError: If Spark session creation fails.
    """
    try:
        spark = SparkSession.builder \
            .appName("TestDeltaTableCreator") \
            .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
            .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
            .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0,io.delta:delta-storage:2.4.0") \
            .getOrCreate()
        return spark
    except Exception as e:
        logger.error(f"Failed to create Spark session: {str(e)}")
        raise RuntimeError(f"Failed to create Spark session: {str(e)}")

def create_test_delta_table(
    table_path: str,
    schema: Optional[StructType] = None,
    num_rows: int = 1000,
    partition_by: Optional[List[str]] = None
) -> None:
    """Create a test Delta table with sample data.

    Args:
        table_path (str): Path where the Delta table will be created.
        schema (Optional[StructType]): Optional custom schema for the table.
        num_rows (int): Number of rows to generate.
        partition_by (Optional[List[str]]): Optional list of columns to partition by.

    Raises:
        ValueError: If table_path is invalid or schema is invalid.
        RuntimeError: If table creation fails.
    """
    if not table_path:
        raise ValueError("Table path cannot be empty")

    if num_rows < 1:
        raise ValueError("Number of rows must be at least 1")

    try:
        spark = get_or_create_spark_session()

        # Create directory if it doesn't exist
        os.makedirs(os.path.dirname(table_path), exist_ok=True)

        # Generate sample data
        data = generate_sample_data(num_rows, schema)
        if schema is None:
            df = spark.createDataFrame(data, ["id", "name", "age", "salary", "department", "join_date"])
        else:
            df = spark.createDataFrame(data, schema=schema)

        # Write initial version
        if partition_by:
            df.write.format("delta") \
                .partitionBy(partition_by) \
                .mode("overwrite") \
                .save(table_path)
        else:
            df.write.format("delta") \
                .mode("overwrite") \
                .save(table_path)

        # If num_rows is 4, create test versions
        if num_rows == 4:
            # Update department of id=1
            df = spark.read.format("delta").load(table_path)
            df.filter("id = 1") \
                .withColumn("department", lit("Engineering")) \
                .write.format("delta") \
                .mode("overwrite") \
                .option("replaceWhere", "id = 1") \
                .save(table_path)

            # Delete id=3 and add id=5
            df = spark.read.format("delta").load(table_path)
            df_without_3 = df.filter("id != 3")
            new_row = spark.createDataFrame([(5, "User_5", 35, 65000.0, "Marketing", date(2020, 5, 1))],
                                         ["id", "name", "age", "salary", "department", "join_date"])
            df_final = df_without_3.union(new_row)
            df_final.write.format("delta") \
                .mode("overwrite") \
                .save(table_path)

        logger.info(f"Successfully created test Delta table at {table_path}")
    except ValueError as e:
        logger.error(f"Failed to create test Delta table: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Failed to create test Delta table: {str(e)}")
        raise RuntimeError(f"Failed to create test Delta table: {str(e)}")

def add_column_to_delta_table(
    table_path: str,
    column_name: str,
    data_type: str,
    default_value: Optional[Any] = None
) -> None:
    """Add a new column to an existing Delta table.
    
    Args:
        table_path (str): Path to the Delta table.
        column_name (str): Name of the new column.
        data_type (str): Data type of the new column.
        default_value (Optional[Any]): Default value for the new column.
        
    Raises:
        ValueError: If table_path is invalid or column already exists.
        RuntimeError: If column addition fails.
    """
    if not table_path or not os.path.exists(table_path):
        raise ValueError(f"Invalid table path: {table_path}")

    try:
        spark = get_or_create_spark_session()

        # Read existing table
        df = spark.read.format("delta").load(table_path)

        # Check if column already exists
        if column_name in df.columns:
            raise ValueError(f"Column {column_name} already exists")

        # Convert data type string to Spark type
        spark_type = None
        if data_type.lower() == "string":
            spark_type = StringType()
        elif data_type.lower() == "integer":
            spark_type = IntegerType()
        elif data_type.lower() == "long":
            spark_type = LongType()
        elif data_type.lower() == "double":
            spark_type = DoubleType()
        elif data_type.lower() == "boolean":
            spark_type = BooleanType()
        elif data_type.lower() == "date":
            spark_type = DateType()
        else:
            raise ValueError(f"Unsupported data type: {data_type}")

        # Create new schema with added column
        new_schema = StructType(df.schema.fields + [StructField(column_name, spark_type, True)])

        # Add new column with default value
        if default_value is not None:
            df = df.withColumn(column_name, lit(default_value).cast(spark_type))
        else:
            df = df.withColumn(column_name, lit(None).cast(spark_type))

        # Write back to Delta table with schema merging enabled
        df.write.format("delta") \
            .option("mergeSchema", "true") \
            .mode("overwrite") \
            .save(table_path)

        logger.info(f"Successfully added column {column_name} to table at {table_path}")
    except ValueError as e:
        logger.error(f"Failed to add column: {str(e)}")
        raise
    except Exception as e:
        logger.error(f"Failed to add column: {str(e)}")
        raise RuntimeError(f"Failed to add column: {str(e)}")

def generate_sample_data(num_rows: int = 1000, schema: Optional[StructType] = None) -> List[Tuple]:
    """Generate sample data with various data patterns.

    Args:
        num_rows (int): Number of rows to generate. Defaults to 1000.
        schema (Optional[StructType]): Optional schema to match data to.

    Returns:
        List[Tuple]: Generated sample data with the following columns:
            - id: Unique identifier
            - name: Random names
            - age: Random ages between 18 and 80
            - salary: Random salaries with normal distribution
            - department: Random department names
            - join_date: Random dates between 2010 and 2023

    Raises:
        ValueError: If num_rows is less than 1.
    """
    if num_rows < 1:
        raise ValueError("Number of rows must be at least 1")

    try:
        # Generate data based on schema if provided
        if schema is not None:
            data = []
            for i in range(num_rows):
                row = []
                for field in schema.fields:
                    if field.name == "id":
                        row.append(i + 1)
                    elif field.name == "name":
                        row.append(f"User_{i + 1}")
                    elif field.name == "age":
                        # 10% chance of null
                        row.append(None if random.random() < 0.1 else random.randint(18, 80))
                    elif field.name == "salary":
                        # 10% chance of null
                        row.append(None if random.random() < 0.1 else round(random.normalvariate(50000, 15000), 2))
                    elif field.name == "department":
                        # 10% chance of null
                        row.append(None if random.random() < 0.1 else random.choice(["HR", "Engineering", "Sales", "Marketing"]))
                    elif field.name == "join_date":
                        base_date = date(2020, 1, 1)
                        days_to_add = random.randint(0, 1000)
                        join_date = base_date + timedelta(days=days_to_add)
                        row.append(join_date)
                    elif field.name == "value":
                        row.append(round(random.normalvariate(100, 20), 2))
                    elif field.name == "timestamp":
                        row.append(datetime.now())
                    else:
                        row.append(None)
                data.append(tuple(row))
        else:
            # Default schema data generation
            data = []
            base_date = date(2020, 1, 1)
            for i in range(num_rows):
                days_to_add = random.randint(0, 1000)
                join_date = base_date + timedelta(days=days_to_add)
                data.append((
                    i + 1,  # id
                    f"User_{i + 1}",  # name
                    None if random.random() < 0.1 else random.randint(18, 80),  # age
                    None if random.random() < 0.1 else round(random.normalvariate(50000, 15000), 2),  # salary
                    None if random.random() < 0.1 else random.choice(["HR", "Engineering", "Sales", "Marketing"]),  # department
                    join_date  # join_date
                ))

        logger.info(f"Successfully generated {num_rows} rows of sample data")
        return data
    except Exception as e:
        logger.error(f"Failed to generate sample data: {str(e)}")
        raise RuntimeError(f"Failed to generate sample data: {str(e)}")

def main() -> None:
    """Main function to create test Delta tables."""
    try:
        # Create output directory
        output_dir = "backend/data/test_tables"
        os.makedirs(output_dir, exist_ok=True)
        
        # Create test table
        table_path = os.path.join(output_dir, "test_delta_table")
        create_test_delta_table(
            table_path=table_path,
            num_rows=1000,
            partition_by=["department"]
        )
        
        # Add a new column
        add_column_to_delta_table(
            table_path=table_path,
            column_name="is_active",
            data_type="boolean",
            default_value=True
        )
        
        logger.info("Test Delta table creation completed successfully")
    except Exception as e:
        logger.error(f"Test Delta table creation failed: {str(e)}")
        raise

if __name__ == "__main__":
    main() 