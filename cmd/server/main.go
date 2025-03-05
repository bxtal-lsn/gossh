package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bxtal-lsn/gossh/pkg/ssh"
)

func main() {
	var (
		privateKeyPath string
		pubKeyPath     string
		err            error
	)

	flag.StringVar(&privateKeyPath, "key", "server.pem", "Path to the server private key")
	flag.StringVar(&pubKeyPath, "authorized-keys", "authorized_keys", "Path to the authorized keys file")
	flag.Parse()

	serverKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Failed to load server key: %v\n", err)
		os.Exit(1)
	}

	authorizedKeysBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		fmt.Printf("Failed to load authorized keys: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting SSH server...")
	if err = ssh.StartServer(serverKeyBytes, authorizedKeysBytes); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
