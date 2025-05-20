package security_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSLIntegration(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "ssl-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := security.SSLConfig{
		Enabled:  true,
		CertFile: filepath.Join(tempDir, "cert.pem"),
		KeyFile:  filepath.Join(tempDir, "key.pem"),
	}

	// Create cert manager
	cm := security.NewCertManager(config)
	require.NotNil(t, cm)

	t.Run("Certificate Generation", func(t *testing.T) {
		// Generate self-signed certificate
		err := cm.GenerateSelfSignedCertForTest()
		require.NoError(t, err)

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
	})
}

func TestSSLWithExistingCertificates(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "ssl-existing-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create SSL config
	config := security.SSLConfig{
		Enabled:  true,
		CertFile: filepath.Join(tempDir, "cert.pem"),
		KeyFile:  filepath.Join(tempDir, "key.pem"),
	}

	// Create first cert manager to generate certificates
	cm1 := security.NewCertManager(config)
	require.NotNil(t, cm1)

	// Generate certificates
	err = cm1.GenerateSelfSignedCertForTest()
	require.NoError(t, err)

	// Create second cert manager to use existing certificates
	cm2 := security.NewCertManager(config)
	require.NotNil(t, cm2)

	// Should still be able to get TLS config
	tlsConfig, err := cm2.GetTLSConfig()
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
}
