"""Script to customize Grafana alerts."""

import json
import argparse
from typing import Dict, Any, List

class AlertCustomizer:
    """Helps customize Grafana alerts."""
    
    def __init__(self, dashboard_path: str):
        self.dashboard_path = dashboard_path
        with open(dashboard_path, 'r') as f:
            self.dashboard = json.load(f)
            
    def get_available_alerts(self) -> List[Dict[str, Any]]:
        """Get list of available alerts in the dashboard."""
        alerts = []
        for panel in self.dashboard.get('panels', []):
            if 'alert' in panel:
                alerts.append({
                    'panel_id': panel['id'],
                    'title': panel['title'],
                    'metric': panel['targets'][0]['expr'],
                    'current_threshold': panel['alert']['conditions'][0]['evaluator']['params'][0]
                })
        return alerts
        
    def update_alert(self, panel_id: int, threshold: float, 
                    notification_channels: List[str] = None) -> None:
        """Update alert threshold and notification channels."""
        for panel in self.dashboard['panels']:
            if panel['id'] == panel_id and 'alert' in panel:
                # Update threshold
                panel['alert']['conditions'][0]['evaluator']['params'][0] = threshold
                
                # Update notification channels if provided
                if notification_channels:
                    panel['alert']['notifications'] = notification_channels
                    
    def save_dashboard(self) -> None:
        """Save updated dashboard to file."""
        with open(self.dashboard_path, 'w') as f:
            json.dump(self.dashboard, f, indent=2)
            
def main():
    parser = argparse.ArgumentParser(description='Customize Grafana alerts')
    parser.add_argument('--dashboard', required=True, help='Path to dashboard JSON file')
    parser.add_argument('--list', action='store_true', help='List available alerts')
    parser.add_argument('--update', nargs=3, metavar=('PANEL_ID', 'THRESHOLD', 'CHANNELS'),
                      help='Update alert threshold and notification channels')
    
    args = parser.parse_args()
    customizer = AlertCustomizer(args.dashboard)
    
    if args.list:
        alerts = customizer.get_available_alerts()
        print("\nAvailable Alerts:")
        for alert in alerts:
            print(f"\nPanel: {alert['title']}")
            print(f"ID: {alert['panel_id']}")
            print(f"Metric: {alert['metric']}")
            print(f"Current Threshold: {alert['current_threshold']}")
            
    elif args.update:
        panel_id = int(args.update[0])
        threshold = float(args.update[1])
        channels = args.update[2].split(',') if args.update[2] else None
        
        customizer.update_alert(panel_id, threshold, channels)
        customizer.save_dashboard()
        print(f"Updated alert for panel {panel_id} with threshold {threshold}")
        
if __name__ == "__main__":
    main() 