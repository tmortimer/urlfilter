// DB Connectors to back DB based filters.
package connectors

// Interface to underlying database connection pool and comand runner.
type Connector interface {
	// Check if the URL is in the database.
	ContainsURL(url string) (bool, error)

	// Add the URL to the database. Only used if this DB is being used as a cache.
	AddURL(url string) error

	// Return the name of this connector. Used for logging.
	Name() string
}
