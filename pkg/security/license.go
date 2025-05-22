// Package security provides license validation and security features
package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// License represents a Nessi license
type License struct {
	Token     string    `json:"token"`
	Issued    time.Time `json:"issued"`
	Expires   time.Time `json:"expires"`
	Customer  string    `json:"customer"`
	Plan      string    `json:"plan"`
	Features  []string  `json:"features"`
	Signature string    `json:"signature"`
}

// LicenseClaims represents the JWT claims in a license token
type LicenseClaims struct {
	jwt.RegisteredClaims
	Customer string   `json:"customer"`
	Plan     string   `json:"plan"`
	Features []string `json:"features"`
}

// ErrLicenseExpired is returned when the license has expired
var ErrLicenseExpired = errors.New("license has expired")

// ErrLicenseInvalid is returned when the license is invalid
var ErrLicenseInvalid = errors.New("license is invalid")

// ErrLicenseNotFound is returned when the license file is not found
var ErrLicenseNotFound = errors.New("license file not found")

// ErrFeatureNotLicensed is returned when a premium feature is not included in the license
var ErrFeatureNotLicensed = errors.New("feature not included in license")

// ErrPlanDowngrade is returned when attempting to downgrade to a lower plan
var ErrPlanDowngrade = errors.New("cannot downgrade to a lower plan")

// ErrTrialExpired is returned when the trial period has expired
var ErrTrialExpired = errors.New("trial period has expired")

// ValidateLicense validates the license file
func ValidateLicense() (*License, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// License file path
	licensePath := filepath.Join(homeDir, ".nessi", "license.json")

	// Check if license file exists
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return nil, ErrLicenseNotFound
	}

	// Read license file
	data, err := os.ReadFile(licensePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read license file: %w", err)
	}

	// Parse license file
	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, fmt.Errorf("failed to parse license file: %w", err)
	}

	// Verify JWT token
	if err := verifyToken(license.Token); err != nil {
		return nil, err
	}

	// Check expiry
	if time.Now().After(license.Expires) {
		return &license, ErrLicenseExpired
	}

	return &license, nil
}

// verifyToken verifies the JWT token in the license
func verifyToken(tokenString string) error {
	// This is a simplified implementation
	// In a real implementation, we would verify the signature using a public key

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &LicenseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// In a real implementation, we would verify the signing method and return the public key
		return []byte("secret"), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return ErrLicenseInvalid
	}

	return nil
}

// HasFeature checks if the license has a specific feature
func HasFeature(license *License, feature string) bool {
	if license == nil {
		return false
	}

	// Check direct feature match
	for _, f := range license.Features {
		if f == feature {
			return true
		}
	}

	// Check if the feature is included in the plan
	if planFeatures, ok := PlanFeatures[license.Plan]; ok {
		for _, f := range planFeatures {
			if f == feature {
				return true
			}
		}
	}

	return false
}

// RequireFeature checks if the license has a specific feature and returns an error if not
func RequireFeature(feature string) error {
	// For development mode, use the standard implementation
	if os.Getenv("NESSI_DEV_MODE") == "true" {
		return standardRequireFeature(feature)
	}

	// In production, use the secure implementation with anti-tampering measures
	return secureRequireFeature(feature)
}

// standardRequireFeature is the original implementation of RequireFeature
func standardRequireFeature(feature string) error {
	license, err := ValidateLicense()
	if err == ErrLicenseNotFound {
		// For community features, allow usage without a license
		if feature == "" || strings.HasPrefix(feature, "community_") {
			return nil
		}

		// Check if trial is active
		if isTrialActive() {
			return nil // Allow all features during trial period
		}

		return fmt.Errorf("%w: %s requires a valid license", err, FeatureDescriptions[feature])
	} else if err != nil {
		return err
	}

	if !HasFeature(license, feature) {
		return fmt.Errorf("%w: %s is only available in %s", ErrFeatureNotLicensed,
			FeatureDescriptions[feature], getMinimumPlanForFeature(feature))
	}

	return nil
}

// getMinimumPlanForFeature returns the minimum plan that includes the feature
func getMinimumPlanForFeature(feature string) string {
	// Order of plans from lowest to highest
	plans := []string{PlanStarter, PlanPro}

	// Check each plan
	// First check if it's a Starter plan feature
	for _, f := range PlanFeatures[PlanStarter] {
		if f == feature {
			return PlanStarter
		}
	}

	// If not in Starter, it must be in Pro
	for _, f := range PlanFeatures[PlanPro] {
		if f == feature {
			return PlanPro
		}
	}

	return "unknown plan"
}

