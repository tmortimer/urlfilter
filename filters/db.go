package filters

import (
	"github.com/tmortimer/urlfilter/connectors"
	"log"
)

// Database based filter. Depending on config can be used as a cache.
// If this is being used as a cache a secondary filter must be set.
type DB struct {
	// Secondary filter in the filter chain.
	next Filter

	// The underlying DB connection pool.
	conn connectors.Connector
}

// Return a new database filter.
func NewDB(conn connectors.Connector) *DB {
	return &DB{
		conn: conn,
	}
}

// Add a secondary filter. Necessary if using this DB as a cache.
func (d *DB) AddSecondaryFilter(filter Filter) {
	d.next = filter
}

// Return true if the URL is found in the Database. If it's not then return false
// if there are no further filters in the chain, otherwise call the next filter.
// If the database generates an error and this is only a cache we can continue down the
// filter chain, since each subsequent level should have better information.
func (d *DB) ContainsURL(url string) (bool, error) {
	//TOM error information is lost here on subsequent steps.
	found, err := d.conn.ContainsURL(url)
	if err != nil {
		log.Printf("%s generated an the error %s when checking for %s.", d.conn.Name(), err.Error(), url)
	}

	if found || d.next == nil {
		return found, err
	}

	// Not found in the cache, try the next filter.
	found, err = d.next.ContainsURL(url)

	if found {
		// Add it to the cache.
		err = d.conn.AddURL(url)
		if err != nil {
			log.Printf("%s generated an the error %s when adding %s.", d.conn.Name(), err.Error(), url)
		}
	}

	return found, err
}
