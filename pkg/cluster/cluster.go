package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"pgedgecli/pkg/jsonutils"
	"pgedgecli/pkg/mqttclient"
)

type ClusterSettings struct {
	ClusterName string   `json:"cluster_name"`
	Nodes       []Node   `json:"node_groups"`
	DB          Database `json:"pgedge"`
}

type Node struct {
	Name       string `json:"name"`
	PublicIP   string `json:"public_ip"`
	PrivateIP  string `json:"private_ip"`
	Port       int    `json:"port"`
	Path       string `json:"path"`
	SSHUser    string `json:"ssh_user"`
	SSHKey     string `json:"ssh_key"`
	IsActive   string `json:"is_active"`
	Backrest   string `json:"backrest"`
	LogLevel   string `json:"log_level"`
	Repo1Type  string `json:"repo1_type"`
	Repo1Path  string `json:"repo1_path"`
	BackupID   string `json:"backup_id"`
	PGVersion  string `json:"pg_version"`
	SpockVer   string `json:"spock_version"`
	AutoStart  bool   `json:"auto_start"`
	AutoDDL    bool   `json:"auto_ddl"`
	LogDir     string `json:"log_directory"`
	SharedLibs string `json:"shared_preload_libraries"`
}

type Database struct {
	PGVersion    string `json:"pg_version"`
	SpockVersion string `json:"spock_version"`
	AutoDDL      string `json:"auto_ddl"`
	AutoStart    string `json:"auto_start"`
}

func initCluster(clusterName string, broker string, clientID string) {
	// Load JSON configuration
	clusterFile := fmt.Sprintf("cluster/%s/%s.json", clusterName, clusterName)
	clusterData, err := jsonutils.LoadJSON(clusterFile)
	if err != nil {
		log.Fatalf("Failed to load cluster JSON: %v", err)
	}

	var settings ClusterSettings
	err = json.Unmarshal(clusterData, &settings)
	if err != nil {
		log.Fatalf("Failed to parse cluster JSON: %v", err)
	}

	// Initialize MQTT Client
	mqttClient, err := mqttclient.NewMQTTClient(broker, clientID)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	// Iterate over nodes and send MQTT messages
	for _, node := range settings.Nodes {
		message := map[string]interface{}{
			"action": "init",
			"node": map[string]interface{}{
				"name":       node.Name,
				"public_ip":  node.PublicIP,
				"private_ip": node.PrivateIP,
				"port":       node.Port,
				"path":       node.Path,
				"pg_version": settings.DB.PGVersion,
				"spock_ver":  settings.DB.SpockVersion,
			},
		}

		topic := fmt.Sprintf("pgedge/cluster/%s/node/%s/init", clusterName, node.Name)
		err = mqttClient.Publish(topic, message)
		if err != nil {
			log.Printf("Failed to send MQTT message for node %s: %v", node.Name, err)
		} else {
			log.Printf("Initialization message sent to node %s on topic %s", node.Name, topic)
		}
	}

	log.Println("Cluster initialization complete.")
}

func AddNode(node string) error {
	if node == "" {
		return errors.New("node name cannot be empty")
	}
	return nil
}

func RemoveNode(node string) error {
	if node == "" {
		return errors.New("node name cannot be empty")
	}
	return nil
}

func ListNodes() ([]string, error) {
	return []string{"node1", "node2", "node3"}, nil
}
