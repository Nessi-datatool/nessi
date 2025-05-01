#!/usr/bin/env python3

"""Sample data generation module for creating test datasets.

This module provides functionality to generate sample data with various patterns
and save it in different formats (Delta, Parquet, CSV).
"""

import os
import argparse
import logging
from datetime import datetime, timedelta, date
from typing import List, Tuple, Optional, Union
import numpy as np
import pandas as pd
from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, LongType, DoubleType, DateType, BooleanType
import random
import shutil

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
            .appName("SampleDataGenerator") \
            .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
            .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
            .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0,io.delta:delta-storage:2.4.0") \
            .getOrCreate()
        return spark
    except Exception as e:
        logger.error(f"Failed to create Spark session: {str(e)}")
        raise RuntimeError(f"Failed to create Spark session: {str(e)}")

def generate_sample_data(num_rows: int = 1000) -> List[Tuple]:
    """Generate sample data with various data patterns.
    
    Args:
        num_rows (int): Number of rows to generate. Defaults to 1000.
        
    Returns:
        List[Tuple]: Generated sample data with the following columns:
            - id: Unique identifier
            - name: Random names
            - age: Random ages between 18 and 80
            - salary: Random salaries with normal distribution
            - department: Random department names
            - join_date: Random dates between 2010 and 2023
            - is_active: Random boolean values
            - rating: Random ratings between 1 and 5
            
    Raises:
        ValueError: If num_rows is less than 1.
    """
    if num_rows < 1:
        raise ValueError("Number of rows must be at least 1")
        
    try:
        # Generate dates as datetime64[ns] objects
        dates = pd.date_range(start='2020-01-01', periods=num_rows, freq='D')
        
        # Create DataFrame with consistent types
        df = pd.DataFrame({
            'id': range(1, num_rows + 1),
            'name': [f'User_{i}' for i in range(1, num_rows + 1)],
            'age': [random.randint(20, 60) for _ in range(num_rows)],
            'salary': [random.uniform(50000, 150000) for _ in range(num_rows)],
            'department': [random.choice(['HR', 'Engineering', 'Marketing', 'Sales', 'Finance']) for _ in range(num_rows)],
            'join_date': dates
        })
        
        # Add specific null values using None instead of NaN
        df.loc[df.index % 10 == 0, 'age'] = None
        df.loc[df.index % 5 == 0, 'salary'] = None
        df.loc[df.index % 7 == 0, 'department'] = None
        
        # Convert DataFrame to list of tuples with Python native types
        data = []
        for row in df.itertuples(index=False, name=None):
            converted_row = []
            for val, col in zip(row, df.columns):
                if pd.isna(val):
                    converted_row.append(None)
                elif col == 'id':
                    converted_row.append(int(val))  # Convert to long
                elif col == 'age' and val is not None:
                    converted_row.append(int(val))  # Convert to int
                elif col == 'salary' and val is not None:
                    converted_row.append(float(val))  # Convert to double
                elif col == 'join_date':
                    converted_row.append(val.date())  # Convert to date
                else:
                    converted_row.append(str(val))  # Convert to string
            data.append(tuple(converted_row))
        
        logger.info(f"Successfully generated {num_rows} rows of sample data")
        return data
    except Exception as e:
        logger.error(f"Failed to generate sample data: {str(e)}")
        raise

def save_as_delta(data: Union[List[Tuple], pd.DataFrame], table_path: str) -> None:
    """Save data as a Delta table.
    
    Args:
        data (Union[List[Tuple], pd.DataFrame]): Data to save.
        table_path (str): Path where the Delta table will be created.
        
    Raises:
        ValueError: If data is empty or table_path is invalid.
        RuntimeError: If saving fails.
    """
    if not table_path:
        raise ValueError("Table path cannot be empty")

    if isinstance(data, list) and len(data) == 0:
        raise ValueError("Data cannot be empty")
    elif isinstance(data, pd.DataFrame) and len(data) == 0:
        raise ValueError("DataFrame cannot be empty")

    try:
        spark = get_or_create_spark_session()
        
        # Create directory if it doesn't exist
        os.makedirs(os.path.dirname(table_path), exist_ok=True)
        
        # Convert data to DataFrame if needed
        if isinstance(data, list):
            df = spark.createDataFrame(data, ["id", "name", "age", "salary", "department", "join_date"])
        else:
            # Convert pandas DataFrame to Spark DataFrame with explicit schema
            schema = StructType([
                StructField("id", LongType(), True),
                StructField("name", StringType(), True),
                StructField("age", LongType(), True),
                StructField("salary", DoubleType(), True),
                StructField("department", StringType(), True),
                StructField("join_date", DateType(), True)
            ])
            
            # Convert numeric columns to appropriate types
            data = data.copy()
            data['id'] = data['id'].astype('int64')
            data['age'] = data['age'].fillna(-1).astype('int64')  # Fill NA with -1 for integers
            data['salary'] = data['salary'].astype('float64')
            
            df = spark.createDataFrame(data, schema=schema)

        # Write to Delta table with schema merging enabled
        df.write.format("delta") \
            .option("mergeSchema", "true") \
            .mode("overwrite") \
            .save(table_path)
            
        logger.info(f"Successfully saved data as Delta table at {table_path}")
    except Exception as e:
        logger.error(f"Failed to save data as Delta table: {str(e)}")
        raise RuntimeError(f"Failed to save data as Delta table: {str(e)}")

