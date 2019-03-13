// Package for parsing and handling urlfilter config.
package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// All config for urlfilter.
type Config struct {
	Port  string      `json:"port"`
	Redis RedisConfig `json:"redis"`
}

// Return Config with default values.
func NewConfig() *Config {
	return &Config{
		Port:  "8080",
		Redis: NewRedisConfig(),
	}
}

// Open the config file at path and parse it.
func ParseConfigFile(path string) *Config {
	configFile, err := os.Open(path)
	defer configFile.Close()
	if path != "" && err != nil {
		log.Fatal(err)
	}

	return ParseConfig(configFile)
}

// Parse the config.
func ParseConfig(reader io.Reader) *Config {
	config := NewConfig()
	jsonParser := json.NewDecoder(reader)
	jsonParser.Decode(&config)
	return config
}
