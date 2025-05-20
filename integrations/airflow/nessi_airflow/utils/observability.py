"""
Observability utilities for Nessi Airflow integration.

This module provides utilities for metrics collection, structured logging,
and tracing to enhance monitoring and debugging capabilities.
"""

import time
import logging
import json
import functools
import uuid
import socket
import os
from typing import Dict, Any, Optional, Callable, List, Union, TypeVar, cast

from airflow.configuration import conf
from airflow.stats import Stats

# Type variables for function decorators
F = TypeVar('F', bound=Callable[..., Any])

# Configure logger
logger = logging.getLogger(__name__)


class MetricsCollector:
    """Collector for Nessi Airflow integration metrics."""
    
    # Metric name prefixes
    PREFIX = "nessi_airflow"
    API_PREFIX = f"{PREFIX}.api"
    HOOK_PREFIX = f"{PREFIX}.hook"
    OPERATOR_PREFIX = f"{PREFIX}.operator"
    SENSOR_PREFIX = f"{PREFIX}.sensor"
    
    @classmethod
    def timer(cls, name: str) -> Callable[[F], F]:
        """
        Decorator to time a function and record the duration as a metric.
        
        Args:
            name: Metric name
            
        Returns:
            Decorated function
        """
        def decorator(func: F) -> F:
            @functools.wraps(func)
            def wrapper(*args, **kwargs):
                start_time = time.time()
                
                try:
                    result = func(*args, **kwargs)
                    # Record successful execution
                    Stats.timing(
                        f"{name}.duration",
                        (time.time() - start_time) * 1000,  # Convert to ms
                    )
                    Stats.incr(f"{name}.success")
                    return result
                except Exception as e:
                    # Record failed execution
                    Stats.timing(
                        f"{name}.duration",
                        (time.time() - start_time) * 1000,  # Convert to ms
                    )
                    Stats.incr(f"{name}.failure")
                    # Re-raise the exception
                    raise
            
            return cast(F, wrapper)
        
        return decorator
    
    @classmethod
    def api_timer(cls, endpoint: str) -> Callable[[F], F]:
        """
        Decorator to time an API call and record the duration as a metric.
        
        Args:
            endpoint: API endpoint name
            
        Returns:
            Decorated function
        """
        return cls.timer(f"{cls.API_PREFIX}.{endpoint}")
    
    @classmethod
    def hook_timer(cls, method: str) -> Callable[[F], F]:
        """
        Decorator to time a hook method and record the duration as a metric.
        
        Args:
            method: Hook method name
            
        Returns:
            Decorated function
        """
        return cls.timer(f"{cls.HOOK_PREFIX}.{method}")
    
    @classmethod
    def operator_timer(cls, operator: str) -> Callable[[F], F]:
        """
        Decorator to time an operator and record the duration as a metric.
        
        Args:
            operator: Operator name
            
        Returns:
            Decorated function
        """
        return cls.timer(f"{cls.OPERATOR_PREFIX}.{operator}")
    
    @classmethod
    def sensor_timer(cls, sensor: str) -> Callable[[F], F]:
        """
        Decorator to time a sensor and record the duration as a metric.
        
        Args:
            sensor: Sensor name
            
        Returns:
            Decorated function
        """
        return cls.timer(f"{cls.SENSOR_PREFIX}.{sensor}")
    
    @classmethod
    def count(cls, name: str, value: int = 1) -> None:
        """
        Increment a counter metric.
        
        Args:
            name: Metric name
            value: Value to increment by
        """
        Stats.incr(f"{cls.PREFIX}.{name}", value)
    
    @classmethod
    def gauge(cls, name: str, value: float) -> None:
        """
        Set a gauge metric.
        
        Args:
            name: Metric name
            value: Gauge value
        """
        Stats.gauge(f"{cls.PREFIX}.{name}", value)


