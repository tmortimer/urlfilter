package main

import (
	"flag"
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/filters"
	"github.com/tmortimer/urlfilter/handlers"
	"github.com/tmortimer/urlfilter/server"
	"log"
	"net/http"
)

// Conifgure and launch URL filtering service which can be used
// to check against known malicious URLs.
func main() {
	configPath := flag.String("config", "", "Path to config file.")
	flag.Parse()
	config, err := config.ParseConfigFile(*configPath)
	if err != nil {
		log.Fatalf("Unable to load config: %s", err)
	}

	filter, err := filters.FilterFactory(config)
	if err != nil {
		log.Fatalf("Unable to configure filter chain: %s", err)
	}

	handlers := []handlers.Handler{
		handlers.NewFilterHandler(filter),
	}

	server.Run(handlers, &http.Server{Addr: config.Host + ":" + config.Port})
}
