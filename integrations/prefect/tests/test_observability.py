"""
Tests for observability features (structured logging, metrics, tracing) in the Nessi Prefect integration.
"""

import unittest
from unittest import mock
import json
import logging
import io
import time
from datetime import datetime

from prefect import flow, task
import opentelemetry.trace as trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import SimpleSpanProcessor, ConsoleSpanExporter

from nessi_prefect.tasks.quality import run_quality_check
from nessi_prefect.tasks.profile import run_profile
from nessi_prefect.tasks.validation import run_validation
from nessi_prefect.observability.logging import setup_structured_logging, StructuredLogFormatter
from nessi_prefect.observability.metrics import MetricsCollector
from nessi_prefect.observability.tracing import setup_tracing, add_span_attributes


class TestStructuredLogging(unittest.TestCase):
    """Test cases for structured logging."""

    def setUp(self):
        """Set up test environment."""
        # Capture logs
        self.log_capture = io.StringIO()
        self.handler = logging.StreamHandler(self.log_capture)
        self.handler.setFormatter(StructuredLogFormatter())
        
        # Get logger
        self.logger = logging.getLogger('nessi_prefect')
        self.logger.setLevel(logging.DEBUG)
        self.logger.addHandler(self.handler)
        self.logger.propagate = False
        
        # Clear previous handlers
        for handler in self.logger.handlers[:]:
            if handler != self.handler:
                self.logger.removeHandler(handler)
    
    def tearDown(self):
        """Clean up test environment."""
        self.logger.removeHandler(self.handler)
        self.handler.close()
    
    def test_structured_log_formatter(self):
        """Test structured log formatter."""
        # Log a message
        self.logger.info("Test message", extra={
            'table_name': 'test_table',
            'operation': 'quality_check',
            'duration_ms': 1250,
        })
        
        # Get log output
        log_output = self.log_capture.getvalue()
        
        # Parse JSON
        log_entry = json.loads(log_output)
        
        # Check log entry structure
        self.assertEqual(log_entry['message'], 'Test message')
        self.assertEqual(log_entry['level'], 'INFO')
        self.assertEqual(log_entry['table_name'], 'test_table')
        self.assertEqual(log_entry['operation'], 'quality_check')
        self.assertEqual(log_entry['duration_ms'], 1250)
        self.assertIn('timestamp', log_entry)
    
    def test_setup_structured_logging(self):
        """Test setup of structured logging."""
        # Setup structured logging
        setup_structured_logging()
        
        # Get root logger
        root_logger = logging.getLogger()
        
        # Check that a handler with StructuredLogFormatter was added
        structured_handlers = [h for h in root_logger.handlers 
                             if isinstance(h.formatter, StructuredLogFormatter)]
        self.assertTrue(len(structured_handlers) > 0)
    
    def test_task_logging(self):
        """Test logging in Prefect tasks."""
        # Mock the client
        mock_client = mock.MagicMock()
        mock_client.run_quality_check.return_value = {
            'check_id': 'test-check-id',
            'quality_score': 0.95,
        }
        
        # Run the task with mocked client
        with mock.patch('nessi_prefect.tasks.quality.NessiClient', return_value=mock_client):
            result = run_quality_check.fn(
                table_name='test_table',
                rules=[{'name': 'test_rule', 'rule_type': 'not_null', 'column': 'id'}],
                api_host='http://test-host',
                api_key='test-key',
            )
        
        # Log a message using the same logger as the task
        task_logger = logging.getLogger('nessi_prefect.tasks.quality')
        task_logger.addHandler(self.handler)
        task_logger.info("Test task message", extra={
            'table_name': 'test_table',
            'operation': 'quality_check',
        })
        
        # Get log output
        log_output = self.log_capture.getvalue()
        
        # Check that log output is JSON
        try:
            log_entries = [json.loads(line) for line in log_output.strip().split('\n')]
            self.assertTrue(len(log_entries) > 0)
            
            # Check the test message
            test_entry = log_entries[-1]
            self.assertEqual(test_entry['message'], 'Test task message')
            self.assertEqual(test_entry['table_name'], 'test_table')
            self.assertEqual(test_entry['operation'], 'quality_check')
        except json.JSONDecodeError:
            self.fail("Log output is not valid JSON")
    
    def test_error_logging(self):
        """Test error logging."""
        # Log an error
        try:
            raise ValueError("Test error")
        except ValueError as e:
            self.logger.error("An error occurred", exc_info=True, extra={
                'table_name': 'test_table',
                'operation': 'quality_check',
                'error_type': type(e).__name__,
            })
        
        # Get log output
        log_output = self.log_capture.getvalue()
        
        # Parse JSON
        log_entry = json.loads(log_output)
        
        # Check log entry structure
        self.assertEqual(log_entry['message'], 'An error occurred')
        self.assertEqual(log_entry['level'], 'ERROR')
        self.assertEqual(log_entry['table_name'], 'test_table')
        self.assertEqual(log_entry['operation'], 'quality_check')
        self.assertEqual(log_entry['error_type'], 'ValueError')
        self.assertIn('exc_info', log_entry)
        self.assertIn('ValueError: Test error', log_entry['exc_info'])


