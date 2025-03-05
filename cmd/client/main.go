package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	var (
		host           string
		port           string
		user           string
		privateKeyPath string
		cmd            string
		err            error
	)

	flag.StringVar(&host, "host", "localhost", "SSH server hostname")
	flag.StringVar(&port, "port", "2022", "SSH server port")
	flag.StringVar(&user, "user", "user", "SSH username")
	flag.StringVar(&privateKeyPath, "key", "id_rsa", "Path to private key")
	flag.StringVar(&cmd, "cmd", "", "Command to execute (optional)")
	flag.Parse()

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Failed to load private key: %v\n", err)
		os.Exit(1)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Failed to parse private key: %v\n", err)
		os.Exit(1)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not recommended for production
		Timeout:         time.Second * 10,
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		os.Exit(1)
	}
	defer session.Close()

	// Set up I/O
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if cmd != "" {
		// Run a specific command
		err = session.Run(cmd)
		if err != nil {
			fmt.Printf("Failed to run command: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Start an interactive shell
		session.Stdin = os.Stdin

		// Request PTY
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			fmt.Printf("Failed to request PTY: %v\n", err)
			os.Exit(1)
		}

		if err := session.Shell(); err != nil {
			fmt.Printf("Failed to start shell: %v\n", err)
			os.Exit(1)
		}

		if err := session.Wait(); err != nil {
			if e, ok := err.(*ssh.ExitError); ok {
				os.Exit(e.ExitStatus())
			} else {
				fmt.Printf("Session error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
