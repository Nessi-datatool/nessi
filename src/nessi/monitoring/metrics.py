"""Custom metric collection and SLA tracking."""
import logging
from typing import Dict, List, Optional, Any
from datetime import datetime, timedelta
import json
import os
from dataclasses import dataclass
from enum import Enum

logger = logging.getLogger(__name__)

class MetricType(Enum):
    """Types of metrics."""
    COUNTER = "counter"
    GAUGE = "gauge"
    HISTOGRAM = "histogram"
    SUMMARY = "summary"

@dataclass
class Metric:
    """Metric definition."""
    name: str
    type: MetricType
    description: str
    labels: Dict[str, str]
    value: float
    timestamp: datetime

@dataclass
class SLA:
    """Service Level Agreement definition."""
    name: str
    metric_name: str
    threshold: float
    window: timedelta
    description: str

class MetricsCollector:
    """Custom metric collection system."""
    def __init__(self, metrics_dir: str = "metrics"):
        self.metrics: Dict[str, List[Metric]] = {}
        self.slas: Dict[str, SLA] = {}
        self.metrics_dir = metrics_dir
        
        # Create metrics directory if it doesn't exist
        os.makedirs(metrics_dir, exist_ok=True)

    def register_metric(self, name: str, metric_type: MetricType,
                       description: str, labels: Dict[str, str]) -> None:
        """Register a new metric.
        
        Args:
            name: Name of the metric
            metric_type: Type of metric
            description: Description of the metric
            labels: Labels for the metric
        """
        try:
            if name not in self.metrics:
                self.metrics[name] = []
                
            logger.info(f"Registered metric: {name}")
            
        except Exception as e:
            logger.error(f"Error registering metric: {str(e)}")
            raise

    def record_metric(self, name: str, value: float,
                     labels: Optional[Dict[str, str]] = None) -> None:
        """Record a metric value.
        
        Args:
            name: Name of the metric
            value: Value to record
            labels: Optional labels for the metric
        """
        try:
            if name not in self.metrics:
                raise ValueError(f"Metric {name} not registered")
                
            metric = Metric(
                name=name,
                type=self._get_metric_type(name),
                description=self._get_metric_description(name),
                labels=labels or {},
                value=value,
                timestamp=datetime.now()
            )
            
            self.metrics[name].append(metric)
            
            # Check SLAs
            self._check_slas(name, value)
            
            # Save metric
            self._save_metric(metric)
            
        except Exception as e:
            logger.error(f"Error recording metric: {str(e)}")
            raise

    def _get_metric_type(self, name: str) -> MetricType:
        """Get type of a registered metric."""
        # In a real implementation, this would return the actual type
        return MetricType.GAUGE

    def _get_metric_description(self, name: str) -> str:
        """Get description of a registered metric."""
        # In a real implementation, this would return the actual description
        return ""

    def _save_metric(self, metric: Metric) -> None:
        """Save metric to disk."""
        try:
            metric_file = os.path.join(self.metrics_dir, f"{metric.name}.json")
            
            metric_data = {
                "name": metric.name,
                "type": metric.type.value,
                "description": metric.description,
                "labels": metric.labels,
                "value": metric.value,
                "timestamp": metric.timestamp.isoformat()
            }
            
            with open(metric_file, "a") as f:
                f.write(json.dumps(metric_data) + "\n")
                
        except Exception as e:
            logger.error(f"Error saving metric: {str(e)}")
            raise

    def define_sla(self, name: str, metric_name: str, threshold: float,
                  window: timedelta, description: str) -> SLA:
        """Define a new SLA.
        
        Args:
            name: Name of the SLA
            metric_name: Name of the metric to track
            threshold: Threshold value
            window: Time window for evaluation
            description: Description of the SLA
            
        Returns:
            Created SLA
        """
        try:
            sla = SLA(
                name=name,
                metric_name=metric_name,
                threshold=threshold,
                window=window,
                description=description
            )
            
            self.slas[name] = sla
            
            logger.info(f"Defined SLA: {name}")
            
            return sla
            
        except Exception as e:
            logger.error(f"Error defining SLA: {str(e)}")
            raise

    def _check_slas(self, metric_name: str, value: float) -> None:
        """Check SLAs for a metric value."""
        try:
            for sla in self.slas.values():
                if sla.metric_name == metric_name:
                    self._evaluate_sla(sla, value)
                    
        except Exception as e:
            logger.error(f"Error checking SLAs: {str(e)}")
            raise

    def _evaluate_sla(self, sla: SLA, value: float) -> None:
        """Evaluate an SLA against a metric value."""
        try:
            # Get recent metrics within the SLA window
            recent_metrics = self._get_recent_metrics(
                sla.metric_name,
                datetime.now() - sla.window
            )
            
            if not recent_metrics:
                return
                
            # Calculate average
            avg_value = sum(m.value for m in recent_metrics) / len(recent_metrics)
            
            # Check if SLA is violated
            if avg_value > sla.threshold:
                logger.warning(
                    f"SLA violation: {sla.name} "
                    f"(threshold: {sla.threshold}, actual: {avg_value})"
                )
                
        except Exception as e:
            logger.error(f"Error evaluating SLA: {str(e)}")
            raise

    def _get_recent_metrics(self, metric_name: str,
                           start_time: datetime) -> List[Metric]:
        """Get recent metrics for a given time range."""
        try:
            if metric_name not in self.metrics:
                return []
                
            return [
                m for m in self.metrics[metric_name]
                if m.timestamp >= start_time
            ]
            
        except Exception as e:
            logger.error(f"Error getting recent metrics: {str(e)}")
            raise

    def get_metric_history(self, name: str, start_time: Optional[datetime] = None,
                          end_time: Optional[datetime] = None) -> List[Metric]:
        """Get metric history.
        
        Args:
            name: Name of the metric
            start_time: Optional start time for filtering
            end_time: Optional end time for filtering
            
        Returns:
            List of metrics
        """
        try:
            if name not in self.metrics:
                return []
                
            metrics = self.metrics[name]
            
            if start_time:
                metrics = [m for m in metrics if m.timestamp >= start_time]
                
            if end_time:
                metrics = [m for m in metrics if m.timestamp <= end_time]
                
            return metrics
            
        except Exception as e:
            logger.error(f"Error getting metric history: {str(e)}")
            raise

    def get_sla_status(self, name: str) -> Dict[str, Any]:
        """Get current status of an SLA.
        
        Args:
            name: Name of the SLA
            
        Returns:
            Dictionary with SLA status
        """
        try:
            if name not in self.slas:
                raise ValueError(f"SLA {name} not found")
                
            sla = self.slas[name]
            recent_metrics = self._get_recent_metrics(
                sla.metric_name,
                datetime.now() - sla.window
            )
            
            if not recent_metrics:
                return {
                    "name": name,
                    "status": "unknown",
                    "message": "No recent metrics"
                }
                
            avg_value = sum(m.value for m in recent_metrics) / len(recent_metrics)
            
            return {
                "name": name,
                "status": "ok" if avg_value <= sla.threshold else "violated",
                "threshold": sla.threshold,
                "current_value": avg_value,
                "window": str(sla.window)
            }
            
        except Exception as e:
            logger.error(f"Error getting SLA status: {str(e)}")
            raise 