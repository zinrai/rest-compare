package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Represents the structure of the configuration file
type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
	Settings  Settings   `yaml:"settings"`
}

// Represents the configuration of an API endpoint
type Endpoint struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Auth string `yaml:"auth,omitempty"` // Authentication is optional
}

// Represents comparison settings
type Settings struct {
	Timeout     int      `yaml:"timeout,omitempty"`
	IgnoredKeys []string `yaml:"ignoredKeys,omitempty"`
	JSONPath    string   `yaml:"jsonPath,omitempty"` // Optional JSON path
}

// Loads the configuration file and converts it to a Config structure
func LoadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// Set default values
	setDefaults(&config)

	return &config, nil
}

// Validates the configuration
func validateConfig(config *Config) error {
	// Check the number of endpoints
	if len(config.Endpoints) < 2 {
		return errors.New("at least two endpoints are required for comparison")
	}

	// Check if endpoint URL is set
	for i, endpoint := range config.Endpoints {
		if endpoint.Name == "" {
			return errors.New("endpoint name is required")
		}
		if endpoint.URL == "" {
			return errors.New("endpoint URL is required")
		}

		// Check for duplicate endpoint names
		for j := i + 1; j < len(config.Endpoints); j++ {
			if endpoint.Name == config.Endpoints[j].Name {
				return errors.New("duplicate endpoint name: " + endpoint.Name)
			}
		}
	}

	return nil
}

// Sets default values
func setDefaults(config *Config) {
	// Default timeout value
	if config.Settings.Timeout <= 0 {
		config.Settings.Timeout = 30
	}

	// Default ignored keys
	if config.Settings.IgnoredKeys == nil {
		config.Settings.IgnoredKeys = []string{}
	}

	// Default JSON path is empty string (compare entire response)
}

// Returns the first two endpoints
func (c *Config) GetDefaultEndpoints() (Endpoint, Endpoint) {
	return c.Endpoints[0], c.Endpoints[1]
}

// Returns the timeout value
func (c *Config) GetTimeout() int {
	return c.Settings.Timeout
}

// Returns the list of keys to ignore
func (c *Config) GetIgnoredKeys() []string {
	return c.Settings.IgnoredKeys
}

// Returns the JSON path
func (c *Config) GetJSONPath() string {
	return c.Settings.JSONPath
}
