package connectors

// Interface to underlying database connection pool from which a bloom filter is loaded.
type Loader interface {
	// Check if the URL is in the database.
	GetURLPage(start int, number int) ([]string, error)
}
