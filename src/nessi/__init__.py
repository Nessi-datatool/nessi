"""Nessi package."""

from .spark_config import get_spark_session
from .decorators import require_valid_license, LicenseError
from .license import LicenseValidator
from .sample_data_generator import generate_sample_data, save_as_parquet, save_as_csv
from .create_test_delta_table import create_test_delta_table, add_column_to_delta_table
from .scanner.table_scanner import TableScanner

__version__ = "0.1.0"
__all__ = ['scanner', 'LicenseValidator', 'require_valid_license', 'LicenseError'] 