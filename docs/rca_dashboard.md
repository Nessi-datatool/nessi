# Root Cause Analysis Dashboard

The Root Cause Analysis (RCA) Dashboard provides an interactive visualization interface for exploring and understanding the underlying causes of data quality anomalies in your Delta Lake tables.

## Features

- **Recent Analyses**: View a list of recent RCA results with key information about primary root causes, confidence scores, and affected tables
- **Run Analysis**: Trigger new RCA analyses for specific anomalies with configurable output formats
- **Insights**: Visualize aggregated statistics about root cause distributions, affected tables, and common root causes
- **Detailed Analysis View**: Explore comprehensive analysis results including:
  - Primary root cause with confidence scoring
  - Secondary potential causes
  - Affected tables and lineage information
  - Recommended actions
  - Links to Grafana monitoring dashboards

## Accessing the Dashboard

The RCA Dashboard is available at the `/rca` endpoint of your Nessi.dev dashboard:

```
http://your-nessi-server:port/rca
```

## Dashboard Sections

### Recent Analyses

This tab displays a table of recent RCA results with the following information:
- Anomaly ID
- Analysis Time
- Primary Root Cause
- Confidence Score
- Affected Tables
- Actions (View button to see detailed analysis)

### Run Analysis

This tab allows you to trigger a new RCA analysis with the following options:
- Anomaly ID selection
- Output format (HTML or JSON)

After running an analysis, the result will be displayed in a modal window and added to the Recent Analyses list.

### Insights

This tab provides visualizations of aggregated RCA data:
- Root Cause Distribution chart (pie chart)
- Affected Tables chart (bar chart)
- Common Root Causes list with occurrence counts and average confidence scores

## API Endpoints

The RCA Dashboard is powered by the following API endpoints:

- `GET /api/rca/recent` - Get recent RCA results
- `POST /api/rca/analyze` - Run a new RCA analysis
- `GET /api/rca/{anomaly_id}` - Get detailed analysis for a specific anomaly
- `GET /api/rca/insights` - Get aggregated insights from all RCA results
- `GET /api/rca/export` - Export RCA data in CSV or JSON format
- `GET /api/rca/{anomaly_id}/export` - Export a specific RCA result in JSON format

## Integration with Other Components

The RCA Dashboard integrates with other Nessi.dev components:

- **Monitoring System**: Uses the monitoring client to retrieve anomaly data
- **Delta Lake Connector**: Analyzes schema changes in Delta tables
- **Alerting System**: Links to alerts related to the analyzed anomalies
- **Grafana Integration**: Provides links to relevant Grafana dashboards

## Security

All RCA Dashboard API endpoints are protected by the same authentication and authorization mechanisms as the rest of the Nessi.dev dashboard. Users must have appropriate permissions to view and run RCA analyses.

## Example Use Cases

1. **Investigating Data Quality Anomalies**:
   - Receive an alert about a data quality anomaly
   - Access the RCA Dashboard to analyze the anomaly
   - View the primary root cause and recommended actions
   - Follow links to affected tables and related anomalies

2. **Trend Analysis**:
   - Use the Insights tab to identify common root causes
   - Identify tables that frequently experience issues
   - Implement preventative measures based on patterns

3. **Documentation and Reporting**:
   - Export RCA results for documentation purposes
   - Include RCA findings in incident reports
   - Track resolution of root causes over time

## Customization

The RCA Dashboard can be customized through the following configuration options:

- Storage path for RCA results
- Default output format
- Maximum analysis depth
- Analysis timeout
- Number of results to display in the Recent Analyses tab

These options can be configured in your Nessi.dev configuration file.
