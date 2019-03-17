package config

// Redis config for urlfilter.
type Redis struct {
	// Hostname of Redis - default "".
	Host string `json:"host"`

	// Port of Redis - default 6379.
	Port string `json:"port"`

	// Password for Redis - default "".
	Password string `json:"password"`

	// Maximum number of idle connections to Redis in the pool - default 10.
	MaxIdle int `json:"maxIdle"`

	// Timeout before closing idle Redis connections in seconds - default 600 (10 Minutes).
	IdleTimeout int `json:"idleTimeout"`

	// Config values for Redis, an array of strings that can be passed to CONFIG SET - default [].
	Config []string `json:"config"`

	// Max number of URLs to bulk insert to Redis in one MSET command - default 100.
	InsertChunkSize int `json:"insertChunkSize"`
}

// Return Redis config with default values.
func NewRedis() Redis {
	return Redis{
		Host:            "",
		Port:            "6379",
		Password:        "",
		MaxIdle:         10,
		IdleTimeout:     600,
		InsertChunkSize: 1000,
	}
}
