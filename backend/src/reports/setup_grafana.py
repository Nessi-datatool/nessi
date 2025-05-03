"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

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
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Script to set up Grafana with example dashboard and data sources."""

import requests
import json
import time
from pathlib import Path
from typing import Dict, Any
import os
from requests.exceptions import ConnectionError

class GrafanaSetup:
    """Sets up Grafana with example dashboard and data sources."""
    
    def __init__(self, base_url: str = "http://localhost:3000", api_key: str = None):
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {"Content-Type": "application/json"}
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
        self.base_dir = Path(__file__).parent
        self.dashboard_path = self.base_dir / "grafana_dashboard.json"
        
    def wait_for_grafana(self, timeout: int = 60) -> bool:
        """Wait for Grafana to be ready."""
        print("Waiting for Grafana to be ready...")
        start_time = time.time()
        retry_count = 0
        
        while time.time() - start_time < timeout:
            try:
                response = requests.get(f"{self.base_url}/api/health")
                if response.status_code == 200:
                    print("Grafana is ready!")
                    return True
            except ConnectionError:
                retry_count += 1
                print(f"Connection error, retrying... (attempt {retry_count})")
                time.sleep(0.1)
                continue
            except Exception as e:
                retry_count += 1
                print(f"Error checking Grafana health: {e} (attempt {retry_count})")
                time.sleep(0.1)
                continue
        
        print("Timeout waiting for Grafana to be ready")
        return False
        
    def setup_prometheus_datasource(self) -> bool:
        """Set up Prometheus data source."""
        try:
            print("Setting up Prometheus data source...")
            datasource = {
                "name": "Prometheus",
                "type": "prometheus",
                "url": "http://prometheus:9090",
                "access": "proxy",
                "isDefault": True
            }
            
            response = requests.post(
                f"{self.base_url}/api/datasources",
                json=datasource,
                headers=self.headers
            )
            
            if response.status_code in [200, 409]:  # 409 means already exists
                print("Prometheus data source configured")
                return True
            else:
                print(f"Failed to configure Prometheus: {response.text}")
                return False
                
        except Exception as e:
            print(f"Error setting up Prometheus: {str(e)}")
            return False
            
    def setup_dashboard(self) -> bool:
        """Set up example dashboard."""
        try:
            print("Setting up example dashboard...")
            
            # Read dashboard JSON
            with open(self.dashboard_path, 'r') as f:
                dashboard = json.load(f)
                
            # Set dashboard properties
            dashboard["dashboard"]["id"] = None
            dashboard["dashboard"]["uid"] = "nessi-reports"
            dashboard["dashboard"]["title"] = "Nessi Reports Dashboard"
            dashboard["overwrite"] = True
            
            # Create dashboard
            response = requests.post(
                f"{self.base_url}/api/dashboards/db",
                json=dashboard,
                headers=self.headers
            )
            
            if response.status_code == 200:
                print("Dashboard created successfully")
                return True
            else:
                print(f"Failed to create dashboard: {response.text}")
                return False
                
        except Exception as e:
            print(f"Error setting up dashboard: {str(e)}")
            return False
            
    def setup_alert_channels(self) -> bool:
        """Set up example alert notification channels."""
        try:
            print("Setting up alert notification channels...")
            
            # Example email channel
            email_channel = {
                "name": "Email Alerts",
                "type": "email",
                "settings": {
                    "addresses": "admin@example.com"
                }
            }
            
            response = requests.post(
                f"{self.base_url}/api/alert-notifications",
                json=email_channel,
                headers=self.headers
            )
            
            if response.status_code in [200, 409]:  # 409 means already exists
                print("Alert notification channels configured")
                return True
            else:
                print(f"Failed to configure alert channels: {response.text}")
                return False
                
        except Exception as e:
            print(f"Error setting up alert channels: {str(e)}")
            return False
            
    def run(self) -> bool:
        """Run the complete setup process."""
        if not self.wait_for_grafana():
            return False
            
        if not self.setup_prometheus_datasource():
            return False
            
        if not self.setup_dashboard():
            return False
            
        if not self.setup_alert_channels():
            return False
            
        print("\nGrafana setup completed successfully!")
        print(f"Dashboard URL: {self.base_url}/d/nessi-reports")
        return True
        
def main():
    """Run the Grafana setup."""
    setup = GrafanaSetup()
    setup.run()
    
if __name__ == "__main__":
    main()