package filters

import (
	"errors"
	"strings"
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

type TestLoader struct {
	db    map[int]string
	maxID int
}

func NewTestLoader() *TestLoader {
	loader := &TestLoader{}
	loader.db = make(map[int]string)

	return loader
}

func (t *TestLoader) GetURLPage(start int, number int) ([]string, int, error) {
	urls := make([]string, 0, number)
	for i := start; i < start+number; i++ {
		if url, ok := t.db[i]; ok {
			urls = append(urls, url)
		} else {
			break
		}
	}
	return urls, start + number - 1, nil
}

func (t *TestLoader) GetMaxID() (int, error) {
	return t.maxID, nil
}

func (t *TestLoader) AddURLs(urls []string) {
	for _, url := range urls {
		t.maxID++
		t.db[t.maxID] = url
	}
}
