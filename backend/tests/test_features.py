import unittest
import os
import shutil
from pathlib import Path
from backend.src.scanner.table_scanner import TableScanner
from backend.src.scanner.sample_data_generator import SampleDataGenerator
from backend.src.scanner.spark_config import get_spark_session

class TestFeatures(unittest.TestCase):
    def setUp(self):
        self.test_data_dir = Path("data/test")
        self.test_data_dir.mkdir(parents=True, exist_ok=True)
        self.parquet_path = self.test_data_dir / "test_parquet"
        self.csv_path = self.test_data_dir / "test_csv"
        
        # Initialize components
        self.data_generator = SampleDataGenerator()
        self.table_scanner = TableScanner()
        
        # Generate and save sample data
        df = self.data_generator.generate_sample_data()
        self.data_generator.save_as_parquet(df, str(self.parquet_path))
        self.data_generator.save_as_csv(df, str(self.csv_path))

    def tearDown(self):
        # Clean up test data
        if self.test_data_dir.exists():
            shutil.rmtree(self.test_data_dir)
        # Stop Spark sessions
        if hasattr(self, 'data_generator') and hasattr(self.data_generator, 'spark'):
            self.data_generator.spark.stop()
        if hasattr(self, 'table_scanner') and hasattr(self.table_scanner, 'spark'):
            self.table_scanner.spark.stop()

    def test_sample_data_generation(self):
        """Test sample data generation"""
        df = self.data_generator.generate_sample_data()
        self.assertEqual(df.count(), 5)  # Check number of rows
        self.assertEqual(len(df.columns), 5)  # Check number of columns
        
        # Verify column names
        expected_columns = {'id', 'name', 'age', 'email', 'created_at'}
        self.assertEqual(set(df.columns), expected_columns)

    def test_parquet_scanning(self):
        """Test scanning Parquet files"""
        # Get the actual parquet file path (Spark creates a directory)
        parquet_file = str(self.parquet_path)
        
        # Scan the parquet table
        result = self.table_scanner.scan_table(parquet_file, "parquet")
        
        # Verify the scan results
        self.assertIn('schema', result)
        self.assertIn('row_count', result)
        self.assertIn('column_stats', result)
        
        self.assertEqual(result['row_count'], 5)
        self.assertEqual(len(result['schema']), 5)
        
        # Check column statistics
        stats = result['column_stats']
        self.assertIn('age', stats)
        self.assertIn('count', stats['age'])
        self.assertIn('mean', stats['age'])

    def test_csv_scanning(self):
        """Test scanning CSV files"""
        # Get the CSV file path
        csv_file = str(self.csv_path)
        
        # Scan the CSV table
        result = self.table_scanner.scan_table(csv_file, "csv")
        
        # Verify the scan results
        self.assertIn('schema', result)
        self.assertIn('row_count', result)
        self.assertIn('column_stats', result)
        
        self.assertEqual(result['row_count'], 5)
        self.assertEqual(len(result['schema']), 5)
        
        # Check column statistics
        stats = result['column_stats']
        self.assertIn('age', stats)
        self.assertIn('count', stats['age'])
        self.assertIn('mean', stats['age'])

    def test_invalid_path(self):
        """Test handling of invalid file paths"""
        with self.assertRaises(ValueError):
            self.table_scanner.scan_table("/nonexistent/path", "parquet")

    def test_invalid_format(self):
        """Test handling of invalid format types"""
        with self.assertRaises(ValueError):
            self.table_scanner.scan_table(str(self.parquet_path), "invalid_format")

if __name__ == '__main__':
    unittest.main() 