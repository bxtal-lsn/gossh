# gossh

A robust SSH utility suite built in Go, providing key generation, secure client connections, and server functionality.

## Overview

gossh is a command-line SSH toolkit that simplifies key management, SSH connections, and running SSH servers. It's designed with a focus on ease of use and proper security practices.

## Features

### Key Management
- Generate RSA key pairs (4096-bit) with proper file permissions
- Command-line interface for key generation

### SSH Client
- Connect to SSH servers with public key authentication
- Execute commands remotely with detailed output
- Interactive shell support with proper terminal handling
- Configurable connection timeouts

### SSH Server
- Public key authentication
- Command execution handling
- Customizable port binding
- Detailed logging capabilities

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/bxtal-lsn/gossh.git
cd gossh

# Build
go build -o gossh

# Or install to your Go bin path
go install
```

## Usage

### Key Generator

```bash
# Generate a default key pair (id_rsa and id_rsa.pub)
gossh keygen

# Specify output files
gossh keygen --private-key mykey.pem --public-key mykey.pub
```

### SSH Client

```bash
# Connect with interactive shell
gossh client --host example.com --port 2022 --user admin --key id_rsa

# Execute a command
gossh client --host example.com --user admin --key id_rsa --cmd "ls -la"

# Execute with timeout
gossh client --host example.com --user admin --key id_rsa --cmd "backup.sh" --timeout 30s
```

### SSH Server

```bash
# Start a basic server
gossh server --key server.pem --authorized-keys authorized_keys

# Configure server options
gossh server --key server.pem --authorized-keys authorized_keys --port 2222 --bind 0.0.0.0

# Run with detailed logging
gossh server --key server.pem --authorized-keys authorized_keys --log-level debug
```

## Project Structure

```
gossh/
├── cmd/                   # Command line interfaces
│   ├── client.go          # SSH client command
│   ├── keygen.go          # Key generation command
│   ├── root.go            # Root command configuration
│   └── server.go          # SSH server command
├── pkg/                   # Core packages
│   └── ssh/               # SSH functionality
│       ├── keygen.go      # Key generation
│       └── server.go      # Server implementation
├── main.go                # Application entry point
└── go.mod                 # Go module definition
```

## Development

### Building

```bash
# Simple build
go build -o gossh

# Run tests
go test ./...
```

### Testing

The project includes unit tests for the core functionality:

```bash
# Run all tests
go test ./...

# Run specific test suite
go test ./pkg/ssh -v
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure the server is running on the specified port
   - Check firewall settings and network connectivity

2. **Authentication Failures**
   - Verify that the private key corresponds to a public key in the authorized keys file
   - Check file permissions (private key should be 0600)

3. **Command Execution Failures**
   - Verify the user has appropriate permissions on the server

### Debugging

Enable verbose logging for troubleshooting:

```bash
# Client debugging
gossh client --host example.com --user admin --key id_rsa --log-level debug

# Server debugging
gossh server --key server.pem --authorized-keys authorized_keys --log-level debug
```
