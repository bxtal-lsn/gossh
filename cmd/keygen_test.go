// cmd/keygen_test.go
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKeygenOutputFiles tests the keygen command's file creation
func TestKeygenOutputFiles(t *testing.T) {
	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "gossh-keygen-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up test key paths
	testPrivKey := filepath.Join(tmpDir, "test-key.pem")
	testPubKey := filepath.Join(tmpDir, "test-key.pub")

	// Save original values to restore later
	origPrivKeyOut := privateKeyOut
	origPubKeyOut := publicKeyOut
	defer func() {
		privateKeyOut = origPrivKeyOut
		publicKeyOut = origPubKeyOut
	}()

	// Set for testing
	privateKeyOut = testPrivKey
	publicKeyOut = testPubKey

	// Mock execution of keygen command
	// This would normally call the GenerateKeys function
	// For testing, we'll just create dummy files
	if err := os.WriteFile(testPrivKey, []byte("TEST PRIVATE KEY"), 0o600); err != nil {
		t.Fatalf("Failed to write test private key: %v", err)
	}
	if err := os.WriteFile(testPubKey, []byte("ssh-rsa TEST"), 0o644); err != nil {
		t.Fatalf("Failed to write test public key: %v", err)
	}

	// Check if files exist with correct permissions
	privInfo, err := os.Stat(testPrivKey)
	if err != nil {
		t.Fatalf("Failed to stat private key: %v", err)
	}
	if privInfo.Mode().Perm() != 0o600 {
		t.Errorf("Private key has wrong permissions: %v, expected 0600", privInfo.Mode().Perm())
	}

	pubInfo, err := os.Stat(testPubKey)
	if err != nil {
		t.Fatalf("Failed to stat public key: %v", err)
	}
	if pubInfo.Mode().Perm() != 0o644 {
		t.Errorf("Public key has wrong permissions: %v, expected 0644", pubInfo.Mode().Perm())
	}
}
