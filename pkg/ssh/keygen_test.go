// pkg/ssh/keygen_test.go
package ssh

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateKeys(t *testing.T) {
	// Test key generation
	privateKey, publicKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Key generation failed: %v", err)
	}

	// Verify private key format
	block, _ := pem.Decode(privateKey)
	if block == nil {
		t.Fatal("Failed to decode private key PEM block")
	}
	if block.Type != "RSA PRIVATE KEY" {
		t.Errorf("Expected RSA PRIVATE KEY, got %s", block.Type)
	}

	// Parse the private key to ensure it's valid
	_, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Generated private key is invalid: %v", err)
	}

	// Verify public key format (should be in authorized_keys format)
	if !bytes.HasPrefix(publicKey, []byte("ssh-rsa ")) {
		t.Errorf("Public key doesn't start with 'ssh-rsa': %s", publicKey)
	}

	// Parse the public key to ensure it's valid
	_, _, _, _, err = ssh.ParseAuthorizedKey(publicKey)
	if err != nil {
		t.Fatalf("Generated public key is invalid: %v", err)
	}
}

func TestGenerateKeysMatchingPair(t *testing.T) {
	// Generate a key pair
	privateKeyBytes, publicKeyBytes, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Key generation failed: %v", err)
	}

	// Parse the private key
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		t.Fatal("Failed to decode private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
	if err != nil {
		t.Fatalf("Failed to parse public key: %v", err)
	}

	// Generate a public key from the private key
	derivedPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to derive public key: %v", err)
	}

	// Compare the keys
	if !bytes.Equal(pubKey.Marshal(), derivedPublicKey.Marshal()) {
		t.Error("Public key doesn't match the one derived from private key")
	}
}
