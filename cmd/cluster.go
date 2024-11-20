package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pgedgecli/pkg/cluster"
)

func NewClusterCommand() *cobra.Command {
	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage PostgreSQL clusters",
		Long:  "Commands to manage PostgreSQL clusters such as adding, removing, and listing nodes.",
	}

	// Subcommands for cluster
	clusterCmd.AddCommand(newAddNodeCommand())
	clusterCmd.AddCommand(newRemoveNodeCommand())
	clusterCmd.AddCommand(newListNodesCommand())

	return clusterCmd
}

func newAddNodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add-node",
		Short: "Add a node to the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			// Parse flags and execute logic
			node := cmd.Flag("node").Value.String()
			if err := cluster.AddNode(node); err != nil {
				fmt.Printf("Failed to add node: %v\n", err)
			} else {
				fmt.Println("Node added successfully.")
			}
		},
	}
}

func newRemoveNodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-node",
		Short: "Remove a node from the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			node := cmd.Flag("node").Value.String()
			if err := cluster.RemoveNode(node); err != nil {
				fmt.Printf("Failed to remove node: %v\n", err)
			} else {
				fmt.Println("Node removed successfully.")
			}
		},
	}
}

func newListNodesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list-nodes",
		Short: "List all nodes in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			nodes, err := cluster.ListNodes()
			if err != nil {
				fmt.Printf("Failed to list nodes: %v\n", err)
			} else {
				fmt.Println("Nodes in the cluster:")
				for _, node := range nodes {
					fmt.Println(node)
				}
			}
		},
	}
}

