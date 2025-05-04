from typing import Dict, List, Optional
from pydantic import BaseModel
import os
from datetime import timedelta

class MonitoringConfig(BaseModel):
    """Monitoring configuration for the application."""
    
    # Prometheus Configuration
    PROMETHEUS_ENABLED: bool = True
    PROMETHEUS_PORT: int = 9090
    PROMETHEUS_PATH: str = "/metrics"
    PROMETHEUS_SCRAPE_INTERVAL: int = 15  # seconds
    
    # Grafana Configuration
    GRAFANA_ENABLED: bool = True
    GRAFANA_PORT: int = 3000
    GRAFANA_ADMIN_USER: str = "admin"
    GRAFANA_ADMIN_PASSWORD: str = os.getenv("GRAFANA_ADMIN_PASSWORD", "admin")
    
    # Alertmanager Configuration
    ALERTMANAGER_ENABLED: bool = True
    ALERTMANAGER_PORT: int = 9093
    ALERTMANAGER_CONFIG_PATH: str = "/etc/alertmanager/config.yml"
    
    # Metrics Configuration
    METRICS: Dict[str, Dict] = {
        "data_quality_score": {
            "type": "gauge",
            "description": "Overall data quality score",
            "labels": ["table", "database"],
            "thresholds": {
                "warning": 0.8,
                "critical": 0.6
            }
        },
        "scan_duration": {
            "type": "histogram",
            "description": "Time taken for table scans",
            "labels": ["table", "database"],
            "buckets": [0.1, 0.5, 1.0, 2.0, 5.0, 10.0]
        },
        "row_count": {
            "type": "gauge",
            "description": "Number of rows in tables",
            "labels": ["table", "database"]
        },
        "error_count": {
            "type": "counter",
            "description": "Number of errors encountered",
            "labels": ["type", "severity"]
        }
    }
    
    # Alert Rules
    ALERT_RULES: List[Dict] = [
        {
            "name": "DataQualityDegradation",
            "expr": "data_quality_score < 0.8",
            "for": "5m",
            "labels": {
                "severity": "warning"
            },
            "annotations": {
                "summary": "Data quality score below threshold",
                "description": "Data quality score for {{ $labels.table }} is {{ $value }}"
            }
        },
        {
            "name": "HighErrorRate",
            "expr": "rate(error_count[5m]) > 0.1",
            "for": "5m",
            "labels": {
                "severity": "critical"
            },
            "annotations": {
                "summary": "High error rate detected",
                "description": "Error rate is {{ $value }} errors per second"
            }
        }
    ]
    
    # Logging Configuration
    LOGGING: Dict = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": {
            "json": {
                "class": "pythonjsonlogger.jsonlogger.JsonFormatter",
                "format": "%(asctime)s %(levelname)s %(name)s %(message)s"
            }
        },
        "handlers": {
            "console": {
                "class": "logging.StreamHandler",
                "formatter": "json",
                "stream": "ext://sys.stdout"
            },
            "file": {
                "class": "logging.handlers.RotatingFileHandler",
                "formatter": "json",
                "filename": "/var/log/nessi/app.log",
                "maxBytes": 10485760,
                "backupCount": 5
            }
        },
        "root": {
            "level": "INFO",
            "handlers": ["console", "file"]
        }
    }
    
    # Health Check Configuration
    HEALTH_CHECK_INTERVAL: int = 30  # seconds
    HEALTH_CHECK_TIMEOUT: int = 10  # seconds
    HEALTH_CHECK_RETRIES: int = 3
    
    # Performance Monitoring
    PERFORMANCE_METRICS: Dict = {
        "cpu_usage": {
            "enabled": True,
            "interval": 60  # seconds
        },
        "memory_usage": {
            "enabled": True,
            "interval": 60  # seconds
        },
        "disk_usage": {
            "enabled": True,
            "interval": 300  # seconds
        }
    }
    
    # Custom Metrics
    CUSTOM_METRICS: List[Dict] = [
        {
            "name": "custom_metric_1",
            "type": "gauge",
            "description": "Custom metric 1",
            "labels": ["label1", "label2"]
        }
    ] 