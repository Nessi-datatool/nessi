import click
from typing import Optional
import os
from pathlib import Path
import logging
from ..monitoring.config import MonitoringConfig
from ..monitoring.metrics_collector import MetricsCollector
from ..monitoring.api import MonitoringAPI
from ..security.user_manager import UserManager
from ..security.rbac import Role
from ..transform.pipeline import PipelineManager
import uvicorn
import json

@click.group()
def cli():
    """Nessi CLI - Data Quality and Monitoring Tool"""
    pass

@cli.group()
def table():
    """Manage Delta Lake tables."""
    pass

@table.command()
@click.argument('path')
@click.option('--format', '-f', type=click.Choice(['delta', 'parquet', 'csv']), required=True)
@click.option('--partition-by', '-p', help='Partition columns')
@click.option('--z-order-by', '-z', help='Z-order columns')
def create(path: str, format: str, partition_by: Optional[str] = None, z_order_by: Optional[str] = None):
    """Create a new table with specified format and optimizations."""
    click.echo(f"Creating {format} table at {path}")
    # TODO: Implement table creation with optimizations

@table.command()
@click.argument('path')
@click.option('--version', '-v', help='Version to time travel to')
def time_travel(path: str, version: Optional[str] = None):
    """Time travel to a specific version of the table."""
    click.echo(f"Time traveling to version {version if version else 'latest'} of {path}")
    # TODO: Implement time travel

@cli.group()
def quality():
    """Manage data quality checks."""
    pass

@quality.command()
@click.argument('table_path')
@click.option('--rules', '-r', help='Path to rules file')
@click.option('--output', '-o', help='Output path for report')
def check(table_path: str, rules: Optional[str] = None, output: Optional[str] = None):
    """Run data quality checks on a table."""
    click.echo(f"Running quality checks on {table_path}")
    # TODO: Implement quality checks

@cli.group()
def monitor():
    """Manage monitoring and alerts."""
    pass

@monitor.command()
@click.option('--config', '-c', help='Path to config file')
def start(config: Optional[str] = None):
    """Start monitoring services."""
    click.echo("Starting monitoring services...")
    if config:
        os.environ['NESSI_CONFIG'] = config
    config = MonitoringConfig()
    metrics_collector = MetricsCollector(config)
    api = MonitoringAPI(metrics_collector, config)
    uvicorn.run(api.get_app(), host="0.0.0.0", port=8000)

@monitor.command()
@click.option('--threshold', '-t', type=float, help='Alert threshold')
@click.option('--metric', '-m', help='Metric to monitor')
def alert(threshold: float, metric: str):
    """Configure alert thresholds."""
    click.echo(f"Setting alert threshold for {metric} to {threshold}")
    # TODO: Implement alert configuration

@cli.group()
def user():
    """Manage users and permissions."""
    pass

@user.command()
@click.option('--username', '-u', required=True, help='Username')
@click.option('--email', '-e', required=True, help='Email')
@click.option('--role', '-r', type=click.Choice(['admin', 'manager', 'analyst', 'viewer']), required=True)
def create(username: str, email: str, role: str):
    """Create a new user."""
    user_manager = UserManager()
    try:
        user = user_manager.create_user(username, email, Role(role))
        click.echo(f"Created user {username} with role {role}")
        click.echo(f"API Key: {user.api_key}")
    except ValueError as e:
        click.echo(f"Error: {str(e)}", err=True)

@cli.group()
def pipeline():
    """Manage data transformation pipelines."""
    pass

@pipeline.command()
@click.argument('name')
def create(name: str):
    """Create a new transformation pipeline."""
    pipeline_manager = PipelineManager()
    try:
        pipeline = pipeline_manager.create_pipeline(name)
        click.echo(f"Created pipeline {name}")
    except ValueError as e:
        click.echo(f"Error: {str(e)}", err=True)

@pipeline.command()
@click.argument('name')
@click.option('--transformation', '-t', help='Transformation type')
@click.option('--config', '-c', help='Transformation config (JSON)')
def add_transformation(name: str, transformation: str, config: str):
    """Add a transformation to a pipeline."""
    pipeline_manager = PipelineManager()
    pipeline = pipeline_manager.get_pipeline(name)
    if not pipeline:
        click.echo(f"Error: Pipeline {name} not found", err=True)
        return
    
    try:
        config_dict = json.loads(config)
        # TODO: Add transformation based on type
        click.echo(f"Added {transformation} transformation to pipeline {name}")
    except json.JSONDecodeError:
        click.echo("Error: Invalid JSON config", err=True)

if __name__ == '__main__':
    cli() 