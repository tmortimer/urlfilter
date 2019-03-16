package filters

import (
	"errors"
	"strings"
	"testing"
)

type TestConnector struct {
	db map[string]bool
}

func NewTestConnector() *TestConnector {
	connector := &TestConnector{}
	connector.db = make(map[string]bool)

	return connector
}

func (t *TestConnector) ContainsURL(url string) (bool, error) {
	if strings.Contains(url, "merp") {
		return false, errors.New("Bad things happened!")
	}

	if strings.Contains(url, "derp") {
		return true, errors.New("Bad things happened!")
	}

	return t.db[url], nil
}

func (t *TestConnector) AddURL(url string) error {
	if strings.Contains(url, "perm") {
		return errors.New("Bad things happened!")
	}

	t.db[url] = true

	return nil
}

func (t *TestConnector) Name() string {
	return "Test"
}

func TestSetsSecondaryFilter(t *testing.T) {
	db := NewDB(NewTestConnector())

	db.AddSecondaryFilter(NewFake())
}

func TestContainsURLFound(t *testing.T) {
	url := "facebook.com"
	conn := NewTestConnector()
	conn.db[url] = true
	db := NewDB(conn)

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	if err != nil {
		t.Errorf("An error was generated when none was expected: %s.", err.Error())
	}
}

func TestContainsURLNotFound(t *testing.T) {
	url := "facebook.com"
	db := NewDB(NewTestConnector())

	found, err := db.ContainsURL(url)

	if found {
		t.Errorf("URL \"%s\" was incorrectly returned by the filter.", url)
	}
	if err != nil {
		t.Errorf("An error was generated when none was expected: %s.", err.Error())
	}
}

func TestContainsURLError(t *testing.T) {
	url := "facebook.com/merp"
	db := NewDB(NewTestConnector())

	found, err := db.ContainsURL(url)

	if found {
		t.Errorf("URL \"%s\" was incorrectly returned by the filter.", url)
	}
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}
}

func TestContainsURLErrorFound(t *testing.T) {
	url := "facebook.com/derp"
	db := NewDB(NewTestConnector())

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}
}

func TestContainsURLFoundWSecondary(t *testing.T) {
	url := "werpwerp.com"
	conn := NewTestConnector()
	conn.db[url] = true
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	if err != nil {
		t.Errorf("An error was generated when none was expected: %s.", err.Error())
	}
}

func TestContainsURLNotFoundWSecondaryNotFound(t *testing.T) {
	url := "werpwerp.com"
	db := NewDB(NewTestConnector())
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if found {
		t.Errorf("URL \"%s\" was incorrectly returned by the filter.", url)
	}
	if err != nil {
		t.Errorf("An error was generated when none was expected: %s.", err.Error())
	}
}

func TestContainsURLNotFoundWSecondaryFound(t *testing.T) {
	url := "facebook.com"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	if err != nil {
		t.Errorf("An error was generated when none was expected: %s.", err.Error())
	}

	// Should be added to the cache now
	if !conn.db[url] {
		t.Errorf("URL \"%s\" was not added to the cache.", url)
	}
}

// Error in primary, found in secondary, added to cache
// Error in primary, found in primary
// Error in primary, error in secondary, found in secondary
// Error in primary, error in secondary, not found in secondary
// Not found in primary, not found in secondary, error in secondary
// not found in primary, found in secondary, error in secondary
// not found in primary, found in secondary, error adding to primary
