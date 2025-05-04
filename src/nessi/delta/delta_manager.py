from typing import Dict, Any, Optional, List, Union
import logging
from datetime import datetime, timedelta
from pyspark.sql import SparkSession, DataFrame
from delta.tables import DeltaTable
import json
from pathlib import Path
from pyspark.sql.types import StructType, DataType, IntegerType, DoubleType, StringType, BooleanType, TimestampType
import pandas as pd

logger = logging.getLogger(__name__)

class DeltaManager:
    """Advanced Delta Lake management with schema evolution and optimization."""
    
    def __init__(self, spark: SparkSession):
        self.spark = spark
        self._setup_delta_config()
    
    def _setup_delta_config(self):
        """Configure Delta Lake settings."""
        self.spark.conf.set("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension")
        self.spark.conf.set("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog")
        self.spark.conf.set("spark.databricks.delta.retentionDurationCheck.enabled", "false")
        self.spark.conf.set("spark.databricks.delta.schema.autoMerge.enabled", "true")
    
    def create_table(self, path: str, schema: str, partition_by: Optional[List[str]] = None):
        """Create a new Delta table with schema."""
        try:
            df = self.spark.createDataFrame([], schema)
            df.write.format("delta").partitionBy(partition_by or []).save(path)
            logger.info(f"Created Delta table at {path}")
        except Exception as e:
            logger.error(f"Error creating Delta table: {str(e)}")
            raise
    
    def time_travel(self, table_path: str, version: Optional[int] = None, 
                   timestamp: Optional[Union[str, datetime]] = None) -> DataFrame:
        """Time travel to a specific version or timestamp of the table.
        
        Args:
            table_path: Path to the Delta table
            version: Version number to travel to
            timestamp: Timestamp to travel to (ISO format string or datetime)
            
        Returns:
            DataFrame at the specified version/timestamp
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        
        if version is not None:
            return delta_table.history(version).select("*")
        elif timestamp is not None:
            if isinstance(timestamp, str):
                timestamp = datetime.fromisoformat(timestamp)
            return delta_table.history(timestamp).select("*")
        else:
            raise ValueError("Either version or timestamp must be specified")
    
    def get_version_history(self, table_path: str) -> List[Dict]:
        """Get the complete version history of a Delta table.
        
        Args:
            table_path: Path to the Delta table
            
        Returns:
            List of version history entries
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        history = delta_table.history().collect()
        
        return [{
            "version": row.version,
            "timestamp": row.timestamp,
            "operation": row.operation,
            "operation_parameters": json.loads(row.operationParameters),
            "operation_metrics": json.loads(row.operationMetrics),
            "user": row.userId,
            "notebook": row.notebookId
        } for row in history]
    
    def restore_version(self, table_path: str, version: int) -> None:
        """Restore a Delta table to a specific version.
        
        Args:
            table_path: Path to the Delta table
            version: Version number to restore to
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        delta_table.restoreToVersion(version)
    
    def optimize_table(self, table_path: str, z_order_by: Optional[List[str]] = None,
                      auto_compact: bool = True) -> Dict:
        """Optimize a Delta table with Z-ordering and compaction.
        
        Args:
            table_path: Path to the Delta table
            z_order_by: List of columns to Z-order by
            auto_compact: Whether to automatically compact small files
            
        Returns:
            Dictionary with optimization results
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        
        # Get current table statistics
        stats = self._get_table_stats(delta_table)
        
        # Optimize with Z-ordering if specified
        if z_order_by:
            delta_table.optimize().executeZOrderBy(z_order_by)
            stats["z_ordered"] = True
            stats["z_order_columns"] = z_order_by
        
        # Auto-compact if enabled
        if auto_compact:
            delta_table.optimize().executeCompaction()
            stats["compacted"] = True
        
        # Get post-optimization statistics
        new_stats = self._get_table_stats(delta_table)
        stats["optimization_metrics"] = {
            "files_before": stats["num_files"],
            "files_after": new_stats["num_files"],
            "size_before": stats["size_bytes"],
            "size_after": new_stats["size_bytes"]
        }
        
        return stats
    
    def _get_table_stats(self, delta_table: DeltaTable) -> Dict:
        """Get detailed statistics about a Delta table.
        
        Args:
            delta_table: DeltaTable object
            
        Returns:
            Dictionary with table statistics
        """
        # Get basic statistics
        stats = {
            "num_files": delta_table.detail().select("numFiles").collect()[0][0],
            "size_bytes": delta_table.detail().select("sizeInBytes").collect()[0][0],
            "partition_columns": delta_table.detail().select("partitionColumns").collect()[0][0],
            "properties": delta_table.detail().select("properties").collect()[0][0]
        }
        
        # Get file statistics
        files = delta_table.detail().select("files").collect()[0][0]
        stats["file_stats"] = {
            "total": len(files),
            "avg_size": sum(f["size"] for f in files) / len(files),
            "min_size": min(f["size"] for f in files),
            "max_size": max(f["size"] for f in files)
        }
        
        # Get partition statistics
        partitions = delta_table.detail().select("partitionStatistics").collect()[0][0]
        stats["partition_stats"] = {
            "num_partitions": len(partitions),
            "avg_rows": sum(p["numRecords"] for p in partitions) / len(partitions),
            "total_rows": sum(p["numRecords"] for p in partitions)
        }
        
        return stats
    
    def get_optimization_hints(self, table_path: str) -> Dict:
        """Get optimization hints for a Delta table.
        
        Args:
            table_path: Path to the Delta table
            
        Returns:
            Dictionary with optimization hints
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        stats = self._get_table_stats(delta_table)
        
        hints = {
            "z_ordering": [],
            "compaction": False,
            "vacuum": False,
            "partitioning": []
        }
        
        # Check if Z-ordering would be beneficial
        if stats["num_files"] > 1000:  # Threshold for Z-ordering
            hints["z_ordering"] = stats["partition_columns"]
        
        # Check if compaction is needed
        if stats["file_stats"]["avg_size"] < 128 * 1024 * 1024:  # 128MB threshold
            hints["compaction"] = True
        
        # Check if vacuum is needed
        if stats["num_files"] > 10000:  # Threshold for vacuum
            hints["vacuum"] = True
        
        # Check if partitioning could be improved
        if stats["partition_stats"]["num_partitions"] < 10:  # Too few partitions
            hints["partitioning"] = ["Consider adding more partition columns"]
        elif stats["partition_stats"]["num_partitions"] > 1000:  # Too many partitions
            hints["partitioning"] = ["Consider reducing partition columns"]
        
        return hints
    
    def get_history(self, path: str) -> List[Dict[str, Any]]:
        """Get the version history of a Delta table."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            return [{
                "version": row.version,
                "timestamp": row.timestamp,
                "operation": row.operation,
                "operationParameters": json.loads(row.operationParameters),
                "operationMetrics": json.loads(row.operationMetrics)
            } for row in delta_table.history().collect()]
        except Exception as e:
            logger.error(f"Error getting history: {str(e)}")
            raise
    
    def vacuum(self, path: str, retention_hours: int = 168):
        """Vacuum old files with configurable retention."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            delta_table.vacuum(retention_hours)
            logger.info(f"Vacuumed Delta table at {path}")
        except Exception as e:
            logger.error(f"Error vacuuming table: {str(e)}")
            raise
    
    def get_metadata(self, path: str) -> Dict[str, Any]:
        """Get detailed metadata about the Delta table."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            detail = delta_table.detail().collect()[0]
            
            return {
                "format": detail.format,
                "id": detail.id,
                "name": detail.name,
                "description": detail.description,
                "location": detail.location,
                "createdAt": detail.createdAt,
                "lastModified": detail.lastModified,
                "partitionColumns": detail.partitionColumns,
                "numFiles": detail.numFiles,
                "sizeInBytes": detail.sizeInBytes,
                "properties": detail.properties,
                "minReaderVersion": detail.minReaderVersion,
                "minWriterVersion": detail.minWriterVersion,
                "currentVersion": detail.currentVersion,
                "schema": detail.schema
            }
        except Exception as e:
            logger.error(f"Error getting metadata: {str(e)}")
            raise
    
    def analyze_table(self, path: str) -> Dict[str, Any]:
        """Analyze table for optimization hints."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            detail = delta_table.detail().collect()[0]
            
            # Get file statistics
            files = self.spark.read.format("delta").load(path).inputFiles()
            file_sizes = [Path(f).stat().st_size for f in files]
            
            # Get partition statistics
            partition_stats = self.spark.sql(f"SELECT * FROM delta.`{path}`").groupBy(
                *detail.partitionColumns
            ).count().collect()
            
            return {
                "file_count": len(files),
                "total_size": sum(file_sizes),
                "avg_file_size": sum(file_sizes) / len(files) if files else 0,
                "partition_stats": [{
                    "partition": dict(zip(detail.partitionColumns, row[:-1])),
                    "count": row[-1]
                } for row in partition_stats],
                "optimization_hints": {
                    "z_order_candidates": self._get_z_order_candidates(path),
                    "partition_candidates": self._get_partition_candidates(path)
                }
            }
        except Exception as e:
            logger.error(f"Error analyzing table: {str(e)}")
            raise
    
    def _get_z_order_candidates(self, path: str) -> List[str]:
        """Get columns that are good candidates for Z-ordering."""
        try:
            df = self.spark.read.format("delta").load(path)
            stats = df.summary()
            
            candidates = []
            for col in df.columns:
                # Check if column has high cardinality and is frequently used in queries
                unique_count = stats.filter(f"summary = 'count'").select(col).collect()[0][0]
                if unique_count > 1000:  # Arbitrary threshold
                    candidates.append(col)
            
            return candidates
        except Exception as e:
            logger.error(f"Error getting Z-order candidates: {str(e)}")
            return []
    
    def _get_partition_candidates(self, path: str) -> List[str]:
        """Get columns that are good candidates for partitioning."""
        try:
            df = self.spark.read.format("delta").load(path)
            stats = df.summary()
            
            candidates = []
            for col in df.columns:
                # Check if column has low cardinality and is frequently used in filters
                unique_count = stats.filter(f"summary = 'count'").select(col).collect()[0][0]
                if unique_count < 100:  # Arbitrary threshold
                    candidates.append(col)
            
            return candidates
        except Exception as e:
            logger.error(f"Error getting partition candidates: {str(e)}")
            return []
    
    def evolve_schema(self, path: str, schema_evolution: Dict[str, Any]):
        """Evolve the schema of a Delta table."""
        try:
            delta_table = DeltaTable.forPath(self.spark, path)
            
            # Apply schema changes
            for col, changes in schema_evolution.items():
                if "type" in changes:
                    delta_table.alterColumn(col, changes["type"])
                if "nullable" in changes:
                    delta_table.alterColumn(col, nullable=changes["nullable"])
                if "comment" in changes:
                    delta_table.alterColumn(col, comment=changes["comment"])
            
            logger.info(f"Evolved schema for Delta table at {path}")
        except Exception as e:
            logger.error(f"Error evolving schema: {str(e)}")
            raise
    
    def get_version_ui(self, table_path: str) -> Dict:
        """Get version control UI data for a Delta table.
        
        Args:
            table_path: Path to the Delta table
            
        Returns:
            Dictionary with version control UI data
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        history = self.get_version_history(table_path)
        
        # Get current schema
        current_schema = delta_table.toDF().schema
        
        # Get version diffs
        version_diffs = []
        for i in range(1, len(history)):
            prev_version = history[i-1]
            curr_version = history[i]
            
            # Get schema changes
            schema_changes = self._get_schema_changes(
                prev_version.get('schema', {}),
                curr_version.get('schema', {})
            )
            
            # Get data changes
            data_changes = self._get_data_changes(
                table_path,
                prev_version['version'],
                curr_version['version']
            )
            
            version_diffs.append({
                'from_version': prev_version['version'],
                'to_version': curr_version['version'],
                'timestamp': curr_version['timestamp'],
                'operation': curr_version['operation'],
                'schema_changes': schema_changes,
                'data_changes': data_changes,
                'user': curr_version['user'],
                'notebook': curr_version['notebook']
            })
        
        return {
            'current_version': history[-1]['version'],
            'total_versions': len(history),
            'schema': str(current_schema),
            'version_diffs': version_diffs,
            'optimization_history': self._get_optimization_history(table_path)
        }

    def _get_schema_changes(self, old_schema: Dict, new_schema: Dict) -> Dict:
        """Get schema changes between versions."""
        changes = {
            'added': [],
            'removed': [],
            'modified': []
        }
        
        # Compare schemas
        old_fields = {f['name']: f for f in old_schema.get('fields', [])}
        new_fields = {f['name']: f for f in new_schema.get('fields', [])}
        
        # Find added fields
        changes['added'] = [name for name in new_fields if name not in old_fields]
        
        # Find removed fields
        changes['removed'] = [name for name in old_fields if name not in new_fields]
        
        # Find modified fields
        for name in set(old_fields.keys()) & set(new_fields.keys()):
            if old_fields[name] != new_fields[name]:
                changes['modified'].append({
                    'field': name,
                    'old': old_fields[name],
                    'new': new_fields[name]
                })
        
        return changes

    def _get_data_changes(self, table_path: str, from_version: int, to_version: int) -> Dict:
        """Get data changes between versions."""
        # Get dataframes for both versions
        df_old = self.time_travel(table_path, version=from_version)
        df_new = self.time_travel(table_path, version=to_version)
        
        # Calculate changes
        changes = {
            'rows_added': df_new.count() - df_old.count(),
            'rows_removed': df_old.count() - df_new.count(),
            'columns_changed': [],
            'summary_stats': {}
        }
        
        # Compare columns
        for col in set(df_old.columns) & set(df_new.columns):
            old_stats = df_old.select(col).summary()
            new_stats = df_new.select(col).summary()
            
            if old_stats != new_stats:
                changes['columns_changed'].append({
                    'column': col,
                    'old_stats': old_stats.collect(),
                    'new_stats': new_stats.collect()
                })
        
        # Calculate summary statistics
        changes['summary_stats'] = {
            'old_row_count': df_old.count(),
            'new_row_count': df_new.count(),
            'old_column_count': len(df_old.columns),
            'new_column_count': len(df_new.columns)
        }
        
        return changes

    def _get_optimization_history(self, table_path: str) -> List[Dict]:
        """Get optimization history for a Delta table."""
        delta_table = DeltaTable.forPath(self.spark, table_path)
        history = delta_table.history().collect()
        
        optimizations = []
        for entry in history:
            if entry.operation == 'OPTIMIZE':
                optimizations.append({
                    'version': entry.version,
                    'timestamp': entry.timestamp,
                    'metrics': json.loads(entry.operationMetrics),
                    'parameters': json.loads(entry.operationParameters)
                })
        
        return optimizations

    def compare_versions(self, table_path: str, version1: int, version2: int) -> Dict:
        """Compare two versions of a Delta table.
        
        Args:
            table_path: Path to the Delta table
            version1: First version to compare
            version2: Second version to compare
            
        Returns:
            Dictionary with comparison results
        """
        # Get dataframes for both versions
        df1 = self.time_travel(table_path, version=version1)
        df2 = self.time_travel(table_path, version=version2)
        
        # Compare schemas
        schema_diff = self._get_schema_changes(
            json.loads(df1.schema.json()),
            json.loads(df2.schema.json())
        )
        
        # Compare data
        data_diff = self._get_data_changes(table_path, version1, version2)
        
        # Get version metadata
        history = self.get_version_history(table_path)
        v1_meta = next(v for v in history if v['version'] == version1)
        v2_meta = next(v for v in history if v['version'] == version2)
        
        return {
            'versions': {
                'from': {
                    'version': version1,
                    'timestamp': v1_meta['timestamp'],
                    'operation': v1_meta['operation']
                },
                'to': {
                    'version': version2,
                    'timestamp': v2_meta['timestamp'],
                    'operation': v2_meta['operation']
                }
            },
            'schema_diff': schema_diff,
            'data_diff': data_diff
        }

    def auto_evolve_schema(self, table_path: str, new_data: DataFrame) -> Dict:
        """Automatically evolve schema based on new data.
        
        Args:
            table_path: Path to the Delta table
            new_data: New data with potential schema changes
            
        Returns:
            Dictionary with evolution results
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        current_schema = delta_table.toDF().schema
        
        # Get schema differences
        schema_diff = self._get_schema_diff(current_schema, new_data.schema)
        
        if not schema_diff['changes']:
            return {'evolved': False, 'message': 'No schema changes needed'}
        
        try:
            # Apply schema changes
            for change in schema_diff['changes']:
                if change['type'] == 'add_column':
                    delta_table.alterColumn(change['name'], change['data_type'])
                elif change['type'] == 'modify_column':
                    delta_table.alterColumn(change['name'], change['new_type'])
                elif change['type'] == 'drop_column':
                    delta_table.alterColumn(change['name'], nullable=True)
            
            return {
                'evolved': True,
                'changes': schema_diff['changes'],
                'message': 'Schema successfully evolved'
            }
        except Exception as e:
            logger.error(f"Error evolving schema: {str(e)}")
            raise

    def validate_schema(self, table_path: str, schema_rules: Dict) -> Dict:
        """Validate schema against rules.
        
        Args:
            table_path: Path to the Delta table
            schema_rules: Dictionary of validation rules
            
        Returns:
            Dictionary with validation results
        """
        delta_table = DeltaTable.forPath(self.spark, table_path)
        schema = delta_table.toDF().schema
        
        validation_results = {
            'valid': True,
            'errors': [],
            'warnings': []
        }
        
        # Validate column types
        for col in schema:
            if col.name in schema_rules.get('required_types', {}):
                expected_type = schema_rules['required_types'][col.name]
                if str(col.dataType) != expected_type:
                    validation_results['valid'] = False
                    validation_results['errors'].append(
                        f"Column {col.name} has type {col.dataType}, expected {expected_type}"
                    )
        
        # Validate required columns
        for req_col in schema_rules.get('required_columns', []):
            if req_col not in [col.name for col in schema]:
                validation_results['valid'] = False
                validation_results['errors'].append(
                    f"Required column {req_col} is missing"
                )
        
        # Validate column constraints
        for col in schema:
            if col.name in schema_rules.get('constraints', {}):
                constraints = schema_rules['constraints'][col.name]
                
                # Check nullable constraint
                if 'nullable' in constraints and col.nullable != constraints['nullable']:
                    validation_results['warnings'].append(
                        f"Column {col.name} nullable constraint mismatch"
                    )
                
                # Check default value
                if 'default' in constraints:
                    default_value = constraints['default']
                    if not self._validate_default_value(col.dataType, default_value):
                        validation_results['errors'].append(
                            f"Invalid default value for column {col.name}"
                        )
        
        return validation_results

    def _get_schema_diff(self, current_schema: StructType, new_schema: StructType) -> Dict:
        """Get differences between current and new schema."""
        current_fields = {f.name: f for f in current_schema.fields}
        new_fields = {f.name: f for f in new_schema.fields}
        
        changes = []
        
        # Find added columns
        for name in set(new_fields.keys()) - set(current_fields.keys()):
            changes.append({
                'type': 'add_column',
                'name': name,
                'data_type': str(new_fields[name].dataType)
            })
        
        # Find modified columns
        for name in set(current_fields.keys()) & set(new_fields.keys()):
            if str(current_fields[name].dataType) != str(new_fields[name].dataType):
                changes.append({
                    'type': 'modify_column',
                    'name': name,
                    'old_type': str(current_fields[name].dataType),
                    'new_type': str(new_fields[name].dataType)
                })
        
        # Find dropped columns
        for name in set(current_fields.keys()) - set(new_fields.keys()):
            changes.append({
                'type': 'drop_column',
                'name': name
            })
        
        return {'changes': changes}

    def _validate_default_value(self, data_type: DataType, value: Any) -> bool:
        """Validate default value against data type."""
        try:
            # Convert value to appropriate type
            if isinstance(data_type, IntegerType):
                int(value)
            elif isinstance(data_type, DoubleType):
                float(value)
            elif isinstance(data_type, StringType):
                str(value)
            elif isinstance(data_type, BooleanType):
                bool(value)
            elif isinstance(data_type, TimestampType):
                pd.to_datetime(value)
            return True
        except (ValueError, TypeError):
            return False

    def auto_optimize_partitions(self, table_path: str) -> Dict[str, Any]:
        """Automatically optimize table partitions based on query patterns.
        
        Args:
            table_path: Path to the Delta table
            
        Returns:
            Dictionary with optimization results
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, table_path)
            current_partitions = delta_table.detail().select("partitionColumns").collect()[0][0]
            
            # Analyze query patterns
            query_patterns = self._analyze_query_patterns(table_path)
            
            # Get partition candidates
            candidates = self._get_partition_candidates(table_path, query_patterns)
            
            # Compare with current partitions
            optimization_needed = set(candidates) != set(current_partitions)
            
            if optimization_needed:
                # Create new table with optimized partitions
                new_path = f"{table_path}_optimized"
                self._repartition_table(delta_table, new_path, candidates)
                
                return {
                    "optimized": True,
                    "old_partitions": current_partitions,
                    "new_partitions": candidates,
                    "new_path": new_path
                }
            
            return {
                "optimized": False,
                "message": "Current partitions are optimal"
            }
            
        except Exception as e:
            logger.error(f"Error optimizing partitions: {str(e)}")
            raise

    def _analyze_query_patterns(self, table_path: str) -> List[Dict[str, Any]]:
        """Analyze query patterns from the table history."""
        try:
            history = self.get_history(table_path)
            query_patterns = []
            
            for entry in history:
                if entry["operation"] == "SELECT":
                    # Extract filter conditions
                    filters = self._extract_filters(entry["operationParameters"])
                    if filters:
                        query_patterns.append({
                            "timestamp": entry["timestamp"],
                            "filters": filters,
                            "frequency": 1
                        })
            
            # Aggregate patterns
            aggregated_patterns = {}
            for pattern in query_patterns:
                key = tuple(sorted(pattern["filters"].items()))
                if key in aggregated_patterns:
                    aggregated_patterns[key]["frequency"] += 1
                else:
                    aggregated_patterns[key] = pattern
            
            return list(aggregated_patterns.values())
            
        except Exception as e:
            logger.error(f"Error analyzing query patterns: {str(e)}")
            return []

    def _extract_filters(self, parameters: Dict[str, Any]) -> Dict[str, Any]:
        """Extract filter conditions from query parameters."""
        filters = {}
        
        if "where" in parameters:
            where_clause = parameters["where"]
            # Parse where clause to extract column filters
            # This is a simplified version - implement proper SQL parsing
            conditions = where_clause.split("AND")
            for condition in conditions:
                parts = condition.strip().split()
                if len(parts) >= 3:
                    column = parts[0]
                    operator = parts[1]
                    value = " ".join(parts[2:])
                    filters[column] = {
                        "operator": operator,
                        "value": value
                    }
        
        return filters

    def _repartition_table(self, delta_table: DeltaTable, new_path: str, 
                          new_partitions: List[str]) -> None:
        """Repartition table with new partition columns."""
        try:
            # Read current data
            df = delta_table.toDF()
            
            # Write with new partitions
            df.write.format("delta") \
                .partitionBy(*new_partitions) \
                .mode("overwrite") \
                .save(new_path)
            
            logger.info(f"Table repartitioned to {new_path}")
            
        except Exception as e:
            logger.error(f"Error repartitioning table: {str(e)}")
            raise

    def dynamic_schema_evolution(self, table_path: str, new_data: DataFrame) -> Dict[str, Any]:
        """Dynamically evolve schema based on new data.
        
        Args:
            table_path: Path to the Delta table
            new_data: New data with potential schema changes
            
        Returns:
            Dictionary with evolution results
        """
        try:
            delta_table = DeltaTable.forPath(self.spark, table_path)
            current_schema = delta_table.toDF().schema
            
            # Compare schemas
            schema_diff = self._compare_schemas(current_schema, new_data.schema)
            
            if schema_diff["changes"]:
                # Apply schema changes
                self._apply_schema_changes(delta_table, schema_diff["changes"])
                
                return {
                    "evolved": True,
                    "changes": schema_diff["changes"],
                    "new_schema": str(new_data.schema)
                }
            
            return {
                "evolved": False,
                "message": "No schema changes needed"
            }
            
        except Exception as e:
            logger.error(f"Error evolving schema: {str(e)}")
            raise

    def _compare_schemas(self, current_schema: StructType, new_schema: StructType) -> Dict[str, Any]:
        """Compare current and new schemas."""
        changes = {
            "added": [],
            "removed": [],
            "modified": []
        }
        
        current_fields = {f.name: f for f in current_schema.fields}
        new_fields = {f.name: f for f in new_schema.fields}
        
        # Find added fields
        for name in set(new_fields.keys()) - set(current_fields.keys()):
            changes["added"].append({
                "name": name,
                "type": str(new_fields[name].dataType),
                "nullable": new_fields[name].nullable
            })
        
        # Find removed fields
        for name in set(current_fields.keys()) - set(new_fields.keys()):
            changes["removed"].append({
                "name": name,
                "type": str(current_fields[name].dataType)
            })
        
        # Find modified fields
        for name in set(current_fields.keys()) & set(new_fields.keys()):
            current_field = current_fields[name]
            new_field = new_fields[name]
            
            if str(current_field.dataType) != str(new_field.dataType):
                changes["modified"].append({
                    "name": name,
                    "old_type": str(current_field.dataType),
                    "new_type": str(new_field.dataType)
                })
            
            if current_field.nullable != new_field.nullable:
                changes["modified"].append({
                    "name": name,
                    "old_nullable": current_field.nullable,
                    "new_nullable": new_field.nullable
                })
        
        return {
            "changes": changes,
            "current_schema": str(current_schema),
            "new_schema": str(new_schema)
        }

    def _apply_schema_changes(self, delta_table: DeltaTable, changes: Dict[str, List[Dict[str, Any]]]) -> None:
        """Apply schema changes to the Delta table."""
        try:
            # Add new columns
            for change in changes["added"]:
                delta_table.alterColumn(
                    change["name"],
                    change["type"],
                    nullable=change["nullable"]
                )
            
            # Modify existing columns
            for change in changes["modified"]:
                if "new_type" in change:
                    delta_table.alterColumn(
                        change["name"],
                        change["new_type"]
                    )
                if "new_nullable" in change:
                    delta_table.alterColumn(
                        change["name"],
                        nullable=change["new_nullable"]
                    )
            
            # Note: Delta Lake doesn't support dropping columns directly
            # Mark removed columns as nullable
            for change in changes["removed"]:
                delta_table.alterColumn(
                    change["name"],
                    nullable=True
                )
            
            logger.info("Schema changes applied successfully")
            
        except Exception as e:
            logger.error(f"Error applying schema changes: {str(e)}")
            raise 