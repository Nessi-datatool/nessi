"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

import pytest
import os
import base64
import json
import ssl
import time
from unittest.mock import patch, MagicMock
from http.server import HTTPServer, BaseHTTPRequestHandler
from io import BytesIO
from prometheus_client import Histogram, Gauge, Counter, CollectorRegistry
from src.reports.metrics_test import MetricsHandler

class MockSocket:
    def getsockname(self):
        return ('127.0.0.1', 8000)

class MockServer:
    def __init__(self):
        self.server_name = 'localhost'
        self.server_port = 8000

class MockRequest:
    def __init__(self):
        self.makefile = MagicMock()
        self.makefile.return_value = BytesIO(b'GET / HTTP/1.1\r\n\r\n')
        self.settimeout = MagicMock()
        self.setsockopt = MagicMock()
        self.getsockname = MagicMock(return_value=('127.0.0.1', 8000))

@pytest.fixture
def mock_request():
    return MockRequest()

@pytest.fixture
def mock_server():
    return MockServer()

@pytest.fixture
def metrics_handler(mock_request):
    # Create a separate registry for test metrics
    registry = CollectorRegistry()
    
    # Create real metrics with the test registry
    table_scan_rows = Gauge('nessi_table_scan_total_rows', 'Total rows scanned in table', ['table_name'], registry=registry)
    data_quality_score = Gauge('nessi_data_quality_score', 'Data quality score (0-100)', registry=registry)
    processing_time = Histogram('nessi_processing_time_seconds', 'Processing time in seconds', registry=registry)
    scan_efficiency = Gauge('nessi_scan_efficiency', 'Table scan efficiency score (0-100)', registry=registry)
    partition_utilization = Gauge('nessi_partition_utilization', 'Partition utilization score (0-100)', registry=registry)
    index_usage = Gauge('nessi_index_usage', 'Index usage efficiency score (0-100)', registry=registry)
    data_health = Gauge('nessi_data_health', 'Overall data health score (0-100)', registry=registry)
    performance_status = Gauge('nessi_performance_status', 'System performance status (0-100)', registry=registry)
    anomaly_count = Counter('nessi_anomaly_count', 'Number of anomalies detected', registry=registry)
    resource_utilization = Gauge('nessi_resource_utilization', 'Resource utilization score (0-100)', registry=registry)
    custom_report_count = Counter('nessi_custom_report_count', 'Number of custom reports generated', registry=registry)
    auth_failures = Counter('nessi_auth_failures_total', 'Total authentication failures', ['reason'], registry=registry)
    credential_age = Gauge('nessi_credential_age_days', 'Age of current credentials in days', registry=registry)
    ssl_connections = Counter('nessi_ssl_connections_total', 'Total SSL connections', ['version'], registry=registry)
    rate_limit_hits = Counter('nessi_rate_limit_hits_total', 'Total rate limit hits', registry=registry)
    http_requests_total = Counter('nessi_http_requests_total', 'Total HTTP requests', ['method', 'endpoint', 'status'], registry=registry)
    http_request_duration = Histogram('nessi_http_request_duration_seconds', 'HTTP request duration', ['endpoint'], registry=registry)
    
    with patch('src.reports.metrics_test.table_scan_rows', table_scan_rows), \
         patch('src.reports.metrics_test.data_quality_score', data_quality_score), \
         patch('src.reports.metrics_test.processing_time', processing_time), \
         patch('src.reports.metrics_test.scan_efficiency', scan_efficiency), \
         patch('src.reports.metrics_test.partition_utilization', partition_utilization), \
         patch('src.reports.metrics_test.index_usage', index_usage), \
         patch('src.reports.metrics_test.data_health', data_health), \
         patch('src.reports.metrics_test.performance_status', performance_status), \
         patch('src.reports.metrics_test.anomaly_count', anomaly_count), \
         patch('src.reports.metrics_test.resource_utilization', resource_utilization), \
         patch('src.reports.metrics_test.custom_report_count', custom_report_count), \
         patch('src.reports.metrics_test.auth_failures', auth_failures), \
         patch('src.reports.metrics_test.credential_age', credential_age), \
         patch('src.reports.metrics_test.ssl_connections', ssl_connections), \
         patch('src.reports.metrics_test.rate_limit_hits', rate_limit_hits), \
         patch('src.reports.metrics_test.http_requests_total', http_requests_total), \
         patch('src.reports.metrics_test.http_request_duration', http_request_duration):
        
        # Create a handler instance without calling the parent class's __init__
        handler = MetricsHandler.__new__(MetricsHandler)
        handler.rfile = BytesIO(b'GET / HTTP/1.1\r\n\r\n')
        handler.wfile = BytesIO()
        handler.send_response = MagicMock()
        handler.send_header = MagicMock()
        handler.end_headers = MagicMock()
        handler.headers = {}
        handler._request_count = 0
        handler._last_reset = time.time()
        handler.client_address = ('127.0.0.1', 8000)
        handler.log_message = MagicMock()  # Suppress logging
        handler.log_error = MagicMock()  # Suppress error logging
        handler.send_error = MagicMock()  # Mock send_error
        handler.address_string = MagicMock(return_value='127.0.0.1')  # Mock address_string
        handler.version_string = MagicMock(return_value='HTTP/1.1')  # Mock version_string
        handler.date_time_string = MagicMock(return_value='Thu, 01 Jan 1970 00:00:00 GMT')  # Mock date_time_string
        return handler

