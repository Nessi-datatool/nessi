"""
Format handling functionality for the Nessi monitoring client.
"""

import json
import os
import tempfile
from typing import Dict, List, Optional, Any, Union, BinaryIO, TextIO

import pandas as pd
import pyarrow as pa
import pyarrow.csv as csv
import pyarrow.parquet as pq
import pyarrow.dataset as ds
from deltalake import DeltaTable

from .format_models import (
    SchemaField, 
    Schema, 
    FormatConfig, 
    FormatDetectionResult, 
    DataBatch
)


class FormatHandler:
    """
    Handler for working with multiple data formats.
    
    This class provides functionality for:
    - Detecting file formats (Delta, Parquet, CSV)
    - Inferring schemas from data
    - Reading data from various formats
    - Converting between formats
    """
    
    def __init__(self, config: Optional[FormatConfig] = None):
        """
        Initialize the format handler.
        
        Args:
            config: Configuration for format handling
        """
        self.config = config or FormatConfig()
    
    def detect_format(self, path: str) -> FormatDetectionResult:
        """
        Detect the format of a file or directory.
        
        Args:
            path: Path to the file or directory
            
        Returns:
            FormatDetectionResult with detected format and confidence
        """
        # Check for Delta format (directory with _delta_log)
        if os.path.isdir(path) and os.path.exists(os.path.join(path, "_delta_log")):
            try:
                # Try to open as Delta table
                table = DeltaTable(path)
                schema = self._convert_arrow_schema_to_schema(table.schema())
                return FormatDetectionResult(
                    format="delta",
                    confidence=1.0,
                    schema=schema,
                    metadata={"version": table.version()}
                )
            except Exception:
                pass
        
        # Check for Parquet file
        if os.path.isfile(path) and path.endswith(".parquet"):
            try:
                # Try to open as Parquet file
                parquet_file = pq.ParquetFile(path)
                schema = self._convert_arrow_schema_to_schema(parquet_file.schema)
                return FormatDetectionResult(
                    format="parquet",
                    confidence=1.0,
                    schema=schema,
                    metadata=parquet_file.metadata.to_dict() if parquet_file.metadata else {}
                )
            except Exception:
                pass
        
        # Check for CSV file
        if os.path.isfile(path) and (path.endswith(".csv") or path.endswith(".txt")):
            try:
                # Try to infer schema from CSV
                csv_options = csv.ParseOptions(delimiter=self.config.csv_delimiter)
                read_options = csv.ReadOptions(
                    skip_rows=1 if self.config.csv_has_header else 0,
                    column_names=None if self.config.csv_has_header else []
                )
                
                # Read a sample of the CSV to infer schema
                with open(path, 'rb') as f:
                    sample_data = f.read(10240)  # Read first 10KB
                
                with tempfile.NamedTemporaryFile(delete=False) as temp_file:
                    temp_file.write(sample_data)
                    temp_path = temp_file.name
                
                try:
                    table = csv.read_csv(temp_path, read_options=read_options, parse_options=csv_options)
                    schema = self._convert_arrow_schema_to_schema(table.schema)
                    return FormatDetectionResult(
                        format="csv",
                        confidence=0.9,
                        schema=schema,
                        metadata={"has_header": self.config.csv_has_header}
                    )
                finally:
                    os.unlink(temp_path)
            except Exception:
                pass
        
        # Default to unknown format
        return FormatDetectionResult(
            format="unknown",
            confidence=0.0,
            schema=None,
            metadata={}
        )
    
    def infer_schema(self, data: Union[str, BinaryIO, TextIO, List[Dict[str, Any]]]) -> Schema:
        """
        Infer schema from data.
        
        Args:
            data: Data to infer schema from (file path, file object, or list of records)
            
        Returns:
            Inferred schema
        """
        if isinstance(data, str):
            # Data is a file path
            return self.detect_format(data).schema
        
        if isinstance(data, list):
            # Data is a list of records
            if not data:
                return Schema(fields=[])
            
            # Convert to pandas DataFrame for schema inference
            df = pd.DataFrame(data)
            return self._infer_schema_from_dataframe(df)
        
        # Try to read as CSV or Parquet
        try:
            with tempfile.NamedTemporaryFile(delete=False) as temp_file:
                if hasattr(data, 'read'):
                    # Data is a file-like object
                    temp_file.write(data.read())
                else:
                    # Data is a string or bytes
                    temp_file.write(data if isinstance(data, bytes) else data.encode('utf-8'))
                
                temp_path = temp_file.name
            
            try:
                # Try as Parquet first
                try:
                    parquet_file = pq.ParquetFile(temp_path)
                    return self._convert_arrow_schema_to_schema(parquet_file.schema)
                except Exception:
                    pass
                
                # Try as CSV
                try:
                    csv_options = csv.ParseOptions(delimiter=self.config.csv_delimiter)
                    read_options = csv.ReadOptions(
                        skip_rows=1 if self.config.csv_has_header else 0,
                        column_names=None if self.config.csv_has_header else []
                    )
                    table = csv.read_csv(temp_path, read_options=read_options, parse_options=csv_options)
                    return self._convert_arrow_schema_to_schema(table.schema)
                except Exception:
                    pass
                
                # Try to read as JSON
                try:
                    with open(temp_path, 'r') as f:
                        json_data = json.load(f)
                    
                    if isinstance(json_data, list):
                        df = pd.DataFrame(json_data)
                        return self._infer_schema_from_dataframe(df)
                except Exception:
                    pass
                
                # Default to empty schema
                return Schema(fields=[])
            finally:
                os.unlink(temp_path)
        except Exception:
            # Default to empty schema
            return Schema(fields=[])
    
    def read_data(self, path: str) -> DataBatch:
        """
        Read data from a file or directory.
        
        Args:
            path: Path to the file or directory
            
        Returns:
            DataBatch with data and schema
        """
        format_result = self.detect_format(path)
        
        if format_result.format == "delta":
            # Read Delta table
            table = DeltaTable(path)
            arrow_table = table.to_pyarrow_table()
            data = self._convert_arrow_table_to_records(arrow_table)
            return DataBatch(
                data=data,
                schema=format_result.schema,
                format="delta",
                row_count=len(data),
                metadata={"version": table.version()}
            )
        
        elif format_result.format == "parquet":
            # Read Parquet file
            table = pq.read_table(path)
            data = self._convert_arrow_table_to_records(table)
            return DataBatch(
                data=data,
                schema=format_result.schema,
                format="parquet",
                row_count=len(data),
                metadata={}
            )
        
        elif format_result.format == "csv":
            # Read CSV file
            csv_options = csv.ParseOptions(delimiter=self.config.csv_delimiter)
            read_options = csv.ReadOptions(
                skip_rows=1 if self.config.csv_has_header else 0,
                column_names=None if self.config.csv_has_header else []
            )
            table = csv.read_csv(path, read_options=read_options, parse_options=csv_options)
            data = self._convert_arrow_table_to_records(table)
            return DataBatch(
                data=data,
                schema=format_result.schema,
                format="csv",
                row_count=len(data),
                metadata={"has_header": self.config.csv_has_header}
            )
        
        else:
            raise ValueError(f"Unsupported format: {format_result.format}")
    
    def _convert_arrow_schema_to_schema(self, arrow_schema: pa.Schema) -> Schema:
        """Convert Arrow schema to Schema model."""
        fields = []
        
        for field in arrow_schema:
            fields.append(SchemaField(
                name=field.name,
                data_type=str(field.type),
                nullable=field.nullable,
                metadata=field.metadata
            ))
        
        return Schema(
            fields=fields,
            metadata=arrow_schema.metadata if arrow_schema.metadata else {}
        )
    
    def _convert_arrow_table_to_records(self, table: pa.Table) -> List[Dict[str, Any]]:
        """Convert Arrow table to list of records."""
        return table.to_pandas().to_dict(orient='records')
    
    def _infer_schema_from_dataframe(self, df: pd.DataFrame) -> Schema:
        """Infer schema from pandas DataFrame."""
        fields = []
        
        for column_name, dtype in df.dtypes.items():
            data_type = str(dtype)
            
            # Map pandas dtypes to more generic types
            if data_type.startswith('int'):
                data_type = 'int'
            elif data_type.startswith('float'):
                data_type = 'float'
            elif data_type == 'bool':
                data_type = 'boolean'
            elif data_type.startswith('datetime'):
                data_type = 'timestamp'
            else:
                data_type = 'string'
            
            fields.append(SchemaField(
                name=column_name,
                data_type=data_type,
                nullable=df[column_name].isna().any()
            ))
        
        return Schema(fields=fields)
