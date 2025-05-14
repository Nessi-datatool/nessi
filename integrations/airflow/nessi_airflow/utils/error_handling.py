"""
Error handling utilities for Nessi Airflow integration.

This module provides utilities for handling errors and implementing retry logic
for API calls to Nessi.dev.
"""

import time
import logging
import functools
from typing import Callable, Any, Dict, Optional, Type, List, Union, Tuple

import requests
from requests.exceptions import RequestException, Timeout, ConnectionError
from airflow.exceptions import AirflowException

logger = logging.getLogger(__name__)

# Error categories
class NessiApiError(AirflowException):
    """Base exception for Nessi API errors."""
    def __init__(self, message: str, status_code: Optional[int] = None, response: Optional[Dict[str, Any]] = None):
        self.status_code = status_code
        self.response = response
        super().__init__(message)


class NessiAuthenticationError(NessiApiError):
    """Exception for authentication errors."""
    pass


class NessiResourceNotFoundError(NessiApiError):
    """Exception for resource not found errors."""
    pass


class NessiServerError(NessiApiError):
    """Exception for server-side errors."""
    pass


class NessiClientError(NessiApiError):
    """Exception for client-side errors."""
    pass


class NessiTimeoutError(NessiApiError):
    """Exception for timeout errors."""
    pass


class NessiConnectionError(NessiApiError):
    """Exception for connection errors."""
    pass


# Error mapping
def map_http_error(status_code: int, response_body: Dict[str, Any]) -> Type[NessiApiError]:
    """
    Maps HTTP status codes to specific error types.
    
    Args:
        status_code: HTTP status code
        response_body: Response body as dictionary
        
    Returns:
        Appropriate error class
    """
    if 400 <= status_code < 500:
        if status_code == 401 or status_code == 403:
            return NessiAuthenticationError
        elif status_code == 404:
            return NessiResourceNotFoundError
        else:
            return NessiClientError
    elif status_code >= 500:
        return NessiServerError
    else:
        return NessiApiError


# Retry decorator with exponential backoff
def with_exponential_backoff(
    max_retries: int = 3,
    initial_backoff: float = 1.0,
    max_backoff: float = 60.0,
    backoff_factor: float = 2.0,
    retry_on_exceptions: Optional[List[Type[Exception]]] = None,
    retry_on_status_codes: Optional[List[int]] = None,
):
    """
    Decorator that implements exponential backoff for retrying functions.
    
    Args:
        max_retries: Maximum number of retries
        initial_backoff: Initial backoff time in seconds
        max_backoff: Maximum backoff time in seconds
        backoff_factor: Factor to multiply backoff time by after each retry
        retry_on_exceptions: List of exceptions to retry on
        retry_on_status_codes: List of HTTP status codes to retry on
        
    Returns:
        Decorated function
    """
    if retry_on_exceptions is None:
        retry_on_exceptions = [
            ConnectionError,
            Timeout,
            NessiServerError,
        ]
        
    if retry_on_status_codes is None:
        retry_on_status_codes = [429, 500, 502, 503, 504]
    
    def decorator(func: Callable) -> Callable:
        @functools.wraps(func)
        def wrapper(*args, **kwargs) -> Any:
            last_exception = None
            backoff_time = initial_backoff
            
            for retry in range(max_retries + 1):
                try:
                    return func(*args, **kwargs)
                except tuple(retry_on_exceptions) as e:
                    last_exception = e
                    if retry == max_retries:
                        logger.error(
                            "Maximum retries reached for %s: %s",
                            func.__name__,
                            str(e),
                        )
                        raise
                    
                    # Check if it's a RequestException with a response
                    if isinstance(e, RequestException) and hasattr(e, 'response') and e.response is not None:
                        status_code = e.response.status_code
                        if status_code not in retry_on_status_codes:
                            # Don't retry if status code is not in retry_on_status_codes
                            raise
                    
                    # Calculate backoff time
                    backoff_time = min(backoff_time * backoff_factor, max_backoff)
                    
                    logger.warning(
                        "Retrying %s in %.2f seconds after error: %s (retry %d/%d)",
                        func.__name__,
                        backoff_time,
                        str(e),
                        retry + 1,
                        max_retries,
                    )
                    
                    time.sleep(backoff_time)
            
            # This should never happen, but just in case
            if last_exception:
                raise last_exception
            
            return None  # This should never be reached
        
        return wrapper
    
    return decorator


# Exception handler for API calls
def handle_api_exceptions(func: Callable) -> Callable:
    """
    Decorator that handles API exceptions and converts them to appropriate error types.
    
    Args:
        func: Function to decorate
        
    Returns:
        Decorated function
    """
    @functools.wraps(func)
    def wrapper(*args, **kwargs) -> Any:
        try:
            return func(*args, **kwargs)
        except requests.exceptions.Timeout as e:
            raise NessiTimeoutError(f"Request timed out: {str(e)}")
        except requests.exceptions.ConnectionError as e:
            raise NessiConnectionError(f"Connection error: {str(e)}")
        except requests.exceptions.RequestException as e:
            # Handle response errors
            if hasattr(e, 'response') and e.response is not None:
                status_code = e.response.status_code
                try:
                    response_body = e.response.json()
                except ValueError:
                    response_body = {"error": e.response.text}
                
                error_class = map_http_error(status_code, response_body)
                error_message = response_body.get('error', str(e))
                
                raise error_class(
                    f"API error ({status_code}): {error_message}",
                    status_code=status_code,
                    response=response_body,
                )
            else:
                # Generic request exception
                raise NessiApiError(f"Request failed: {str(e)}")
    
    return wrapper
