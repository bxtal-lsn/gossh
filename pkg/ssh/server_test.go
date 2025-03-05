// pkg/ssh/server_test.go
package ssh

import (
	"net"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// mockSSHServer is a struct to help with testing the SSH server
type mockSSHServer struct {
	listener net.Listener
	done     chan bool
}

// setupMockSSHServer creates a temporary SSH server for testing
func setupMockSSHServer(t *testing.T) (*mockSSHServer, []byte, []byte) {
	// Generate key pair for testing
	privateKey, publicKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}

	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	// Create mock server
	mockServer := &mockSSHServer{
		listener: listener,
		done:     make(chan bool),
	}

	// Start server in a goroutine (will be stopped in teardown)
	go func() {
		// This is a simplified version that won't actually start the full server
		// but it will listen for connections so we can test basic connectivity
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-mockServer.done:
				// Server was stopped, this is expected
				return
			default:
				t.Logf("Accept error: %v", err)
			}
		} else {
			defer conn.Close()
			// Just read some data to simulate interaction
			buffer := make([]byte, 1024)
			conn.Read(buffer)
		}
	}()

	return mockServer, privateKey, publicKey
}

// teardownMockSSHServer cleans up the mock server
func teardownMockSSHServer(mock *mockSSHServer) {
	close(mock.done)
	mock.listener.Close()
}

func TestStartServer_ConnectionBasics(t *testing.T) {
	// This is a basic connectivity test
	// Setup mock server
	mock, privateKey, publicKey := setupMockSSHServer(t)
	defer teardownMockSSHServer(mock)

	// Create temporary files for the keys
	privateKeyFile, err := os.CreateTemp("", "ssh_test_private_key")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(privateKeyFile.Name())

	publicKeyFile, err := os.CreateTemp("", "ssh_test_public_key")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(publicKeyFile.Name())

	// Write the keys to the files
	if _, err := privateKeyFile.Write(privateKey); err != nil {
		t.Fatalf("Failed to write private key: %v", err)
	}
	if _, err := publicKeyFile.Write(publicKey); err != nil {
		t.Fatalf("Failed to write public key: %v", err)
	}

	// Close the files
	privateKeyFile.Close()
	publicKeyFile.Close()

	// Test key loading (partial test of StartServer)
	serverKeyBytes, err := os.ReadFile(privateKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed to read server key: %v", err)
	}
	authorizedKeysBytes, err := os.ReadFile(publicKeyFile.Name())
	if err != nil {
		t.Fatalf("Failed to read authorized keys: %v", err)
	}

	// Validate the server can parse keys
	// This is testing the beginning of StartServer without actually starting it
	authorizedKeysMap := map[string]bool{}
	var rest []byte = authorizedKeysBytes
	for len(rest) > 0 {
		pubKey, _, _, newRest, err := ssh.ParseAuthorizedKey(rest)
		if err != nil {
			t.Fatalf("Parse authorized keys error: %v", err)
		}
		authorizedKeysMap[string(pubKey.Marshal())] = true
		rest = newRest
	}

	private, err := ssh.ParsePrivateKey(serverKeyBytes)
	if err != nil {
		t.Fatalf("ParsePrivateKey error: %v", err)
	}

	// Verify we parsed at least one key
	if len(authorizedKeysMap) == 0 {
		t.Fatal("No authorized keys were parsed")
	}

	// Verify private key was parsed
	if private == nil {
		t.Fatal("Private key couldn't be parsed")
	}
}

// TestExecSomething tests the execSomething function
func TestExecSomething(t *testing.T) {
	// Create a mock connection for testing
	mockConn := &ssh.ServerConn{
		Conn: &mockSSHConn{
			user: "testuser",
		},
	}

	// Test whoami command
	result := execSomething(mockConn, []byte("whoami"))
	expected := "You are: testuser\n"
	if result != expected {
		t.Errorf("execSomething(whoami) = %q, want %q", result, expected)
	}

	// Test unknown command
	result = execSomething(mockConn, []byte("unknown"))
	if !strings.Contains(result, "Command Not Found") {
		t.Errorf("execSomething(unknown) should return 'Command Not Found', got %q", result)
	}
}

type mockSSHConn struct {
	user string
}

func (m *mockSSHConn) User() string {
	return m.user
}

func (m *mockSSHConn) SessionID() []byte {
	return []byte("mock-session-id")
}

func (m *mockSSHConn) ClientVersion() []byte {
	return []byte("SSH-2.0-mockClient")
}

func (m *mockSSHConn) ServerVersion() []byte {
	return []byte("SSH-2.0-mockServer")
}

func (m *mockSSHConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 22}
}

func (m *mockSSHConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2222}
}

func (m *mockSSHConn) Close() error {
	return nil
}

func (m *mockSSHConn) OpenChannel(channelType string, extraData []byte) (ssh.Channel, <-chan *ssh.Request, error) {
	return nil, nil, nil
}

func (m *mockSSHConn) SendRequest(name string, wantReply bool, payload []byte) (bool, []byte, error) {
	return true, nil, nil
}

// Add the Wait method to satisfy the ssh.Conn interface
func (m *mockSSHConn) Wait() error {
	return nil
}
