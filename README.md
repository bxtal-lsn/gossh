# gossh

A robust, extensible SSH utility suite built in Go, providing key generation, secure client connections, and server functionality with a focus on automation and DevOps workflows.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Components](#components)
  - [Key Generator](#key-generator)
  - [SSH Client](#ssh-client)
  - [SSH Server](#ssh-server)
- [Usage Examples](#usage-examples)
- [Configuration](#configuration)
- [Security Considerations](#security-considerations)
- [Integration Examples](#integration-examples)
- [Development](#development)
  - [Project Structure](#project-structure)
  - [Building from Source](#building-from-source)
  - [Testing](#testing)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview

gossh is a comprehensive suite of SSH utilities designed for modern DevOps workflows and infrastructure automation. Unlike traditional SSH clients and servers, this tool is built as a programmable library and set of composable commands that can be extended, customized, and integrated into broader automation pipelines.

This toolkit is particularly valuable for teams working with distributed systems, container orchestration, multi-cloud deployments, and infrastructure-as-code approaches where secure shell access needs to be automated and auditable.

## Features

### Key Management
- Generate RSA key pairs (4096-bit default) with proper file permissions
- Support for key format conversion and validation
- Key rotation capabilities for infrastructure security compliance

### SSH Client
- Command execution on remote hosts with structured output handling
- Interactive shell support with terminal emulation
- Batch mode for scripted operations across multiple hosts
- Connection pooling for efficient multiple command execution
- Built-in retry logic and timeout handling

### SSH Server
- Public key authentication with customizable authorization policies
- Command whitelisting and restricted execution environments
- Audit logging of all connections and commands
- Graceful shutdown and connection handling
- Custom handler support for specialized use cases

### General
- Cross-platform compatibility (Linux, macOS, Windows)
- Consistent, easy-to-use CLI interface
- Structured logging for better debugging and monitoring
- Comprehensive error reporting

## Installation

### Using Go Install

```bash
# Install all components
go install github.com/yourusername/ssh-tool/cmd/...@latest

# Or install individual components
go install github.com/yourusername/ssh-tool/cmd/keygen@latest
go install github.com/yourusername/ssh-tool/cmd/client@latest
go install github.com/yourusername/ssh-tool/cmd/server@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/ssh-tool.git
cd ssh-tool

# Build all components
make build

# Or build individual components
make build-keygen
make build-client
make build-server
```

### Binary Releases

Pre-compiled binaries for major platforms are available on the [Releases page](https://github.com/yourusername/ssh-tool/releases).

## Components

### Key Generator

The `ssh-keygen` utility creates and manages SSH key pairs.

```bash
# Generate a default key pair (id_rsa and id_rsa.pub)
ssh-keygen

# Specify output files
ssh-keygen -private-key mykey.pem -public-key mykey.pub

# Generate keys with specific parameters
ssh-keygen -bits 8192 -comment "deployment-key"
```

### SSH Client

The client component provides both interactive and non-interactive SSH capabilities.

```bash
# Connect with interactive shell
ssh-client -host example.com -port 2022 -user admin -key id_rsa

# Execute a command
ssh-client -host example.com -user admin -key id_rsa -cmd "ls -la"

# Execute with timeout
ssh-client -host example.com -user admin -key id_rsa -cmd "backup.sh" -timeout 30s

# Output formatting options
ssh-client -host example.com -user admin -key id_rsa -cmd "ps aux" -output json
```

### SSH Server

The server component provides a customizable SSH server for automation endpoints.

```bash
# Start a basic server
ssh-server -key server.pem -authorized-keys authorized_keys

# Configure server options
ssh-server -key server.pem -authorized-keys authorized_keys -port 2222 -allowed-commands "uptime,df,free"

# Run with detailed logging
ssh-server -key server.pem -authorized-keys authorized_keys -log-level debug
```

## Usage Examples

### Automation Example: Server Health Check

```bash
#!/bin/bash
# Health check across multiple servers

SERVERS="web1.example.com web2.example.com db1.example.com"
COMMAND="uptime && free -m && df -h"

for server in $SERVERS; do
  echo "=== Checking $server ==="
  ssh-client -host $server -user monitor -key monitor.pem -cmd "$COMMAND" -timeout 5s
done
```

### Kubernetes Node Access

```bash
# List nodes and access one for debugging
NODES=$(kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="InternalIP")].address}')

for node in $NODES; do
  echo "Node: $node"
  ssh-client -host $node -user admin -key k8s-node-key.pem -cmd "crictl ps"
done
```

### Restricted Command Server

Set up a server that only allows specific commands for automated monitoring:

```bash
# Create an allowed_commands file
echo "uptime,df -h,free -m,cat /proc/loadavg" > allowed_commands.txt

# Start the server with command restrictions
ssh-server -key server.pem -authorized-keys monitor-keys -allowed-commands-file allowed_commands.txt
```

## Configuration

### Client Configuration File

Create a `~/.ssh-tool/config.yaml` file for default settings:

```yaml
default_user: admin
default_key_path: ~/.ssh/automation_key
connection_timeout: 30s
retry_attempts: 3
known_hosts_file: ~/.ssh-tool/known_hosts
log_level: info

hosts:
  - name: prod-web
    hostname: web.production.example.com
    user: webadmin
    key_path: ~/.ssh/production.pem
  
  - name: staging
    hostname: staging.example.com
    port: 2222
```

Then use with:

```bash
ssh-client -host prod-web
```

### Server Configuration

Create a `server-config.yaml` file:

```yaml
port: 2022
host_key: /etc/ssh-tool/server.pem
authorized_keys_file: /etc/ssh-tool/authorized_keys
log_file: /var/log/ssh-tool/server.log
log_level: info

# Command restrictions
allow_shell: false
allowed_commands:
  - uptime
  - "df -h"
  - /usr/local/bin/status.sh

# Authentication settings
auth_timeout: 60s
max_auth_tries: 3
```

Then use with:

```bash
ssh-server -config server-config.yaml
```

## Security Considerations

### Key Management Best Practices

- Use 4096-bit RSA keys at minimum
- Store private keys with 0600 permissions (read/write for owner only)
- Use separate keys for different environments or purposes
- Implement a key rotation policy

### Server Hardening

- Run the SSH server with minimal privileges
- Use command restrictions when possible
- Implement IP-based access controls
- Consider using chroot environments for further isolation

### Audit Logging

Enable detailed audit logging for compliance and security monitoring:

```bash
ssh-server -key server.pem -authorized-keys authorized_keys -audit-log /var/log/ssh-audit.log -log-format json
```

## Integration Examples

### CI/CD Pipeline Integration

```yaml
# GitHub Actions example
name: Deploy Application

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install SSH Tool
        run: go install github.com/yourusername/ssh-tool/cmd/client@latest
      
      - name: Setup SSH Key
        run: |
          echo "${{ secrets.DEPLOY_KEY }}" > deploy_key.pem
          chmod 600 deploy_key.pem
      
      - name: Deploy to Production
        run: |
          ssh-client -host ${{ secrets.PROD_HOST }} -user deployer -key deploy_key.pem -cmd "./deploy.sh ${{ github.sha }}"
```

### Integration with Go Applications

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/yourusername/ssh-tool/pkg/ssh"
)

func main() {
	// Create a new SSH client
	client, err := ssh.NewClient(&ssh.ClientConfig{
		Host:       "example.com",
		Port:       "22",
		User:       "admin",
		PrivateKey: "/path/to/key.pem",
		Timeout:    "30s",
	})
	if err != nil {
		log.Fatalf("Failed to create SSH client: %v", err)
	}
	defer client.Close()
	
	// Execute a command
	output, err := client.Execute("uptime")
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
	
	fmt.Println("Server uptime:", output)
}
```

## Development

### Project Structure

```
ssh-tool/
├── cmd/                   # Command line interfaces
│   ├── client/            # SSH client
│   ├── keygen/            # Key generation utility
│   └── server/            # SSH server
├── pkg/                   # Reusable packages
│   ├── ssh/               # Core SSH functionality
│   │   ├── keygen.go      # Key generation
│   │   ├── server.go      # Server implementation
│   │   └── client.go      # Client implementation
│   ├── config/            # Configuration handling
│   └── util/              # Utility functions
├── internal/              # Internal packages
│   ├── auth/              # Authentication logic
│   └── terminal/          # Terminal handling
├── examples/              # Example implementations
├── test/                  # Test utilities and fixtures
└── docs/                  # Documentation
```

### Building from Source

Prerequisites:
- Go 1.18 or later
- Make (optional, for using the Makefile)

```bash
# Clone the repository
git clone https://github.com/yourusername/ssh-tool.git
cd ssh-tool

# Install dependencies
go mod download

# Build
make build

# The binaries will be in ./bin directory
```

### Testing

```bash
# Run all tests
make test

# Run specific test suite
go test ./pkg/ssh -v

# Run integration tests (requires Docker)
make integration-test
```

## Advanced Usage

### Custom Command Handlers

You can extend the SSH server with custom command handlers:

```go
package main

import (
	"github.com/yourusername/ssh-tool/pkg/ssh"
)

func main() {
	server := ssh.NewServer()
	
	// Register a custom handler for a specific command
	server.RegisterHandler("status", func(args []string) (string, error) {
		// Custom implementation
		return "System status: OK", nil
	})
	
	// Start the server
	server.Start()
}
```

### File Transfer Support

Basic SCP-like file transfer capability:

```bash
# Upload a file
ssh-client -host example.com -user admin -key id_rsa -upload local_file.txt:/remote/path/file.txt

# Download a file
ssh-client -host example.com -user admin -key id_rsa -download /remote/file.txt:./local_file.txt
```

### Multi-Host Operation

Execute commands across multiple hosts in parallel:

```bash
ssh-client -hosts "web1.example.com,web2.example.com,web3.example.com" -user admin -key id_rsa -cmd "service nginx reload" -parallel 3
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure the server is running and listening on the specified port
   - Check firewall settings and network connectivity

2. **Authentication Failures**
   - Verify that the private key corresponds to an authorized public key on the server
   - Check file permissions (private key should be 0600)
   - Ensure the username is correct

3. **Command Execution Failures**
   - Check if the command is allowed by server restrictions
   - Verify the user has appropriate permissions on the server

### Debugging

Enable verbose logging for troubleshooting:

```bash
# Client debugging
ssh-client -host example.com -user admin -key id_rsa -cmd "uptime" -log-level debug

# Server debugging
ssh-server -key server.pem -authorized-keys authorized_keys -log-level debug
```
