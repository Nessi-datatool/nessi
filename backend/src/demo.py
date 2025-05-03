"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

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
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""
Demo script to showcase Nessi features.
Generates sample data, scans tables, and creates a Grafana dashboard.
"""

from pathlib import Path
import time
import os
import requests
from src.scanner.sample_data_generator import SampleDataGenerator
from src.scanner.table_scanner import TableScanner
from src.scanner.grafana_dashboard import GrafanaDashboard
from src.scanner.spark_config import get_spark_session, stop_spark_session

def wait_for_grafana(url: str, max_retries: int = 30, delay: int = 2) -> bool:
    """Wait for Grafana to be ready."""
    for _ in range(max_retries):
        try:
            response = requests.get(f"{url}/api/health")
            if response.status_code == 200:
                return True
        except requests.exceptions.ConnectionError:
            pass
        time.sleep(delay)
    return False

def main():
    # Create demo directory
    demo_dir = Path("data/demo")
    demo_dir.mkdir(parents=True, exist_ok=True)

    # Initialize Spark session
    spark = get_spark_session("NessiDemo")

    try:
        # Generate sample data
        print("Generating sample data...")
        generator = SampleDataGenerator()
        df = generator.generate_sample_data()
        
        # Save in different formats
        parquet_path = str(demo_dir / "sample.parquet")
        csv_path = str(demo_dir / "sample.csv")
        
        print(f"Saving data to {parquet_path} and {csv_path}...")
        generator.save_as_parquet(df, parquet_path)
        generator.save_as_csv(df, csv_path)

        # Scan tables
        print("\nScanning tables...")
        scanner = TableScanner(spark)
        parquet_results = scanner.scan_table(parquet_path, "parquet")
        csv_results = scanner.scan_table(csv_path, "csv")

        print("\nParquet table statistics:")
        print(f"- Row count: {parquet_results['row_count']}")
        print(f"- Columns: {', '.join(field['name'] for field in parquet_results['schema'])}")
        print("\nColumn statistics:")
        for col, stats in parquet_results['column_stats'].items():
            print(f"- {col}:")
            for stat_name, stat_value in stats.items():
                print(f"  - {stat_name}: {stat_value}")

        # Create Grafana dashboard
        print("\nSetting up Grafana dashboard...")
        grafana_url = os.getenv("GRAFANA_URL", "http://localhost:3000")
        grafana_admin_password = os.getenv("GF_SECURITY_ADMIN_PASSWORD", "admin")
        
        # Wait for Grafana to be ready
        print("Waiting for Grafana to be ready...")
        if not wait_for_grafana(grafana_url):
            print("Warning: Grafana is not responding. Skipping dashboard creation.")
            return

        try:
            # Create dashboard
            print("Creating Grafana dashboard...")
            dashboard = GrafanaDashboard(grafana_url, grafana_admin_password)
            dashboard_response = dashboard.create_dashboard("Nessi Demo Dashboard")
            dashboard_uid = dashboard_response["dashboard"]["uid"]

            # Add panels
            print("Adding dashboard panels...")
            dashboard.add_table_metrics_panel(dashboard_uid, parquet_results)
            dashboard.add_column_stats_panel(dashboard_uid, parquet_results["column_stats"])

            # Export dashboard
            dashboard_path = str(demo_dir / "dashboard.json")
            print(f"\nExporting dashboard to {dashboard_path}...")
            dashboard.export_dashboard(dashboard_uid, dashboard_path)

            print("\nDemo completed successfully!")
            print("Generated files:")
            print(f"1. Sample Parquet file: {parquet_path}")
            print(f"2. Sample CSV file: {csv_path}")
            print(f"3. Grafana dashboard: {dashboard_path}")
            print(f"\nYou can now access the dashboard at: {grafana_url}")
            print("Login credentials:")
            print("Username: admin")
            print(f"Password: {grafana_admin_password}")

        except Exception as e:
            print(f"\nError creating Grafana dashboard: {str(e)}")
            print("Continuing with data generation and scanning results...")
    finally:
        # Stop Spark session
        stop_spark_session()

if __name__ == "__main__":
    main()