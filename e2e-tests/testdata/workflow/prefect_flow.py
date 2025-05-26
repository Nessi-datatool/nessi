from prefect import flow, task
from prefect.task_runners import SequentialTaskRunner
import subprocess
import json
import os

@task
def scan_table(table_name, output_file):
    """Scan a Delta Lake table and save results to a file."""
    result = subprocess.run(
        ["nessi", "scan", table_name, "--output-format=json", f"--output-file={output_file}"],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to scan table {table_name}: {result.stderr}")
    return output_file

@task
def validate_schema(table_name, schema_file, output_file):
    """Validate a table's schema against a schema file."""
    result = subprocess.run(
        ["nessi", "validate", "schema", table_name, f"--schema-file={schema_file}", 
         "--output-format=json", f"--output-file={output_file}"],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to validate schema for {table_name}: {result.stderr}")
    return output_file

@task
def apply_rules(table_name, rules_file, output_file):
    """Apply validation rules to a table."""
    result = subprocess.run(
        ["nessi", "validate", "rules", table_name, f"--rules-file={rules_file}", 
         "--output-format=json", f"--output-file={output_file}"],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to apply rules to {table_name}: {result.stderr}")
    return output_file

@task
def generate_report(input_files, output_file, format="html"):
    """Generate a quality report from input files."""
    input_files_str = ",".join(input_files)
    result = subprocess.run(
        ["nessi", "report", "generate", f"--input-files={input_files_str}", 
         f"--output-format={format}", f"--output-file={output_file}"],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to generate report: {result.stderr}")
    return output_file

@flow(name="Nessi Data Quality Flow", task_runner=SequentialTaskRunner())
def data_quality_flow(table_name, schema_file, rules_file, output_dir="/tmp"):
    """Prefect flow for Nessi data quality checks."""
    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)
    
    # Define output file paths
    scan_output = f"{output_dir}/{table_name}_scan_results.json"
    schema_output = f"{output_dir}/{table_name}_schema_validation.json"
    rules_output = f"{output_dir}/{table_name}_rule_validation.json"
    report_output = f"{output_dir}/{table_name}_quality_report.html"
    
    # Execute tasks in sequence
    scan_result = scan_table(table_name, scan_output)
    schema_result = validate_schema(table_name, schema_file, schema_output)
    rules_result = apply_rules(table_name, rules_file, rules_output)
    
    # Generate final report
    report = generate_report(
        [scan_result, schema_result, rules_result],
        report_output
    )
    
    return {
        "scan_output": scan_result,
        "schema_output": schema_result,
        "rules_output": rules_result,
        "report_output": report
    }

if __name__ == "__main__":
    data_quality_flow(
        table_name="users",
        schema_file="/path/to/schemas/user_schema_v1.json",
        rules_file="/path/to/configs/validation_rules.json"
    )
