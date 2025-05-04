"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.
"""

import json
import logging
import requests
from typing import Dict, List, Optional, Union, Any
from pathlib import Path
from datetime import datetime, timedelta
import pandas as pd
import numpy as np
from scipy import stats
from dataclasses import dataclass
from statsmodels.tsa.seasonal import seasonal_decompose
from statsmodels.tsa.statespace.sarimax import SARIMAX

@dataclass
class PanelTemplate:
    """Grafana panel template configuration."""
    title: str
    type: str
    datasource: str
    targets: List[Dict]
    options: Dict
    field_config: Dict
    transformations: List[Dict]

class GrafanaIntegration:
    def __init__(self, base_url: str = "http://localhost:3000", api_key: Optional[str] = None):
        self.base_url = base_url
        self.api_key = api_key
        self.logger = logging.getLogger(__name__)
        self.headers = {
            "Content-Type": "application/json",
            "Accept": "application/json"
        }
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"

    def create_dashboard(self, config: Dict) -> Dict:
        """Create a new Grafana dashboard."""
        try:
            dashboard = {
                "dashboard": {
                    "title": config["title"],
                    "panels": config["panels"],
                    "timezone": "browser",
                    "schemaVersion": 16,
                    "version": 0
                },
                "overwrite": False
            }

            response = requests.post(
                f"{self.base_url}/api/dashboards/db",
                headers=self.headers,
                json=dashboard
            )
            response.raise_for_status()
            return response.json()["dashboard"]

        except Exception as e:
            self.logger.error(f"Error creating Grafana dashboard: {str(e)}")
            raise

    def create_alert(self, config: Dict) -> Dict:
        """Create a new Grafana alert."""
        try:
            alert = {
                "name": config["name"],
                "condition": config["condition"],
                "data": config["data"],
                "notifications": config["notifications"],
                "frequency": "1m",
                "for": "5m"
            }

            response = requests.post(
                f"{self.base_url}/api/alert-notifications",
                headers=self.headers,
                json=alert
            )
            response.raise_for_status()
            return response.json()

        except Exception as e:
            self.logger.error(f"Error creating Grafana alert: {str(e)}")
            raise

    def get_dashboard(self, uid: str) -> Dict:
        """Get a Grafana dashboard by UID."""
        try:
            response = requests.get(
                f"{self.base_url}/api/dashboards/uid/{uid}",
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()["dashboard"]

        except Exception as e:
            self.logger.error(f"Error getting Grafana dashboard: {str(e)}")
            raise

    def delete_dashboard(self, uid: str) -> None:
        """Delete a Grafana dashboard by UID."""
        try:
            response = requests.delete(
                f"{self.base_url}/api/dashboards/uid/{uid}",
                headers=self.headers
            )
            response.raise_for_status()

        except Exception as e:
            self.logger.error(f"Error deleting Grafana dashboard: {str(e)}")
            raise

    def update_dashboard(self, uid: str, config: Dict) -> Dict:
        """Update an existing Grafana dashboard."""
        try:
            dashboard = self.get_dashboard(uid)
            dashboard.update(config)
            dashboard["version"] += 1

            response = requests.post(
                f"{self.base_url}/api/dashboards/db",
                headers=self.headers,
                json={"dashboard": dashboard, "overwrite": True}
            )
            response.raise_for_status()
            return response.json()["dashboard"]

        except Exception as e:
            self.logger.error(f"Error updating Grafana dashboard: {str(e)}")
            raise

    def list_dashboards(self) -> List[Dict]:
        """List all Grafana dashboards."""
        try:
            response = requests.get(
                f"{self.base_url}/api/search",
                headers=self.headers,
                params={"type": "dash-db"}
            )
            response.raise_for_status()
            return response.json()

        except Exception as e:
            self.logger.error(f"Error listing Grafana dashboards: {str(e)}")
            raise

    def create_custom_panel(self, dashboard_uid: str, panel_template: PanelTemplate) -> Dict:
        """Create a custom panel in a dashboard using a template.
        
        Args:
            dashboard_uid: UID of the target dashboard
            panel_template: Panel template configuration
            
        Returns:
            Created panel configuration
        """
        # Get current dashboard
        dashboard = self.get_dashboard(dashboard_uid)
        
        # Create panel from template
        panel = {
            "title": panel_template.title,
            "type": panel_template.type,
            "datasource": panel_template.datasource,
            "targets": panel_template.targets,
            "options": panel_template.options,
            "fieldConfig": panel_template.field_config,
            "transformations": panel_template.transformations
        }
        
        # Add panel to dashboard
        dashboard['dashboard']['panels'].append(panel)
        
        # Update dashboard
        self.update_dashboard(dashboard_uid, dashboard)
        
        return panel

    def analyze_trends(self, metric_name: str, time_range: str = '24h',
                      analysis_type: str = 'all') -> Dict:
        """Perform advanced trend analysis on a metric.
        
        Args:
            metric_name: Name of the metric to analyze
            time_range: Time range for analysis (e.g., '24h', '7d')
            analysis_type: Type of analysis to perform ('all', 'seasonal', 'anomaly')
            
        Returns:
            Dictionary with trend analysis results
        """
        # Get metric data
        data = self._get_metric_data(metric_name, time_range)
        
        # Convert to pandas DataFrame
        df = pd.DataFrame(data)
        df['time'] = pd.to_datetime(df['time'])
        df.set_index('time', inplace=True)
        
        results = {}
        
        # Basic trend analysis
        if analysis_type in ['all', 'basic']:
            results['basic'] = self._analyze_basic_trends(df)
        
        # Seasonal analysis
        if analysis_type in ['all', 'seasonal']:
            results['seasonal'] = self._analyze_seasonal_patterns(df)
        
        # Anomaly detection
        if analysis_type in ['all', 'anomaly']:
            results['anomaly'] = self._detect_anomalies(df)
        
        # Predictive analysis
        if analysis_type in ['all', 'predictive']:
            results['predictive'] = self._predict_future_trends(df)
        
        return results

    def _analyze_basic_trends(self, df: pd.DataFrame) -> Dict:
        """Analyze basic trends in the data."""
        # Calculate basic statistics
        stats = {
            'mean': df['value'].mean(),
            'std': df['value'].std(),
            'min': df['value'].min(),
            'max': df['value'].max(),
            'trend': self._calculate_trend(df['value'])
        }
        
        # Calculate moving averages
        stats['moving_avg_1h'] = df['value'].rolling('1H').mean().tolist()
        stats['moving_avg_24h'] = df['value'].rolling('24H').mean().tolist()
        
        # Calculate percentiles
        stats['percentiles'] = {
            '25': df['value'].quantile(0.25),
            '50': df['value'].quantile(0.5),
            '75': df['value'].quantile(0.75),
            '95': df['value'].quantile(0.95)
        }
        
        return stats

    def _analyze_seasonal_patterns(self, df: pd.DataFrame) -> Dict:
        """Analyze seasonal patterns in the data."""
        # Resample to hourly data
        hourly = df['value'].resample('H').mean()
        
        # Calculate daily patterns
        daily_pattern = hourly.groupby(hourly.index.hour).mean()
        
        # Calculate weekly patterns
        weekly_pattern = hourly.groupby(hourly.index.dayofweek).mean()
        
        # Calculate monthly patterns
        monthly_pattern = hourly.groupby(hourly.index.month).mean()
        
        # Perform seasonal decomposition
        decomposition = seasonal_decompose(hourly, model='additive', period=24)
        
        return {
            'daily_pattern': daily_pattern.to_dict(),
            'weekly_pattern': weekly_pattern.to_dict(),
            'monthly_pattern': monthly_pattern.to_dict(),
            'seasonality_score': self._calculate_seasonality_score(hourly),
            'trend_component': decomposition.trend.tolist(),
            'seasonal_component': decomposition.seasonal.tolist(),
            'residual_component': decomposition.resid.tolist()
        }

    def _detect_anomalies(self, df: pd.DataFrame) -> Dict:
        """Detect anomalies in the data."""
        # Calculate z-scores
        z_scores = np.abs(stats.zscore(df['value']))
        
        # Identify anomalies (z-score > 3)
        anomalies = df[z_scores > 3]
        
        # Calculate anomaly severity
        severity = z_scores[z_scores > 3]
        
        # Group anomalies by time periods
        hourly_anomalies = anomalies.groupby(anomalies.index.hour).size()
        daily_anomalies = anomalies.groupby(anomalies.index.dayofweek).size()
        
        return {
            'anomaly_count': len(anomalies),
            'anomaly_times': anomalies.index.tolist(),
            'anomaly_values': anomalies['value'].tolist(),
            'severity_scores': severity.tolist(),
            'hourly_distribution': hourly_anomalies.to_dict(),
            'daily_distribution': daily_anomalies.to_dict(),
            'anomaly_rate': len(anomalies) / len(df) * 100
        }

    def _predict_future_trends(self, df: pd.DataFrame) -> Dict:
        """Predict future trends using time series forecasting."""
        # Prepare data for forecasting
        series = df['value'].resample('H').mean()
        
        # Fit SARIMA model
        model = SARIMAX(series, order=(1,1,1), seasonal_order=(1,1,1,24))
        results = model.fit()
        
        # Generate forecasts
        forecast = results.get_forecast(steps=24)
        forecast_mean = forecast.predicted_mean
        forecast_ci = forecast.conf_int()
        
        return {
            'forecast': forecast_mean.tolist(),
            'confidence_intervals': {
                'lower': forecast_ci.iloc[:, 0].tolist(),
                'upper': forecast_ci.iloc[:, 1].tolist()
            },
            'model_metrics': {
                'aic': results.aic,
                'bic': results.bic,
                'hqic': results.hqic
            }
        }

    def _calculate_trend(self, series: pd.Series) -> str:
        """Calculate the overall trend of a series."""
        # Fit linear regression
        x = np.arange(len(series))
        slope, _, r_value, _, _ = stats.linregress(x, series)
        
        if slope > 0.1 and r_value > 0.5:
            return 'strongly_increasing'
        elif slope > 0.1:
            return 'increasing'
        elif slope < -0.1 and r_value > 0.5:
            return 'strongly_decreasing'
        elif slope < -0.1:
            return 'decreasing'
        else:
            return 'stable'

    def _calculate_seasonality_score(self, series: pd.Series) -> float:
        """Calculate a seasonality score for the data."""
        # Calculate autocorrelation
        autocorr = pd.Series(series).autocorr(lag=24)  # 24-hour lag
        
        # Calculate seasonal strength
        decomposition = seasonal_decompose(series, model='additive', period=24)
        seasonal_strength = 1 - np.var(decomposition.resid) / np.var(decomposition.seasonal + decomposition.resid)
        
        return (abs(autocorr) + seasonal_strength) / 2  # Combined score

    def _get_metric_data(self, metric_name: str, time_range: str) -> List[Dict]:
        """Get metric data from Prometheus."""
        # This is a placeholder - implement actual Prometheus query
        return []

    def create_advanced_alert(self, config: Dict) -> Dict:
        """Create an advanced alert with multiple conditions and notifications.
        
        Args:
            config: Alert configuration dictionary
            
        Returns:
            Created alert configuration
        """
        alert = {
            "name": config["name"],
            "conditions": config["conditions"],
            "evaluate_for": config.get("evaluate_for", "5m"),
            "no_data_state": config.get("no_data_state", "no_data"),
            "execution_error_state": config.get("execution_error_state", "alerting"),
            "notifications": config["notifications"],
            "tags": config.get("tags", []),
            "annotations": config.get("annotations", {})
        }

        try:
            response = requests.post(
                f"{self.base_url}/api/alert-notifications",
                headers=self.headers,
                json=alert
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            self.logger.error(f"Error creating advanced alert: {str(e)}")
            raise

    def create_alert_rule(self, config: Dict) -> Dict:
        """Create an alert rule with multiple conditions.
        
        Args:
            config: Alert rule configuration
            
        Returns:
            Created alert rule
        """
        rule = {
            "name": config["name"],
            "condition": self._build_alert_condition(config["conditions"]),
            "data": config["data"],
            "interval": config.get("interval", "1m"),
            "for": config.get("for", "5m"),
            "annotations": config.get("annotations", {}),
            "labels": config.get("labels", {})
        }

        try:
            response = requests.post(
                f"{self.base_url}/api/alert-rules",
                headers=self.headers,
                json=rule
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            self.logger.error(f"Error creating alert rule: {str(e)}")
            raise

    def create_notification_channel(self, config: Dict) -> Dict:
        """Create a notification channel.
        
        Args:
            config: Notification channel configuration
            
        Returns:
            Created notification channel
        """
        channel = {
            "name": config["name"],
            "type": config["type"],
            "settings": config["settings"],
            "send_reminder": config.get("send_reminder", False),
            "frequency": config.get("frequency", "1h")
        }

        try:
            response = requests.post(
                f"{self.base_url}/api/alert-notifications",
                headers=self.headers,
                json=channel
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            self.logger.error(f"Error creating notification channel: {str(e)}")
            raise

    def _build_alert_condition(self, conditions: List[Dict]) -> str:
        """Build alert condition from multiple conditions."""
        condition_parts = []
        
        for condition in conditions:
            if condition["type"] == "threshold":
                part = f"{condition['metric']} {condition['operator']} {condition['value']}"
            elif condition["type"] == "relative":
                part = f"{condition['metric']} {condition['operator']} {condition['value']}%"
            elif condition["type"] == "anomaly":
                part = f"{condition['metric']} is anomalous"
            
            condition_parts.append(part)
        
        return " AND ".join(condition_parts)

    def get_alert_history(self, alert_id: str, time_range: str = "24h") -> List[Dict]:
        """Get alert history for a specific alert.
        
        Args:
            alert_id: ID of the alert
            time_range: Time range for history
            
        Returns:
            List of alert history entries
        """
        try:
            response = requests.get(
                f"{self.base_url}/api/alerts/{alert_id}/history",
                headers=self.headers,
                params={"from": time_range}
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            self.logger.error(f"Error getting alert history: {str(e)}")
            raise

    def get_alert_stats(self, time_range: str = "24h") -> Dict:
        """Get alert statistics.
        
        Args:
            time_range: Time range for statistics
            
        Returns:
            Dictionary with alert statistics
        """
        try:
            response = requests.get(
                f"{self.base_url}/api/alerts/stats",
                headers=self.headers,
                params={"from": time_range}
            )
            response.raise_for_status()
            return response.json()
        except Exception as e:
            self.logger.error(f"Error getting alert stats: {str(e)}")
            raise

    def mute_alert(self, alert_id: str, duration: str = "1h") -> None:
        """Mute an alert for a specified duration.
        
        Args:
            alert_id: ID of the alert
            duration: Duration to mute the alert
        """
        try:
            response = requests.post(
                f"{self.base_url}/api/alerts/{alert_id}/mute",
                headers=self.headers,
                json={"duration": duration}
            )
            response.raise_for_status()
        except Exception as e:
            self.logger.error(f"Error muting alert: {str(e)}")
            raise

    def unmute_alert(self, alert_id: str) -> None:
        """Unmute an alert.
        
        Args:
            alert_id: ID of the alert
        """
        try:
            response = requests.post(
                f"{self.base_url}/api/alerts/{alert_id}/unmute",
                headers=self.headers
            )
            response.raise_for_status()
        except Exception as e:
            self.logger.error(f"Error unmuting alert: {str(e)}")
            raise

    def create_custom_dashboard(self, template_name: str, variables: Dict[str, Any]) -> Dict[str, Any]:
        """Create a custom dashboard from a template.
        
        Args:
            template_name: Name of the dashboard template
            variables: Variables to substitute in the template
            
        Returns:
            Dictionary with dashboard creation results
        """
        try:
            # Load template
            template = self._load_dashboard_template(template_name)
            
            # Substitute variables
            dashboard_json = self._substitute_template_variables(template, variables)
            
            # Create dashboard
            response = self._create_dashboard(dashboard_json)
            
            return {
                "created": True,
                "dashboard": response,
                "template": template_name
            }
            
        except Exception as e:
            self.logger.error(f"Error creating custom dashboard: {str(e)}")
            raise

    def _load_dashboard_template(self, template_name: str) -> Dict[str, Any]:
        """Load dashboard template from file."""
        try:
            template_path = f"templates/dashboards/{template_name}.json"
            with open(template_path, "r") as f:
                return json.load(f)
                
        except Exception as e:
            self.logger.error(f"Error loading dashboard template: {str(e)}")
            raise

    def _substitute_template_variables(self, template: Dict[str, Any], 
                                     variables: Dict[str, Any]) -> Dict[str, Any]:
        """Substitute variables in dashboard template."""
        try:
            template_str = json.dumps(template)
            
            for key, value in variables.items():
                template_str = template_str.replace(f"${key}", str(value))
            
            return json.loads(template_str)
            
        except Exception as e:
            self.logger.error(f"Error substituting template variables: {str(e)}")
            raise

    def _create_dashboard(self, dashboard_json: Dict[str, Any]) -> Dict[str, Any]:
        """Create dashboard in Grafana."""
        try:
            response = requests.post(
                f"{self.base_url}/api/dashboards/db",
                json=dashboard_json,
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
            
        except Exception as e:
            self.logger.error(f"Error creating dashboard: {str(e)}")
            raise

    def create_alert_route(self, route_config: Dict[str, Any]) -> Dict[str, Any]:
        """Create an advanced alert routing configuration.
        
        Args:
            route_config: Alert routing configuration
            
        Returns:
            Dictionary with route creation results
        """
        try:
            # Validate route config
            self._validate_route_config(route_config)
            
            # Create route
            response = self._create_alert_route(route_config)
            
            return {
                "created": True,
                "route": response,
                "config": route_config
            }
            
        except Exception as e:
            self.logger.error(f"Error creating alert route: {str(e)}")
            raise

    def _validate_route_config(self, config: Dict[str, Any]) -> None:
        """Validate alert route configuration."""
        required_fields = ["name", "receiver", "matchers"]
        
        for field in required_fields:
            if field not in config:
                raise ValueError(f"Missing required field: {field}")
            
        if not isinstance(config["matchers"], list):
            raise ValueError("matchers must be a list")
            
        for matcher in config["matchers"]:
            if not all(k in matcher for k in ["name", "value", "isRegex"]):
                raise ValueError("Invalid matcher format")

    def _create_alert_route(self, config: Dict[str, Any]) -> Dict[str, Any]:
        """Create alert route in Grafana."""
        try:
            response = requests.post(
                f"{self.base_url}/api/alertmanager/grafana/api/v2/routes",
                json=config,
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
            
        except Exception as e:
            self.logger.error(f"Error creating alert route: {str(e)}")
            raise

    def get_alert_routes(self) -> List[Dict[str, Any]]:
        """Get all alert routes."""
        try:
            response = requests.get(
                f"{self.base_url}/api/alertmanager/grafana/api/v2/routes",
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
            
        except Exception as e:
            self.logger.error(f"Error getting alert routes: {str(e)}")
            raise

    def update_alert_route(self, route_id: str, config: Dict[str, Any]) -> Dict[str, Any]:
        """Update an alert route."""
        try:
            response = requests.put(
                f"{self.base_url}/api/alertmanager/grafana/api/v2/routes/{route_id}",
                json=config,
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
            
        except Exception as e:
            self.logger.error(f"Error updating alert route: {str(e)}")
            raise

    def delete_alert_route(self, route_id: str) -> Dict[str, Any]:
        """Delete an alert route."""
        try:
            response = requests.delete(
                f"{self.base_url}/api/alertmanager/grafana/api/v2/routes/{route_id}",
                headers=self.headers
            )
            response.raise_for_status()
            return {"deleted": True, "route_id": route_id}
            
        except Exception as e:
            self.logger.error(f"Error deleting alert route: {str(e)}")
            raise 