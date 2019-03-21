package filters

import (
	"errors"
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
func (f *Fake) AddSecondaryFilter(filter Filter) error {
	return nil
}

// Returns true if the url contains facebook anywhere in it,
// because that's as good as anything to block.
func (f *Fake) ContainsURL(url string) (bool, error) {
	if strings.Contains(url, "facebook") {
		return true, nil
	}

	if strings.Contains(url, "bookface") {
		return false, errors.New("Bad things happened!")
	}

	if strings.Contains(url, "faceface") {
		return true, errors.New("Bad things happened!")
	}

	return false, nil
}
