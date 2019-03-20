package config

// MySQL config for urlfilter.
type MySQL struct {
	// Hostname of MySQL - default "".
	Host string `json:"host"`

	// Port of MySQL - default 6379.
	Port string `json:"port"`

	// Username for MySQL - default "".
	Username string `json:"username"`

	// Password for MySQL - default "".
	Password string `json:"password"`
}

// Return MySQL config with default values.
func NewMySQL() MySQL {
	return MySQL{
		Host:     "",
		Port:     "3306",
		Username: "",
		Password: "",
	}
}
