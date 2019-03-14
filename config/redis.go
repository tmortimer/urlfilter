package config

// Redis config for urlfilter.
type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// Return a Redis config with default values.
func NewRedis() Redis {
	return Redis{
		Host:     "",
		Port:     "6379",
		Password: "",
	}
}
