// Package security provides license validation and security features
package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetMachineID returns a unique identifier for the current machine
// This is used to limit the number of trials per machine
func GetMachineID() (string, error) {
	var id string
	var err error

	switch runtime.GOOS {
	case "darwin":
		// On macOS, use the hardware UUID
		id, err = getMacOSMachineID()
	case "linux":
		// On Linux, use machine-id
		id, err = getLinuxMachineID()
	case "windows":
		// On Windows, use the MachineGuid
		id, err = getWindowsMachineID()
	default:
		// Fallback to network interfaces
		id, err = getNetworkBasedMachineID()
	}

	if err != nil {
		// Fallback to network interfaces if platform-specific method fails
		id, err = getNetworkBasedMachineID()
		if err != nil {
			return "", fmt.Errorf("failed to generate machine ID: %w", err)
		}
	}

	// Hash the ID to ensure consistent length and format
	hash := sha256.Sum256([]byte(id))
	return hex.EncodeToString(hash[:]), nil
}

// getMacOSMachineID gets the hardware UUID on macOS
func getMacOSMachineID() (string, error) {
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				uuid := strings.TrimSpace(parts[1])
				uuid = strings.Trim(uuid, "\"")
				return uuid, nil
			}
		}
	}

	return "", fmt.Errorf("could not find IOPlatformUUID")
}

// getLinuxMachineID gets the machine-id on Linux
func getLinuxMachineID() (string, error) {
	// Try /etc/machine-id first
	if id, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(id)), nil
	}

	// Try /var/lib/dbus/machine-id as fallback
	if id, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		return strings.TrimSpace(string(id)), nil
	}

	return "", fmt.Errorf("could not read machine-id")
}

// getWindowsMachineID gets the MachineGuid on Windows
func getWindowsMachineID() (string, error) {
	cmd := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Cryptography", "/v", "MachineGuid")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "MachineGuid") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[len(parts)-1], nil
			}
		}
	}

	return "", fmt.Errorf("could not find MachineGuid")
}

// getNetworkBasedMachineID generates a machine ID based on network interfaces
// This is a fallback method if platform-specific methods fail
func getNetworkBasedMachineID() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	h := sha256.New()
	for _, iface := range interfaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			io.WriteString(h, iface.HardwareAddr.String())
		}
	}

	// Add hostname for additional uniqueness
	hostname, err := os.Hostname()
	if err == nil {
		io.WriteString(h, hostname)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
