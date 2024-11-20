package mqttclient

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient wraps the MQTT connection
type MQTTClient struct {
	Client mqtt.Client
}

// NewMQTTClient initializes a new MQTT client
func NewMQTTClient(broker, clientID string) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTClient{Client: client}, nil
}

// Publish sends a message to the MQTT broker
func (m *MQTTClient) Publish(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	token := m.Client.Publish(topic, 0, false, data)
	token.Wait()
	return token.Error()
}

// Subscribe listens to messages on the given topic
func (m *MQTTClient) Subscribe(topic string, callback func([]byte)) error {
	token := m.Client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg.Payload())
	})
	token.Wait()
	return token.Error()
}

// Disconnect safely disconnects the MQTT client
func (m *MQTTClient) Disconnect() {
	m.Client.Disconnect(250)
}
