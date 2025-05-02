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