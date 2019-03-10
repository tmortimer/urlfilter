package handlers

import (
	"log"
	"net/http"
)

// Empty type to implete interfaces such as handlers.Initializer.
type Filter struct {}


// Function which handles URL filtering requests.
func filterHandler(w http.ResponseWriter, r *http.Request) {
	// Magic happens here
}

// Initialize URL filter API.
func (f Filter) Init() {
	http.HandleFunc("/urlinfo/1/", filterHandler)
}