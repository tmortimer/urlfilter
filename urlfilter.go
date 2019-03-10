package main

import (
	"github.com/tmortimer/urlfilter/handlers"
	"github.com/tmortimer/urlfilter/server"
)

// Conifgure and launch URL filtering service which can be used
// to check against known malicious URLs.
func main() {
	handlers := []handlers.Initializer {
		handlers.Filter {},
	}
    server.Run(handlers)
}
