package handlers

import (
	"github.com/tmortimer/urlfilter/filters"
	"log"
	"net/http"
)

const FILTER_ENDPOINT = "/urlinfo/1/"

// Holds a filter which the handler uses to check for banned URLs.
type FilterHandler struct {
	// The chain of filters used by this handler to see if a URL is flagged.
	filter filters.Filter
}

// Create a FilterHandler instance with the underlying filters.Filter chain.
func NewFilterHandler(filter filters.Filter) *FilterHandler {
	return &FilterHandler{filter: filter}
}

// Handles URL filtering requests.
func (f *FilterHandler) filterHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RequestURI()[len(FILTER_ENDPOINT):]
	found, err := f.filter.ContainsURL(url)

	// If we generated an error but the URL was found we can still act on
	// that information. If an error was generated but the URL was not found
	// we have to let the requester know we're unable to answer their request.
	if err != nil && !found {
		//TOM needs something more usefule here.
		w.WriteHeader(http.StatusInternalServerError)
	} else if found {
		// Return negative response, URL is banned.
		//TOM Not settled on API yet, just use this code to indicate error.
		w.WriteHeader(http.StatusLocked)
	} else {
		log.Printf("URL %s not found in the filter.", url)
		// Return positive response.
	}
}

// Initialize URL filter API.
func (f *FilterHandler) Init() {
	http.HandleFunc(FILTER_ENDPOINT, f.filterHandler)
}
