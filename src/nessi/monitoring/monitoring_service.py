"""Service for monitoring data quality and performance metrics."""

from typing import Dict, Any, Optional
import time
import logging
from pathlib import Path
from datetime import datetime
import threading

from .monitor import Monitor
from .metrics_collector import MetricsCollector
from .config import MonitoringConfig

logger = logging.getLogger(__name__)

class MonitoringService:
    """Service for monitoring data quality and performance metrics."""

    def __init__(self, config: MonitoringConfig):
        """Initialize the monitoring service.

        Args:
            config: Configuration for the monitoring service
        """
        self.config = config
        self.monitor = Monitor(config)
        self.metrics_collector = MetricsCollector(config)
        self._monitoring_thread = None
        self._stop_monitoring = False

    def start(self) -> None:
        """Start the monitoring service."""
        if self._monitoring_thread is not None and self._monitoring_thread.is_alive():
            logger.warning("Monitoring service is already running")
            return

        self._stop_monitoring = False
        self._monitoring_thread = threading.Thread(target=self._monitoring_loop)
        self._monitoring_thread.daemon = True
        self._monitoring_thread.start()
        logger.info("Started monitoring service")

    def stop(self) -> None:
        """Stop the monitoring service."""
        if self._monitoring_thread is None or not self._monitoring_thread.is_alive():
            logger.warning("Monitoring service is not running")
            return

        self._stop_monitoring = True
        self._monitoring_thread.join()
        logger.info("Stopped monitoring service")

    def _monitoring_loop(self) -> None:
        """Main monitoring loop."""
        while not self._stop_monitoring:
            try:
                # Collect metrics
                metrics = self.metrics_collector.collect_metrics()
                
                # Monitor data quality
                quality_metrics = self.monitor.check_data_quality()
                
                # Combine metrics
                all_metrics = {**metrics, **quality_metrics}
                
                # Store metrics
                self._store_metrics(all_metrics)
                
                # Wait for next collection interval
                time.sleep(self.config.collection_interval)
            except Exception as e:
                logger.error(f"Error in monitoring loop: {e}", exc_info=True)
                time.sleep(60)  # Wait before retrying

    def _store_metrics(self, metrics: Dict[str, Any]) -> None:
        """Store collected metrics.

        Args:
            metrics: Dictionary of metrics to store
        """
        try:
            # Add timestamp
            metrics["timestamp"] = datetime.utcnow().isoformat()
            
            # Store metrics (implement storage logic here)
            logger.debug(f"Stored metrics: {metrics}")
        except Exception as e:
            logger.error(f"Error storing metrics: {e}", exc_info=True) 