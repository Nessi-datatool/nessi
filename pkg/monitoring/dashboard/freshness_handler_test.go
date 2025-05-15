package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/freshness"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFreshnessManager is a mock implementation of the FreshnessManager interface
type MockFreshnessManager struct {
	mock.Mock
}

func (m *MockFreshnessManager) GetSLA(tableName string) (*freshness.SLAConfig, error) {
	args := m.Called(tableName)
	return args.Get(0).(*freshness.SLAConfig), args.Error(1)
}

func (m *MockFreshnessManager) ListSLAs() []*freshness.SLAConfig {
	args := m.Called()
	return args.Get(0).([]*freshness.SLAConfig)
}

func (m *MockFreshnessManager) SetSLA(config *freshness.SLAConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockFreshnessManager) DeleteSLA(tableName string) error {
	args := m.Called(tableName)
	return args.Error(0)
}

func (m *MockFreshnessManager) CheckFreshness(tableName string) (*freshness.FreshnessStatus, error) {
	args := m.Called(tableName)
	return args.Get(0).(*freshness.FreshnessStatus), args.Error(1)
}

func (m *MockFreshnessManager) CheckAllFreshness() ([]*freshness.FreshnessStatus, error) {
	args := m.Called()
	return args.Get(0).([]*freshness.FreshnessStatus), args.Error(1)
}

func (m *MockFreshnessManager) GetTableTrends(tableName string) (*freshness.FreshnessTrends, error) {
	args := m.Called(tableName)
	return args.Get(0).(*freshness.FreshnessTrends), args.Error(1)
}

func (m *MockFreshnessManager) GetAllTablesTrends() (*freshness.FreshnessTrends, error) {
	args := m.Called()
	return args.Get(0).(*freshness.FreshnessTrends), args.Error(1)
}

func TestDashboard_HandleFreshnessDashboard(t *testing.T) {
	// Create a test dashboard with a mock freshness manager
	mockFreshnessManager := new(MockFreshnessManager)
	dash := &Dashboard{
		freshnessManager: mockFreshnessManager,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/freshness", nil)
	w := httptest.NewRecorder()

	// Set up templates
	var err error
	dash.templates, err = parseTestTemplates()
	require.NoError(t, err)

	// Call the handler
	dash.handleFreshnessDashboard(w, req)

	// Check response
	resp := w.Result()
	defer resp.Body.Close()

	// Should return 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDashboard_HandleFreshnessStatusAPI(t *testing.T) {
	// Create a test dashboard with a mock freshness manager
	mockFreshnessManager := new(MockFreshnessManager)
	dash := &Dashboard{
		freshnessManager: mockFreshnessManager,
	}

	// Create test data
	now := time.Now()
	config := &freshness.SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	status := &freshness.FreshnessStatus{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		LastUpdateTime:    now.Add(-30 * time.Minute),
		TimeSinceUpdate:   30 * time.Minute,
		ExpectedFrequency: time.Hour,
		Status:            freshness.SLALevelInfo,
		NextExpectedUpdate: now.Add(30 * time.Minute),
		SLAConfig:         config,
	}

	// Set up mock behavior
	mockFreshnessManager.On("CheckFreshness", "test_table").Return(status, nil)
	mockFreshnessManager.On("CheckAllFreshness").Return([]*freshness.FreshnessStatus{status}, nil)

	// Test getting status for a specific table
	t.Run("Get status for specific table", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/status?table=test_table", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessStatusAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var statuses []*freshness.FreshnessStatus
		err := json.NewDecoder(resp.Body).Decode(&statuses)
		require.NoError(t, err)
		assert.Len(t, statuses, 1)
		assert.Equal(t, "test_table", statuses[0].TableName)
	})

	// Test getting status for all tables
	t.Run("Get status for all tables", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/status", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessStatusAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var statuses []*freshness.FreshnessStatus
		err := json.NewDecoder(resp.Body).Decode(&statuses)
		require.NoError(t, err)
		assert.Len(t, statuses, 1)
		assert.Equal(t, "test_table", statuses[0].TableName)
	})

	// Verify all mock expectations were met
	mockFreshnessManager.AssertExpectations(t)
}

