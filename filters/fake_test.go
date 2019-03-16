package filters

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

var containsURL = []string{
	"www.facebook.com/wjwjw/wdqwd",
	"facebook.com/pewpewpew",
	"www.google.ca/facebook/facechew",
	"www.google.ca/facebook",
}

var doesNotContainURL = []string{
	"www.facehook.com",
	"cisco.com/facehok",
	"www.netapp.com",
}

var errorsURL = []string{
	"www.bookface.com/wjwjw/wdqwd",
	"bookface.com/pewpewpew",
	"www.google.ca/bookface/facechew",
	"www.google.ca/bookface",
}

var errorsContainsURL = []string{
	"www.faceface.com/wjwjw/wdqwd",
	"faceface.com/pewpewpew",
	"www.google.ca/faceface/facechew",
	"www.google.ca/faceface",
}

func TestAddSecondaryFilterDoesNothing(t *testing.T) {
	f := NewFake()
	f.AddSecondaryFilter(NewFake())
	if !cmp.Equal(*f, *NewFake()) {
		t.Fatalf("Fake AddSecondaryFilter affected the filter... somehow.")
		t.Fatal(cmp.Diff(*f, *NewFake()))
	}
}

func TestContainsURLContains(t *testing.T) {
	f := NewFake()
	for _, url := range containsURL {
		contains, err := f.ContainsURL(url)
		if !contains {
			t.Errorf("URL \"%s\" was incorrectly missed by the filter.", url)
		}
		if err != nil {
			t.Errorf("An error was generated when none was expected: %s.", err.Error())
		}
	}
}

func TestContainsURLDoesNotContain(t *testing.T) {
	f := NewFake()
	for _, url := range doesNotContainURL {
		contains, err := f.ContainsURL(url)
		if contains {
			t.Errorf("URL \"%s\" was incorrectly flagged by the filter.", url)
		}
		if err != nil {
			t.Errorf("An error was generated when none was expected: %s.", err.Error())
		}
	}
}

func TestContainsURLGeneratesError(t *testing.T) {
	f := NewFake()
	for _, url := range errorsURL {
		contains, err := f.ContainsURL(url)
		if contains {
			t.Errorf("URL \"%s\" was incorrectly flagged by the filter.", url)
		}
		if err == nil {
			t.Error("An error was not generated when one was expected.")
		}
	}
}

func TestContainsURLGeneratesErrorStillFound(t *testing.T) {
	f := NewFake()
	for _, url := range errorsContainsURL {
		contains, err := f.ContainsURL(url)
		if !contains {
			t.Errorf("URL \"%s\" was incorrectly missed by the filter.", url)
		}
		if err == nil {
			t.Error("An error was not generated when one was expected.")
		}
	}
}
