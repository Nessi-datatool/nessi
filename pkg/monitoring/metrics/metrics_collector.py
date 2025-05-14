"""
Metrics Collector for Nessi.dev.
"""

from typing import Dict, List, Any, Optional
from datetime import datetime


class MetricsCollector:
    """Metrics Collector for Nessi.dev."""

    def __init__(self):
        """Initialize the Metrics Collector."""
        self.metrics = {}

    def record_metric(self, name: str, value: Any, tags: Optional[Dict[str, str]] = None) -> Dict[str, Any]:
        """Record a metric.

        Args:
            name: Metric name.
            value: Metric value.
            tags: Optional tags.

        Returns:
            Dict with metric information.
        """
        timestamp = datetime.now().isoformat()
        tags = tags or {}
        
        metric = {
            'name': name,
            'value': value,
            'tags': tags,
            'timestamp': timestamp,
        }
        
        if name not in self.metrics:
            self.metrics[name] = []
        
        self.metrics[name].append(metric)
        return metric

    def get_metrics(self, name: str, tags: Optional[Dict[str, str]] = None) -> List[Dict[str, Any]]:
        """Get metrics.

        Args:
            name: Metric name.
            tags: Optional tags filter.

        Returns:
            List of metrics.
        """
        if name not in self.metrics:
            return []
        
        metrics = self.metrics[name]
        
        if tags:
            filtered_metrics = []
            for metric in metrics:
                match = True
                for key, value in tags.items():
                    if key not in metric['tags'] or metric['tags'][key] != value:
                        match = False
                        break
                if match:
                    filtered_metrics.append(metric)
            return filtered_metrics
        
        return metrics
