# Data Lineage Visualization and Comparison

This document provides an overview of the Data Lineage Visualization and Comparison feature in Nessi.dev.

## Overview

Data lineage is a critical component of data governance and quality management. It tracks the flow of data from its source to consumption, providing visibility into how data moves and transforms throughout its lifecycle. Nessi.dev's lineage feature allows you to:

1. **Visualize Data Flows**: Create interactive graph visualizations showing data dependencies
2. **Track Data Evolution**: Capture snapshots of lineage graphs to track how they change over time
3. **Compare Versions**: Compare different lineage snapshots to identify changes
4. **Analyze Impact**: Assess the impact of schema or process changes on downstream systems
5. **Integrate with Data Catalogs**: Import lineage information from popular data catalogs
6. **Export in Multiple Formats**: Share lineage information in various formats (HTML, JSON, SVG, DOT)

## Core Concepts

### Lineage Graph

A lineage graph represents the flow of data through various systems and processes. It consists of:

- **Nodes**: Represent data assets (tables, views, files) or processes (ETL jobs, applications)
- **Edges**: Represent relationships between nodes (read, write, depends_on)
- **Properties**: Additional metadata associated with nodes and edges

### Node Types

Nessi.dev supports the following node types:

- `table`: Database tables
- `view`: Database views
- `query`: SQL queries
- `process`: Data processing jobs (ETL, transformations)
- `dataset`: File-based datasets
- `api`: API endpoints
- `application`: Applications that produce or consume data

### Edge Types

Edges represent relationships between nodes:

- `read`: Source node is read by target node
- `write`: Source node writes to target node
- `depends_on`: Source node depends on target node
- `produces`: Source node produces target node

### Snapshots

Snapshots capture the state of a lineage graph at a specific point in time. They are useful for:

- Historical record-keeping
- Auditing
- Comparing changes over time

## Using the Lineage Feature

### Via API

The lineage feature is accessible through a RESTful API:

#### Graph Management

```
GET    /lineage/graphs              # List all lineage graphs
POST   /lineage/graphs              # Create a new lineage graph
GET    /lineage/graphs/{id}         # Get a specific lineage graph
DELETE /lineage/graphs/{id}         # Delete a lineage graph
POST   /lineage/graphs/{id}/nodes   # Add a node to a graph
POST   /lineage/graphs/{id}/edges   # Add an edge to a graph
```

#### Snapshot Management

```
GET    /lineage/graphs/{id}/snapshots              # List snapshots for a graph
POST   /lineage/graphs/{id}/snapshots              # Create a new snapshot
GET    /lineage/snapshots/{id}                     # Get a specific snapshot
GET    /lineage/graphs/{id}/snapshots/time/{time}  # Get snapshot at a specific time
```

#### Visualization and Comparison

```
GET    /lineage/graphs/{id}/visualize          # Visualize a lineage graph
GET    /lineage/graphs/{id}/export/{format}    # Export lineage in a specific format
POST   /lineage/compare                        # Compare two snapshots
GET    /lineage/diffs/{id}/impact              # Analyze impact of changes
```

#### Integration

```
POST   /lineage/import/catalog    # Import lineage from a data catalog
```

### Via CLI

The lineage feature is also accessible through the Nessi CLI:

```bash
# Graph management
nessi lineage graph list
nessi lineage graph create --name "My Graph" --description "My lineage graph"
nessi lineage graph get <graph-id>
nessi lineage graph delete <graph-id>
nessi lineage graph add-node <graph-id> --name "My Table" --type "table" --description "My table"
nessi lineage graph add-edge <graph-id> --source <source-id> --target <target-id> --type "read"

# Snapshot management
nessi lineage snapshot list <graph-id>
nessi lineage snapshot create <graph-id> --comment "Initial snapshot"
nessi lineage snapshot get <snapshot-id>
nessi lineage snapshot get-at-time <graph-id> "2023-01-01T00:00:00Z"

# Visualization and comparison
nessi lineage visualize <graph-id> --format html --output lineage.html
nessi lineage export <graph-id> --format json --output lineage.json
nessi lineage compare --snapshot1 <snapshot-id-1> --snapshot2 <snapshot-id-2> --output diff.json

# Integration
nessi lineage import catalog --catalog "aws-glue" --database "my-db" --table "my-table"
```

