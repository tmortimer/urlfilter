package config

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestNewRedisDefaults(t *testing.T) {
	redis := NewRedis()

	if redis.Host != "" {
		t.Errorf("Redis.Hostname should be empty but was %s.", redis.Host)
	}

	if redis.Port != "6379" {
		t.Errorf("Redis.Port should be 6379 but was %s.", redis.Port)
	}

	if redis.Password != "" {
		t.Errorf("Redis.Password should be empty but was %s.", redis.Password)
	}

	if redis.MaxIdle != 10 {
		t.Errorf("Redis.MaxIdle should be 10 but was %d.", redis.MaxIdle)
	}

	if redis.IdleTimeout != 600 {
		t.Errorf("Redis.MaxIdle should be 600 but was %d.", redis.IdleTimeout)
	}

	if len(redis.Config) > 0 {
		t.Errorf("Redis.Config should be empty but was not should be 600 but was %v.", redis.Config)
	}

	if redis.InsertChunkSize != 1000 {
		t.Errorf("Redis.InsertChunkSize should be 1000 but was %d.", redis.InsertChunkSize)
	}

}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.Host != "" {
		t.Errorf("The Host should be empty but was %s.", config.Port)
	}

	if config.Port != "8080" {
		t.Errorf("The Port should be 8080 but was %s.", config.Port)
	}

	redis := NewRedis()
	if !cmp.Equal(config.Redis, redis) {
		t.Error("The default config options had non-default Redis config.")
		t.Error(cmp.Diff(config.Redis, redis))
	}
}

func TestParseConfig(t *testing.T) {
	config := NewConfig()
	config.Host = "google.ca"
	config.Port = "6060"
	config.Redis.Host = "google.ca"
	config.Redis.Port = "444"
	config.Redis.Password = "Changeme"

	configBytes, err := json.Marshal(config)
	if err != nil {
		t.Fatal("Failed to create JSON string from config.Config")
	}
	parsedConfig := ParseConfig(strings.NewReader(string(configBytes)))

	if !cmp.Equal(parsedConfig, config) {
		t.Error("The parsed config did not match the input config.")
		t.Error(cmp.Diff(parsedConfig, config))
	}
}

func TestParseConfigFileReturnsDefaultConfigEmptyPath(t *testing.T) {
	config := NewConfig()
	parsedConfig := ParseConfigFile("")
	if !cmp.Equal(parsedConfig, config) {
		t.Error("The parsed config did not match the input config.")
		t.Error(cmp.Diff(parsedConfig, config))
	}
}
