package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
)

// SSLConfig represents SSL configuration
type SSLConfig struct {
	Enabled      bool   `json:"enabled"`
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
	AutoGenerate bool   `json:"auto_generate"`
}

// CertManager handles SSL certificates
type CertManager struct {
	config SSLConfig
}

// NewCertManager creates a new CertManager
func NewCertManager(config SSLConfig) *CertManager {
	return &CertManager{
		config: config,
	}
}

// GetTLSConfig returns a TLS configuration
func (cm *CertManager) GetTLSConfig() (*tls.Config, error) {
	if !cm.config.Enabled {
		return nil, fmt.Errorf("SSL is not enabled")
	}

	// Check if certificate files exist
	certExists, keyExists := cm.certificateFilesExist()

	// Generate self-signed certificate if needed
	if cm.config.AutoGenerate && (!certExists || !keyExists) {
		if err := cm.generateSelfSignedCert(); err != nil {
			return nil, fmt.Errorf("failed to generate self-signed certificate: %w", err)
		}
	}

	// Load certificate
	cert, err := tls.LoadX509KeyPair(cm.config.CertFile, cm.config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// certificateFilesExist checks if certificate files exist
func (cm *CertManager) certificateFilesExist() (bool, bool) {
	certExists := false
	keyExists := false

	if _, err := os.Stat(cm.config.CertFile); err == nil {
		certExists = true
	}

	if _, err := os.Stat(cm.config.KeyFile); err == nil {
		keyExists = true
	}

	return certExists, keyExists
}

// generateSelfSignedCert generates a self-signed certificate
func (cm *CertManager) generateSelfSignedCert() error {
	// Create directories if they don't exist
	certDir := filepath.Dir(cm.config.CertFile)
	keyDir := filepath.Dir(cm.config.KeyFile)

	if err := os.MkdirAll(certDir, 0755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	if err := os.MkdirAll(keyDir, 0755); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Nessi Data Profiling System"},
			CommonName:   "localhost",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Add IP addresses and DNS names
	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	template.DNSNames = []string{"localhost"}

	// Create certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write certificate to file
	certOut, err := os.Create(cm.config.CertFile)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Write private key to file
	keyOut, err := os.OpenFile(cm.config.KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyOut.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}

	logging.Info(fmt.Sprintf("Generated self-signed certificate: %s", cm.config.CertFile))
	return nil
}

// StartHTTPSServer starts an HTTPS server
func (cm *CertManager) StartHTTPSServer(addr string, handler http.Handler) error {
	if !cm.config.Enabled {
		return fmt.Errorf("SSL is not enabled")
	}

	// Get TLS config
	tlsConfig, err := cm.GetTLSConfig()
	if err != nil {
		return err
	}

	// Create server
	server := &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	logging.Info(fmt.Sprintf("Starting HTTPS server on %s", addr))
	return server.ListenAndServeTLS("", "")
}
