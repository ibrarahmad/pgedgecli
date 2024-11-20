package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pgedgecli/pkg/cluster"
	"pgedgecli/pkg/jsonutils"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// NewClusterCommand creates the "cluster" command with its subcommands
func NewClusterCommand() *cobra.Command {
	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage PostgreSQL clusters",
		Long:  "Commands to manage PostgreSQL clusters such as creating, adding, and listing nodes.",
	}

	// Add subcommands to the cluster command
	clusterCmd.AddCommand(newCreateJSONCommand())
	clusterCmd.AddCommand(newAddNodeCommand())
	clusterCmd.AddCommand(newRemoveNodeCommand())
	clusterCmd.AddCommand(newListNodesCommand())

	return clusterCmd
}

// newCreateJSONCommand creates the "create-json" subcommand
func newCreateJSONCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-json [Cluster Name] [Number of Nodes] [Database Name] [User] [Password]",
		Short: "Create a cluster JSON configuration file",
		Args:  cobra.ExactArgs(5),
		RunE:  runCreateJSONCommand,
	}
}

// runCreateJSONCommand is the execution function for the "create-json" subcommand
func runCreateJSONCommand(cmd *cobra.Command, args []string) error {
	clusterName := args[0]
	numNodes, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid number of nodes: %v", err)
	}
	dbName := args[2]
	dbUser := args[3]
	dbPassword := args[4]

	reader := bufio.NewReader(os.Stdin)

	pgVersion := promptUser(reader, "Enter PostgreSQL version (default: 16): ", "16")
	startPort, err := promptUserForInt(reader, "Enter starting port number (default: 5432): ", 5432)
	if err != nil {
		return fmt.Errorf("invalid port number: %v", err)
	}

	nodeGroups := gatherNodeGroups(reader, clusterName, numNodes, startPort)

	clusterConfig := &jsonutils.PGEConfig{
		JsonVersion: jsonutils.JsonVersion,
		ClusterName: clusterName,
		LogLevel:    "debug",
		UpdateDate:  "",
		PgEdge: jsonutils.PgEdge{
			PgVersion: pgVersion,
			AutoStart: "off",
			Spock: jsonutils.Spock{
				SpockVersion: "",
				AutoDDL:      "off",
			},
			Databases: []jsonutils.DbInfo{
				{
					DbName:     dbName,
					DbUser:     dbUser,
					DbPassword: dbPassword,
				},
			},
		},
		NodeGroups: nodeGroups,
	}

	if err := jsonutils.JSONValidate(clusterConfig); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if err := jsonutils.WriteClusterJSON(clusterName, clusterConfig); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	fmt.Printf("\nCluster JSON file created successfully at: %s\n", clusterConfig.ClusterName)
	return nil
}

// promptUser prompts the user for input with a default value
func promptUser(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

// promptUserForInt prompts the user for an integer input with a default value
func promptUserForInt(reader *bufio.Reader, prompt string, defaultValue int) (int, error) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(input)
}

// gatherNodeGroups collects information about nodes interactively
func gatherNodeGroups(reader *bufio.Reader, clusterName string, numNodes, startPort int) []jsonutils.NodeGroup {
	var nodeGroups []jsonutils.NodeGroup
	for i := 1; i <= numNodes; i++ {
		fmt.Printf("\nNode %d details:\n", i)
		nodeName := promptUser(reader, fmt.Sprintf("Node name (default: n%d): ", i), fmt.Sprintf("n%d", i))
		publicIP := promptUser(reader, "Enter public IP (default: 127.0.0.1): ", "127.0.0.1")
		nodePath := promptUser(reader, fmt.Sprintf("Path (default: /home/%s/%s/%s): ", os.Getenv("USER"), clusterName, nodeName), fmt.Sprintf("/home/%s/%s/%s", os.Getenv("USER"), clusterName, nodeName))
		subNodes := gatherSubNodes(reader, nodeName)

		nodeGroups = append(nodeGroups, jsonutils.NodeGroup{
			Name:      nodeName,
			IsActive:  "on",
			PublicIP:  publicIP,
			PrivateIP: publicIP,
			Port:      fmt.Sprintf("%d", startPort+i-1),
			Path:      nodePath,
			SubNodes:  subNodes,
		})
	}
	return nodeGroups
}

// gatherSubNodes interactively collects information about subnodes for a node
func gatherSubNodes(reader *bufio.Reader, parentNodeName string) []jsonutils.NodeGroup {
	var subNodes []jsonutils.NodeGroup
	confirm := promptUser(reader, fmt.Sprintf("\nAdd sub-nodes for %s? (yes/no, default: no): ", parentNodeName), "no")
	if strings.ToLower(confirm) == "yes" || strings.ToLower(confirm) == "y" {
		for {
			fmt.Printf("\nSub-node details for %s:\n", parentNodeName)
			subNodeName := promptUser(reader, "Sub-node name: ", "")
			if subNodeName == "" {
				fmt.Println("Sub-node name cannot be empty. Try again.")
				continue
			}
			publicIP := promptUser(reader, "Enter public IP (default: 127.0.0.1): ", "127.0.0.1")
			subNodes = append(subNodes, jsonutils.NodeGroup{
				Name:      subNodeName,
				IsActive:  "off",
				PublicIP:  publicIP,
				PrivateIP: publicIP,
				Port:      "5432",
				Path:      fmt.Sprintf("/home/%s/%s/%s", os.Getenv("USER"), parentNodeName, subNodeName),
			})
			another := promptUser(reader, "Add another sub-node? (yes/no, default: no): ", "no")
			if strings.ToLower(another) != "yes" && strings.ToLower(another) != "y" {
				break
			}
		}
	}
	return subNodes
}

// newAddNodeCommand creates the "add-node" subcommand
func newAddNodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add-node",
		Short: "Add a node to the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			node := cmd.Flag("node").Value.String()
			if err := cluster.AddNode(node); err != nil {
				fmt.Printf("Failed to add node: %v\n", err)
			} else {
				fmt.Println("Node added successfully.")
			}
		},
	}
}

// newRemoveNodeCommand creates the "remove-node" subcommand
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

// newListNodesCommand creates the "list-nodes" subcommand
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
