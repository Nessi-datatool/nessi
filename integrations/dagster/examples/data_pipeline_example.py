"""
Example Dagster pipeline demonstrating Nessi.dev integration.

This example shows how to use Nessi.dev ops and jobs to perform data quality checks,
validation, and lineage tracking as part of a Dagster pipeline.
"""

import os
import pandas as pd
from datetime import datetime, timedelta

from dagster import (
    job, op, In, Out, resource, ResourceDefinition, 
    DagsterType, TypeCheck, Field, String, Int, Bool, 
    make_values_resource, file_relative_path
)

from nessi_dagster.resources import nessi_resource
from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
from nessi_dagster.ops.validation_ops import run_validation, wait_for_validation_results
from nessi_dagster.ops.lineage_ops import get_lineage


# Define custom types
def is_dataframe_type_check(_, value):
    return isinstance(value, pd.DataFrame)

DataFrame = DagsterType(
    name="DataFrame",
    description="A pandas DataFrame",
    type_check_fn=is_dataframe_type_check,
)


# Define ops for the example pipeline
@op(
    name="generate_sample_data",
    description="Generate sample data for testing",
    ins={
        "rows": In(int, description="Number of rows to generate", default_value=1000),
    },
    out=Out(DataFrame, description="Generated sample data"),
)
def generate_sample_data(context, rows: int = 1000) -> pd.DataFrame:
    """
    Generate sample data for testing.

    Args:
        context: Dagster execution context
        rows: Number of rows to generate

    Returns:
        Sample data as DataFrame
    """
    context.log.info(f"Generating {rows} rows of sample data")
    
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
    
    context.log.info(f"Sample data generated with {len(df)} rows and {len(df.columns)} columns")
    return df


@op(
    name="save_to_table",
    description="Save DataFrame to a table",
    ins={
        "df": In(DataFrame, description="DataFrame to save"),
        "table_name": In(str, description="Name of the table"),
    },
    out=Out(str, description="Name of the saved table"),
)
def save_to_table(context, df: pd.DataFrame, table_name: str) -> str:
    """
    Save DataFrame to a table.

    Args:
        context: Dagster execution context
        df: DataFrame to save
        table_name: Name of the table

    Returns:
        Name of the saved table
    """
    context.log.info(f"Saving DataFrame to table: {table_name}")
    
    # In a real implementation, this would save to a database
    # For now, just log the action
    context.log.info(f"Saved {len(df)} rows to table {table_name}")
    
    return table_name


@op(
    name="process_sample_data",
    description="Process sample data based on quality results",
    ins={
        "df": In(DataFrame, description="DataFrame to process"),
        "quality_results": In(dict, description="Quality check results"),
    },
    out=Out(DataFrame, description="Processed DataFrame"),
)
def process_sample_data(context, df: pd.DataFrame, quality_results: dict) -> pd.DataFrame:
    """
    Process sample data based on quality results.

    Args:
        context: Dagster execution context
        df: DataFrame to process
        quality_results: Quality check results

    Returns:
        Processed DataFrame
    """
    context.log.info(f"Processing DataFrame with {len(df)} rows")
    
    # Get quality score
    quality_score = quality_results.get("quality_score", 0)
    context.log.info(f"Input data quality score: {quality_score:.2f}")
    
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
    
    context.log.info(f"Processed DataFrame has {len(df)} rows and {len(df.columns)} columns")
    return df


# Define the pipeline
@job(
    name="nessi_sample_pipeline",
    description="Sample data pipeline with Nessi.dev integration",
    resource_defs={
        "nessi": nessi_resource.configured({
            "api_host": "http://localhost:8080",
            "api_key": "your-api-key",
            "api_secret": "",
            "timeout": 300,
        }),
    },
)
def sample_pipeline():
    """
    Sample data pipeline with Nessi.dev integration.
    
    This pipeline demonstrates how to use Nessi.dev ops to perform
    data quality checks, validation, and lineage tracking as part of a 
    Dagster pipeline.
    """
    # Generate sample data
    df = generate_sample_data()
    
    # Save to input table
    input_table = save_to_table(df=df, table_name="sales_data")
    
    # Run quality check on input table
    quality_response = run_quality_check(
        table_name=input_table,
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
        profile=True,
    )
    
    # Wait for quality results
    quality_results = wait_for_quality_results(
        check_id=quality_response["check_id"],
        quality_threshold=0.8,
        fail_on_rule_failure=True,
    )
    
    # Process data
    processed_df = process_sample_data(df=df, quality_results=quality_results)
    
    # Save to output table
    output_table = save_to_table(df=processed_df, table_name="sales_data_processed")
    
    # Run validation on output table
    validation_response = run_validation(
        table_name=output_table,
        rules=[
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
    )
    
    # Wait for validation results
    validation_results = wait_for_validation_results(
        validation_id=validation_response["validation_id"],
        fail_on_validation_failure=True,
    )
    
    # Get lineage for output table
    lineage = get_lineage(table_name=output_table)


if __name__ == "__main__":
    result = sample_pipeline.execute_in_process()
