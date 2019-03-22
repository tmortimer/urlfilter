package filters

import (
	"errors"
	"github.com/tmortimer/urlfilter/connectors"
	"log"
	"sync/atomic"
	"time"
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
	pageLoadInterval time.Duration

	// The ID of the last URL loaded from the DB.
	lastIdLoaded int

	// The number of URLs loaded into the Bloom Filter.
	numURLs int

	// Store whether the bloom filter is ready or not
	ready int32

	// Timer for refreshing the Bloom Filter and picking up new entries.
	ticker *time.Ticker
}

// Return a new database filter.
func NewBloom(conn connectors.Connector, loader connectors.Loader, pageLoadSize int, pageLoadInterval int) *Bloom {
	bloom := &Bloom{
		conn:             conn,
		loader:           loader,
		pageLoadSize:     pageLoadSize,
		pageLoadInterval: time.Duration(pageLoadInterval),
		lastIdLoaded:     0,
		numURLs:          0,
		ready:            0,
		ticker:           time.NewTicker(time.Duration(pageLoadInterval) * time.Minute),
	}

	go func() {
		bloom.Load()
		atomic.StoreInt32(&(bloom.ready), 1)

		for _ = range bloom.ticker.C {
			bloom.Load()
		}
	}()

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
	log.Printf("The Bloom Filter loaded %d urls for a total of %d.", count, b.numURLs)
}

// Stop the Bloom Filter's background loading task.
func (b *Bloom) StopLoading() {
	b.ticker.Stop()
}

// Add a secondary filter. Required for Bloom Filters.
func (b *Bloom) AddSecondaryFilter(filter Filter) error {
	if filter == nil {
		return errors.New("Bloom Filter can't be configured without a secondary Filter.")
	}
	b.next = filter
	return nil
}

// Check the Bloom Filter for the URL. If the URL is found in the
// Bloom Filter we have to then check the next next filter in the
// chain because Bloom Filters can return false positives. If
// it's not found then we can return right away as a negative
// result is final.
// If the Bloom Filter has not yet been loaded, skip it.
func (b *Bloom) ContainsURL(url string) (bool, error) {
	if atomic.LoadInt32(&(b.ready)) == 0 {
		log.Printf("%s Bloom Filter is not yet loaded, checking the next filter.", b.conn.Name())
		return b.next.ContainsURL(url)
	}

	found, err := b.conn.ContainsURL(url)
	if found || err != nil {
		if err == nil {
			log.Printf("URL %s found in %s Bloom Filter, checking the next filter.", url, b.conn.Name())
		} else {
			log.Printf("%s Bloom Filter generated an error, %s, when checking for %s.", b.conn.Name(), err.Error(), url)
		}
		return b.next.ContainsURL(url)
	}

	// Not found. Nothing to see here.
	log.Printf("URL %s not found in %s Bloom Filter.", url, b.conn.Name())
	return false, nil
}
