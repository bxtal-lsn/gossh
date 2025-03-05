package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bxtal-lsn/gossh/pkg/ssh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	serverKeyPath string
	pubKeyPath    string
	serverPort    string
	bindAddress   string
	allowedCmds   string
	noColor       bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an SSH server",
	Long: `The server component provides a customizable SSH server for automation endpoints.

Examples:
  # Start a basic server
  gossh server --key server.pem --authorized-keys authorized_keys

  # Configure server options
  gossh server --key server.pem --authorized-keys authorized_keys --port 2222

  # Run with detailed logging
  gossh server --key server.pem --authorized-keys authorized_keys --log-level debug`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure colors based on the noColor flag
		if noColor {
			color.NoColor = true
		}

		// Create colored output helpers
		successColor := color.New(color.FgGreen, color.Bold).SprintFunc()
		infoColor := color.New(color.FgCyan).SprintFunc()
		errorColor := color.New(color.FgRed, color.Bold).SprintFunc()

		// Log the start of server initialization
		log.Info("Initializing SSH server...")

		// Read the server key
		log.Debug("Reading private key from: ", serverKeyPath)
		serverKeyBytes, err := os.ReadFile(serverKeyPath)
		if err != nil {
			log.Error("Failed to load server key: ", err)
			fmt.Println(errorColor("✗ Failed to load server key: ") + err.Error())
			os.Exit(1)
		}
		fmt.Println(successColor("✓ ") + "Server key loaded from " + infoColor(serverKeyPath))

		// Read the authorized keys
		log.Debug("Reading authorized keys from: ", pubKeyPath)
		authorizedKeysBytes, err := os.ReadFile(pubKeyPath)
		if err != nil {
			log.Error("Failed to load authorized keys: ", err)
			fmt.Println(errorColor("✗ Failed to load authorized keys: ") + err.Error())
			os.Exit(1)
		}
		fmt.Println(successColor("✓ ") + "Authorized keys loaded from " + infoColor(pubKeyPath))

		// Print allowed commands if specified
		if allowedCmds != "" {
			fmt.Println(infoColor("ℹ ") + "Restricted to commands: " + allowedCmds)
		} else {
			fmt.Println(infoColor("ℹ ") + "No command restrictions applied")
		}

		// Print server configuration
		fmt.Println()
		fmt.Println(successColor("→ ") + "Starting SSH server with configuration:")
		fmt.Printf("  • Bind Address: %s\n", infoColor(bindAddress))
		fmt.Printf("  • Port: %s\n", infoColor(serverPort))
		fmt.Printf("  • Private Key: %s\n", infoColor(serverKeyPath))
		fmt.Printf("  • Authorized Keys: %s\n", infoColor(pubKeyPath))
		fmt.Println()

		// Simulate server startup countdown for visual appeal
		fmt.Print("Starting server in: ")
		for i := 3; i > 0; i-- {
			fmt.Print(color.YellowString("%d... ", i))
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println(successColor("Launched!"))

		// Actually start the server
		log.Info("SSH server starting on ", bindAddress, ":", serverPort)
		if err = ssh.StartServer(serverKeyBytes, authorizedKeysBytes); err != nil {
			log.Error("Server error: ", err)
			fmt.Println(errorColor("\n✗ Server failed: ") + err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Define flags for the server command
	serverCmd.Flags().StringVarP(&serverKeyPath, "key", "k", "server.pem", "Path to the server private key")
	serverCmd.Flags().StringVarP(&pubKeyPath, "authorized-keys", "a", "authorized_keys", "Path to the authorized keys file")
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "2022", "Port for the SSH server to listen on")
	serverCmd.Flags().StringVarP(&bindAddress, "bind", "b", "0.0.0.0", "Address to bind the SSH server to")
	serverCmd.Flags().StringVar(&allowedCmds, "allowed-commands", "", "Comma-separated list of allowed commands (empty for unrestricted)")
	serverCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")

	// Mark required flags
	serverCmd.MarkFlagRequired("key")
	serverCmd.MarkFlagRequired("authorized-keys")
}