// getExistingTrials returns the number of trials used for a specific machine ID
func getExistingTrials(machineID string) (int, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Registry file path
	registryPath := filepath.Join(homeDir, ".nessi", "trial_registry.json")

	// Check if registry file exists
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		// Registry doesn't exist yet, create it
		registry := map[string]int{}

		// Create signed registry
		signedData, err := CreateSignedData(registry)
		if err != nil {
			return 0, fmt.Errorf("failed to create signed registry: %w", err)
		}

		data, err := json.MarshalIndent(signedData, "", "  ")
		if err != nil {
			return 0, fmt.Errorf("failed to marshal registry: %w", err)
		}

		if err := os.WriteFile(registryPath, data, 0644); err != nil {
			return 0, fmt.Errorf("failed to write registry file: %w", err)
		}

		return 0, nil
	}

	// Read registry file
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read registry file: %w", err)
	}

	// Parse registry file with signature verification
	var signedData SignedData
	if err := json.Unmarshal(data, &signedData); err != nil {
		// Try parsing as legacy format (without signature)
		var registry map[string]int
		if err := json.Unmarshal(data, &registry); err != nil {
			return 0, fmt.Errorf("failed to parse registry file: %w", err)
		}

		// Convert to signed format for future reads
		if err := updateMachineTrialsRegistry(machineID, registry[machineID]); err != nil {
			// Non-fatal error, just log it
			fmt.Printf("Warning: Failed to convert registry to signed format: %v\n", err)
		}

		return registry[machineID], nil
	}

	// Verify signature and unmarshal data
	var registry map[string]int
	if err := VerifyAndUnmarshal(&signedData, &registry); err != nil {
		return 0, fmt.Errorf("registry verification failed: %w", err)
	}

	// Return the number of trials for this machine ID
	return registry[machineID], nil
}

// updateMachineTrialsRegistry updates the number of trials used for a specific machine ID
func updateMachineTrialsRegistry(machineID string, count int) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Registry file path
	registryPath := filepath.Join(homeDir, ".nessi", "trial_registry.json")

	// Read existing registry or create new one
	var registry map[string]int
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		registry = map[string]int{}
	} else {
		data, err := os.ReadFile(registryPath)
		if err != nil {
			return fmt.Errorf("failed to read registry file: %w", err)
		}

		// Try to parse with signature verification
		var signedData SignedData
		if err := json.Unmarshal(data, &signedData); err != nil {
			// Try parsing as legacy format (without signature)
			if err := json.Unmarshal(data, &registry); err != nil {
				return fmt.Errorf("failed to parse registry file: %w", err)
			}
		} else {
			// Verify signature and unmarshal data
			if err := VerifyAndUnmarshal(&signedData, &registry); err != nil {
				return fmt.Errorf("registry verification failed: %w", err)
			}
		}
	}

	// Update registry
	registry[machineID] = count

	// Create signed registry
	signedData, err := CreateSignedData(registry)
	if err != nil {
		return fmt.Errorf("failed to create signed registry: %w", err)
	}

	// Write updated registry with signature
	data, err := json.MarshalIndent(signedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	return nil
}

// isTrialActive checks if the free trial period is active
func isTrialActive() bool {
	// Get trial start time from file
	trialInfo, err := getTrialInfo()
	if err != nil {
		// If there's an error reading trial info, no trial exists
		fmt.Printf("Debug: Trial info error: %v\n", err)
		return false // No automatic trial creation
	}

	// Check if trial has expired (1 month = 30 days)
	trialEndTime := trialInfo.StartTime.AddDate(0, 1, 0) // Add 1 month
	isActive := time.Now().Before(trialEndTime) && trialInfo.Active
	fmt.Printf("Debug: Trial active check: before=%v, active=%v, result=%v\n",
		time.Now().Before(trialEndTime), trialInfo.Active, isActive)
	return isActive
}

// TrialInfo stores information about the free trial
type TrialInfo struct {
	StartTime   time.Time `json:"start_time"`
	Active      bool      `json:"active"`
	MachineID   string    `json:"machine_id"`
	TrialCount  int       `json:"trial_count"`
	LastRenewal time.Time `json:"last_renewal"`
}

