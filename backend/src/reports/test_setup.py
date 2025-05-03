"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

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
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Comprehensive test script for the monitoring setup."""

import subprocess
import time
import requests
import json
from pathlib import Path
from typing import Dict, Any, List

class MonitoringTest:
    """Tests the entire monitoring setup."""
    
    def __init__(self):
        self.base_dir = Path(__file__).parent.parent.parent
        self.dashboard_path = self.base_dir / "src/reports/grafana_dashboard.json"
        self.docker_compose_path = self.base_dir / "docker-compose.test.yml"
        
    def start_services(self) -> bool:
        """Start Docker services for testing."""
        try:
            print("Starting Docker services...")
            subprocess.run([
                "docker-compose", "-f", str(self.docker_compose_path), "up", "-d"
            ], check=True)
            print("Services started successfully")
            return True
        except subprocess.CalledProcessError as e:
            print(f"Error starting services: {str(e)}")
            return False
            
    def wait_for_services(self, timeout: int = 60) -> bool:
        """Wait for services to be ready."""
        print("Waiting for services to be ready...")
        start_time = time.time()
        services_ready = False
        
        while time.time() - start_time < timeout:
            try:
                # Check metrics endpoint
                metrics_response = requests.get("http://localhost:8000/metrics")
                # Check Prometheus
                prometheus_response = requests.get("http://localhost:9090/-/ready")
                # Check Grafana
                grafana_response = requests.get("http://localhost:3000/api/health")
                
                if (metrics_response.status_code == 200 and
                    prometheus_response.status_code == 200 and
                    grafana_response.status_code == 200):
                    services_ready = True
                    break
                    
            except requests.exceptions.ConnectionError:
                pass
                
            time.sleep(5)
            
        if services_ready:
            print("All services are ready")
            return True
        else:
            print("Timeout waiting for services")
            return False
            
    def test_metrics_generation(self) -> bool:
        """Test if metrics are being generated."""
        try:
            print("Testing metrics generation...")
            response = requests.get("http://localhost:8000/metrics")
            if response.status_code == 200:
                metrics = response.text
                required_metrics = [
                    "nessi_table_scan_total_rows",
                    "nessi_data_quality_score",
                    "nessi_processing_time_seconds",
                    "nessi_custom_report_count"
                ]
                if all(metric in metrics for metric in required_metrics):
                    print("All required metrics are being generated")
                    return True
            print("Some metrics are missing")
            return False
        except Exception as e:
            print(f"Error testing metrics: {str(e)}")
            return False
            
    def test_prometheus_scraping(self) -> bool:
        """Test if Prometheus is scraping metrics."""
        try:
            print("Testing Prometheus scraping...")
            response = requests.get(
                "http://localhost:9090/api/v1/query",
                params={"query": "nessi_table_scan_total_rows"}
            )
            if response.status_code == 200:
                data = response.json()
                if data["status"] == "success" and len(data["data"]["result"]) > 0:
                    print("Prometheus is successfully scraping metrics")
                    return True
            print("Prometheus is not scraping metrics correctly")
            return False
        except Exception as e:
            print(f"Error testing Prometheus: {str(e)}")
            return False
            
    def test_grafana_dashboard(self) -> bool:
        """Test if Grafana dashboard is working."""
        try:
            print("Testing Grafana dashboard...")
            # Login to Grafana
            session = requests.Session()
            login_response = session.post(
                "http://localhost:3000/api/login",
                json={"user": "admin", "password": "admin"}
            )
            if login_response.status_code != 200:
                print("Failed to login to Grafana")
                return False
                
            # Check dashboard
            dashboard_response = session.get(
                "http://localhost:3000/api/dashboards/uid/nessi-reports"
            )
            if dashboard_response.status_code == 200:
                print("Dashboard is accessible")
                return True
            print("Dashboard is not accessible")
            return False
        except Exception as e:
            print(f"Error testing Grafana: {str(e)}")
            return False
            
    def test_alerts(self) -> bool:
        """Test if alerts are configured correctly."""
        try:
            print("Testing alert configurations...")
            session = requests.Session()
            session.post(
                "http://localhost:3000/api/login",
                json={"user": "admin", "password": "admin"}
            )
            
            # Check alert rules
            alerts_response = session.get("http://localhost:3000/api/alert-notifications")
            if alerts_response.status_code == 200:
                alerts = alerts_response.json()
                if len(alerts) > 0:
                    print("Alerts are configured")
                    return True
            print("No alerts are configured")
            return False
        except Exception as e:
            print(f"Error testing alerts: {str(e)}")
            return False
            
    def cleanup(self) -> None:
        """Clean up Docker services."""
        try:
            print("Cleaning up services...")
            subprocess.run([
                "docker-compose", "-f", str(self.docker_compose_path), "down"
            ], check=True)
            print("Services cleaned up successfully")
        except subprocess.CalledProcessError as e:
            print(f"Error cleaning up services: {str(e)}")
            
    def run_tests(self) -> Dict[str, bool]:
        """Run all tests."""
        results = {}
        
        try:
            # Start services
            if not self.start_services():
                return {"services_start": False}
            results["services_start"] = True
            
            # Wait for services
            if not self.wait_for_services():
                return results
            results["services_ready"] = True
            
            # Run individual tests
            results["metrics_generation"] = self.test_metrics_generation()
            results["prometheus_scraping"] = self.test_prometheus_scraping()
            results["grafana_dashboard"] = self.test_grafana_dashboard()
            results["alerts"] = self.test_alerts()
            
            return results
            
        finally:
            self.cleanup()
            
def main():
    """Run the test suite."""
    test = MonitoringTest()
    print("Starting comprehensive monitoring test...")
    results = test.run_tests()
    
    print("\nTest Results:")
    for test_name, passed in results.items():
        print(f"{test_name}: {'✓' if passed else '✗'}")
        
    if all(results.values()):
        print("\nAll tests passed! Monitoring system is working correctly.")
    else:
        print("\nSome tests failed. Please check the logs above for details.")
        
if __name__ == "__main__":
    main()