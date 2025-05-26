from dagster import job, op, In, Out, Nothing
import subprocess
import json
import os

@op(out=Out(str))
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

@op(out=Out(str))
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

@op(out=Out(str))
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

@op(ins={
    "scan_result": In(str),
    "schema_result": In(str),
    "rules_result": In(str)
})
def generate_report(context, scan_result, schema_result, rules_result):
    """Generate a quality report from input files."""
    output_file = f"/tmp/{context.op_config['table_name']}_quality_report.html"
    input_files_str = f"{scan_result},{schema_result},{rules_result}"
    
    result = subprocess.run(
        ["nessi", "report", "generate", f"--input-files={input_files_str}", 
         "--output-format=html", f"--output-file={output_file}"],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to generate report: {result.stderr}")
    
    context.log.info(f"Report generated at {output_file}")
    return output_file

@job
def nessi_data_quality_job():
    """Dagster job for Nessi data quality checks."""
    # Configuration will be provided at runtime
    table_name = "users"
    schema_file = "/path/to/schemas/user_schema_v1.json"
    rules_file = "/path/to/configs/validation_rules.json"
    
    # Define output file paths
    output_dir = "/tmp"
    scan_output = f"{output_dir}/{table_name}_scan_results.json"
    schema_output = f"{output_dir}/{table_name}_schema_validation.json"
    rules_output = f"{output_dir}/{table_name}_rule_validation.json"
    
    # Execute ops in sequence
    scan_result = scan_table(table_name, scan_output)
    schema_result = validate_schema(table_name, schema_file, schema_output)
    rules_result = apply_rules(table_name, rules_file, rules_output)
    
    # Generate final report
    generate_report(scan_result, schema_result, rules_result)

if __name__ == "__main__":
    result = nessi_data_quality_job.execute_in_process()
