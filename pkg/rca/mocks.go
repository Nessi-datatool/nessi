package rca

import (
	"fmt"
	"time"
)

// MockMonitoringClient is a mock implementation of the monitoring client for testing
type MockMonitoringClient struct {
	anomalies map[string]*AnomalyInfo
}

// NewMockMonitoringClient creates a new mock monitoring client
func NewMockMonitoringClient() *MockMonitoringClient {
	return &MockMonitoringClient{
		anomalies: make(map[string]*AnomalyInfo),
	}
}

// AddMockAnomaly adds a mock anomaly to the client
func (m *MockMonitoringClient) AddMockAnomaly(anomaly *AnomalyInfo) {
	m.anomalies[anomaly.ID] = anomaly
}

// GetAnomaly gets an anomaly by ID
func (m *MockMonitoringClient) GetAnomaly(id string) (*AnomalyInfo, error) {
	if anomaly, ok := m.anomalies[id]; ok {
		return anomaly, nil
	}
	return nil, fmt.Errorf("anomaly not found: %s", id)
}

// GetRecentAnomalies gets recent anomalies between start and end time
func (m *MockMonitoringClient) GetRecentAnomalies(start, end time.Time) ([]string, error) {
	var result []string
	for id, anomaly := range m.anomalies {
		if (anomaly.Timestamp.Equal(start) || anomaly.Timestamp.After(start)) &&
			(anomaly.Timestamp.Equal(end) || anomaly.Timestamp.Before(end)) {
			result = append(result, id)
		}
	}
	return result, nil
}

// MockDeltaConnector is a mock implementation of the Delta Lake connector for testing
type MockDeltaConnector struct {
	schemaChanges map[string][]*SchemaChange
}

// SchemaChange represents a schema change in the Delta Lake
type SchemaChange struct {
	TablePath    string
	Timestamp    time.Time
	ColumnName   string
	PreviousType string
	CurrentType  string
	ChangeAuthor string
}

// NewMockDeltaConnector creates a new mock Delta Lake connector
func NewMockDeltaConnector() *MockDeltaConnector {
	return &MockDeltaConnector{
		schemaChanges: make(map[string][]*SchemaChange),
	}
}

// AddMockSchemaChange adds a mock schema change to the connector
func (d *MockDeltaConnector) AddMockSchemaChange(change *SchemaChange) {
	if _, ok := d.schemaChanges[change.TablePath]; !ok {
		d.schemaChanges[change.TablePath] = make([]*SchemaChange, 0)
	}
	d.schemaChanges[change.TablePath] = append(d.schemaChanges[change.TablePath], change)
}

// GetRecentSchemaChanges gets recent schema changes for a table
func (d *MockDeltaConnector) GetRecentSchemaChanges(tablePath string, since time.Time) ([]*SchemaChange, error) {
	if changes, ok := d.schemaChanges[tablePath]; ok {
		var result []*SchemaChange
		for _, change := range changes {
			if change.Timestamp.Equal(since) || change.Timestamp.After(since) {
				result = append(result, change)
			}
		}
		return result, nil
	}
	return nil, nil
}
