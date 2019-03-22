package filters

import (
	"sync/atomic"
	"testing"
	"time"
)

var urls = []string{
	"facebook.com",
	"google.ca/facebook",
	"eeeh.com/facebook/what",
	"cisco.com",
}

var updatedURLs = []string{
	"cars.com",
	"subarusarecool.com",
	"vrooom.com",
	"pewpewpew.net",
}

func NewBloomFilter() *Bloom {
	connector := &TestConnector{}
	connector.db = make(map[string]bool)

	loader := NewTestLoader()
	loader.AddURLs(urls)

	bloom := NewBloom(connector, loader, 1, 1)
	for atomic.LoadInt32(&(bloom.ready)) == 0 {
		time.Sleep(1 * time.Second)
	}
	return bloom
}

func NewBloomFilterNoBackgroundLoading() *Bloom {
	bloom := NewBloomFilter()
	bloom.StopLoading()
	return bloom
}

func TestRequiresSecondaryFilter(t *testing.T) {
	bloom := NewBloomFilterNoBackgroundLoading()
	err := bloom.AddSecondaryFilter(nil)
	if err == nil {
		t.Fatal("Adding a nil secondary filter did not return an error when one was expected.")
	}
}

func TestLoadsExistingURLS(t *testing.T) {
	bloom := NewBloomFilterNoBackgroundLoading()
	_ = bloom.AddSecondaryFilter(NewFake())
	for _, url := range urls {
		found, _ := bloom.conn.ContainsURL(url)
		if !found {
			t.Errorf("URL %s was not found in the Bloom Filter when it was supposed to be.", url)
		}
	}
}

func TestLoadsNewURLS(t *testing.T) {
	bloom := NewBloomFilter()
	_ = bloom.AddSecondaryFilter(NewFake())
	for _, url := range updatedURLs {
		found, _ := bloom.conn.ContainsURL(url)
		if found {
			t.Errorf("URL %s was found in the Bloom Filter when it was not supposed to be.", url)
		}
	}

	bloom.loader.(*TestLoader).AddURLs(updatedURLs)
	// Sadness
	time.Sleep(61 * time.Second)

	for _, url := range updatedURLs {
		found, _ := bloom.conn.ContainsURL(url)
		if !found {
			t.Errorf("URL %s was not found in the Bloom Filter when it was supposed to be.", url)
		}
	}
	bloom.StopLoading()
}

func TestFalsePositive(t *testing.T) {
	url := urls[3]
	bloom := NewBloomFilterNoBackgroundLoading()
	_ = bloom.AddSecondaryFilter(NewFake())
	found, _ := bloom.ContainsURL(url)
	if found {
		t.Errorf("URL %s was found in the filter chain when it was not supposed to be.", url)
	}
}

func TestPositive(t *testing.T) {
	url := urls[2]
	bloom := NewBloomFilterNoBackgroundLoading()
	_ = bloom.AddSecondaryFilter(NewFake())
	found, _ := bloom.ContainsURL(url)
	if !found {
		t.Errorf("URL %s was not found in the filter chain when it was supposed to be.", url)
	}
}

func TestNegative(t *testing.T) {
	url := "chickens.com/facebook"
	bloom := NewBloomFilterNoBackgroundLoading()
	_ = bloom.AddSecondaryFilter(NewFake())
	found, _ := bloom.ContainsURL(url)
	if found {
		t.Errorf("URL %s was found in the filter chain when it was not supposed to be.", url)
	}
}

func TestErrorChecksNextFilter(t *testing.T) {
	url := "facebook.com/merp"
	bloom := NewBloomFilterNoBackgroundLoading()
	_ = bloom.AddSecondaryFilter(NewFake())
	found, _ := bloom.ContainsURL(url)
	if !found {
		t.Errorf("URL %s was not found in the filter chain when it was supposed to be.", url)
	}
}
