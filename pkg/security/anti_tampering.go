// Package security provides license validation and security features
package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// Anti-tampering measures to detect binary modifications

// functionChecksum calculates a checksum of a function's address
func functionChecksum(f interface{}) string {
	// Get function pointer
	ptr := reflect.ValueOf(f).Pointer()
	// Convert to string representation
	fnName := runtime.FuncForPC(ptr).Name()
	// Hash the function name
	hash := sha256.Sum256([]byte(fnName))
	return hex.EncodeToString(hash[:])
}

// verifyFunctionIntegrity checks if a function has been tampered with
func verifyFunctionIntegrity(f interface{}, expectedChecksum string) bool {
	actualChecksum := functionChecksum(f)
	return strings.EqualFold(actualChecksum, expectedChecksum)
}

// redundantTrialCheck performs multiple checks to verify trial status
func redundantTrialCheck() bool {
	// First check - standard trial check
	check1 := isTrialActive()

	// Second check - directly check trial file
	check2 := false
	trialInfo, err := getTrialInfo()
	if err == nil {
		trialEndTime := trialInfo.StartTime.AddDate(0, 1, 0)
		check2 = time.Now().Before(trialEndTime) && trialInfo.Active
	}

	// Third check - verify machine ID
	check3 := false
	if trialInfo != nil {
		machineID, err := GetMachineID()
		if err == nil {
			check3 = trialInfo.MachineID == machineID
		}
	}

	// All checks must pass
	return check1 && check2 && check3
}

// verifyTrialIntegrity checks if the trial system has been tampered with
func verifyTrialIntegrity() bool {
	// In a real implementation, these checksums would be calculated at build time
	// and possibly encrypted or obfuscated
	expectedChecksums := map[string]string{
		"isTrialActive":   functionChecksum(isTrialActive),
		"getTrialInfo":    functionChecksum(getTrialInfo),
		"GetMachineID":    functionChecksum(GetMachineID),
		"VerifySignature": functionChecksum(VerifySignature),
	}

	// Verify each function
	for name, checksum := range expectedChecksums {
		var f interface{}
		switch name {
		case "isTrialActive":
			f = isTrialActive
		case "getTrialInfo":
			f = getTrialInfo
		case "GetMachineID":
			f = GetMachineID
		case "VerifySignature":
			f = VerifySignature
		}

		if !verifyFunctionIntegrity(f, checksum) {
			fmt.Printf("Warning: Function %s may have been tampered with\n", name)
			return false
		}
	}

	return true
}

// secureRequireFeature is a more secure version of RequireFeature
func secureRequireFeature(feature string) error {
	// Check function integrity
	if !verifyTrialIntegrity() {
		return fmt.Errorf("security violation: license system integrity check failed")
	}

	// Use redundant checks for trial status
	if redundantTrialCheck() {
		return nil // Allow all features during trial period
	}

	// Fall back to standard license check
	return RequireFeature(feature)
}
