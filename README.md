# pgEdgeCLI

`pgEdgeCLI` is a command-line tool to manage `pgEdge`'s PostgreSQL clusters, providing functionalities for cluster management, replication setup, and more.

---

## Features

- **Cluster Management**: Manage PostgreSQL clusters with ease.
- **Spock Replication**: Configure and manage logical replication using Spock.
- **MQTT Integration**: Efficient communication using MQTT for cluster operations.
- **Versioned Binary**: Includes version information in the CLI binary.
- **Cross-Platform**: Supports Linux, macOS, and Windows.

---

## Requirements

- **Go**: Version 1.20 or later.
- **Dependencies**: 
  - [`cobra`](https://github.com/spf13/cobra): For CLI commands.
  - [`paho.mqtt.golang`](https://github.com/eclipse/paho.mqtt.golang): For MQTT integration.

---

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/pgedgecli.git
   cd pgedgecli
   make build

## Usage

$ ./pgedgecli 
pgEdgeCLI is a command-line tool to manage pgEdge's PostgreSQL cluster
  ```bash
   Usage:
      pgedgecli [command]

    Available Commands:
      cluster     Manage PostgreSQL clusters
      help        Help about any command
      spock       Manage Spock replication
      version     Print the version number of pgedgecli

  Flags:
      -h, --help   help for pgedgecli

      Use "pgedgecli [command] --help" for more information about a command.
  ```
