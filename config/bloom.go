package config

// Bloom Filter config for urlfilter.
type Bloom struct {
	// Redis specific config.
	Redis Redis `json:"redis"`

	// Page size of entries to load at a time - default 1000.
	LoadPageSize int `json:"loadpagesize"`
}

// Return Bloom Filter config with default values.
func NewBloom() Bloom {
	return Bloom{
		Redis:        NewRedis(),
		LoadPageSize: 1000,
	}
}