func TestDashboard_HandleFreshnessSLAAPI(t *testing.T) {
	// Create a test dashboard with a mock freshness manager
	mockFreshnessManager := new(MockFreshnessManager)
	dash := &Dashboard{
		freshnessManager: mockFreshnessManager,
	}

	// Create test data
	config1 := &freshness.SLAConfig{
		TableName:         "test_table1",
		TablePath:         "/path/to/test_table1",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	config2 := &freshness.SLAConfig{
		TableName:         "test_table2",
		TablePath:         "/path/to/test_table2",
		ExpectedFrequency: 24 * time.Hour,
		WarningThreshold:  125,
		CriticalThreshold: 150,
		Enabled:           true,
	}

	// Set up mock behavior
	mockFreshnessManager.On("GetSLA", "test_table1").Return(config1, nil)
	mockFreshnessManager.On("ListSLAs").Return([]*freshness.SLAConfig{config1, config2}, nil)
	mockFreshnessManager.On("SetSLA", mock.AnythingOfType("*freshness.SLAConfig")).Return(nil)
	mockFreshnessManager.On("DeleteSLA", "test_table1").Return(nil)

	// Test getting SLA for a specific table
	t.Run("Get SLA for specific table", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/sla?table=test_table1", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessSLAAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var configs []*freshness.SLAConfig
		err := json.NewDecoder(resp.Body).Decode(&configs)
		require.NoError(t, err)
		assert.Len(t, configs, 1)
		assert.Equal(t, "test_table1", configs[0].TableName)
	})

	// Test getting all SLA configurations
	t.Run("Get all SLA configurations", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/sla", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessSLAAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var configs []*freshness.SLAConfig
		err := json.NewDecoder(resp.Body).Decode(&configs)
		require.NoError(t, err)
		assert.Len(t, configs, 2)
	})

	// Test creating/updating an SLA configuration
	t.Run("Create/update SLA configuration", func(t *testing.T) {
		configJSON := `{
			"table_name": "new_table",
			"table_path": "/path/to/new_table",
			"expected_frequency": 3600000000000,
			"warning_threshold": 150,
			"critical_threshold": 200,
			"enabled": true
		}`

		req := httptest.NewRequest("POST", "/api/freshness/sla", strings.NewReader(configJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		dash.handleFreshnessSLAAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test deleting an SLA configuration
	t.Run("Delete SLA configuration", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/freshness/sla?table=test_table1", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessSLAAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Verify all mock expectations were met
	mockFreshnessManager.AssertExpectations(t)
}

func TestDashboard_HandleFreshnessTrendsAPI(t *testing.T) {
	// Create a test dashboard with a mock freshness manager
	mockFreshnessManager := new(MockFreshnessManager)
	dash := &Dashboard{
		freshnessManager: mockFreshnessManager,
	}

	// Create test data
	now := time.Now()
	historyEntry := freshness.FreshnessHistoryEntry{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		Timestamp:         now.Add(-1 * time.Hour),
		LastUpdateTime:    now.Add(-2 * time.Hour),
		TimeSinceUpdate:   time.Hour,
		ExpectedFrequency: time.Hour,
		Status:            "info",
	}

	compliance := freshness.FreshnessCompliance{
		InfoCount:      1,
		WarningCount:   0,
		CriticalCount:  0,
		TotalCount:     1,
		ComplianceRate: 100.0,
	}

	trends := &freshness.FreshnessTrends{
		History:    []freshness.FreshnessHistoryEntry{historyEntry},
		Compliance: compliance,
	}

	// Set up mock behavior
	mockFreshnessManager.On("GetTableTrends", "test_table").Return(trends, nil)
	mockFreshnessManager.On("GetAllTablesTrends").Return(trends, nil)

	// Test getting trends for a specific table
	t.Run("Get trends for specific table", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/trends?table=test_table", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessTrendsAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var trendsData freshness.FreshnessTrends
		err := json.NewDecoder(resp.Body).Decode(&trendsData)
		require.NoError(t, err)
		assert.Len(t, trendsData.History, 1)
		assert.Equal(t, "test_table", trendsData.History[0].TableName)
		assert.Equal(t, float64(100), trendsData.Compliance.ComplianceRate)
	})

	// Test getting trends for all tables
	t.Run("Get trends for all tables", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/trends", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessTrendsAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		var trendsData freshness.FreshnessTrends
		err := json.NewDecoder(resp.Body).Decode(&trendsData)
		require.NoError(t, err)
		assert.Len(t, trendsData.History, 1)
	})

	// Verify all mock expectations were met
	mockFreshnessManager.AssertExpectations(t)
}

func TestDashboard_HandleFreshnessExportAPI(t *testing.T) {
	// Create a test dashboard with a mock freshness manager
	mockFreshnessManager := new(MockFreshnessManager)
	dash := &Dashboard{
		freshnessManager: mockFreshnessManager,
	}

	// Create test data
	now := time.Now()
	config := &freshness.SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	status := &freshness.FreshnessStatus{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		LastUpdateTime:    now.Add(-30 * time.Minute),
		TimeSinceUpdate:   30 * time.Minute,
		ExpectedFrequency: time.Hour,
		Status:            freshness.SLALevelInfo,
		NextExpectedUpdate: now.Add(30 * time.Minute),
		SLAConfig:         config,
	}

	// Set up mock behavior
	mockFreshnessManager.On("CheckAllFreshness").Return([]*freshness.FreshnessStatus{status}, nil)

	// Test exporting as JSON
	t.Run("Export as JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/export?format=json", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessExportAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Contains(t, resp.Header.Get("Content-Disposition"), "freshness_data.json")

		var statuses []*freshness.FreshnessStatus
		err := json.NewDecoder(resp.Body).Decode(&statuses)
		require.NoError(t, err)
		assert.Len(t, statuses, 1)
		assert.Equal(t, "test_table", statuses[0].TableName)
	})

	// Test exporting as CSV
	t.Run("Export as CSV", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/freshness/export?format=csv", nil)
		w := httptest.NewRecorder()

		dash.handleFreshnessExportAPI(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/csv", resp.Header.Get("Content-Type"))
		assert.Contains(t, resp.Header.Get("Content-Disposition"), "freshness_data.csv")
	})

	// Verify all mock expectations were met
	mockFreshnessManager.AssertExpectations(t)
}

// Helper function to parse test templates
func parseTestTemplates() (*http.ServeMux, error) {
	// Return a simple mux that can handle the freshness.html template
	mux := http.NewServeMux()
	mux.HandleFunc("/freshness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Freshness Dashboard</body></html>"))
	})
	return mux, nil
}
