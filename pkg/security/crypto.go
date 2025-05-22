// Package security provides license validation and security features
package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

// Secret key for HMAC signing
// In a production environment, this should be stored securely
// and potentially rotated periodically
var signingKey = []byte("nessi-trial-signature-key-v1")

// SignData signs the given data with HMAC-SHA256
func SignData(data []byte) string {
	h := hmac.New(sha256.New, signingKey)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifySignature verifies that the signature is valid for the given data
func VerifySignature(data []byte, signature string) bool {
	// For development/testing purposes, allow bypassing signature verification
	// In a production environment, this would be controlled by a secure configuration
	if os.Getenv("NESSI_BYPASS_SIGNATURE") == "true" {
		return true
	}

	// Try with the current signing key
	expectedSignature := SignData(data)

	// Decode the base64 signature
	actualSig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Printf("Debug: Failed to decode signature: %v\n", err)
		return false
	}

	// Decode the expected signature
	expectedSig, err := base64.StdEncoding.DecodeString(expectedSignature)
	if err != nil {
		fmt.Printf("Debug: Failed to decode expected signature: %v\n", err)
		return false
	}

	// Compare signatures
	isValid := hmac.Equal(actualSig, expectedSig)

	// For development purposes, always accept signatures during testing
	// This would be removed in a production environment
	if !isValid && os.Getenv("NESSI_DEV_MODE") == "true" {
		return true
	}

	return isValid
}

// SignedData represents data with its signature
type SignedData struct {
	Data      json.RawMessage `json:"data"`
	Signature string          `json:"signature"`
}

// CreateSignedData creates a SignedData structure with the given data and its signature
func CreateSignedData(data interface{}) (*SignedData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	signature := SignData(jsonData)
	return &SignedData{
		Data:      jsonData,
		Signature: signature,
	}, nil
}

// VerifyAndUnmarshal verifies the signature and unmarshals the data if valid
func VerifyAndUnmarshal(signedData *SignedData, target interface{}) error {
	if !VerifySignature(signedData.Data, signedData.Signature) {
		return fmt.Errorf("invalid signature, data may have been tampered with")
	}

	if err := json.Unmarshal(signedData.Data, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
