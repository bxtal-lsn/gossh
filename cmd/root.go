package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global logger instance
var log = logrus.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gossh",
	Short: "A comprehensive SSH toolset built in Go",
	Long: `gossh is a robust, extensible SSH utility suite providing key generation,
secure client connections, and server functionality with a focus on
automation and DevOps workflows.

Complete documentation is available at https://github.com/bxtal-lsn/gossh`,
	// This will run before any subcommand
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Print a fancy header
		color.New(color.FgHiCyan, color.Bold).Println("┌─────────────────────────────┐")
		color.New(color.FgHiCyan, color.Bold).Println("│        GoSSH Toolset        │")
		color.New(color.FgHiCyan, color.Bold).Println("└─────────────────────────────┘")
		fmt.Println()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		color.Red("Error: %s", err)
		os.Exit(1)
	}
}

func init() {
	// Configure logging
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Set the default log output to stdout
	log.SetOutput(os.Stdout)

	// Define persistent flags for root command
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress all output except errors")

	// Set up a hook to adjust log level based on flag
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Get the log level from the flag
	logLevel, _ := rootCmd.PersistentFlags().GetString("log-level")

	// Set the log level
	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Debug logging enabled")
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Check if quiet mode is enabled
	quiet, _ := rootCmd.PersistentFlags().GetBool("quiet")
	if quiet {
		// In quiet mode, only show errors
		log.SetLevel(logrus.ErrorLevel)
		// Also disable the header by setting PersistentPreRun to nil
		rootCmd.PersistentPreRun = nil
	}
}
