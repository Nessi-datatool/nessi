from abc import ABC, abstractmethod
from typing import Dict, Any, Optional, List
import pandas as pd
import pyarrow as pa
import pyarrow.parquet as pq
import json
import logging
from pathlib import Path
from delta.tables import DeltaTable
from pyspark.sql import SparkSession

logger = logging.getLogger(__name__)

class DataFormatHandler(ABC):
    """Abstract base class for data format handlers."""
    
    @abstractmethod
    def read(self, path: str, **kwargs) -> pd.DataFrame:
        """Read data from the specified path."""
        pass
    
    @abstractmethod
    def write(self, df: pd.DataFrame, path: str, **kwargs) -> None:
        """Write data to the specified path."""
        pass
    
    @abstractmethod
    def infer_schema(self, path: str) -> Dict[str, Any]:
        """Infer schema from the data at the specified path."""
        pass

class DeltaHandler(DataFormatHandler):
    """Handler for Delta Lake format."""
    
    def __init__(self, spark: SparkSession):
        self.spark = spark
    
    def read(self, path: str, **kwargs) -> pd.DataFrame:
        """Read Delta Lake table."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            df = delta_table.toDF()
            return df.toPandas()
        except Exception as e:
            logger.error(f"Error reading Delta table: {str(e)}")
            raise
    
    def write(self, df: pd.DataFrame, path: str, **kwargs) -> None:
        """Write to Delta Lake table."""
        try:
            spark_df = self.spark.createDataFrame(df)
            spark_df.write.format("delta").mode("overwrite").save(path)
        except Exception as e:
            logger.error(f"Error writing Delta table: {str(e)}")
            raise
    
    def infer_schema(self, path: str) -> Dict[str, Any]:
        """Infer schema from Delta Lake table."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            schema = delta_table.toDF().schema
            return {
                "fields": [
                    {
                        "name": field.name,
                        "type": str(field.dataType),
                        "nullable": field.nullable
                    }
                    for field in schema.fields
                ]
            }
        except Exception as e:
            logger.error(f"Error inferring Delta schema: {str(e)}")
            raise

class ParquetHandler(DataFormatHandler):
    """Handler for Parquet format."""
    
    def read(self, path: str, **kwargs) -> pd.DataFrame:
        """Read Parquet file."""
        try:
            return pd.read_parquet(path, **kwargs)
        except Exception as e:
            logger.error(f"Error reading Parquet file: {str(e)}")
            raise
    
    def write(self, df: pd.DataFrame, path: str, **kwargs) -> None:
        """Write to Parquet file."""
        try:
            df.to_parquet(path, **kwargs)
        except Exception as e:
            logger.error(f"Error writing Parquet file: {str(e)}")
            raise
    
    def infer_schema(self, path: str) -> Dict[str, Any]:
        """Infer schema from Parquet file."""
        try:
            schema = pq.read_schema(path)
            return {
                "fields": [
                    {
                        "name": field.name,
                        "type": str(field.type),
                        "nullable": field.nullable
                    }
                    for field in schema
                ]
            }
        except Exception as e:
            logger.error(f"Error inferring Parquet schema: {str(e)}")
            raise

class CSVHandler(DataFormatHandler):
    """Handler for CSV format."""
    
    def read(self, path: str, **kwargs) -> pd.DataFrame:
        """Read CSV file."""
        try:
            return pd.read_csv(path, **kwargs)
        except Exception as e:
            logger.error(f"Error reading CSV file: {str(e)}")
            raise
    
    def write(self, df: pd.DataFrame, path: str, **kwargs) -> None:
        """Write to CSV file."""
        try:
            df.to_csv(path, **kwargs)
        except Exception as e:
            logger.error(f"Error writing CSV file: {str(e)}")
            raise
    
    def infer_schema(self, path: str) -> Dict[str, Any]:
        """Infer schema from CSV file."""
        try:
            df = pd.read_csv(path, nrows=1000)  # Read first 1000 rows for schema inference
            return {
                "fields": [
                    {
                        "name": col,
                        "type": str(df[col].dtype),
                        "nullable": True  # CSV doesn't enforce nullability
                    }
                    for col in df.columns
                ]
            }
        except Exception as e:
            logger.error(f"Error inferring CSV schema: {str(e)}")
            raise

class DataFormatManager:
    """Manager for handling different data formats."""
    
    def __init__(self, spark: Optional[SparkSession] = None):
        self.handlers = {
            "delta": DeltaHandler(spark) if spark else None,
            "parquet": ParquetHandler(),
            "csv": CSVHandler()
        }
    
    def get_handler(self, format: str) -> DataFormatHandler:
        """Get handler for the specified format."""
        handler = self.handlers.get(format.lower())
        if not handler:
            raise ValueError(f"Unsupported format: {format}")
        if format.lower() == "delta" and not handler:
            raise ValueError("Spark session required for Delta format")
        return handler
    
    def read(self, path: str, format: str, **kwargs) -> pd.DataFrame:
        """Read data from the specified path and format."""
        handler = self.get_handler(format)
        return handler.read(path, **kwargs)
    
    def write(self, df: pd.DataFrame, path: str, format: str, **kwargs) -> None:
        """Write data to the specified path and format."""
        handler = self.get_handler(format)
        handler.write(df, path, **kwargs)
    
    def infer_schema(self, path: str, format: str) -> Dict[str, Any]:
        """Infer schema from the specified path and format."""
        handler = self.get_handler(format)
        return handler.infer_schema(path)
    
    def get_supported_formats(self) -> list:
        """Get list of supported formats."""
        return list(self.handlers.keys()) 