from typing import List, Dict, Any, Optional
import logging
from pyspark.sql import DataFrame, SparkSession
from pyspark.sql.functions import col, when, lit
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, DoubleType
import pandas as pd
import numpy as np

class Transformation:
    """Base class for data transformations."""
    
    def __init__(self, name: str):
        self.name = name
        self.logger = logging.getLogger(__name__)
    
    def apply(self, df: DataFrame) -> DataFrame:
        """Apply the transformation to the DataFrame."""
        raise NotImplementedError

class FilterTransformation(Transformation):
    """Filter rows based on conditions."""
    
    def __init__(self, condition: str):
        super().__init__("filter")
        self.condition = condition
    
    def apply(self, df: DataFrame) -> DataFrame:
        return df.filter(self.condition)

class RenameTransformation(Transformation):
    """Rename columns."""
    
    def __init__(self, column_mapping: Dict[str, str]):
        super().__init__("rename")
        self.column_mapping = column_mapping
    
    def apply(self, df: DataFrame) -> DataFrame:
        for old_name, new_name in self.column_mapping.items():
            df = df.withColumnRenamed(old_name, new_name)
        return df

class TypeCastTransformation(Transformation):
    """Cast columns to specified types."""
    
    def __init__(self, type_mapping: Dict[str, str]):
        super().__init__("type_cast")
        self.type_mapping = type_mapping
    
    def apply(self, df: DataFrame) -> DataFrame:
        for column, target_type in self.type_mapping.items():
            df = df.withColumn(column, col(column).cast(target_type))
        return df

class FillNullTransformation(Transformation):
    """Fill null values with specified values."""
    
    def __init__(self, fill_values: Dict[str, Any]):
        super().__init__("fill_null")
        self.fill_values = fill_values
    
    def apply(self, df: DataFrame) -> DataFrame:
        for column, value in self.fill_values.items():
            df = df.fillna({column: value})
        return df

class Pipeline:
    """Data transformation pipeline."""
    
    def __init__(self, name: str):
        self.name = name
        self.transformations: List[Transformation] = []
        self.logger = logging.getLogger(__name__)
    
    def add_transformation(self, transformation: Transformation) -> None:
        """Add a transformation to the pipeline."""
        self.transformations.append(transformation)
    
    def apply(self, df: DataFrame) -> DataFrame:
        """Apply all transformations in sequence."""
        for transformation in self.transformations:
            self.logger.info(f"Applying transformation: {transformation.name}")
            df = transformation.apply(df)
        return df
    
    def validate(self, df: DataFrame) -> bool:
        """Validate the pipeline against a DataFrame."""
        try:
            self.apply(df)
            return True
        except Exception as e:
            self.logger.error(f"Pipeline validation failed: {str(e)}")
            return False
    
    def to_dict(self) -> Dict:
        """Convert pipeline to dictionary representation."""
        return {
            "name": self.name,
            "transformations": [t.name for t in self.transformations]
        }

class PipelineManager:
    """Manages multiple transformation pipelines."""
    
    def __init__(self):
        self.pipelines: Dict[str, Pipeline] = {}
        self.logger = logging.getLogger(__name__)
    
    def create_pipeline(self, name: str) -> Pipeline:
        """Create a new pipeline."""
        if name in self.pipelines:
            raise ValueError(f"Pipeline {name} already exists")
        
        pipeline = Pipeline(name)
        self.pipelines[name] = pipeline
        return pipeline
    
    def get_pipeline(self, name: str) -> Optional[Pipeline]:
        """Get a pipeline by name."""
        return self.pipelines.get(name)
    
    def delete_pipeline(self, name: str) -> None:
        """Delete a pipeline."""
        if name in self.pipelines:
            del self.pipelines[name]
    
    def list_pipelines(self) -> List[Dict]:
        """List all pipelines."""
        return [p.to_dict() for p in self.pipelines.values()] 