package main

import (
	"flag"
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/filters"
	"github.com/tmortimer/urlfilter/handlers"
	"github.com/tmortimer/urlfilter/server"
	"net/http"
)

// Conifgure and launch URL filtering service which can be used
// to check against known malicious URLs.
func main() {
	configPath := flag.String("config", "", "Path to config file.")
	flag.Parse()
	config := config.ParseConfigFile(*configPath)

	handlers := []handlers.Handler{
		handlers.NewFilterHandler(filters.NewFake()),
	}

	server.Run(handlers, &http.Server{Addr: config.Host + ":" + config.Port})
}
