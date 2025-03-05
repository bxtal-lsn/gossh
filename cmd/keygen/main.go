package main

import (
	"fmt"
	"os"

	"github.com/bxtal-lsn/gossh/pkg/ssh"
)

func main() {
	var (
		privateKey []byte
		publicKey  []byte
		err        error
	)

	if privateKey, publicKey, err = ssh.GenerateKeys(); err != nil {
		fmt.Printf("Error generating keys: %s\n", err)
		os.Exit(1)
	}

	if err = os.WriteFile("id_rsa", privateKey, 0o600); err != nil {
		fmt.Printf("Error writing private key: %s\n", err)
		os.Exit(1)
	}

	if err = os.WriteFile("id_rsa.pub", publicKey, 0o644); err != nil {
		fmt.Printf("Error writing public key: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("SSH key pair generated successfully:")
	fmt.Println("Private key: id_rsa")
	fmt.Println("Public key: id_rsa.pub")
}