def save_as_parquet(data: Union[List[Tuple], pd.DataFrame], table_path: str) -> None:
    """Save data as a Parquet file.
    
    Args:
        data (Union[List[Tuple], pd.DataFrame]): Data to save.
        table_path (str): Path where the Parquet file will be created.
        
    Raises:
        ValueError: If data is empty or table_path is invalid.
        RuntimeError: If saving fails.
    """
    if not table_path:
        raise ValueError("Parquet path cannot be empty")

    if isinstance(data, list) and len(data) == 0:
        raise ValueError("Data cannot be empty")
    elif isinstance(data, pd.DataFrame) and len(data) == 0:
        raise ValueError("DataFrame cannot be empty")

    try:
        spark = get_or_create_spark_session()
        
        # Convert data to DataFrame if needed
        if isinstance(data, list):
            df = spark.createDataFrame(data, ["id", "name", "age", "salary", "department", "join_date"])
        else:
            df = spark.createDataFrame(data)
        
        # Write to Parquet
        df.write.parquet(table_path, mode="overwrite")
        logger.info(f"Successfully saved data as Parquet at {table_path}")
    except Exception as e:
        logger.error(f"Failed to save data as Parquet: {str(e)}")
        raise RuntimeError(f"Failed to save data as Parquet: {str(e)}")

def save_as_csv(data: Union[List[Tuple], pd.DataFrame], file_path: str) -> None:
    """Save data as a CSV file.
    
    Args:
        data (Union[List[Tuple], pd.DataFrame]): Data to save.
        file_path (str): Path where the CSV file will be created.
        
    Raises:
        ValueError: If data is empty or file_path is invalid.
        RuntimeError: If saving fails.
    """
    if not file_path:
        raise ValueError("File path cannot be empty")
    
    try:
        # Convert data to DataFrame if needed
        if isinstance(data, list):
            if len(data) == 0:
                raise ValueError("Data cannot be empty")
            df = pd.DataFrame(data, columns=["id", "name", "age", "salary", "department", "join_date"])
        else:
            if len(data) == 0:
                raise ValueError("DataFrame cannot be empty")
            df = data
        
        # Write to CSV
        df.to_csv(file_path, index=False)
        logger.info(f"Successfully saved data as CSV at {file_path}")
    except Exception as e:
        logger.error(f"Failed to save data as CSV: {str(e)}")
        raise RuntimeError(f"Failed to save data as CSV: {str(e)}")

def main() -> None:
    """Main function to generate and save sample data.
    
    This function:
    1. Parses command line arguments
    2. Generates sample data
    3. Saves data in specified format
    4. Cleans up temporary files
    """
    parser = argparse.ArgumentParser(description="Nessi Sample Data Generator")
    parser.add_argument("--format", choices=["delta", "parquet", "csv", "all"], 
                      default="all", help="Output format (default: all)")
    parser.add_argument("--rows", type=int, default=1000, 
                      help="Number of rows to generate (default: 1000)")
    parser.add_argument("--output-dir", default="backend/data/test_tables",
                      help="Output directory (default: backend/data/test_tables)")
    args = parser.parse_args()

    try:
        # Generate sample data
        data = generate_sample_data(args.rows)
        
        # Create output directory if it doesn't exist
        os.makedirs(args.output_dir, exist_ok=True)
        
        # Save in specified format(s)
        if args.format in ["delta", "all"]:
            save_as_delta(data, os.path.join(args.output_dir, "sample_delta_table"))
        
        if args.format in ["parquet", "all"]:
            save_as_parquet(data, os.path.join(args.output_dir, "sample_parquet_table"))
        
        if args.format in ["csv", "all"]:
            save_as_csv(data, os.path.join(args.output_dir, "sample_csv_table.csv"))
            
        logger.info("Sample data generation completed successfully")
    except Exception as e:
        logger.error(f"Sample data generation failed: {str(e)}")
        raise

if __name__ == "__main__":
    main() 