package connectors

// Interface to underlying database connection pool from which a bloom filter is loaded.
type Loader interface {
	// Check if the URL is in the database. Returns the set of URLs and
	// the highest ID loaded from the DB.
	GetURLPage(start int, number int) ([]string, int, error)

	// Get the current max ID in the DB.
	GetMaxID() (int, error)
}
