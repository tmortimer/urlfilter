package filters

import (
	"strings"
)

// A Fake filter used for testing and setup.
type Fake struct {}

// Since this filter is full of lies it does nothing
// with any secondary filter chaning anyways.
func (f *Fake) AddSecondaryFilter(filter Filter) {}

// Returns true if the url contains facebook.com anywhere in it,
// because that's as good as anything to block.
func (f *Fake) ContainsURL(url string) bool {
	if strings.Contains(url, "facebook.com") {
		return true
	}

	return false
}