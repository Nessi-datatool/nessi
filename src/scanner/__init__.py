"""Scanner module for analyzing data files."""

from ..decorators import require_valid_license

@require_valid_license
def scan_table(table_path, format='delta'):
    """Scan a table and generate a report."""
    # Implementation will be added later
    pass

@require_valid_license
def generate_sample_data(output_path, format='delta', rows=1000):
    """Generate sample data for testing."""
    # Implementation will be added later
    pass

__all__ = ['scan_table', 'generate_sample_data'] 