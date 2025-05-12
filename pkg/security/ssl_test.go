package security

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertManager(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "certs-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := SSLConfig{
		Enabled:      true,
		CertFile:     filepath.Join(tempDir, "cert.pem"),
		KeyFile:      filepath.Join(tempDir, "key.pem"),
		AutoGenerate: true,
	}

	// Create cert manager
	cm := NewCertManager(config)
	require.NotNil(t, cm)

	// Test GetTLSConfig
	tlsConfig, err := cm.GetTLSConfig()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)

	// Verify certificate files were created
	_, err = os.Stat(config.CertFile)
	assert.NoError(t, err)
	_, err = os.Stat(config.KeyFile)
	assert.NoError(t, err)

	// Test with SSL disabled
	disabledConfig := SSLConfig{
		Enabled: false,
	}
	disabledCM := NewCertManager(disabledConfig)
	_, err = disabledCM.GetTLSConfig()
	assert.Error(t, err)
}

func TestStartHTTPSServer(t *testing.T) {
	// Skip this test in automated testing environments
	t.Skip("This test starts an actual HTTPS server")

	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "certs-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := SSLConfig{
		Enabled:      true,
		CertFile:     filepath.Join(tempDir, "cert.pem"),
		KeyFile:      filepath.Join(tempDir, "key.pem"),
		AutoGenerate: true,
	}

	// Create cert manager
	cm := NewCertManager(config)
	require.NotNil(t, cm)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HTTPS works!"))
	})

	// Start HTTPS server in a goroutine
	go func() {
		err := cm.StartHTTPSServer(":8443", testHandler)
		if err != nil {
			t.Logf("HTTPS server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Create HTTP client that skips certificate verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Make HTTPS request
	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		t.Logf("HTTPS request error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Verify response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGenerateSelfSignedCert(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "certs-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := SSLConfig{
		Enabled:      true,
		CertFile:     filepath.Join(tempDir, "cert.pem"),
		KeyFile:      filepath.Join(tempDir, "key.pem"),
		AutoGenerate: true,
	}

	// Create cert manager
	cm := NewCertManager(config)
	require.NotNil(t, cm)

	// Generate self-signed certificate
	err = cm.generateSelfSignedCert()
	require.NoError(t, err)

	// Verify certificate files were created
	_, err = os.Stat(config.CertFile)
	assert.NoError(t, err)
	_, err = os.Stat(config.KeyFile)
	assert.NoError(t, err)

	// Load certificate to verify it's valid
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}
