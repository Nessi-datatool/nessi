"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

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
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Test script to verify Prometheus metrics are being exposed correctly."""

import time
import random
import logging
import os
import base64
import json
import ssl
from http.server import HTTPServer, BaseHTTPRequestHandler
from prometheus_client import start_http_server, Gauge, Counter, Histogram
from datetime import datetime

# Configure logging
logging.basicConfig(
    level=os.getenv('LOG_LEVEL', 'INFO'),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Initialize Prometheus metrics
table_scan_rows = Gauge('nessi_table_scan_total_rows', 'Total rows scanned in table', ['table_name'])
data_quality_score = Gauge('nessi_data_quality_score', 'Data quality score (0-100)')
processing_time = Histogram('nessi_processing_time_seconds', 'Processing time in seconds')
custom_report_count = Counter('nessi_custom_report_count', 'Number of custom reports generated')

# Security metrics
auth_failures = Counter('nessi_auth_failures_total', 'Total authentication failures', ['reason'])
credential_age = Gauge('nessi_credential_age_days', 'Age of current credentials in days')
ssl_connections = Counter('nessi_ssl_connections_total', 'Total SSL connections', ['version'])
rate_limit_hits = Counter('nessi_rate_limit_hits_total', 'Total rate limit hits')

# New detailed analysis metrics
scan_efficiency = Gauge('nessi_scan_efficiency', 'Table scan efficiency score (0-100)')
partition_utilization = Gauge('nessi_partition_utilization', 'Partition utilization score (0-100)')
index_usage = Gauge('nessi_index_usage', 'Index usage efficiency score (0-100)')
data_health = Gauge('nessi_data_health', 'Overall data health score (0-100)')
performance_status = Gauge('nessi_performance_status', 'System performance status (0-100)')
anomaly_count = Counter('nessi_anomaly_count', 'Number of anomalies detected')
resource_utilization = Gauge('nessi_resource_utilization', 'Resource utilization score (0-100)')

# Protocol metrics
http_requests_total = Counter('nessi_http_requests_total', 'Total HTTP requests', ['method', 'endpoint', 'status'])
http_request_duration = Histogram('nessi_http_request_duration_seconds', 'HTTP request duration', ['endpoint'])

class MetricsHandler(BaseHTTPRequestHandler):
    protocol_version = 'HTTP/1.1'
    
    def __init__(self, *args, **kwargs):
        self._request_count = 0
        self._last_reset = time.time()
        super().__init__(*args, **kwargs)
    
    def do_GET(self):
        start_time = time.time()
        try:
            # Check rate limiting
            if self.check_rate_limit():
                self.send_response(429)
                self.send_header('Content-type', 'text/plain')
                self.end_headers()
                self.wfile.write(b'Too Many Requests')
                return
            
            if self.path == '/health':
                self.handle_health()
            elif self.path == '/metrics':
                self.handle_metrics()
            else:
                self.send_error(404, 'Not Found')
        except Exception as e:
            logger.error(f"Error handling request: {str(e)}")
            self.send_error(500, 'Internal Server Error')
        finally:
            duration = time.time() - start_time
            http_request_duration.labels(endpoint=self.path).observe(duration)
    
    def check_rate_limit(self):
        """Check if request should be rate limited."""
        current_time = time.time()
        
        # Reset counter if more than 1 second has passed
        if current_time - self._last_reset > 1:
            self._request_count = 0
            self._last_reset = current_time
        
        # Increment request count
        self._request_count += 1
        
        # Check if rate limit exceeded (10 requests per second)
        if self._request_count > 10:
            rate_limit_hits.inc()
            return True
        
        return False
    
    def handle_health(self):
        """Handle health check requests."""
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        response = {
            'status': 'healthy',
            'timestamp': datetime.utcnow().isoformat(),
            'version': '1.0.0',
            'security': {
                'authentication': os.getenv('AUTH_ENABLED', 'true'),
                'ssl': os.getenv('SSL_ENABLED', 'true')
            }
        }
        self.wfile.write(json.dumps(response).encode())
        http_requests_total.labels(method='GET', endpoint='/health', status='200').inc()
    
    def handle_metrics(self):
        """Handle metrics requests with authentication."""
        # Check authentication
        if not self.authenticate():
            return
        
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        
        # Write metrics
        self.write_metric('table_scan_rows', table_scan_rows)
        self.write_metric('data_quality_score', data_quality_score)
        self.write_metric('processing_time_seconds', processing_time)
        self.write_metric('scan_efficiency', scan_efficiency)
        self.write_metric('partition_utilization', partition_utilization)
        self.write_metric('index_usage', index_usage)
        self.write_metric('data_health', data_health)
        self.write_metric('performance_status', performance_status)
        self.write_metric('anomaly_count', anomaly_count)
        self.write_metric('resource_utilization', resource_utilization)
        self.write_metric('custom_report_count', custom_report_count)
        
        # Write security metrics
        self.write_metric('auth_failures_total', auth_failures)
        self.write_metric('credential_age_days', credential_age)
        self.write_metric('ssl_connections_total', ssl_connections)
        self.write_metric('rate_limit_hits_total', rate_limit_hits)
        
        http_requests_total.labels(method='GET', endpoint='/metrics', status='200').inc()
    
    def authenticate(self):
        """Authenticate the request."""
        if not os.getenv('AUTH_ENABLED', 'true').lower() == 'true':
            return True
        
        auth_header = self.headers.get('Authorization')
        if not auth_header or not auth_header.startswith('Basic '):
            self.send_response(401)
            self.send_header('WWW-Authenticate', 'Basic realm="Metrics"')
            self.end_headers()
            auth_failures.labels(reason='missing_auth').inc()
            http_requests_total.labels(method='GET', endpoint='/metrics', status='401').inc()
            return False
        
        try:
            auth = base64.b64decode(auth_header[6:]).decode()
            username, password = auth.split(':')
            
            if (username != os.getenv('METRICS_USERNAME', 'test_user') or
                password != os.getenv('METRICS_PASSWORD', 'test_pass')):
                self.send_response(401)
                self.send_header('WWW-Authenticate', 'Basic realm="Metrics"')
                self.end_headers()
                auth_failures.labels(reason='invalid_credentials').inc()
                http_requests_total.labels(method='GET', endpoint='/metrics', status='401').inc()
                return False
        except Exception as e:
            logger.error(f"Authentication error: {str(e)}")
            self.send_response(401)
            self.send_header('WWW-Authenticate', 'Basic realm="Metrics"')
            self.end_headers()
            auth_failures.labels(reason='auth_error').inc()
            http_requests_total.labels(method='GET', endpoint='/metrics', status='401').inc()
            return False
        
        return True
    
    def write_metric(self, name, metric):
        """Write a metric to the response."""
        self.wfile.write(f'# HELP nessi_{name} {metric._documentation}\n'.encode())
        self.wfile.write(f'# TYPE nessi_{name} {metric._type}\n'.encode())
        
        if metric._type == 'histogram':
            # Handle histogram metrics
            self.wfile.write(f'nessi_{name}_count {metric._count()}\n'.encode())
            self.wfile.write(f'nessi_{name}_sum {metric._sum()}\n'.encode())
            for bucket, value in metric._buckets.items():
                self.wfile.write(f'nessi_{name}_bucket{{le="{bucket}"}} {value}\n'.encode())
        elif hasattr(metric, '_labelnames') and metric._labelnames:
            # Handle labeled metrics
            for labels, value in metric._metrics.items():
                label_str = ','.join(f'{k}="{v}"' for k, v in zip(metric._labelnames, labels))
                self.wfile.write(f'nessi_{name}{{{label_str}}} {value._value}\n'.encode())
        else:
            # Handle simple metrics
            self.wfile.write(f'nessi_{name} {metric._value}\n'.encode())

def update_normal_metrics():
    """Update metrics for normal operation scenario."""
    table_scan_rows.labels(table_name='users').set(random.uniform(1000, 5000))
    table_scan_rows.labels(table_name='orders').set(random.uniform(5000, 10000))
    table_scan_rows.labels(table_name='products').set(random.uniform(1000, 3000))
    data_quality_score.set(random.uniform(85, 95))
    processing_time.observe(random.uniform(0.1, 0.5))
    scan_efficiency.set(random.uniform(80, 90))
    partition_utilization.set(random.uniform(60, 80))
    index_usage.set(random.uniform(70, 85))
    data_health.set(random.uniform(85, 95))
    performance_status.set(random.uniform(80, 90))
    resource_utilization.set(random.uniform(40, 60))

def update_high_load_metrics():
    """Update metrics for high load scenario."""
    table_scan_rows.labels(table_name='users').set(random.uniform(5000, 10000))
    table_scan_rows.labels(table_name='orders').set(random.uniform(10000, 20000))
    table_scan_rows.labels(table_name='products').set(random.uniform(3000, 6000))
    data_quality_score.set(random.uniform(75, 85))
    processing_time.observe(random.uniform(0.5, 1.0))
    scan_efficiency.set(random.uniform(60, 75))
    partition_utilization.set(random.uniform(80, 95))
    index_usage.set(random.uniform(50, 70))
    data_health.set(random.uniform(75, 85))
    performance_status.set(random.uniform(60, 75))
    resource_utilization.set(random.uniform(70, 90))

def update_data_quality_metrics():
    """Update metrics for data quality issues scenario."""
    table_scan_rows.labels(table_name='users').set(random.uniform(2000, 4000))
    table_scan_rows.labels(table_name='orders').set(random.uniform(4000, 8000))
    table_scan_rows.labels(table_name='products').set(random.uniform(1000, 2000))
    data_quality_score.set(random.uniform(50, 70))
    processing_time.observe(random.uniform(0.3, 0.7))
    scan_efficiency.set(random.uniform(50, 65))
    partition_utilization.set(random.uniform(40, 60))
    index_usage.set(random.uniform(30, 50))
    data_health.set(random.uniform(40, 60))
    performance_status.set(random.uniform(50, 65))
    resource_utilization.set(random.uniform(30, 50))

def update_performance_metrics():
    """Update metrics for performance degradation scenario."""
    table_scan_rows.labels(table_name='users').set(random.uniform(3000, 6000))
    table_scan_rows.labels(table_name='orders').set(random.uniform(6000, 12000))
    table_scan_rows.labels(table_name='products').set(random.uniform(2000, 4000))
    data_quality_score.set(random.uniform(65, 80))
    processing_time.observe(random.uniform(1.0, 2.0))
    scan_efficiency.set(random.uniform(40, 60))
    partition_utilization.set(random.uniform(70, 90))
    index_usage.set(random.uniform(40, 60))
    data_health.set(random.uniform(60, 75))
    performance_status.set(random.uniform(40, 60))
    resource_utilization.set(random.uniform(80, 95))

def update_custom_metrics():
    """Update metrics for custom analysis scenario."""
    table_scan_rows.labels(table_name='users').set(random.uniform(1000, 3000))
    table_scan_rows.labels(table_name='orders').set(random.uniform(3000, 6000))
    table_scan_rows.labels(table_name='products').set(random.uniform(1000, 2000))
    data_quality_score.set(random.uniform(90, 100))
    processing_time.observe(random.uniform(0.1, 0.3))
    scan_efficiency.set(random.uniform(90, 100))
    partition_utilization.set(random.uniform(50, 70))
    index_usage.set(random.uniform(80, 95))
    data_health.set(random.uniform(90, 100))
    performance_status.set(random.uniform(90, 100))
    resource_utilization.set(random.uniform(30, 50))

def test_metrics():
    """Generate test metrics for demonstration."""
    logger.info("Starting metrics test server...")
    logger.info("Metrics will be available at: http://localhost:8000/metrics")
    
    # Initialize credential age
    credential_age.set(0)
    
    while True:
        try:
            # Get current scenario from environment variable
            scenario = os.getenv('SCENARIO', 'normal_operation')
            
            if scenario == 'normal_operation':
                update_normal_metrics()
            elif scenario == 'high_load':
                update_high_load_metrics()
            elif scenario == 'data_quality_issues':
                update_data_quality_metrics()
            elif scenario == 'performance_degradation':
                update_performance_metrics()
            elif scenario == 'custom_analysis':
                update_custom_metrics()
            
            # Increment custom report count
            custom_report_count.inc()
            
            # Randomly increment anomaly count (30% chance)
            if random.random() < 0.3:
                anomaly_count.inc()
            
            # Update credential age (increment by 1 day every 24 hours)
            credential_age.inc(1/24)
            
            logger.info(f"Generated metrics for scenario: {scenario}")
            time.sleep(float(os.getenv('METRICS_INTERVAL', '5')))
        except Exception as e:
            logger.error(f"Error generating metrics: {str(e)}")
            time.sleep(5)

if __name__ == "__main__":
    # Start Prometheus metrics server
    server = HTTPServer(('0.0.0.0', 8000), MetricsHandler)
    
    # Configure SSL if enabled
    if os.getenv('SSL_ENABLED', 'true').lower() == 'true':
        context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        context.load_cert_chain('/app/ssl/cert.pem', '/app/ssl/key.pem')
        server.socket = context.wrap_socket(server.socket, server_side=True)
    
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        server.server_close()