package config

// Config for a Redis backed Bloom Filter that loads
//from MySQL.
type RedisMySQLBloom struct {
	// Redis specific config.
	Redis Redis `json:"redis"`

	// Redis specific config.
	MySQL MySQL `json:"mysql"`

	// Page size of entries to load at a time - default 1000.
	PageLoadSize int `json:"pageloadsize"`

	// The interval, in minutes, at which we check MySQL for new entries - default 5.
	PageLoadInterval int `json:"pageloadinterval"`
}

// Return Bloom Filter config with default values.
func NewRedisMySQLBloom() RedisMySQLBloom {
	return RedisMySQLBloom{
		Redis:            NewRedis(),
		MySQL:            NewMySQL(),
		PageLoadSize:     1000,
		PageLoadInterval: 5,
	}
}
