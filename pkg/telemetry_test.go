package pkg

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestTelemetryManager_SendEvent_Disabled(t *testing.T) {
	tm := &TelemetryManager{Enabled: false, Endpoint: ""}
	event := TelemetryEvent{Command: "test", Timestamp: time.Now()}
	tm.SendEvent(event) // Should do nothing, no panic
}

func TestTelemetryManager_SendEvent_Enabled(t *testing.T) {
	var called int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(200)
	}))
	defer ts.Close()

	tm := &TelemetryManager{Enabled: true, Endpoint: ts.URL}
	event := TelemetryEvent{Command: "test", Timestamp: time.Now()}
	tm.SendEvent(event)
	time.Sleep(100 * time.Millisecond) // Allow goroutine to finish
	if atomic.LoadInt32(&called) == 0 {
		t.Fatal("expected telemetry endpoint to be called")
	}
}

func TestTelemetryManager_SendEvent_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	tm := &TelemetryManager{Enabled: true, Endpoint: ts.URL}
	event := TelemetryEvent{Command: "fail", Timestamp: time.Now()}
	tm.SendEvent(event)
	// No panic, error is logged
}
