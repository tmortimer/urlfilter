// REST API handlers used by the urlfilter server.
package handlers

// REST API handler interface.
type Handler interface {
	// Initialize REST API endpoints.
	Init()
}
