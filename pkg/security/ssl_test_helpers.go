package security

// GenerateSelfSignedCertForTest is a helper method for tests to generate a self-signed certificate
func (cm *CertManager) GenerateSelfSignedCertForTest() error {
	return cm.generateSelfSignedCert()
}
