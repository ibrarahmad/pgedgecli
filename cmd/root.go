package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// Version of the CLI
const Version = "1.0.0"

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pgedgecli",
		Short: "A CLI for managing PostgreSQL cluster",
		Long:  `pgEdgeCLI is a command-line tool to manage pgEdge's PostgreSQL cluster`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Println("Initializing CLI...")
		},
	}

	// Disable the default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands in the desired order
	rootCmd.AddCommand(NewClusterCommand())
	rootCmd.AddCommand(NewSpockCommand())
	rootCmd.AddCommand(NewVersionCommand())

	return rootCmd
}

// NewVersionCommand adds the `version` command to display CLI version
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of pgedgecli",
		Long:  `All software has versions. This is pgedgecli's version.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("pgedgecli version %s\n", Version)
		},
	}
}