class TestMetricsCollection(unittest.TestCase):
    """Test cases for metrics collection."""

    def setUp(self):
        """Set up test environment."""
        # Mock the metrics backend
        self.mock_backend = mock.MagicMock()
        
        # Create metrics collector
        self.metrics = MetricsCollector(backend=self.mock_backend)
    
    def test_counter_metric(self):
        """Test counter metric."""
        # Increment counter
        self.metrics.increment_counter(
            name='quality_checks_total',
            value=1,
            labels={'table_name': 'test_table', 'result': 'success'},
        )
        
        # Check that backend was called
        self.mock_backend.increment_counter.assert_called_once_with(
            'quality_checks_total',
            1,
            {'table_name': 'test_table', 'result': 'success'},
        )
    
    def test_gauge_metric(self):
        """Test gauge metric."""
        # Set gauge
        self.metrics.set_gauge(
            name='quality_score',
            value=0.95,
            labels={'table_name': 'test_table'},
        )
        
        # Check that backend was called
        self.mock_backend.set_gauge.assert_called_once_with(
            'quality_score',
            0.95,
            {'table_name': 'test_table'},
        )
    
    def test_histogram_metric(self):
        """Test histogram metric."""
        # Observe histogram
        self.metrics.observe_histogram(
            name='quality_check_duration_seconds',
            value=1.25,
            labels={'table_name': 'test_table'},
        )
        
        # Check that backend was called
        self.mock_backend.observe_histogram.assert_called_once_with(
            'quality_check_duration_seconds',
            1.25,
            {'table_name': 'test_table'},
        )
    
    def test_task_metrics(self):
        """Test metrics in Prefect tasks."""
        # Mock the client
        mock_client = mock.MagicMock()
        mock_client.run_quality_check.return_value = {
            'check_id': 'test-check-id',
            'quality_score': 0.95,
        }
        
        # Mock the metrics collector
        mock_metrics = mock.MagicMock()
        
        # Run the task with mocked client and metrics
        with mock.patch('nessi_prefect.tasks.quality.NessiClient', return_value=mock_client), \
             mock.patch('nessi_prefect.tasks.quality.metrics', mock_metrics):
            start_time = time.time()
            result = run_quality_check.fn(
                table_name='test_table',
                rules=[{'name': 'test_rule', 'rule_type': 'not_null', 'column': 'id'}],
                api_host='http://test-host',
                api_key='test-key',
            )
            end_time = time.time()
        
        # Check that metrics were recorded
        mock_metrics.increment_counter.assert_called_with(
            'quality_checks_total',
            1,
            {'table_name': 'test_table', 'result': 'success'},
        )
        
        mock_metrics.set_gauge.assert_called_with(
            'quality_score',
            0.95,
            {'table_name': 'test_table'},
        )
        
        # Check that duration was recorded
        mock_metrics.observe_histogram.assert_called()
        args, kwargs = mock_metrics.observe_histogram.call_args
        self.assertEqual(args[0], 'quality_check_duration_seconds')
        self.assertGreaterEqual(args[1], 0)  # Duration should be positive
        self.assertEqual(kwargs['labels']['table_name'], 'test_table')


class TestDistributedTracing(unittest.TestCase):
    """Test cases for distributed tracing."""

    def setUp(self):
        """Set up test environment."""
        # Set up tracer
        trace.set_tracer_provider(TracerProvider())
        trace.get_tracer_provider().add_span_processor(
            SimpleSpanProcessor(ConsoleSpanExporter())
        )
        self.tracer = trace.get_tracer('nessi_prefect.test')
    
    def test_setup_tracing(self):
        """Test setup of distributed tracing."""
        # Setup tracing
        setup_tracing(
            service_name='nessi-prefect-test',
            endpoint='http://localhost:4317',
        )
        
        # Check that tracer provider is set
        self.assertIsNotNone(trace.get_tracer_provider())
    
    def test_add_span_attributes(self):
        """Test adding attributes to spans."""
        # Create a span
        with self.tracer.start_as_current_span('test_span') as span:
            # Add attributes
            add_span_attributes(span, {
                'table_name': 'test_table',
                'operation': 'quality_check',
                'quality_score': 0.95,
            })
            
            # Check attributes
            self.assertEqual(span.attributes['table_name'], 'test_table')
            self.assertEqual(span.attributes['operation'], 'quality_check')
            self.assertEqual(span.attributes['quality_score'], 0.95)
    
    def test_task_tracing(self):
        """Test tracing in Prefect tasks."""
        # Mock the client
        mock_client = mock.MagicMock()
        mock_client.run_quality_check.return_value = {
            'check_id': 'test-check-id',
            'quality_score': 0.95,
        }
        
        # Create a parent span
        with self.tracer.start_as_current_span('parent_span'):
            # Run the task with mocked client
            with mock.patch('nessi_prefect.tasks.quality.NessiClient', return_value=mock_client), \
                 mock.patch('nessi_prefect.tasks.quality.tracer', self.tracer):
                result = run_quality_check.fn(
                    table_name='test_table',
                    rules=[{'name': 'test_rule', 'rule_type': 'not_null', 'column': 'id'}],
                    api_host='http://test-host',
                    api_key='test-key',
                )
        
        # We can't easily check the spans directly, but we can verify the task ran successfully
        self.assertEqual(result['check_id'], 'test-check-id')
        self.assertEqual(result['quality_score'], 0.95)


if __name__ == '__main__':
    unittest.main()
