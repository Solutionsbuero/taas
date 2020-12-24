package ttrn

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// Config describes the structure of parameters for ttrn.
type Config struct {
	// Port is the port where the page is served on localhost.
	Port int `json:"port"`
	// Db is the path to the sqlite db.
	Db string `json:"db"`
	// MqttHost is the hostname of the the MQTT broker.
	MqttHost string `json:"mqtt_host"`
	// MqttPort is the port of the MQTT broker.
	MqttPort int `json:"mqtt_port"`
	// MqttUser is the user name for the MQTT broker.
	MqttUser string `json:"mqtt_user"`
	// MqttPassword is the password for the MQTT broker.
	MqttPassword string `json:"mqtt_password"`
}

// DefaultConfig returns a new instance of the Config struct with default values.
func DefaultConfig() Config {
	return Config{
		Port:         8000,
		Db:           "data.db",
		MqttHost:     "127.0.0.1",
		MqttPort:     8883,
		MqttUser:     "mqtt-usr",
		MqttPassword: "mqtt-pwd",
	}
}

// OpenConfig loads the configuration from a JSON file and returns the result.
func OpenConfig(path string) Config {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("error loading configuration form %s, %s", path, err)
	}
	var rsl Config
	if err := json.Unmarshal(raw, &rsl); err != nil {
		logrus.Fatalf("error unmarshalling configuration from %s, %s", path, err)
	}
	return rsl
}

// SaveConfig writes a configuration to the given path as a JSON.
func (c Config) SaveConfig(path string) {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		logrus.Fatalf("couldn't marshal config, %s", err)
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		logrus.Fatalf("cloudn't write config to %s, %s", path, data)
	}
}
