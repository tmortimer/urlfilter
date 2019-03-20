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

func TestNewMySQLDefaults(t *testing.T) {
	mysql := NewMySQL()

	if mysql.Host != "" {
		t.Errorf("MySQL.Hostname should be empty but was %s.", mysql.Host)
	}

	if mysql.Port != "3306" {
		t.Errorf("MySQL.Port should be 6379 but was %s.", mysql.Port)
	}

	if mysql.Username != "" {
		t.Errorf("MySQL.Password should be empty but was %s.", mysql.Password)
	}

	if mysql.Password != "" {
		t.Errorf("MySQL.Password should be empty but was %s.", mysql.Password)
	}
}

func TestNewBloomDefaults(t *testing.T) {
	bloom := NewBloom()

	if bloom.LoadPageSize != 1000 {
		t.Errorf("Bloom.LoadPageSize should be 1000 but was %d.", bloom.LoadPageSize)
	}

	redis := NewRedis()
	if !cmp.Equal(bloom.Redis, redis) {
		t.Error("The default bloom config options had non-default Redis config.")
		t.Error(cmp.Diff(bloom.Redis, redis))
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

	if strings.Compare(config.Filters[0], "redis") != 0 || len(config.Filters) != 1 {
		t.Errorf("The Filters list should be [\"redis\"] but was %v.", config.Filters)
	}

	redis := NewRedis()
	if !cmp.Equal(config.Redis, redis) {
		t.Error("The default config options had non-default Redis config.")
		t.Error(cmp.Diff(config.Redis, redis))
	}

	mysql := NewMySQL()
	if !cmp.Equal(config.MySQL, mysql) {
		t.Error("The default config options had non-default MySQL config.")
		t.Error(cmp.Diff(config.MySQL, mysql))
	}

	bloom := NewBloom()
	if !cmp.Equal(config.Bloom, bloom) {
		t.Error("The default config options had non-default Bloom config.")
		t.Error(cmp.Diff(config.Bloom, bloom))
	}
}

func TestParseConfig(t *testing.T) {
	config := NewConfig()
	config.Host = "google.ca"
	config.Port = "6060"
	config.Filters = []string{"redis", "fake"}

	config.Redis.Host = "google.ca"
	config.Redis.Port = "444"
	config.Redis.Password = "Changeme"
	config.Redis.MaxIdle = 100
	config.Redis.IdleTimeout = 4
	config.Redis.Config = []string{"run some things"}
	config.Redis.InsertChunkSize = 1

	config.MySQL.Host = "google.ca"
	config.MySQL.Port = "444"
	config.MySQL.Username = "used"
	config.MySQL.Password = "Changeme"

	config.Bloom.LoadPageSize = 66
	config.Bloom.Redis.Host = "google.ca"
	config.Bloom.Redis.Port = "444"
	config.Bloom.Redis.Password = "Changeme"
	config.Bloom.Redis.MaxIdle = 100
	config.Bloom.Redis.IdleTimeout = 4
	config.Bloom.Redis.Config = []string{"run some things"}
	config.Bloom.Redis.InsertChunkSize = 1

	configBytes, err := json.Marshal(config)
	if err != nil {
		t.Fatal("Failed to create JSON string from config.Config")
	}

	parsedConfig, err := ParseConfig(strings.NewReader(string(configBytes)))
	if err != nil {
		t.Fatalf("Config validation failed: %s", err)
	}

	if !cmp.Equal(parsedConfig, config) {
		t.Error("The parsed config did not match the input config.")
		t.Error(cmp.Diff(parsedConfig, config))
	}
}

func TestValidateConfigFailure(t *testing.T) {
	config := NewConfig()
	config.Filters = []string{"redis", "merp", "fake"}

	configBytes, err := json.Marshal(config)
	if err != nil {
		t.Fatal("Failed to create JSON string from config.Config")
	}

	_, err = ParseConfig(strings.NewReader(string(configBytes)))
	if err == nil {
		t.Fatalf("Config validation failed was supposed to fail but didn't.")
	}
}

func TestParseConfigFileReturnsDefaultConfigEmptyPath(t *testing.T) {
	config := NewConfig()
	parsedConfig, _ := ParseConfigFile("")
	if !cmp.Equal(parsedConfig, config) {
		t.Error("The parsed config did not match the input config.")
		t.Error(cmp.Diff(parsedConfig, config))
	}
}
