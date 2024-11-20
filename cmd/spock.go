package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pgedgecli/pkg/spock"
)

func NewSpockCommand() *cobra.Command {
	spockCmd := &cobra.Command{
		Use:   "spock",
		Short: "Manage Spock replication",
		Long:  "Commands to manage Spock logical replication such as creating, checking, and dropping replication sets.",
	}

	// Subcommands for spock
	spockCmd.AddCommand(newCreateSetCommand())
	spockCmd.AddCommand(newDropSetCommand())
	spockCmd.AddCommand(newCheckSetCommand())

	return spockCmd
}

func newCreateSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-set",
		Short: "Create a replication set",
		Run: func(cmd *cobra.Command, args []string) {
			set := cmd.Flag("set").Value.String()
			if err := spock.CreateSet(set); err != nil {
				fmt.Printf("Failed to create replication set: %v\n", err)
			} else {
				fmt.Println("Replication set created successfully.")
			}
		},
	}
}

func newDropSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "drop-set",
		Short: "Drop a replication set",
		Run: func(cmd *cobra.Command, args []string) {
			set := cmd.Flag("set").Value.String()
			if err := spock.DropSet(set); err != nil {
				fmt.Printf("Failed to drop replication set: %v\n", err)
			} else {
				fmt.Println("Replication set dropped successfully.")
			}
		},
	}
}

func newCheckSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check-set",
		Short: "Check the status of a replication set",
		Run: func(cmd *cobra.Command, args []string) {
			set := cmd.Flag("set").Value.String()
			status, err := spock.CheckSet(set)
			if err != nil {
				fmt.Printf("Failed to check replication set: %v\n", err)
			} else {
				fmt.Printf("Replication set status: %s\n", status)
			}
		},
	}
}

