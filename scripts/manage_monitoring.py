#!/usr/bin/env python3

import argparse
import subprocess
import sys
import os
import time
import logging
from pathlib import Path
import json
import requests

logger = logging.getLogger(__name__)

def setup_logging():
    """Set up logging configuration."""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )

def check_docker_compose():
    """Check if Docker Compose is installed."""
    try:
        subprocess.run(['docker-compose', '--version'], check=True, capture_output=True)
        return True
    except subprocess.CalledProcessError:
        logger.error("Docker Compose is not installed")
        return False

def start_monitoring(config_path: str = None):
    """Start the monitoring stack."""
    try:
        cmd = ['docker-compose', '-f', 'docker-compose.monitoring.yml', 'up', '-d']
        subprocess.run(cmd, check=True)
        logger.info("Monitoring stack started successfully")
        
        # Wait for services to be ready
        time.sleep(10)
        
        # Check service health
        check_services_health()
        
    except subprocess.CalledProcessError as e:
        logger.error(f"Failed to start monitoring stack: {str(e)}")
        sys.exit(1)

def stop_monitoring():
    """Stop the monitoring stack."""
    try:
        cmd = ['docker-compose', '-f', 'docker-compose.monitoring.yml', 'down']
        subprocess.run(cmd, check=True)
        logger.info("Monitoring stack stopped successfully")
    except subprocess.CalledProcessError as e:
        logger.error(f"Failed to stop monitoring stack: {str(e)}")
        sys.exit(1)

def check_services_health():
    """Check the health of monitoring services."""
    services = {
        'Prometheus': 'http://localhost:9090/-/healthy',
        'Grafana': 'http://localhost:3000/api/health',
        'Alertmanager': 'http://localhost:9093/-/healthy'
    }
    
    for service, url in services.items():
        try:
            response = requests.get(url)
            if response.status_code == 200:
                logger.info(f"{service} is healthy")
            else:
                logger.warning(f"{service} health check failed: {response.status_code}")
        except requests.exceptions.RequestException as e:
            logger.error(f"Failed to check {service} health: {str(e)}")

def create_grafana_datasource():
    """Create Prometheus datasource in Grafana."""
    url = 'http://localhost:3000/api/datasources'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Basic YWRtaW46bmVzc2kxMjM='  # admin:nessi123
    }
    data = {
        'name': 'Prometheus',
        'type': 'prometheus',
        'url': 'http://prometheus:9090',
        'access': 'proxy',
        'isDefault': True
    }
    
    try:
        response = requests.post(url, headers=headers, json=data)
        if response.status_code in [200, 409]:  # 409 means datasource already exists
            logger.info("Prometheus datasource configured in Grafana")
        else:
            logger.warning(f"Failed to create Grafana datasource: {response.status_code}")
    except requests.exceptions.RequestException as e:
        logger.error(f"Failed to configure Grafana datasource: {str(e)}")

def import_dashboards():
    """Import dashboards into Grafana."""
    dashboard_dir = Path('grafana/dashboards')
    if not dashboard_dir.exists():
        logger.warning("Dashboard directory not found")
        return
    
    url = 'http://localhost:3000/api/dashboards/db'
    headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Basic YWRtaW46bmVzc2kxMjM='  # admin:nessi123
    }
    
    for dashboard_file in dashboard_dir.glob('*.json'):
        try:
            with open(dashboard_file, 'r') as f:
                dashboard = json.load(f)
            
            data = {
                'dashboard': dashboard,
                'overwrite': True
            }
            
            response = requests.post(url, headers=headers, json=data)
            if response.status_code == 200:
                logger.info(f"Imported dashboard: {dashboard_file.name}")
            else:
                logger.warning(f"Failed to import dashboard {dashboard_file.name}: {response.status_code}")
        except Exception as e:
            logger.error(f"Error importing dashboard {dashboard_file.name}: {str(e)}")

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description='Manage Nessi monitoring services')
    parser.add_argument('action', choices=['start', 'stop', 'status'], help='Action to perform')
    parser.add_argument('--config', help='Path to configuration file')
    
    args = parser.parse_args()
    
    setup_logging()
    
    if not check_docker_compose():
        sys.exit(1)
    
    if args.action == 'start':
        start_monitoring(args.config)
        create_grafana_datasource()
        import_dashboards()
    elif args.action == 'stop':
        stop_monitoring()
    elif args.action == 'status':
        check_services_health()

if __name__ == '__main__':
    main() 