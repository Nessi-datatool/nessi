"""
Example Prefect flow demonstrating Nessi.dev integration.

This example shows how to use Nessi.dev tasks and flows to perform data quality checks,
validation, and lineage tracking as part of a Prefect data pipeline.
"""

import os
import pandas as pd
from datetime import datetime, timedelta

from prefect import flow, task, get_run_logger
from prefect.task_runners import SequentialTaskRunner

from nessi_prefect.flows.data_pipeline_flow import data_pipeline_flow
from nessi_prefect.flows.data_quality_flow import data_quality_flow


@task(name="Generate Sample Data")
def generate_sample_data(output_path: str, rows: int = 1000) -> str:
    """
    Generate sample data for testing.

    Args:
        output_path: Path to save the data
        rows: Number of rows to generate

    Returns:
        Path to the generated data
    """
    logger = get_run_logger()
    logger.info(f"Generating {rows} rows of sample data to {output_path}")
    
    # Generate sample data
    data = {
        "customer_id": [f"CUST-{i:05d}" for i in range(rows)],
        "transaction_date": [(datetime.now() - timedelta(days=i % 30)).strftime("%Y-%m-%d") for i in range(rows)],
        "product_code": [f"PR-{i % 100:04d}" for i in range(rows)],
        "sales_amount": [round(100 + (i % 100) * 1.5, 2) for i in range(rows)],
        "region": [["North", "South", "East", "West"][i % 4] for i in range(rows)],
    }
    
    # Create DataFrame
    df = pd.DataFrame(data)
    
    # Save to CSV
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    df.to_csv(output_path, index=False)
    
    logger.info(f"Sample data generated and saved to {output_path}")
    return output_path


@task(name="Process Sample Data")
def process_sample_data(input_path: str, output_path: str, quality_results: dict) -> str:
    """
    Process sample data based on quality results.

    Args:
        input_path: Path to the input data
        output_path: Path to save the processed data
        quality_results: Quality check results

    Returns:
        Path to the processed data
    """
    logger = get_run_logger()
    logger.info(f"Processing data from {input_path} to {output_path}")
    
    # Get quality score
    quality_score = quality_results.get("quality_score", 0)
    logger.info(f"Input data quality score: {quality_score:.2f}")
    
    # Load data
    df = pd.read_csv(input_path)
    
    # Apply transformations
    # 1. Filter out negative sales amounts
    df = df[df["sales_amount"] > 0]
    
    # 2. Standardize region names
    df["region"] = df["region"].str.upper()
    
    # 3. Add a new column for transaction month
    df["transaction_month"] = pd.to_datetime(df["transaction_date"]).dt.strftime("%Y-%m")
    
    # 4. Add a new column for sales category
    df["sales_category"] = pd.cut(
        df["sales_amount"],
        bins=[0, 50, 100, 200, float("inf")],
        labels=["Low", "Medium", "High", "Premium"]
    )
    
    # Save processed data
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    df.to_csv(output_path, index=False)
    
    logger.info(f"Processed data saved to {output_path}")
    return output_path


@flow(name="Sample Data Pipeline", task_runner=SequentialTaskRunner())
def sample_data_pipeline():
    """
    Sample data pipeline with Nessi.dev integration.
    
    This flow demonstrates how to use Nessi.dev tasks and flows to perform
    data quality checks, validation, and lineage tracking as part of a 
    Prefect data pipeline.
    """
    logger = get_run_logger()
    logger.info("Starting sample data pipeline with Nessi.dev integration")
    
    # Configuration
    api_host = os.environ.get("NESSI_API_HOST", "http://localhost:8080")
    api_key = os.environ.get("NESSI_API_KEY", "your-api-key")
    
    # Generate sample data
    input_path = generate_sample_data(
        output_path="./data/raw/sales_data.csv",
        rows=1000
    )
    
    # Run quality check on input data
    quality_results = data_quality_flow(
        table_name="sales_data",
        rules=[
            {
                "name": "sales_amount_positive",
                "description": "Sales amount should be positive",
                "rule_type": "range",
                "column": "sales_amount",
                "min_value": 0,
            },
            {
                "name": "transaction_date_valid",
                "description": "Transaction date should be valid",
                "rule_type": "not_null",
                "column": "transaction_date",
            },
            {
                "name": "customer_id_unique",
                "description": "Customer ID should be unique",
                "rule_type": "unique",
                "column": "customer_id",
            },
        ],
        quality_threshold=0.8,
        api_host=api_host,
        api_key=api_key,
    )
    
    # Process data
    output_path = process_sample_data(
        input_path=input_path,
        output_path="./data/processed/sales_data_processed.csv",
        quality_results=quality_results,
    )
    
    # Run full data pipeline with validation
    pipeline_results = data_pipeline_flow(
        input_table="sales_data",
        output_table="sales_data_processed",
        validation_rules=[
            {
                "name": "region_in_allowed_values",
                "description": "Region should be in allowed values",
                "rule_type": "enum",
                "column": "region",
                "allowed_values": ["NORTH", "SOUTH", "EAST", "WEST"],
            },
            {
                "name": "sales_category_not_null",
                "description": "Sales category should not be null",
                "rule_type": "not_null",
                "column": "sales_category",
            },
        ],
        api_host=api_host,
        api_key=api_key,
    )
    
    logger.info("Sample data pipeline completed successfully")
    return {
        "input_path": input_path,
        "output_path": output_path,
        "quality_results": quality_results,
        "pipeline_results": pipeline_results,
    }


if __name__ == "__main__":
    sample_data_pipeline()
