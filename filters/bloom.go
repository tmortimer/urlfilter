package filters

import (
	"fmt"
	"github.com/tmortimer/urlfilter/connectors"
	"log"
)

// Bloom Filter Filter (ahahaha). Bloom filters are used to check
// for existance in a large data set. Their nature is that a negative
// result is final, but a positive result could be false. This means
// that a secondary filter must be set.
type Bloom struct {
	// Secondary filter in the filter chain.
	next Filter

	// The underlying DB connection pool.
	conn connectors.Connector
}

// Return a new database filter.
func NewBloom(conn connectors.Connector) *Bloom {
	return &Bloom{
		conn: conn,
	}
}

// Add a secondary filter. Necessary if using this DB as a cache.
func (b *Bloom) AddSecondaryFilter(filter Filter) {
	b.next = filter
}

// If the URL is found in the Bloom Filter we have to then check the next
// next filter in the chain because Bloom Filters can return false
// positives. If it's not found then we can return right away as
// a negative result is final.
func (b *Bloom) ContainsURL(url string) (bool, error) {
	//TOM error information is lost here on subsequent steps.
	found, err := b.conn.ContainsURL(url)
	if err != nil {
		log.Printf("%s Bloom Filter generated an the error %s when checking for %s.", b.conn.Name(), err.Error(), url)
		return false, err
	}

	if found {
		if b.next != nil {
			log.Printf("URL %s found in %s Bloom Filter, checking the next filter.", url, b.conn.Name())
			return b.next.ContainsURL(url)
		} else {
			log.Printf("URL %s found in %s Bloom Filter, but no seconary filter configured.", url, b.conn.Name())
			return false, fmt.Errorf("No secondary filter configured for %s Bloom Filter.", b.conn.Name())
		}
	}

	// Not found. Nothing to see here.
	log.Printf("URL %s not found in %s Bloom Filter.", url, b.conn.Name())
	return false, nil
}
