package handlers

// Interface to initialize REST API handlers.
type Initializer interface {
	// Initialize REST API endpoints.
	Init()
}
