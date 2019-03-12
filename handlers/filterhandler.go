package handlers

import (
	"github.com/tmortimer/urlfilter/filters"
	"net/http"
)

const FILTER_ENDPOINT = "/urlinfo/1/"

// Holds a filter which the handler uses to check for banned URLs.
type FilterHandler struct {
	filter filters.Filter
}

func NewFilterHandler(filter filters.Filter) *FilterHandler {
	return &FilterHandler{filter: filter}
}

// Handles URL filtering requests.
func (f *FilterHandler) filterHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()[len(FILTER_ENDPOINT):]
	if f.filter.ContainsURL(url) {
		// Return negative response, URL is banned.
		// Not settled on API yet, just use this code to indicate error.
		w.WriteHeader(http.StatusLocked)
	}
	// Return positive response.
}

// Initialize URL filter API.
func (f *FilterHandler) Init() {
	http.HandleFunc(FILTER_ENDPOINT, f.filterHandler)
}
