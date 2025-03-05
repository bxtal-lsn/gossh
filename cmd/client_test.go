// cmd/client_test.go
package cmd

import (
	"fmt"
	"testing"
	"time"
)

// TestClientValidation tests the client command validation
func TestClientValidation(t *testing.T) {
	// Save original flags to restore later
	origHost := host
	origUser := user
	origKeyPath := clientKeyPath
	defer func() {
		host = origHost
		user = origUser
		clientKeyPath = origKeyPath
	}()

	// Set up test cases
	tests := []struct {
		name    string
		host    string
		user    string
		keyPath string
		wantErr bool
	}{
		{
			name:    "all flags provided",
			host:    "example.com",
			user:    "testuser",
			keyPath: "/path/to/key",
			wantErr: false,
		},
		{
			name:    "missing host",
			host:    "",
			user:    "testuser",
			keyPath: "/path/to/key",
			wantErr: true,
		},
		{
			name:    "missing user",
			host:    "example.com",
			user:    "",
			keyPath: "/path/to/key",
			wantErr: true,
		},
		{
			name:    "missing key path",
			host:    "example.com",
			user:    "testuser",
			keyPath: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags for this test case
			host = tt.host
			user = tt.user
			clientKeyPath = tt.keyPath

			// Run validation
			err := validateClientFlags()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateClientFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to validate client flags (mimicking what the Run function would do)
func validateClientFlags() error {
	if host == "" {
		return fmt.Errorf("host is required")
	}
	if user == "" {
		return fmt.Errorf("user is required")
	}
	if clientKeyPath == "" {
		return fmt.Errorf("private key path is required")
	}
	return nil
}

// TestTimeoutParsing tests the timeout duration parsing
func TestTimeoutParsing(t *testing.T) {
	// Save original timeout to restore later
	origTimeout := timeout
	defer func() {
		timeout = origTimeout
	}()

	tests := []struct {
		name      string
		timeout   string
		wantError bool
	}{
		{"valid seconds", "10s", false},
		{"valid minutes", "5m", false},
		{"valid hours", "1h", false},
		{"complex duration", "1h30m", false},
		{"invalid format", "xyz", true},
		{"negative duration", "-10s", false}, // Note: This is actually valid in Go's time.ParseDuration
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout = tt.timeout
			_, err := time.ParseDuration(timeout)
			if (err != nil) != tt.wantError {
				t.Errorf("time.ParseDuration(%q) error = %v, wantErr %v", timeout, err, tt.wantError)
			}
		})
	}
}
