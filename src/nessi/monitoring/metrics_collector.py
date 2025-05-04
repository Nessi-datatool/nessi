from typing import Dict, List, Optional
import logging
from datetime import datetime
import time
from prometheus_client import start_http_server, Gauge, Counter, Histogram
from prometheus_client.core import CollectorRegistry
import threading
import queue
import json
import os
from dataclasses import dataclass
from typing import Dict, List, Optional, Any
import asyncio
from fastapi import FastAPI, HTTPException, Security, Depends
from fastapi.security import APIKeyHeader
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
from pydantic import BaseModel
import psutil
from .config import MonitoringConfig

@dataclass
class MetricConfig:
    """Configuration for a metric."""
    name: str
    description: str
    type: str  # gauge, counter, histogram
    labels: List[str]
    thresholds: Dict[str, float]
    alert_rules: Dict[str, str]

class MetricsCollector:
    """Collects and exposes metrics for monitoring."""
    
    def __init__(self, config: MonitoringConfig):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.metrics: Dict = {}
        self._initialize_metrics()
        
    def _initialize_metrics(self):
        """Initialize all metrics based on configuration."""
        # Initialize standard metrics
        for name, metric_config in self.config.METRICS.items():
            metric_type = metric_config["type"]
            if metric_type == "gauge":
                self.metrics[name] = Gauge(
                    name,
                    metric_config["description"],
                    metric_config.get("labels", [])
                )
            elif metric_type == "counter":
                self.metrics[name] = Counter(
                    name,
                    metric_config["description"],
                    metric_config.get("labels", [])
                )
            elif metric_type == "histogram":
                self.metrics[name] = Histogram(
                    name,
                    metric_config["description"],
                    metric_config.get("labels", []),
                    buckets=metric_config.get("buckets", [])
                )
        
        # Initialize custom metrics
        for metric_config in self.config.CUSTOM_METRICS:
            name = metric_config["name"]
            metric_type = metric_config["type"]
            if metric_type == "gauge":
                self.metrics[name] = Gauge(
                    name,
                    metric_config["description"],
                    metric_config.get("labels", [])
                )
    
    def start(self):
        """Start the metrics server and background collection."""
        if self.config.PROMETHEUS_ENABLED:
            start_http_server(self.config.PROMETHEUS_PORT)
            self.logger.info(f"Started metrics server on port {self.config.PROMETHEUS_PORT}")
        
        # Start background collection for performance metrics
        if any(metric["enabled"] for metric in self.config.PERFORMANCE_METRICS.values()):
            self._start_performance_collection()
    
    def _start_performance_collection(self):
        """Start background collection of performance metrics."""
        def collect_performance_metrics():
            while True:
                try:
                    if self.config.PERFORMANCE_METRICS["cpu_usage"]["enabled"]:
                        self.metrics["cpu_usage"].set(psutil.cpu_percent())
                    
                    if self.config.PERFORMANCE_METRICS["memory_usage"]["enabled"]:
                        memory = psutil.virtual_memory()
                        self.metrics["memory_usage"].set(memory.percent)
                    
                    if self.config.PERFORMANCE_METRICS["disk_usage"]["enabled"]:
                        disk = psutil.disk_usage('/')
                        self.metrics["disk_usage"].set(disk.percent)
                except Exception as e:
                    self.logger.error(f"Error collecting performance metrics: {e}")
                
                time.sleep(60)  # Collect every minute
        
        thread = threading.Thread(target=collect_performance_metrics, daemon=True)
        thread.start()
    
    def update_metric(self, name: str, value: float, labels: Optional[Dict[str, str]] = None):
        """Update a metric value with optional labels."""
        if name not in self.metrics:
            self.logger.warning(f"Metric {name} not found")
            return
        
        metric = self.metrics[name]
        if labels:
            metric.labels(**labels).set(value)
        else:
            metric.set(value)
    
    def increment_counter(self, name: str, labels: Optional[Dict[str, str]] = None):
        """Increment a counter metric."""
        if name not in self.metrics:
            self.logger.warning(f"Counter {name} not found")
            return
        
        metric = self.metrics[name]
        if labels:
            metric.labels(**labels).inc()
        else:
            metric.inc()
    
    def observe_histogram(self, name: str, value: float, labels: Optional[Dict[str, str]] = None):
        """Observe a value in a histogram metric."""
        if name not in self.metrics:
            self.logger.warning(f"Histogram {name} not found")
            return
        
        metric = self.metrics[name]
        if labels:
            metric.labels(**labels).observe(value)
        else:
            metric.observe(value)
    
    def check_thresholds(self, name: str, value: float) -> Dict[str, bool]:
        """Check if a metric value exceeds configured thresholds."""
        if name not in self.config.METRICS:
            return {}
        
        metric_config = self.config.METRICS[name]
        if "thresholds" not in metric_config:
            return {}
        
        thresholds = metric_config["thresholds"]
        return {
            "warning": value < thresholds.get("warning", float("inf")),
            "critical": value < thresholds.get("critical", float("inf"))
        }

class MetricsAPI:
    """FastAPI application for metrics and alerts."""
    
    def __init__(self, metrics_collector: MetricsCollector, api_key: str):
        self.app = FastAPI(title="Nessi Metrics API")
        self.metrics_collector = metrics_collector
        self.api_key = api_key
        
        # Add CORS middleware
        self.app.add_middleware(
            CORSMiddleware,
            allow_origins=["*"],
            allow_credentials=True,
            allow_methods=["*"],
            allow_headers=["*"],
        )
        
        # Add security
        self.api_key_header = APIKeyHeader(name="X-API-Key")
        
        # Add routes
        self.app.get("/metrics")(self.get_metrics)
        self.app.get("/alerts")(self.get_alerts)
        self.app.post("/update")(self.update_metric)
    
    async def get_metrics(self, api_key: str = Security(APIKeyHeader(name="X-API-Key"))):
        """Get current metrics."""
        if api_key != self.api_key:
            raise HTTPException(status_code=403, detail="Invalid API key")
        return self.metrics_collector.get_metrics()
    
    async def get_alerts(self, api_key: str = Security(APIKeyHeader(name="X-API-Key"))):
        """Get recent alerts."""
        if api_key != self.api_key:
            raise HTTPException(status_code=403, detail="Invalid API key")
        alerts = []
        while not self.metrics_collector.alert_queue.empty():
            alerts.append(self.metrics_collector.alert_queue.get())
        return alerts
    
    async def update_metric(self, metric_data: Dict, api_key: str = Security(APIKeyHeader(name="X-API-Key"))):
        """Update a metric value."""
        if api_key != self.api_key:
            raise HTTPException(status_code=403, detail="Invalid API key")
        
        self.metrics_collector.update_metric(
            metric_data['name'],
            metric_data['value'],
            metric_data['labels']
        )
        return {"status": "success"}
    
    def run(self, host: str = "0.0.0.0", port: int = 8001):
        """Run the FastAPI application."""
        uvicorn.run(self.app, host=host, port=port) 