// getTrialInfo reads trial information from file
func getTrialInfo() (*TrialInfo, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Trial info file path
	trialPath := filepath.Join(homeDir, ".nessi", "trial.json")

	// Check if trial file exists
	if _, err := os.Stat(trialPath); os.IsNotExist(err) {
		return nil, ErrLicenseNotFound
	}

	// Read trial file
	data, err := os.ReadFile(trialPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read trial file: %w", err)
	}

	// Parse trial file with signature verification
	var signedData SignedData
	if err := json.Unmarshal(data, &signedData); err != nil {
		// Try parsing as legacy format (without signature)
		var trialInfo TrialInfo
		if err := json.Unmarshal(data, &trialInfo); err != nil {
			return nil, fmt.Errorf("failed to parse trial file: %w", err)
		}
		// Convert to signed format for future reads
		if err := saveTrialInfo(&trialInfo); err != nil {
			// Non-fatal error, just log it
			fmt.Printf("Warning: Failed to convert trial info to signed format: %v\n", err)
		}
		return &trialInfo, nil
	}

	// Verify signature and unmarshal data
	var trialInfo TrialInfo
	if err := VerifyAndUnmarshal(&signedData, &trialInfo); err != nil {
		return nil, fmt.Errorf("trial data verification failed: %w", err)
	}

	return &trialInfo, nil
}

// saveTrialInfo saves trial information to file with signature
func saveTrialInfo(trialInfo *TrialInfo) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create .nessi directory if it doesn't exist
	nessiDir := filepath.Join(homeDir, ".nessi")
	if err := os.MkdirAll(nessiDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Trial info file path
	trialPath := filepath.Join(nessiDir, "trial.json")

	// Create signed data
	signedData, err := CreateSignedData(trialInfo)
	if err != nil {
		return fmt.Errorf("failed to create signed trial data: %w", err)
	}

	// Marshal signed data to JSON
	data, err := json.MarshalIndent(signedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal signed trial data: %w", err)
	}

	// Write trial file
	if err := os.WriteFile(trialPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write trial file: %w", err)
	}

	return nil
}

// initializeTrial creates a new trial
func initializeTrial() error {
	// Get machine ID
	machineID, err := GetMachineID()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Check for existing machine trials
	existingTrials, err := getExistingTrials(machineID)
	if err == nil && existingTrials >= 2 {
		return fmt.Errorf("maximum number of trials (2) reached for this machine")
	}

	// Create trial info
	trialInfo := TrialInfo{
		StartTime:   time.Now(),
		Active:      true,
		MachineID:   machineID,
		TrialCount:  existingTrials + 1,
		LastRenewal: time.Now(),
	}

	// Save trial info with signature
	if err := saveTrialInfo(&trialInfo); err != nil {
		return fmt.Errorf("failed to save trial info: %w", err)
	}

	// Update machine trials registry
	if err := updateMachineTrialsRegistry(machineID, existingTrials+1); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("Warning: Failed to update machine trials registry: %v\n", err)
	}

	return nil
}

// StartTrial initiates a free trial for premium features
func StartTrial() error {
	// Check if trial already exists
	trialInfo, err := getTrialInfo()
	if err == nil && trialInfo.Active {
		// Trial already exists and is active
		return fmt.Errorf("trial already active, expires on %s",
			trialInfo.StartTime.AddDate(0, 1, 0).Format("2006-01-02"))
	}

	// Get machine ID
	machineID, err := GetMachineID()
	if err != nil {
		return fmt.Errorf("failed to get machine ID: %w", err)
	}

	// Check if this machine has exceeded trial limit
	trialCount, err := getTrialCountForMachine(machineID)
	if err != nil {
		// If we can't read the registry, assume this is the first trial
		trialCount = 0
	}

	if trialCount >= 2 {
		return fmt.Errorf("maximum number of trials (2) already used on this machine")
	}

	// Create new trial info
	newTrialInfo := TrialInfo{
		StartTime:   time.Now(),
		Active:      true,
		MachineID:   machineID,
		TrialCount:  trialCount + 1,
		LastRenewal: time.Now(),
	}

	// Save trial info
	err = saveTrialInfo(newTrialInfo)
	if err != nil {
		return fmt.Errorf("failed to save trial info: %w", err)
	}

	// Update machine trials registry
	err = updateMachineTrialsRegistry(machineID, trialCount+1)
	if err != nil {
		fmt.Printf("Warning: Failed to update machine trials registry: %v\n", err)
	}

	return nil
}

// SetTrialStartTimeForTesting sets the trial start time to a specific date for testing purposes
// This function should only be used for testing
func SetTrialStartTimeForTesting(startTime time.Time) error {
	// Get current trial info
	trialInfo, err := getTrialInfo()
	if err != nil {
		return fmt.Errorf("no active trial found: %w", err)
	}

	// Update start time
	trialInfo.StartTime = startTime

	// Save modified trial info
	err = saveTrialInfo(trialInfo)
	if err != nil {
		return fmt.Errorf("failed to save modified trial info: %w", err)
	}

	return nil
}
