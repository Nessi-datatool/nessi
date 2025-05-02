"""Script to set up Grafana with example dashboard and data sources."""

import requests
import json
import time
from pathlib import Path
from typing import Dict, Any

class GrafanaSetup:
    """Sets up Grafana with example dashboard and data sources."""
    
    def __init__(self):
        self.base_url = "http://localhost:3000"
        self.auth = ("admin", "admin")
        self.base_dir = Path(__file__).parent
        self.dashboard_path = self.base_dir / "grafana_dashboard.json"
        
    def wait_for_grafana(self, timeout: int = 60) -> bool:
        """Wait for Grafana to be ready."""
        print("Waiting for Grafana to be ready...")
        start_time = time.time()
        
        while time.time() - start_time < timeout:
            try:
                response = requests.get(f"{self.base_url}/api/health")
                if response.status_code == 200:
                    print("Grafana is ready!")
                    return True
            except requests.exceptions.ConnectionError:
                pass
            time.sleep(5)
            
        print("Timeout waiting for Grafana")
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
                auth=self.auth
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
                auth=self.auth
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
                auth=self.auth
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