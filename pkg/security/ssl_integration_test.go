package security_test

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSLIntegration(t *testing.T) {
	if security.ShouldSkipIntegrationTests(t) {
		return
	}
	// Run tests in parallel for faster execution
	t.Parallel()
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "ssl-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := security.SSLConfig{
		Enabled:      true,
		CertFile:     filepath.Join(tempDir, "cert.pem"),
		KeyFile:      filepath.Join(tempDir, "key.pem"),
		AutoGenerate: true,
	}

	// Create cert manager
	cm := security.NewCertManager(config)
	require.NotNil(t, cm)

	t.Run("Certificate Generation", func(t *testing.T) {
		// Get TLS config
		tlsConfig, err := cm.GetTLSConfig()
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)

		// Verify certificate files were created
		_, err = os.Stat(config.CertFile)
		assert.NoError(t, err)
		_, err = os.Stat(config.KeyFile)
		assert.NoError(t, err)

		// Verify certificate content
		certData, err := os.ReadFile(config.CertFile)
		require.NoError(t, err)
		assert.Contains(t, string(certData), "BEGIN CERTIFICATE")

		keyData, err := os.ReadFile(config.KeyFile)
		require.NoError(t, err)
		assert.Contains(t, string(keyData), "PRIVATE KEY")
	})

	t.Run("HTTPS Server", func(t *testing.T) {
		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("HTTPS works!"))
		})

		// Create test server
		server := httptest.NewUnstartedServer(testHandler)
		defer server.Close()

		// Get TLS config
		tlsConfig, err := cm.GetTLSConfig()
		require.NoError(t, err)

		// Configure server with TLS
		server.TLS = tlsConfig
		server.StartTLS()

		// Create HTTP client that skips certificate verification
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		// Make HTTPS request
		resp, err := client.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "HTTPS works!", string(body))
	})

	t.Run("Disabled SSL", func(t *testing.T) {
		// Create disabled config
		disabledConfig := security.SSLConfig{
			Enabled: false,
		}

		// Create cert manager
		disabledCM := security.NewCertManager(disabledConfig)
		require.NotNil(t, disabledCM)

		// Get TLS config should fail
		_, err := disabledCM.GetTLSConfig()
		assert.Error(t, err)

		// Start HTTPS server should fail
		err = disabledCM.StartHTTPSServer(":0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		assert.Error(t, err)
	})
}

func TestSSLWithExistingCertificates(t *testing.T) {
	if security.ShouldSkipIntegrationTests(t) {
		return
	}
	// Run tests in parallel for faster execution
	t.Parallel()
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "ssl-existing-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := security.SSLConfig{
		Enabled:      true,
		CertFile:     filepath.Join(tempDir, "cert.pem"),
		KeyFile:      filepath.Join(tempDir, "key.pem"),
		AutoGenerate: true,
	}

	// Create first cert manager to generate certificates
	cm1 := security.NewCertManager(config)
	require.NotNil(t, cm1)

	// Generate certificates
	_, err = cm1.GetTLSConfig()
	require.NoError(t, err)

	// Create second cert manager with auto-generate disabled
	config.AutoGenerate = false
	cm2 := security.NewCertManager(config)
	require.NotNil(t, cm2)

	// Should still be able to get TLS config
	tlsConfig, err := cm2.GetTLSConfig()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Using existing certificates"))
	})

	// Create test server
	server := httptest.NewUnstartedServer(testHandler)
	defer server.Close()

	// Configure server with TLS
	server.TLS = tlsConfig
	server.StartTLS()

	// Create HTTP client that skips certificate verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Make HTTPS request
	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "Using existing certificates", string(body))
}
