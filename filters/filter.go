// Chainable filters that can be used by handlers.FilterHandler.
package filters

// Represents a chainable filter to identify malicious URLs.
type Filter interface {
	// Add secondary filter. It is up to this filter how the
	// secondary filter is used, if at all..
	AddSecondaryFilter(filter Filter) error

	// Check if the URL is contained in the filter.
	ContainsURL(url string) (bool, error)
}
