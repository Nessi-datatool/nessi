package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golang-jwt/jwt/v4"
)

// LicenseManager handles license operations
type LicenseManager struct{}

// StartTrial starts a free trial for premium features
func (lm *LicenseManager) StartTrial() error {
	// Check if a license is already active
	license, err := ValidateLicense()
	if err == nil {
		// License exists and is valid
		return fmt.Errorf("you already have an active license with plan %s", license.Plan)
	} else if err != ErrLicenseNotFound && err != ErrLicenseExpired {
		// Some other error occurred
		return err
	}

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

	// Check if trial already exists
	trialInfo, err := getTrialInfo()
	if err == nil {
		// Trial exists, check if it's still active
		trialEndTime := trialInfo.StartTime.AddDate(0, 1, 0) // Add 1 month
		if time.Now().Before(trialEndTime) {
			// Trial is still active
			daysRemaining := int(trialEndTime.Sub(time.Now()).Hours() / 24)
			return fmt.Errorf("you already have an active trial with %d days remaining", daysRemaining)
		}
		// Trial has expired
		return ErrTrialExpired
	}

	// Initialize new trial
	if err := initializeTrial(); err != nil {
		return fmt.Errorf("failed to start trial: %w", err)
	}

	return nil
}

// NewLicenseManager creates a new license manager
func NewLicenseManager() *LicenseManager {
	return &LicenseManager{}
}

// ActivateLicense activates a license using a license key
func (lm *LicenseManager) ActivateLicense(licenseKey string) error {
	// Validate license key format
	if !isValidLicenseKeyFormat(licenseKey) {
		return fmt.Errorf("invalid license key format")
	}

	// Parse license key to get claims
	token, err := jwt.ParseWithClaims(licenseKey, &LicenseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// In a real implementation, we would verify the signing method and return the public key
		return []byte("secret"), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse license key: %w", err)
	}

	claims, ok := token.Claims.(*LicenseClaims)
	if !ok || !token.Valid {
		return ErrLicenseInvalid
	}

	// Create license object
	license := &License{
		Token:     licenseKey,
		Issued:    time.Unix(claims.IssuedAt.Unix(), 0),
		Expires:   time.Unix(claims.ExpiresAt.Unix(), 0),
		Customer:  claims.Customer,
		Plan:      claims.Plan,
		Features:  claims.Features,
		Signature: licenseKey[strings.LastIndex(licenseKey, ".")+1:],
	}

	// Save license to file
	return saveLicense(license)
}

// GetLicenseInfo returns information about the current license
func (lm *LicenseManager) GetLicenseInfo() (*LicenseInfo, error) {
	license, err := ValidateLicense()
	if err == ErrLicenseNotFound {
		// Check if trial is active
		trialInfo, trialErr := getTrialInfo()
		if trialErr == nil {
			// Check if trial has expired
			trialEndTime := trialInfo.StartTime.AddDate(0, 1, 0) // Add 1 month
			if time.Now().Before(trialEndTime) && trialInfo.Active {
				// Trial is active
				daysRemaining := int(trialEndTime.Sub(time.Now()).Hours() / 24)

				return &LicenseInfo{
					Status:        "Trial",
					Plan:          PlanEnterprise, // Trial includes all features
					Expires:       trialEndTime,
					DaysRemaining: daysRemaining,
					Message:       fmt.Sprintf("Free trial active. %d days remaining.", daysRemaining),
					IsTrial:       true,
				}, nil
			} else {
				// Trial has expired
				return &LicenseInfo{
					Status:  "Expired",
					Plan:    PlanCommunity,
					Expires: trialEndTime,
					Message: "Your trial has expired. Please purchase a license to continue using premium features.",
				}, nil
			}
		}

		return &LicenseInfo{
			Status:  "Not Licensed",
			Plan:    PlanCommunity,
			Message: "Using Community Edition",
		}, nil
	} else if err == ErrLicenseExpired {
		return &LicenseInfo{
			Status:   "Expired",
			Plan:     license.Plan,
			Customer: license.Customer,
			Expires:  license.Expires,
			Message:  "Your license has expired. Please renew to continue using premium features.",
		}, nil
	} else if err != nil {
		return nil, err
	}

	// Calculate days remaining
	daysRemaining := int(license.Expires.Sub(time.Now()).Hours() / 24)

	return &LicenseInfo{
		Status:        "Active",
		Plan:          license.Plan,
		Customer:      license.Customer,
		Expires:       license.Expires,
		DaysRemaining: daysRemaining,
		Features:      license.Features,
	}, nil
}

// DeactivateLicense removes the current license
func (lm *LicenseManager) DeactivateLicense() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	licensePath := filepath.Join(homeDir, ".nessi", "license.json")
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return ErrLicenseNotFound
	}

	if err := os.Remove(licensePath); err != nil {
		return fmt.Errorf("failed to remove license file: %w", err)
	}

	return nil
}