## Visualization Options

When visualizing a lineage graph, you can customize the output with the following options:

- **Format**: HTML, JSON, SVG, DOT
- **Width/Height**: Dimensions of the visualization
- **Focus Node**: Center the visualization on a specific node
- **Max Depth**: Limit the depth of the graph from the focus node
- **Layout Algorithm**: Force-directed, hierarchical, or radial
- **Node/Edge Filtering**: Include only specific node or edge types

Example:

```bash
nessi lineage visualize <graph-id> \
  --format html \
  --width 1200 \
  --height 800 \
  --focus-node <node-id> \
  --max-depth 3 \
  --layout force \
  --output lineage.html
```

## Comparing Lineage Snapshots

Comparing lineage snapshots helps you understand how data flows have changed over time. The comparison identifies:

- Added nodes
- Removed nodes
- Modified nodes
- Added edges
- Removed edges

This information is crucial for:

- Impact analysis
- Troubleshooting data quality issues
- Ensuring compliance with data governance policies
- Understanding the evolution of your data ecosystem

Example:

```bash
nessi lineage compare \
  --snapshot1 <snapshot-id-1> \
  --snapshot2 <snapshot-id-2> \
  --output diff.json
```

## Data Catalog Integration

Nessi.dev can import lineage information from popular data catalogs:

- AWS Glue Data Catalog
- Azure Purview
- Google Cloud Data Catalog

This integration allows you to leverage existing metadata and lineage information from your cloud providers.

Example:

```bash
nessi lineage import catalog \
  --catalog "aws-glue" \
  --database "my-db" \
  --table "my-table"
```

## Export Formats

Lineage information can be exported in various formats:

- **HTML**: Interactive visualization for web browsers
- **JSON**: Machine-readable format for integration with other tools
- **SVG**: Vector graphics for inclusion in documentation
- **DOT**: Graphviz format for further customization

Example:

```bash
nessi lineage export <graph-id> --format html --output lineage.html
```

## Best Practices

1. **Create snapshots regularly**: Capture lineage snapshots at key milestones or on a regular schedule
2. **Use meaningful node names**: Clear naming helps understand the lineage graph
3. **Add detailed descriptions**: Include context about each node and edge
4. **Focus on critical data assets**: Start with your most important data assets
5. **Integrate with data catalogs**: Leverage existing metadata from your data catalogs
6. **Compare before and after changes**: Create snapshots before and after significant changes
7. **Export for documentation**: Include lineage visualizations in your data documentation

## Use Cases

### Data Governance and Compliance

- Track data flows for regulatory compliance (GDPR, CCPA, etc.)
- Document data lineage for audit purposes
- Ensure data handling follows organizational policies

### Impact Analysis

- Assess the impact of schema changes on downstream systems
- Identify affected reports and outputs when source data changes
- Plan data migrations with full understanding of dependencies

### Troubleshooting

- Trace data quality issues to their source
- Understand how errors propagate through systems
- Identify affected downstream systems when issues occur

### Documentation

- Create visual documentation of data flows
- Onboard new team members with clear data flow visualizations
- Share data lineage with stakeholders

## Technical Implementation

The lineage feature is implemented with the following components:

- **Core Model**: Defines the structure of lineage graphs, nodes, and edges
- **Storage**: Persists lineage graphs and snapshots
- **Visualization**: Renders interactive graph visualizations
- **Comparison**: Compares lineage snapshots and analyzes impact
- **API**: Exposes lineage functionality through REST endpoints
- **CLI**: Provides command-line access to lineage features
- **Catalog Integration**: Connects to external data catalogs

For more details on the implementation, see the code in the `pkg/lineage` directory.
