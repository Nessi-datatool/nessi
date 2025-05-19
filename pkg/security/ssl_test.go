package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertManager(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "certs-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test with SSL enabled
	t.Run("SSL Enabled", func(t *testing.T) {
		// Create SSL config
		config := SSLConfig{
			Enabled:  true,
			CertFile: filepath.Join(tempDir, "cert.pem"),
			KeyFile:  filepath.Join(tempDir, "key.pem"),
		}

		// Create cert manager
		cm := NewCertManager(config)
		require.NotNil(t, cm)

		// Generate certificate
		err = cm.GenerateSelfSignedCertForTest()
		require.NoError(t, err)

		// Verify certificate files were created
		_, err = os.Stat(config.CertFile)
		assert.NoError(t, err)
		_, err = os.Stat(config.KeyFile)
		assert.NoError(t, err)

		// Now TLS config should work
		tlsConfig, err := cm.GetTLSConfig()
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)
	})

	// Test with SSL disabled
	disabledConfig := SSLConfig{
		Enabled: false,
	}
	disabledCM := NewCertManager(disabledConfig)
	_, err = disabledCM.GetTLSConfig()
	assert.Error(t, err)
}

func TestGenerateSelfSignedCert(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "certs-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := SSLConfig{
		Enabled:  true,
		CertFile: filepath.Join(tempDir, "cert.pem"),
		KeyFile:  filepath.Join(tempDir, "key.pem"),
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

	// Now TLS config should work
	tlsConfig, err := cm.GetTLSConfig()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
}