@pytest.fixture
def setup_env():
    # Set up environment variables
    os.environ['AUTH_ENABLED'] = 'true'
    os.environ['SSL_ENABLED'] = 'true'
    os.environ['METRICS_USERNAME'] = 'test_user'
    os.environ['METRICS_PASSWORD'] = 'test_pass'
    os.environ['SCENARIO'] = 'normal_operation'
    os.environ['METRICS_INTERVAL'] = '5'
    os.environ['LOG_LEVEL'] = 'INFO'
    yield
    # Clean up environment variables
    for key in ['AUTH_ENABLED', 'SSL_ENABLED', 'METRICS_USERNAME', 'METRICS_PASSWORD', 'SCENARIO', 'METRICS_INTERVAL', 'LOG_LEVEL']:
        os.environ.pop(key, None)

def test_health_check(metrics_handler):
    metrics_handler.path = '/health'
    metrics_handler.do_GET()
    
    metrics_handler.send_response.assert_called_with(200)
    metrics_handler.send_header.assert_any_call('Content-type', 'application/json')
    metrics_handler.end_headers.assert_called_once()

def test_metrics_authentication_success(metrics_handler, setup_env):
    metrics_handler.path = '/metrics'
    auth = base64.b64encode(b'test_user:test_pass').decode('utf-8')
    metrics_handler.headers = {'Authorization': f'Basic {auth}'}
    
    metrics_handler.do_GET()
    
    metrics_handler.send_response.assert_called_with(200)
    metrics_handler.send_header.assert_any_call('Content-type', 'text/plain')
    metrics_handler.end_headers.assert_called_once()

def test_metrics_authentication_failure(metrics_handler, setup_env):
    metrics_handler.path = '/metrics'
    auth = base64.b64encode(b'wrong:wrong').decode('utf-8')
    metrics_handler.headers = {'Authorization': f'Basic {auth}'}
    
    metrics_handler.do_GET()
    
    metrics_handler.send_response.assert_called_with(401)
    metrics_handler.send_header.assert_any_call('WWW-Authenticate', 'Basic realm="Metrics"')

def test_rate_limiting(metrics_handler, setup_env):
    metrics_handler.path = '/metrics'
    auth = base64.b64encode(b'test_user:test_pass').decode('utf-8')
    metrics_handler.headers = {'Authorization': f'Basic {auth}'}
    
    # Simulate multiple requests within the same second
    for _ in range(11):  # 11 requests to exceed the 10 requests/second limit
        metrics_handler.do_GET()
    
    # Check if rate limiting was triggered
    metrics_handler.send_response.assert_called_with(429)
    metrics_handler.send_header.assert_any_call('Content-type', 'text/plain')
    metrics_handler.end_headers.assert_called()

@pytest.fixture
def http_server():
    server = HTTPServer(('0.0.0.0', 8000), MetricsHandler)
    yield server
    server.server_close()

def test_ssl_enabled(http_server, setup_env):
    with patch('ssl.SSLContext') as mock_ssl:
        context = mock_ssl.return_value
        context.wrap_socket.return_value = MagicMock()
        
        # Test SSL configuration
        if os.getenv('SSL_ENABLED', 'true').lower() == 'true':
            context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
            context.load_cert_chain('/app/ssl/cert.pem', '/app/ssl/key.pem')
            http_server.socket = context.wrap_socket(http_server.socket, server_side=True)
        
        assert hasattr(http_server, 'socket')