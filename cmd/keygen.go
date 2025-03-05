package cmd

import (
	"fmt"
	"os"

	"github.com/bxtal-lsn/gossh/pkg/ssh"
	"github.com/spf13/cobra"
)

var (
	privateKeyOut string
	publicKeyOut  string
	keyBits       int
	keyComment    string
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate SSH key pairs",
	Long: `The keygen utility creates and manages SSH key pairs.

Examples:
  # Generate a default key pair (id_rsa and id_rsa.pub)
  gossh keygen

  # Specify output files
  gossh keygen --private-key mykey.pem --public-key mykey.pub

  # Generate keys with specific parameters
  gossh keygen --private-key server.pem --public-key server.pub --comment "server-key"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating SSH key pair...")

		// Generate the keys
		privateKey, publicKey, err := ssh.GenerateKeys()
		if err != nil {
			fmt.Printf("Error generating keys: %s\n", err)
			os.Exit(1)
		}

		// Save the private key
		if err = os.WriteFile(privateKeyOut, privateKey, 0o600); err != nil {
			fmt.Printf("Error writing private key: %s\n", err)
			os.Exit(1)
		}

		// Save the public key
		if err = os.WriteFile(publicKeyOut, publicKey, 0o644); err != nil {
			fmt.Printf("Error writing public key: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("SSH key pair generated successfully:")
		fmt.Printf("Private key: %s\n", privateKeyOut)
		fmt.Printf("Public key: %s\n", publicKeyOut)
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	// Define flags for the keygen command
	keygenCmd.Flags().StringVarP(&privateKeyOut, "private-key", "k", "id_rsa", "Output file for private key")
	keygenCmd.Flags().StringVarP(&publicKeyOut, "public-key", "p", "id_rsa.pub", "Output file for public key")
	keygenCmd.Flags().IntVarP(&keyBits, "bits", "b", 4096, "Number of bits in the key")
	keygenCmd.Flags().StringVarP(&keyComment, "comment", "c", "", "Comment to include in the public key")

	// Note: The current implementation doesn't use keyBits and keyComment yet,
	// but they are included here for future enhancement
}
