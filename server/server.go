// Main urlfilter server.
package server

import (
	"github.com/tmortimer/urlfilter/handlers"
	"log"
)

// HTTP Server Interface, backs the main server instance.
type HTTPServer interface {
	// Start the server.
	ListenAndServe() error
}

// Initializes REST API handlers and launch the server.
func Run(handlers []handlers.Handler, s HTTPServer) {
	for _, handler := range handlers {
		handler.Init()
	}

	// ListenAndServe launches a goroutine for each connection,
	// so no additional handling necessary to get some concurrency.git
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
