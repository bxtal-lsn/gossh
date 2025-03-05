// cmd/server_test.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestServerValidation tests the server command validation
func TestServerValidation(t *testing.T) {
	// Save original flags to restore later
	origKeyPath := serverKeyPath
	origAuthKeysPath := pubKeyPath
	defer func() {
		serverKeyPath = origKeyPath
		pubKeyPath = origAuthKeysPath
	}()

	// Set up test cases
	tests := []struct {
		name        string
		keyPath     string
		authKeyPath string
		wantErr     bool
	}{
		{
			name:        "all flags provided",
			keyPath:     "/path/to/key",
			authKeyPath: "/path/to/authorized_keys",
			wantErr:     false,
		},
		{
			name:        "missing key path",
			keyPath:     "",
			authKeyPath: "/path/to/authorized_keys",
			wantErr:     true,
		},
		{
			name:        "missing auth keys path",
			keyPath:     "/path/to/key",
			authKeyPath: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags for this test case
			serverKeyPath = tt.keyPath
			pubKeyPath = tt.authKeyPath

			// Run validation
			err := validateServerFlags()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateServerFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to validate server flags (mimicking what the Run function would do)
func validateServerFlags() error {
	if serverKeyPath == "" {
		return fmt.Errorf("server key path is required")
	}
	if pubKeyPath == "" {
		return fmt.Errorf("authorized keys path is required")
	}
	return nil
}

// TestKeyFileReadability tests the ability to read key files
func TestKeyFileReadability(t *testing.T) {
	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "gossh-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	validKeyPath := filepath.Join(tmpDir, "valid-key.pem")
	validAuthKeysPath := filepath.Join(tmpDir, "valid-auth-keys")

	// Write sample content to files
	if err := os.WriteFile(validKeyPath, []byte("TEST PRIVATE KEY"), 0o600); err != nil {
		t.Fatalf("Failed to write test key file: %v", err)
	}
	if err := os.WriteFile(validAuthKeysPath, []byte("ssh-rsa TEST"), 0o600); err != nil {
		t.Fatalf("Failed to write test auth keys file: %v", err)
	}

	// Test reading the files
	_, err = os.ReadFile(validKeyPath)
	if err != nil {
		t.Errorf("Failed to read key file: %v", err)
	}

	_, err = os.ReadFile(validAuthKeysPath)
	if err != nil {
		t.Errorf("Failed to read auth keys file: %v", err)
	}
}
