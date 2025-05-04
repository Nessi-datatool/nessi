from typing import Dict, Any, Optional
import time
import logging
from prometheus_client import start_http_server, Gauge, Counter, Histogram
from datetime import datetime
import threading
from pathlib import Path
import json

logger = logging.getLogger(__name__)

class Monitor:
    """Real-time monitoring service for Nessi.dev."""
    
    def __init__(self, config_path: str = "config/monitoring.json"):
        self.config = self._load_config(config_path)
        self._setup_metrics()
        self._start_metrics_server()
    
    def _load_config(self, config_path: str) -> Dict[str, Any]:
        """Load monitoring configuration."""
        try:
            with open(config_path, 'r') as f:
                return json.load(f)
        except Exception as e:
            logger.error(f"Error loading monitoring config: {str(e)}")
            return {}
    
    def _setup_metrics(self):
        """Set up Prometheus metrics."""
        # Counters
        self.scan_requests = Counter(
            'nessi_scan_requests_total',
            'Total number of scan requests'
        )
        self.validation_requests = Counter(
            'nessi_validation_requests_total',
            'Total number of validation requests'
        )
        self.errors = Counter(
            'nessi_errors_total',
            'Total number of errors',
            ['type']
        )
        
        # Gauges
        self.active_scans = Gauge(
            'nessi_active_scans',
            'Number of currently active scans'
        )
        self.data_quality_score = Gauge(
            'nessi_data_quality_score',
            'Current data quality score',
            ['source']
        )
        self.anomaly_count = Gauge(
            'nessi_anomaly_count',
            'Number of detected anomalies',
            ['source']
        )
        
        # Histograms
        self.scan_duration = Histogram(
            'nessi_scan_duration_seconds',
            'Time taken for scans',
            ['source_type']
        )
        self.validation_duration = Histogram(
            'nessi_validation_duration_seconds',
            'Time taken for validations',
            ['rule_type']
        )
    
    def _start_metrics_server(self):
        """Start the Prometheus metrics server."""
        try:
            port = self.config.get('metrics_port', 9090)
            start_http_server(port)
            logger.info(f"Metrics server started on port {port}")
        except Exception as e:
            logger.error(f"Error starting metrics server: {str(e)}")
    
    def record_scan_start(self, source_type: str):
        """Record the start of a scan operation."""
        self.scan_requests.inc()
        self.active_scans.inc()
        return time.time()
    
    def record_scan_end(self, start_time: float, source_type: str):
        """Record the end of a scan operation."""
        duration = time.time() - start_time
        self.scan_duration.labels(source_type=source_type).observe(duration)
        self.active_scans.dec()
    
    def record_validation_start(self, rule_type: str):
        """Record the start of a validation operation."""
        self.validation_requests.inc()
        return time.time()
    
    def record_validation_end(self, start_time: float, rule_type: str):
        """Record the end of a validation operation."""
        duration = time.time() - start_time
        self.validation_duration.labels(rule_type=rule_type).observe(duration)
    
    def record_error(self, error_type: str):
        """Record an error occurrence."""
        self.errors.labels(type=error_type).inc()
    
    def update_quality_score(self, source: str, score: float):
        """Update the data quality score for a source."""
        self.data_quality_score.labels(source=source).set(score)
    
    def update_anomaly_count(self, source: str, count: int):
        """Update the anomaly count for a source."""
        self.anomaly_count.labels(source=source).set(count)
    
    def generate_alert(self, alert_type: str, message: str, severity: str = "warning"):
        """Generate an alert."""
        alert = {
            "timestamp": datetime.utcnow().isoformat(),
            "type": alert_type,
            "message": message,
            "severity": severity
        }
        self._write_alert(alert)
        return alert
    
    def _write_alert(self, alert: Dict[str, Any]):
        """Write alert to alert log."""
        try:
            alert_path = Path(self.config.get('alert_log_path', 'logs/alerts.log'))
            alert_path.parent.mkdir(parents=True, exist_ok=True)
            
            with open(alert_path, 'a') as f:
                f.write(json.dumps(alert) + '\n')
        except Exception as e:
            logger.error(f"Error writing alert: {str(e)}")

    def update_table_metrics(self, path: str, metrics: Dict[str, Any]):
        """Update table-related metrics."""
        try:
            self.table_size.set(metrics.get('size', 0))
            self.table_row_count.set(metrics.get('num_rows', 0))
            self.table_file_count.set(metrics.get('num_files', 0))
        except Exception as e:
            logger.error(f"Error updating table metrics: {str(e)}")
    
    def update_quality_metrics(self, metrics: Dict[str, Any]):
        """Update quality-related metrics."""
        try:
            self.quality_score.set(metrics.get('quality_score', 0))
            self.completeness_score.set(metrics.get('completeness_score', 0))
            self.consistency_score.set(metrics.get('consistency_score', 0))
        except Exception as e:
            logger.error(f"Error updating quality metrics: {str(e)}")
    
    def record_optimization_duration(self, duration: float):
        """Record the duration of an optimization operation."""
        self.optimization_duration.observe(duration)
    
    def increment_optimization_errors(self):
        """Increment the optimization error counter."""
        self.optimization_errors.inc()
    
    def generate_grafana_dashboard(self, output_path: str):
        """Generate a Grafana dashboard configuration."""
        try:
            dashboard = {
                "dashboard": {
                    "title": "Nessi Data Quality Dashboard",
                    "panels": [
                        {
                            "title": "Table Health",
                            "type": "row",
                            "panels": [
                                {
                                    "title": "Table Size",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_table_size_bytes"}]
                                },
                                {
                                    "title": "Row Count",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_table_row_count"}]
                                },
                                {
                                    "title": "File Count",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_table_file_count"}]
                                }
                            ]
                        },
                        {
                            "title": "Data Quality",
                            "type": "row",
                            "panels": [
                                {
                                    "title": "Quality Score",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_quality_score"}]
                                },
                                {
                                    "title": "Completeness Score",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_completeness_score"}]
                                },
                                {
                                    "title": "Consistency Score",
                                    "type": "gauge",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "nessi_consistency_score"}]
                                }
                            ]
                        },
                        {
                            "title": "Performance",
                            "type": "row",
                            "panels": [
                                {
                                    "title": "Scan Duration",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_scan_duration_seconds_sum[5m]) / rate(nessi_scan_duration_seconds_count[5m])"}]
                                },
                                {
                                    "title": "Validation Duration",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_validation_duration_seconds_sum[5m]) / rate(nessi_validation_duration_seconds_count[5m])"}]
                                },
                                {
                                    "title": "Optimization Duration",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_optimization_duration_seconds_sum[5m]) / rate(nessi_optimization_duration_seconds_count[5m])"}]
                                }
                            ]
                        },
                        {
                            "title": "Errors",
                            "type": "row",
                            "panels": [
                                {
                                    "title": "Scan Errors",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_scan_errors_total[5m])"}]
                                },
                                {
                                    "title": "Validation Errors",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_validation_errors_total[5m])"}]
                                },
                                {
                                    "title": "Optimization Errors",
                                    "type": "graph",
                                    "datasource": "Prometheus",
                                    "targets": [{"expr": "rate(nessi_optimization_errors_total[5m])"}]
                                }
                            ]
                        }
                    ]
                }
            }
            
            with open(output_path, 'w') as f:
                json.dump(dashboard, f, indent=2)
            
            logger.info(f"Generated Grafana dashboard configuration at {output_path}")
        
        except Exception as e:
            logger.error(f"Error generating Grafana dashboard: {str(e)}")
            raise 