package filters

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

var containsURL = []string {
	"www.facebook.com/wjwjw/wdqwd",
	"facebook.com/pewpewpew",
	"www.google.ca/facebook/facechew",
	"www.google.ca/facebook",
}

var doesNotContainURL = []string {
	"www.facehook.com",
	"cisco.com/facehok",
	"www.netapp.com",
}

func TestAddSecondaryFilterDoesNothing(t *testing.T) {
	f := Fake{}
	f.AddSecondaryFilter(&Fake{})
	if !cmp.Equal(f, Fake{}) {
		t.Fatalf("Fake AddSecondaryFilter affected the filter... somehow.")
	}
}

func TestContainsURLContains(t *testing.T) {
	f := Fake{}
	for _, url := range containsURL {
		if !f.ContainsURL(url) {
			t.Errorf("URL \"%s\" was incorrectly missed by the filter.", url)
		}
	}
}

func TestContainsURLDoesNotContain(t *testing.T) {
	f := Fake{}
	for _, url := range doesNotContainURL {
		if f.ContainsURL(url) {
			t.Errorf("URL \"%s\" was incorrectly flagged by the filter.", url)
		}
	}
}
