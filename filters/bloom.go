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

	// The DB connector used to load the Bloom Filter.
	loader connectors.Loader

	// The number of URLs we load at one time while populating the Bloom Filter.
	pageLoadSize int

	// The frequency at which we look for new URLs to add to the Bloom Filter.
	pageLoadInterval int

	// The ID of the last URL loaded from the DB.
	lastIdLoaded int

	// The number of URLs loaded into the Bloom Filter.
	numURLs int
}

// Return a new database filter.
func NewBloom(conn connectors.Connector, loader connectors.Loader, pageLoadSize int, pageLoadInterval int) *Bloom {
	bloom := &Bloom{
		conn:             conn,
		loader:           loader,
		pageLoadSize:     pageLoadSize,
		pageLoadInterval: pageLoadInterval,
		lastIdLoaded:     0,
		numURLs:          0,
	}
	bloom.Load()

	return bloom
}

// Load the bloom filter from the backing data store provided by the loader.
func (b *Bloom) Load() {
	maxID, err := b.loader.GetMaxID()
	if err != nil {
		log.Printf("Failed to load Bloom Filter %s.", err)
		return
	}

	count := 0
	for b.lastIdLoaded < maxID {
		urls, lastIdLoaded, err := b.loader.GetURLPage(b.lastIdLoaded+1, b.pageLoadSize)
		if err != nil {
			log.Printf("Failed to load Bloom Filter %s.", err)
			return
		}
		for _, url := range urls {
			b.conn.AddURL(url)
		}
		count += len(urls)
		b.lastIdLoaded = lastIdLoaded
	}

	b.numURLs += count
	log.Printf("The Bloom Filder loaded an additional %d urls for a total of %d.", count, b.numURLs)
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
		log.Printf("%s Bloom Filter generated an error, %s, when checking for %s.", b.conn.Name(), err.Error(), url)
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
