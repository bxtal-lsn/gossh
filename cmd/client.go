package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	host          string
	port          string
	user          string
	clientKeyPath string
	command       string
	timeout       string
	noSpinner     bool
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to an SSH server",
	Long: `The client component provides both interactive and non-interactive SSH capabilities.

Examples:
  # Connect with interactive shell
  gossh client --host example.com --port 2022 --user admin --key id_rsa

  # Execute a command
  gossh client --host example.com --user admin --key id_rsa --cmd "ls -la"

  # Execute with timeout
  gossh client --host example.com --user admin --key id_rsa --cmd "backup.sh" --timeout 30s`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create colored output helpers
		titleColor := color.New(color.FgBlue, color.Bold).SprintFunc()
		successColor := color.New(color.FgGreen, color.Bold).SprintFunc()
		infoColor := color.New(color.FgCyan).SprintFunc()
		errorColor := color.New(color.FgRed, color.Bold).SprintFunc()
		warningColor := color.New(color.FgYellow).SprintFunc()

		// Print header
		fmt.Println(titleColor("SSH CLIENT CONNECTION"))
		fmt.Println(infoColor("⟹ ") + fmt.Sprintf("Connecting to %s@%s:%s",
			color.CyanString(user),
			color.CyanString(host),
			color.CyanString(port)))

		// Log connection details
		log.Info("Initiating SSH connection")
		log.WithFields(logrus.Fields{
			"host":    host,
			"port":    port,
			"user":    user,
			"key":     clientKeyPath,
			"timeout": timeout,
		}).Debug("Connection parameters")

		// Parse timeout duration
		timeoutDuration, err := time.ParseDuration(timeout)
		if err != nil {
			log.Error("Invalid timeout format: ", err)
			fmt.Println(errorColor("✗ Invalid timeout format: ") + err.Error())
			os.Exit(1)
		}

		// Read the private key
		log.Debug("Reading private key from: ", clientKeyPath)
		privateKeyBytes, err := os.ReadFile(clientKeyPath)
		if err != nil {
			log.Error("Failed to load private key: ", err)
			fmt.Println(errorColor("✗ Failed to load private key: ") + err.Error())
			os.Exit(1)
		}

		// Parse the private key
		log.Debug("Parsing private key")
		signer, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			log.Error("Failed to parse private key: ", err)
			fmt.Println(errorColor("✗ Failed to parse private key: ") + err.Error())
			os.Exit(1)
		}

		// Display a connection warning about host key verification
		fmt.Println(warningColor("⚠ ") + "Warning: Using InsecureIgnoreHostKey() - host won't be verified")

		// Set up SSH client configuration
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Not secure for production
			Timeout:         timeoutDuration,
		}

		// Start a spinner for connection process
		var s *spinner.Spinner
		if !noSpinner {
			s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Suffix = " Establishing connection..."
			s.Color("cyan")
			s.Start()
		}

		// Connect to the SSH server
		addr := fmt.Sprintf("%s:%s", host, port)
		log.Info("Dialing SSH server at ", addr)
		client, err := ssh.Dial("tcp", addr, config)

		// Stop the spinner regardless of connection result
		if !noSpinner {
			s.Stop()
		}

		if err != nil {
			log.Error("Failed to connect: ", err)
			fmt.Println(errorColor("✗ Connection failed: ") + err.Error())
			os.Exit(1)
		}
		fmt.Println(successColor("✓ ") + "Connected successfully to " + infoColor(addr))

		defer client.Close()

		// Create a session
		log.Debug("Creating new SSH session")
		session, err := client.NewSession()
		if err != nil {
			log.Error("Failed to create session: ", err)
			fmt.Println(errorColor("✗ Failed to create session: ") + err.Error())
			os.Exit(1)
		}
		defer session.Close()

		// Set up I/O
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if command != "" {
			// Run a specific command
			fmt.Println(infoColor("⟹ ") + "Executing command: " + color.HiWhiteString(command))
			log.Info("Executing command: ", command)

			if !noSpinner {
				s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				s.Suffix = " Running command..."
				s.Color("cyan")
				s.Start()
			}

			err = session.Run(command)

			if !noSpinner {
				s.Stop()
			}

			if err != nil {
				log.Error("Command execution failed: ", err)
				fmt.Println(errorColor("✗ Command execution failed: ") + err.Error())
				os.Exit(1)
			}
			fmt.Println(successColor("✓ ") + "Command executed successfully")
		} else {
			// Start an interactive shell
			session.Stdin = os.Stdin

			// Request PTY
			log.Debug("Requesting PTY for interactive session")
			modes := ssh.TerminalModes{
				ssh.ECHO:          1,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}

			if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
				log.Error("Failed to request PTY: ", err)
				fmt.Println(errorColor("✗ Failed to request PTY: ") + err.Error())
				os.Exit(1)
			}

			fmt.Println(infoColor("⟹ ") + "Starting interactive shell session")
			fmt.Println(infoColor("ℹ ") + "Press Ctrl+D or type 'exit' to close the connection")
			fmt.Println(strings.Repeat("─", 50))

			if err := session.Shell(); err != nil {
				log.Error("Failed to start shell: ", err)
				fmt.Println(errorColor("✗ Failed to start shell: ") + err.Error())
				os.Exit(1)
			}

			if err := session.Wait(); err != nil {
				if e, ok := err.(*ssh.ExitError); ok {
					log.Warn("Session ended with exit code: ", e.ExitStatus())
					os.Exit(e.ExitStatus())
				} else {
					log.Error("Session error: ", err)
					fmt.Println(errorColor("✗ Session error: ") + err.Error())
					os.Exit(1)
				}
			}

			// Print end of session message
			fmt.Println(strings.Repeat("─", 50))
			fmt.Println(successColor("✓ ") + "Session closed")
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Define flags for the client command
	clientCmd.Flags().StringVarP(&host, "host", "H", "localhost", "SSH server hostname")
	clientCmd.Flags().StringVarP(&port, "port", "p", "22", "SSH server port")
	clientCmd.Flags().StringVarP(&user, "user", "u", "", "SSH username")
	clientCmd.Flags().StringVarP(&clientKeyPath, "key", "k", "", "Path to private key")
	clientCmd.Flags().StringVarP(&command, "cmd", "c", "", "Command to execute (optional)")
	clientCmd.Flags().StringVarP(&timeout, "timeout", "t", "10s", "Connection timeout duration")
	clientCmd.Flags().BoolVar(&noSpinner, "no-spinner", false, "Disable spinner animation")

	// Mark required flags
	clientCmd.MarkFlagRequired("host")
	clientCmd.MarkFlagRequired("user")
	clientCmd.MarkFlagRequired("key")
}
