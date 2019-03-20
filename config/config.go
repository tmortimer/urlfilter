// Package for parsing and handling urlfilter config.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// All config for urlfilter.
type Config struct {
	// Host to bind server to - default "".
	Host string `json:"host"`

	// Port to bind server to - default 8080.
	Port string `json:"port"`

	// Filter chain. Filters are called left to right - default ["redis"].
	// Valid options are: redis, mysql, bloom and fake.
	Filters []string `json:"filters"`

	// Config for Redis.
	Redis Redis `json:"redis"`

	// Config for MySQL.
	MySQL MySQL `json:"mysql"`

	// Config for Redis Bloom Filter.
	RedisMySQLBloom RedisMySQLBloom `json:"redismysqlbloom"`
}

// Valid Filters to use as Cache
var validFilters = map[string]bool{"mysql": true, "redis": true, "fake": true}

// Return Config with default values.
func NewConfig() *Config {
	return &Config{
		Host:            "",
		Port:            "8080",
		Filters:         []string{"redis"},
		Redis:           NewRedis(),
		MySQL:           NewMySQL(),
		RedisMySQLBloom: NewRedisMySQLBloom(),
	}
}

// Open the config file at path and parse it.
func ParseConfigFile(path string) (*Config, error) {
	configFile, err := os.Open(path)
	defer configFile.Close()
	if path != "" && err != nil {
		return nil, err
	}

	return ParseConfig(configFile)
}

// Validate the config.
func ValidateConfig(config *Config) error {
	for i := 0; i < len(config.Filters); i++ {
		if !validFilters[config.Filters[i]] {
			//TOM only returns the first error.
			return fmt.Errorf("%s is not a valid filter, the only valid options are %v", config.Filters[i], validFilters)
		}
	}
	return nil
}

// Parse the config.
func ParseConfig(reader io.Reader) (*Config, error) {
	config := NewConfig()
	jsonParser := json.NewDecoder(reader)
	jsonParser.Decode(&config)

	err := ValidateConfig(config)
	return config, err
}
