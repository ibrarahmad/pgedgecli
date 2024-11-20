package jsonutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	BaseDir     = "cluster"
	JsonVersion = "1.1"
)

// ClusterInfo holds the directory and file paths for a cluster
type ClusterInfo struct {
	Dir  string
	File string
}

// PGEConfig represents the structure of a pgedge cluster JSON configuration
type PGEConfig struct {
	JsonVersion string      `json:"json_version"`
	ClusterName string      `json:"cluster_name"`
	PgEdge      PgEdge      `json:"pgedge"`
	NodeGroups  []NodeGroup `json:"node_groups"`
	UpdateDate  string      `json:"update_date"`
	LogLevel    string      `json:"log_level"`
}

// PgEdge holds database settings
type PgEdge struct {
	PgVersion string   `json:"pg_version"`
	AutoStart string   `json:"auto_start"`
	Spock     Spock    `json:"spock"`
	Databases []DbInfo `json:"databases"`
}

// Spock holds Spock-related settings
type Spock struct {
	SpockVersion string `json:"spock_version"`
	AutoDDL      string `json:"auto_ddl"`
}

// DbInfo represents a database configuration
type DbInfo struct {
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

type NodeGroup struct {
	Name      string      `json:"name"`
	IsActive  string      `json:"is_active"`
	PublicIP  string      `json:"public_ip"`
	PrivateIP string      `json:"private_ip"`
	Port      string      `json:"port"`
	Path      string      `json:"path"`
	SubNodes  []NodeGroup `json:"sub_nodes,omitempty"`
}

// GetClusterInfo returns the cluster directory and file path
func GetClusterInfo(clusterName string, create bool) (*ClusterInfo, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get home directory: %v", err)
	}

	clusterDir := filepath.Join(homeDir, BaseDir, clusterName)
	clusterFile := filepath.Join(clusterDir, fmt.Sprintf("%s.json", clusterName))

	if create {
		if err := os.MkdirAll(clusterDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cluster directory: %v", err)
		}
	}

	return &ClusterInfo{Dir: clusterDir, File: clusterFile}, nil
}

// GetClusterJSON loads the cluster JSON configuration
func GetClusterJSON(clusterName string) (*PGEConfig, error) {
	clusterInfo, err := GetClusterInfo(clusterName, false)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(clusterInfo.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read cluster JSON file: %v", err)
	}

	var config PGEConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse cluster JSON: %v", err)
	}

	return &config, nil
}

// WriteClusterJSON writes the cluster configuration to a file
func WriteClusterJSON(clusterName string, config *PGEConfig) error {
	clusterInfo, err := GetClusterInfo(clusterName, true)
	if err != nil {
		return err
	}

	config.UpdateDate = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cluster JSON: %v", err)
	}

	if err := ioutil.WriteFile(clusterInfo.File, data, 0644); err != nil {
		return fmt.Errorf("failed to write cluster JSON file: %v", err)
	}

	return nil
}

// JSONTemplate creates a JSON template for the cluster
func JSONTemplate(clusterName, dbName, dbUser, dbPassword string, pgVersion string, numNodes int, port int) (*PGEConfig, error) {
	clusterConfig := &PGEConfig{
		JsonVersion: JsonVersion,
		ClusterName: clusterName,
		LogLevel:    "debug",
		UpdateDate:  time.Now().Format(time.RFC3339),
		PgEdge: PgEdge{
			PgVersion: pgVersion,
			AutoStart: "off",
			Spock: Spock{
				SpockVersion: "",
				AutoDDL:      "off",
			},
			Databases: []DbInfo{
				{
					DbName:     dbName,
					DbUser:     dbUser,
					DbPassword: dbPassword,
				},
			},
		},
	}

	for i := 1; i <= numNodes; i++ {
		clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, NodeGroup{
			Name:     fmt.Sprintf("n%d", i),
			IsActive: "on",
			PublicIP: "127.0.0.1",
			Port:     fmt.Sprintf("%d", port+i-1),
			Path:     fmt.Sprintf("/home/%s/%s/n%d", os.Getenv("USER"), clusterName, i),
		})
	}

	return clusterConfig, nil
}

// JSONValidate validates the structure of a cluster JSON file
func JSONValidate(config *PGEConfig) error {
	if config.JsonVersion != JsonVersion {
		return fmt.Errorf("invalid JSON version: expected %s, got %s", JsonVersion, config.JsonVersion)
	}

	if config.ClusterName == "" || config.PgEdge.PgVersion == "" {
		return fmt.Errorf("missing required fields in cluster JSON")
	}

	for _, db := range config.PgEdge.Databases {
		if db.DbName == "" || db.DbUser == "" || db.DbPassword == "" {
			return fmt.Errorf("invalid database configuration: %+v", db)
		}
	}

	return nil
}
