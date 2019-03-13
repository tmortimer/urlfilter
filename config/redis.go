package config

// Redis config for urlfilter.
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// Return a RedisConfig with default values.
func NewRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
	}
}
