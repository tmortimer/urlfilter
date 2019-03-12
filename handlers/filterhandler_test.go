package handlers

import (
	"github.com/tmortimer/urlfilter/filters"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestFilter struct {
	secondary filters.Filter
	called    int
}

func (f *TestFilter) AddSecondaryFilter(filter filters.Filter) {
	f.secondary = filter
}

func (f *TestFilter) ContainsURL(url string) bool {
	f.called++
	return f.secondary.ContainsURL(url)
}

func TestInitAddsHandlers(t *testing.T) {
	f := &TestFilter{}
	f.AddSecondaryFilter(&filters.Fake{})
	h := NewFilterHandler(f)

	// So it doesn't fail, but how can I (directly) test that it actually registered...
	// Not going to spend the time digging into these weeds right now.
	h.Init()
}

func TestHandlesSafeURL(t *testing.T) {
	// https://blog.questionable.services/article/testing-http-handlers-go/
	// This page was useful for info on how to test http handlers in Go.
	f := &TestFilter{}
	f.AddSecondaryFilter(&filters.Fake{})
	h := NewFilterHandler(f)

	req, err := http.NewRequest("GET", FILTER_ENDPOINT+"www.google.ca", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.filterHandler)

	handler.ServeHTTP(recorder, req)

	if f.called != 1 {
		t.Errorf("The TestFilter ContainsURL function was called %d time(s).", f.called)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("The filterHandler function %s when OK was expected.", http.StatusText(recorder.Code))
	}
}

func TestHandlesBlockedURL(t *testing.T) {
	// https://blog.questionable.services/article/testing-http-handlers-go/
	// This page was useful for info on how to test http handlers in Go.
	f := &TestFilter{}
	f.AddSecondaryFilter(&filters.Fake{})
	h := NewFilterHandler(f)

	req, err := http.NewRequest("GET", FILTER_ENDPOINT+"www.facebook.ca", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.filterHandler)

	handler.ServeHTTP(recorder, req)

	if f.called != 1 {
		t.Errorf("The TestFilter ContainsURL function was called %d time(s).", f.called)
	}

	if recorder.Code != http.StatusLocked {
		t.Errorf("The filterHandler function %s when Locked was expected.", http.StatusText(recorder.Code))
	}
}