// PrintLicenseInfo prints license information in a formatted way
func (lm *LicenseManager) PrintLicenseInfo() error {
	info, err := lm.GetLicenseInfo()
	if err != nil {
		return err
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(bold("Nessi License Information"))
	fmt.Println(strings.Repeat("-", 50))

	// Status with color
	statusColor := green
	if info.Status == "Expired" {
		statusColor = red
	} else if info.Status == "Not Licensed" {
		statusColor = yellow
	} else if info.Status == "Trial" {
		statusColor = cyan
	}
	fmt.Printf("%s: %s\n", bold("Status"), statusColor(info.Status))

	// Plan
	fmt.Printf("%s: %s\n", bold("Plan"), info.Plan)

	// Customer if available
	if info.Customer != "" {
		fmt.Printf("%s: %s\n", bold("Customer"), info.Customer)
	}

	// Expiration if available
	if !info.Expires.IsZero() {
		expiresStr := info.Expires.Format("January 2, 2006")
		if info.Status == "Active" || info.Status == "Trial" {
			fmt.Printf("%s: %s (%d days remaining)\n", bold("Expires"), expiresStr, info.DaysRemaining)
		} else {
			fmt.Printf("%s: %s\n", bold("Expired"), expiresStr)
		}
	}

	// Trial count if on trial
	if info.IsTrial {
		// Get machine ID
		machineID, err := GetMachineID()
		if err == nil {
			// Get trial count for this machine
			trialCount, err := getExistingTrials(machineID)
			if err == nil {
				fmt.Printf("%s: %d of 2\n", bold("Trial Count"), trialCount)
			}
		}
	}

	// Features if available
	if len(info.Features) > 0 {
		fmt.Printf("%s:\n", bold("Licensed Features"))
		for _, feature := range info.Features {
			if desc, ok := FeatureDescriptions[feature]; ok {
				fmt.Printf("  - %s: %s\n", feature, desc)
			} else {
				fmt.Printf("  - %s\n", feature)
			}
		}
	}

	// Plan features or trial features
	if info.IsTrial {
		// For trial, show all premium features
		fmt.Printf("%s:\n", bold("Trial Features (All Premium Features)"))
		// Show all features from enterprise plan
		for _, feature := range PlanFeatures[PlanEnterprise] {
			if desc, ok := FeatureDescriptions[feature]; ok {
				fmt.Printf("  - %s: %s\n", feature, desc)
			} else {
				fmt.Printf("  - %s\n", feature)
			}
		}
	} else if planFeatures, ok := PlanFeatures[info.Plan]; len(planFeatures) > 0 && ok {
		fmt.Printf("%s:\n", bold("Plan Features"))
		for _, feature := range planFeatures {
			if desc, ok := FeatureDescriptions[feature]; ok {
				fmt.Printf("  - %s: %s\n", feature, desc)
			} else {
				fmt.Printf("  - %s\n", feature)
			}
		}
	}

	// Message if available
	if info.Message != "" {
		fmt.Println()
		fmt.Println(info.Message)
	}

	// If not on trial, show trial info
	if !info.IsTrial && info.Status != "Active" {
		// Check if trials are available for this machine
		machineID, err := GetMachineID()
		if err == nil {
			trialCount, err := getExistingTrials(machineID)
			if err == nil && trialCount < 2 {
				fmt.Println()
				fmt.Println(bold("Free Trial"))
				fmt.Println("Try all premium features free for 1 month.")
				fmt.Println("Run 'nessi license trial' to start your free trial.")
				fmt.Printf("Trial count: %d of 2 used for this machine.\n", trialCount)
			} else if trialCount >= 2 {
				fmt.Println()
				fmt.Println(bold("Free Trial"))
				fmt.Println("You have used the maximum number of trials (2) for this machine.")
				fmt.Println("Please purchase a license to continue using premium features.")
			}
		}
	}

	return nil
}

// LicenseInfo represents information about a license
type LicenseInfo struct {
	Status        string
	Plan          string
	Customer      string
	Expires       time.Time
	DaysRemaining int
	Features      []string
	Message       string
	IsTrial       bool
}

// Helper functions

// isValidLicenseKeyFormat checks if a license key has a valid format
func isValidLicenseKeyFormat(key string) bool {
	// Simple check: JWT format has 2 dots
	parts := strings.Split(key, ".")
	return len(parts) == 3 && len(key) > 20
}

// saveLicense saves a license to the license file
func saveLicense(license *License) error {
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

	// License file path
	licensePath := filepath.Join(nessiDir, "license.json")

	// Marshal license to JSON
	data, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal license: %w", err)
	}

	// Write license file
	if err := os.WriteFile(licensePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write license file: %w", err)
	}

	return nil
}
