"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

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
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Script to verify dashboard and metrics functionality."""

import time
import requests
import json
from typing import Dict, Any

class DashboardVerifier:
    """Verifies dashboard and metrics functionality."""
    
    def __init__(self):
        self.grafana_url = "http://localhost:3000"
        self.prometheus_url = "http://localhost:9090"
        self.metrics_url = "http://localhost:8000"
        self.auth = ("admin", "admin")
        
    def verify_metrics_endpoint(self) -> bool:
        """Verify metrics endpoint is accessible and returning data."""
        try:
            response = requests.get(f"{self.metrics_url}/metrics")
            if response.status_code == 200:
                metrics = response.text
                required_metrics = [
                    "nessi_table_scan_total_rows",
                    "nessi_data_quality_score",
                    "nessi_processing_time_seconds",
                    "nessi_custom_report_count"
                ]
                return all(metric in metrics for metric in required_metrics)
            return False
        except Exception as e:
            print(f"Error verifying metrics endpoint: {str(e)}")
            return False
            
    def verify_prometheus(self) -> bool:
        """Verify Prometheus is scraping metrics."""
        try:
            response = requests.get(f"{self.prometheus_url}/api/v1/query", 
                                 params={"query": "nessi_table_scan_total_rows"})
            if response.status_code == 200:
                data = response.json()
                return data["status"] == "success" and len(data["data"]["result"]) > 0
            return False
        except Exception as e:
            print(f"Error verifying Prometheus: {str(e)}")
            return False
            
    def verify_grafana_dashboard(self) -> bool:
        """Verify Grafana dashboard is accessible and configured."""
        try:
            # Verify dashboard exists
            response = requests.get(f"{self.grafana_url}/api/dashboards/uid/nessi-reports",
                                 auth=self.auth)
            if response.status_code != 200:
                return False
                
            # Verify data source is configured
            response = requests.get(f"{self.grafana_url}/api/datasources",
                                 auth=self.auth)
            if response.status_code == 200:
                data_sources = response.json()
                return any(ds["type"] == "prometheus" for ds in data_sources)
            return False
        except Exception as e:
            print(f"Error verifying Grafana dashboard: {str(e)}")
            return False
            
    def verify_alerts(self) -> bool:
        """Verify alerts are configured and working."""
        try:
            response = requests.get(f"{self.grafana_url}/api/alert-notifications",
                                 auth=self.auth)
            if response.status_code == 200:
                notifications = response.json()
                return len(notifications) > 0
            return False
        except Exception as e:
            print(f"Error verifying alerts: {str(e)}")
            return False
            
    def run_verification(self) -> Dict[str, Any]:
        """Run all verification checks."""
        results = {
            "metrics_endpoint": self.verify_metrics_endpoint(),
            "prometheus": self.verify_prometheus(),
            "grafana_dashboard": self.verify_grafana_dashboard(),
            "alerts": self.verify_alerts()
        }
        return results

if __name__ == "__main__":
    verifier = DashboardVerifier()
    print("Starting dashboard verification...")
    results = verifier.run_verification()
    
    print("\nVerification Results:")
    for check, status in results.items():
        print(f"{check}: {'✓' if status else '✗'}")
        
    if all(results.values()):
        print("\nAll checks passed! Dashboard is working correctly.")
    else:
        print("\nSome checks failed. Please check the logs above for details.")