class StructuredLogger:
    """
    Logger that outputs structured logs in JSON format.
    
    This logger adds context information to logs, such as correlation IDs,
    to help with debugging and tracing.
    """
    
    def __init__(self, logger_name: str = "nessi_airflow"):
        """
        Initialize the structured logger.
        
        Args:
            logger_name: Name of the logger
        """
        self.logger = logging.getLogger(logger_name)
        self.context = {}
        
        # Generate a unique ID for this logger instance
        self.instance_id = str(uuid.uuid4())
        
        # Add host information
        self.hostname = socket.gethostname()
        
        # Add environment information
        self.environment = os.environ.get("AIRFLOW_ENV", "development")
    
    def with_context(self, **kwargs) -> "StructuredLogger":
        """
        Create a new logger with additional context.
        
        Args:
            **kwargs: Context key-value pairs
            
        Returns:
            New logger instance with updated context
        """
        new_logger = StructuredLogger(self.logger.name)
        new_logger.context = {**self.context, **kwargs}
        return new_logger
    
    def _format_log(self, level: str, message: str, **kwargs) -> str:
        """
        Format a log message as JSON.
        
        Args:
            level: Log level
            message: Log message
            **kwargs: Additional log data
            
        Returns:
            JSON-formatted log string
        """
        log_data = {
            "timestamp": time.time(),
            "level": level,
            "message": message,
            "logger": self.logger.name,
            "instance_id": self.instance_id,
            "hostname": self.hostname,
            "environment": self.environment,
            **self.context,
            **kwargs,
        }
        
        return json.dumps(log_data)
    
    def debug(self, message: str, **kwargs) -> None:
        """
        Log a debug message.
        
        Args:
            message: Log message
            **kwargs: Additional log data
        """
        if self.logger.isEnabledFor(logging.DEBUG):
            self.logger.debug(self._format_log("DEBUG", message, **kwargs))
    
    def info(self, message: str, **kwargs) -> None:
        """
        Log an info message.
        
        Args:
            message: Log message
            **kwargs: Additional log data
        """
        if self.logger.isEnabledFor(logging.INFO):
            self.logger.info(self._format_log("INFO", message, **kwargs))
    
    def warning(self, message: str, **kwargs) -> None:
        """
        Log a warning message.
        
        Args:
            message: Log message
            **kwargs: Additional log data
        """
        if self.logger.isEnabledFor(logging.WARNING):
            self.logger.warning(self._format_log("WARNING", message, **kwargs))
    
    def error(self, message: str, **kwargs) -> None:
        """
        Log an error message.
        
        Args:
            message: Log message
            **kwargs: Additional log data
        """
        if self.logger.isEnabledFor(logging.ERROR):
            self.logger.error(self._format_log("ERROR", message, **kwargs))
    
    def critical(self, message: str, **kwargs) -> None:
        """
        Log a critical message.
        
        Args:
            message: Log message
            **kwargs: Additional log data
        """
        if self.logger.isEnabledFor(logging.CRITICAL):
            self.logger.critical(self._format_log("CRITICAL", message, **kwargs))
    
    def exception(self, message: str, exc_info: Optional[Exception] = None, **kwargs) -> None:
        """
        Log an exception.
        
        Args:
            message: Log message
            exc_info: Exception information
            **kwargs: Additional log data
        """
        if exc_info:
            kwargs["exception"] = str(exc_info)
            kwargs["exception_type"] = type(exc_info).__name__
        
        if self.logger.isEnabledFor(logging.ERROR):
            self.logger.exception(self._format_log("ERROR", message, **kwargs))


# Create a default structured logger
structured_logger = StructuredLogger()


def with_correlation_id(func: F) -> F:
    """
    Decorator to add a correlation ID to the function context.
    
    This helps with tracing requests across multiple components.
    
    Args:
        func: Function to decorate
        
    Returns:
        Decorated function
    """
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        # Generate a correlation ID if not provided
        correlation_id = kwargs.pop("correlation_id", str(uuid.uuid4()))
        
        # Create a logger with the correlation ID
        logger = structured_logger.with_context(correlation_id=correlation_id)
        
        # Log the function call
        logger.debug(
            f"Calling {func.__name__}",
            args=str(args),
            kwargs=str(kwargs),
        )
        
        try:
            # Call the function
            result = func(*args, **kwargs, correlation_id=correlation_id)
            
            # Log the function result
            logger.debug(
                f"Completed {func.__name__}",
                result_type=type(result).__name__,
            )
            
            return result
        except Exception as e:
            # Log the exception
            logger.exception(
                f"Error in {func.__name__}",
                exc_info=e,
            )
            raise
    
    return cast(F, wrapper)


def trace_api_call(endpoint: str) -> Callable[[F], F]:
    """
    Decorator to trace an API call.
    
    This combines metrics collection, structured logging, and correlation IDs.
    
    Args:
        endpoint: API endpoint name
        
    Returns:
        Decorated function
    """
    def decorator(func: F) -> F:
        # Apply decorators in order: correlation ID, metrics, function
        return with_correlation_id(MetricsCollector.api_timer(endpoint)(func))
    
    return decorator
