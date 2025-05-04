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

from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBasic, HTTPBasicCredentials
from pyspark.sql import SparkSession
from typing import Dict, List, Optional
import logging
import secrets
from datetime import datetime
import uvicorn

from nessi.security.ssl_manager import SSLManager
from nessi.security.middleware import SecurityMiddleware
from nessi.security.rate_limiter import RateLimit
from nessi.monitoring.config import MonitoringConfig
from nessi.monitoring.metrics_collector import MetricsCollector
from nessi.security.ip_manager import IPManager
from nessi.security.audit_logger import AuditLogger
from nessi.config.security_config import SecurityConfig

from nessi.scanner.delta_scanner import DeltaTableScanner
from nessi.scanner.parquet_scanner import ParquetTableScanner
from nessi.scanner.csv_scanner import CSVScanner
from nessi.scanner.report_generator import ReportGenerator

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize security components
security_config = SecurityConfig()
ip_manager = IPManager(security_config.ip_allowlist_file)
audit_logger = AuditLogger(security_config.audit_log_dir)

# Create FastAPI app
app = FastAPI(
    title="Nessi.dev",
    description="Data quality monitoring platform",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Add security middleware
app.add_middleware(
    SecurityMiddleware,
    ip_manager=ip_manager,
    audit_logger=audit_logger,
    exclude_paths=security_config.exclude_paths
)

# Initialize SSL manager
ssl_manager = SSLManager(cert_dir="/app/config/ssl")

# Global variables
spark = None
delta_scanner = None
parquet_scanner = None
csv_scanner = None
report_generator = None

# Models
class ScanRequest:
    def __init__(self, table_path: str, format: str):
        self.table_path = table_path
        self.format = format

class SchemaValidationRequest:
    def __init__(self, table_path: str, format: str, expected_schema: List[Dict]):
        self.table_path = table_path
        self.format = format
        self.expected_schema = expected_schema

class QualityCheckRequest:
    def __init__(self, table_path: str, format: str, rules: Dict):
        self.table_path = table_path
        self.format = format
        self.rules = rules

class ProfileRequest:
    def __init__(self, table_path: str, format: str):
        self.table_path = table_path
        self.format = format

# Startup event
@app.on_event("startup")
async def startup_event():
    global spark, delta_scanner, parquet_scanner, csv_scanner, report_generator
    try:
        # Initialize Spark session
        spark = SparkSession.builder \
            .appName("NessiDataQuality") \
            .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
            .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
            .getOrCreate()
        
        # Initialize scanners
        delta_scanner = DeltaTableScanner(spark)
        parquet_scanner = ParquetTableScanner(spark)
        csv_scanner = CSVScanner(spark)
        
        # Initialize report generator
        report_generator = ReportGenerator()
        
        logger.info("Application started successfully")
    except Exception as e:
        logger.error(f"Error during startup: {str(e)}")
        raise

# Shutdown event
@app.on_event("shutdown")
async def shutdown_event():
    global spark
    try:
        if spark:
            spark.stop()
        logger.info("Application shut down successfully")
    except Exception as e:
        logger.error(f"Error during shutdown: {str(e)}")

# Helper function for authentication
def verify_credentials(credentials: HTTPBasicCredentials = Depends(security)):
    correct_username = secrets.compare_digest(credentials.username, "admin")
    correct_password = secrets.compare_digest(credentials.password, "admin")
    if not (correct_username and correct_password):
        raise HTTPException(
            status_code=401,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Basic"},
        )
    return credentials.username

# Health check endpoint
@app.get("/health")
async def health_check():
    return {"status": "healthy"}

# Scan endpoint
@app.post("/scan")
async def scan_data(request: ScanRequest, username: str = Depends(verify_credentials)):
    try:
        if request.format.lower() == "delta":
            result = delta_scanner.scan_table(request.table_path)
        elif request.format.lower() == "parquet":
            result = parquet_scanner.scan_table(request.table_path)
        elif request.format.lower() == "csv":
            result = csv_scanner.scan_table(request.table_path)
        else:
            raise HTTPException(status_code=400, detail="Unsupported format")
        
        # Generate report
        report_path = report_generator.generate_report(result)
        
        return {
            "status": "success",
            "data": result,
            "report_path": report_path
        }
    except Exception as e:
        logger.error(f"Error scanning data: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

# Schema validation endpoint
@app.post("/validate-schema")
async def validate_schema(request: SchemaValidationRequest, username: str = Depends(verify_credentials)):
    try:
        if request.format.lower() == "delta":
            result = delta_scanner.validate_schema(request.table_path, request.expected_schema)
        elif request.format.lower() == "parquet":
            result = parquet_scanner.validate_schema(request.table_path, request.expected_schema)
        elif request.format.lower() == "csv":
            result = csv_scanner.validate_schema(request.table_path, request.expected_schema)
        else:
            raise HTTPException(status_code=400, detail="Unsupported format")
        
        return {
            "status": "success",
            "data": result
        }
    except Exception as e:
        logger.error(f"Error validating schema: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

# Quality check endpoint
@app.post("/check-quality")
async def check_quality(request: QualityCheckRequest, username: str = Depends(verify_credentials)):
    try:
        if request.format.lower() == "delta":
            result = delta_scanner.check_data_quality(request.table_path, request.rules)
        elif request.format.lower() == "parquet":
            result = parquet_scanner.check_data_quality(request.table_path, request.rules)
        elif request.format.lower() == "csv":
            result = csv_scanner.check_data_quality(request.table_path, request.rules)
        else:
            raise HTTPException(status_code=400, detail="Unsupported format")
        
        # Generate quality report
        report_path = report_generator.generate_quality_report(result)
        
        return {
            "status": "success",
            "data": result,
            "report_path": report_path
        }
    except Exception as e:
        logger.error(f"Error checking data quality: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

# Profile endpoint
@app.post("/profile")
async def profile_data(request: ProfileRequest, username: str = Depends(verify_credentials)):
    try:
        if request.format.lower() == "delta":
            result = delta_scanner.scan_table(request.table_path)
        elif request.format.lower() == "parquet":
            result = parquet_scanner.scan_table(request.table_path)
        elif request.format.lower() == "csv":
            result = csv_scanner.scan_table(request.table_path)
        else:
            raise HTTPException(status_code=400, detail="Unsupported format")
        
        # Generate detailed profile report
        report_path = report_generator.generate_profile_report(result)
        
        return {
            "status": "success",
            "data": result,
            "report_path": report_path
        }
    except Exception as e:
        logger.error(f"Error profiling data: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8443,
        ssl_keyfile="/app/config/ssl/key.pem",
        ssl_certfile="/app/config/ssl/cert.pem"
    )