#!/usr/bin/env python3
import argparse
import sys
import logging
from pathlib import Path
import json
from typing import Dict, Any, Optional

# Add the backend directory to the Python path
sys.path.append(str(Path(__file__).parent.parent))

from src.scanner import Scanner
from src.monitoring.monitoring_service import MonitoringService

def setup_logging():
    """Set up logging configuration."""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )

def load_config(config_path: str) -> Dict[str, Any]:
    """Load configuration from a JSON file."""
    try:
        with open(config_path, 'r') as f:
            return json.load(f)
    except Exception as e:
        logging.error(f"Error loading configuration: {str(e)}")
        sys.exit(1)

def main():
    """Main entry point for the monitoring script."""
    parser = argparse.ArgumentParser(description='Nessi Data Quality Monitor')
    subparsers = parser.add_subparsers(dest='command', help='Command to execute')
    
    # Start command
    start_parser = subparsers.add_parser('start', help='Start the monitoring service')
    start_parser.add_argument('--config', type=str, required=True, help='Path to configuration file')
    start_parser.add_argument('--port', type=int, default=8000, help='Port for metrics server')
    start_parser.add_argument('--interval', type=int, default=60, help='Monitoring interval in seconds')
    
    # Stop command
    stop_parser = subparsers.add_parser('stop', help='Stop the monitoring service')
    stop_parser.add_argument('--config', type=str, required=True, help='Path to configuration file')
    
    # Dashboard command
    dashboard_parser = subparsers.add_parser('dashboard', help='Generate Grafana dashboard')
    dashboard_parser.add_argument('--config', type=str, required=True, help='Path to configuration file')
    dashboard_parser.add_argument('--output', type=str, required=True, help='Output path for dashboard configuration')
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        sys.exit(1)
    
    setup_logging()
    logger = logging.getLogger(__name__)
    
    try:
        config = load_config(args.config)
        scanner = Scanner(config)
        monitoring_service = MonitoringService(scanner, port=args.port)
        
        if args.command == 'start':
            monitoring_service.start(interval=args.interval)
            logger.info("Monitoring service started")
            # Keep the script running
            try:
                while True:
                    pass
            except KeyboardInterrupt:
                monitoring_service.stop()
                logger.info("Monitoring service stopped")
        
        elif args.command == 'stop':
            monitoring_service.stop()
            logger.info("Monitoring service stopped")
        
        elif args.command == 'dashboard':
            monitoring_service.generate_dashboard(args.output)
            logger.info(f"Generated dashboard configuration at {args.output}")
    
    except Exception as e:
        logger.error(f"Error: {str(e)}")
        sys.exit(1)

if __name__ == '__main__':
    main() 