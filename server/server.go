package server

import (
	"log"
	"net/http"
	"github.com/tmortimer/urlfilter/handlers"
)

// Initializes REST API handlers and launch the server.
func Run(handlers []handlers.Initializer) {
	for _, handler := range handlers {
		handler.Init()
	}

	// ListenAndServe launches a goroutine for each connection,
	// so no additional handling necessary to get some concurrency.git
	log.Fatal(http.ListenAndServe(":8080", nil))
}
