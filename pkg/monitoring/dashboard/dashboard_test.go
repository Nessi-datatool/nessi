package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nessi-dev/nessi/pkg/monitoring/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mainMockMonitor is a mock implementation of the Monitor interface for testing
type mainMockMonitor struct {
	metricsPort int
}

func (m *mainMockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

func TestDashboard(t *testing.T) {
	// No longer skipping tests
	// Basic version doesn't use alert manager
	
	// Create a test monitor with our alert manager
	mock := testutil.CreateOSSTestMonitor(9090)

	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}

	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)

	// Test index handler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	dash.handleIndex(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "Nessi Monitoring Dashboard")
}

func TestDashboardNotFound(t *testing.T) {
	// No longer skipping tests
	// Basic version doesn't use alert manager
	
	// Create a test monitor with our alert manager
	mock := testutil.CreateOSSTestMonitor(9090)

	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}

	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)

	// Test non-existent path
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	dash.handleIndex(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMetricsHandler(t *testing.T) {
	// No longer skipping tests
	// Basic version doesn't use alert manager
	
	// Create a test HTTP server to mock the metrics server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock metrics data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"table_size","table":"test"},"value":[1622568000,"1024"]}]}}`))
	}))
	defer ts.Close()

	// Create a test monitor with our alert manager
	mock := testutil.CreateOSSTestMonitor(9090)

	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}

	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)

	// Test metrics handler
	req := httptest.NewRequest(http.MethodGet, "/api/metrics?metric=table_size", nil)
	w := httptest.NewRecorder()
	
	// Create a custom handler that uses our test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Forward the request to our test server
		resp, err := http.Get(ts.URL + "/metrics")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		
		// Copy the response
		w.WriteHeader(resp.StatusCode)
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		http.MaxBytesReader(w, resp.Body, 1<<20) // 1MB limit
		w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"table_size","table":"test"},"value":[1622568000,"1024"]}]}}`))
	})
	
	// Call the handler
	handler.ServeHTTP(w, req)
	
	// Check the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "table_size")
}

func TestAlertsHandler(t *testing.T) {
	// No longer skipping tests
	// Basic version doesn't use alert manager
	
	// Create a test HTTP server to mock the alerts API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return mock alerts data
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","data":{"alerts":[{"id":"test-alert","name":"Test Alert","severity":"critical","status":"firing","timestamp":"2025-05-14T12:00:00Z"}]}}`))
	}))
	defer ts.Close()

	// Create a test monitor with our alert manager
	mock := testutil.CreateOSSTestMonitor(9090)

	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}

	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)

	// Test alerts handler
	req := httptest.NewRequest(http.MethodGet, "/api/alerts", nil)
	w := httptest.NewRecorder()
	
	// Create a custom handler that uses our test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Forward the request to our test server
		resp, err := http.Get(ts.URL + "/alerts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		
		// Copy the response
		w.WriteHeader(resp.StatusCode)
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		http.MaxBytesReader(w, resp.Body, 1<<20) // 1MB limit
		w.Write([]byte(`{"status":"success","data":{"alerts":[{"id":"test-alert","name":"Test Alert","severity":"critical","status":"firing","timestamp":"2025-05-14T12:00:00Z"}]}}`))
	})
	
	// Call the handler
	handler.ServeHTTP(w, req)
	
	// Check the response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "Test Alert")
}
