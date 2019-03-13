package filters

import (
	"strings"
)

// A Fake filter used for testing and setup.
type Fake struct{}

// Create a new Fake filter.
func NewFake() *Fake {
	return &Fake{}
}

// Since this filter is full of lies it does nothing
// with any secondary filter chaning anyways.
func (f *Fake) AddSecondaryFilter(filter Filter) {}

// Returns true if the url contains facebook anywhere in it,
// because that's as good as anything to block.
func (f *Fake) ContainsURL(url string) bool {
	if strings.Contains(url, "facebook") {
		return true
	}

	return false
}
