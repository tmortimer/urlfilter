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
	// Host to bind server to - default "".
	Host string `json:"host"`

	// Port to bind server to - default 8080.
	Port string `json:"port"`

	// Filter chain. Filters are called left to right - default ["redis"].
	// Valid options are: redis, fake, redisCache. The right most item
	// can not be a "Cache", anything earlier in the list _must_ be a
	// "Cache"
	Filters []string `json:"filters"`

	// Config for Redis.
	Redis Redis `json:"redis"`
}

// Return Config with default values.
func NewConfig() *Config {
	return &Config{
		Host:  "",
		Port:  "8080",
		Redis: NewRedis(),
		Filters: []string{"redis"},
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
