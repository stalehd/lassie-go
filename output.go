package lassie

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Output contains status, logs and configuration for an
// application output.
type Output struct {
	EUI    string                 `json:"eui"`
	AppEUI string                 `json:"appEUI"`
	Config map[string]interface{} `json:"config"`
	Log    []OutputLogEntry       `json:"logs"`
	Status string                 `json:"status"`
}

// OutputLogEntry contains a single log entry from an application
// output. This log entries are kept just for simple diagnostics. If you are
// having problems with the output inspect the output's log entries.
type OutputLogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

// CreateOutput creates a new application output.
func (c *Client) CreateOutput(appEUI string, config OutputConfig) (Output, error) {
	output := Output{Config: config.Map()}
	err := c.create(fmt.Sprintf("/applications/%s/outputs", appEUI), &output)
	return output, err
}

// UpdateOutput updates an application output with a new configuration
func (c *Client) UpdateOutput(appEUI, outputEUI string, config OutputConfig) (Output, error) {
	output := Output{Config: config.Map()}
	err := c.update(fmt.Sprintf("/applications/%s/outputs/%s", appEUI, outputEUI), &output)
	return output, err
}

// DeleteOutput removes (and stops) an application output
func (c *Client) DeleteOutput(appEUI, outputEUI string) error {
	return c.delete(fmt.Sprintf("/applications/%s/outputs/%s", appEUI, outputEUI))
}

// Output returns a single output
func (c *Client) Output(appEUI, outputEUI string) (Output, error) {
	var output Output
	err := c.get(fmt.Sprintf("/applications/%s/outputs/%s", appEUI, outputEUI), &output)
	return output, err
}

// Outputs returns all outputs for the application
func (c *Client) Outputs(appEUI string) ([]Output, error) {
	var outputs struct {
		Outputs []Output `json:"outputs"`
	}
	err := c.get(fmt.Sprintf("/applications/%s/outputs", appEUI), &outputs)
	return outputs.Outputs, err
}

// OutputConfig is a configuration for application outputs. Use the actual
// output types
type OutputConfig interface {
	Map() map[string]interface{}
	Type() string
}

// MQTTConfig is an application output configuration for the MQTT outputs.
type MQTTConfig struct {
	Endpoint         string `json:"endpoint"`
	Port             int    `json:"port"`
	TLS              bool   `json:"tls"`
	CertificateCheck bool   `json:"certCheck"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	ClientID         string `json:"clientid"`
	TopicName        string `json:"topicName"`
}

// MQTTConfigType is the string idenfier for MQTT configs
const MQTTConfigType = "mqtt"

// Create a MQTT config from a configuration map. Returns nil if
// something went wrong.
func newMQTTConfig(config map[string]interface{}) (MQTTConfig, error) {
	buf, err := json.Marshal(config)
	if err != nil {
		return MQTTConfig{}, err
	}
	ret := MQTTConfig{}
	if err := json.Unmarshal(buf, &ret); err != nil {
		return MQTTConfig{}, err
	}
	return ret, nil
}

// Map returns the configuration as a map
func (m *MQTTConfig) Map() map[string]interface{} {
	var msg map[string]interface{}
	buf, _ := json.Marshal(m)
	json.Unmarshal(buf, &msg)
	msg["type"] = MQTTConfigType
	return msg
}

// Type returns a string identifying the output
func (m *MQTTConfig) Type() string {
	return MQTTConfigType
}

// MQTTConfig returns the MQTTConfig for the output.
func (o *Output) MQTTConfig() (MQTTConfig, error) {
	if o.Config["type"] == MQTTConfigType {
		return newMQTTConfig(o.Config)
	}
	return MQTTConfig{}, errors.New("not a MQTT configuration")
}

// AWSIoTConfig returns the AWS IoT config for the output
func (o *Output) AWSIoTConfig() (AWSIoTConfig, error) {
	if o.Config["type"] == AWSIoTConfigType {
		return newAWSIoTConfig(o.Config)
	}
	return AWSIoTConfig{}, errors.New("not an AWS IoT configuration")
}

// AWSIoTConfig is an application output configuration for AWS IoT.
type AWSIoTConfig struct {
	Endpoint          string `json:"endpoint"`
	ClientID          string `json:"clientid"`
	ClientCertificate string `json:"clientCertificate"`
	PrivateKey        string `json:"privateKey"`
}

// AWSIoTConfigType is the string identifier for AWS IoT configs.
const AWSIoTConfigType = "awsiot"

func newAWSIoTConfig(config map[string]interface{}) (AWSIoTConfig, error) {
	buf, err := json.Marshal(config)
	if err != nil {
		return AWSIoTConfig{}, err
	}
	ret := AWSIoTConfig{}
	if err := json.Unmarshal(buf, &ret); err != nil {
		return AWSIoTConfig{}, err
	}
	return ret, nil
}

// Map returns the configuration as a map
func (a *AWSIoTConfig) Map() map[string]interface{} {
	var msg map[string]interface{}
	buf, _ := json.Marshal(a)
	json.Unmarshal(buf, &msg)
	msg["type"] = AWSIoTConfigType
	return msg
}

// Type returns the string identifying the output
func (a *AWSIoTConfig) Type() string {
	return AWSIoTConfigType
}
