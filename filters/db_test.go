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

func TestContainsURLNotFoundWNotFoundSecondary(t *testing.T) {
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

func TestContainsURLNotFoundWFoundSecondary(t *testing.T) {
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

func TestContainsURLNotFoundErrorWFoundSecondary(t *testing.T) {
	url := "facebook.com/merp"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, _ := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	/* Need to not lose error information here on subsequent stems
	//TOM
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}*/

	// Should be added to the cache now
	if !conn.db[url] {
		t.Errorf("URL \"%s\" was not added to the cache.", url)
	}
}

func TestContainsURLNotFoundErrorWFoundSecondaryError(t *testing.T) {
	url := "facebook.com/derp"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}

	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}
}

func TestContainsURLNotFoundErrorWErrorInSecondary(t *testing.T) {
	url := "bookface.com/merp"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if found {
		t.Errorf("URL \"%s\" was incorrectly returned by the filter.", url)
	}
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}
}

func TestContainsURLNotFoundWErrorInSecondary(t *testing.T) {
	url := "bookface.com"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if found {
		t.Errorf("URL \"%s\" was incorrectly returned by the filter.", url)
	}
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}

	// Should not be in the cache
	if conn.db[url] {
		t.Errorf("URL \"%s\" added to the cache when it was not supposed to be.", url)
	}
}

func TestContainsURLNotFoundWFoundSecondaryErrorInSecondary(t *testing.T) {
	url := "faceface.com"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, _ := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	/* Need to not lose error information here on subsequent stems
	//TOM
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}*/

	// Should be added to the cache now
	if !conn.db[url] {
		t.Errorf("URL \"%s\" was not added to the cache.", url)
	}
}

func TestContainsURLNotFoundWFoundSecondaryErrorAddingToCache(t *testing.T) {
	url := "facebook.com/perm"
	conn := NewTestConnector()
	db := NewDB(conn)
	db.AddSecondaryFilter(NewFake())

	found, err := db.ContainsURL(url)

	if !found {
		t.Errorf("URL \"%s\" was not returned by the filter.", url)
	}
	if err == nil {
		t.Error("An error was not generated when one was expected.")
	}

	// Should not be in the cache
	if conn.db[url] {
		t.Errorf("URL \"%s\" added to the cache when it was not supposed to be.", url)
	}
}
