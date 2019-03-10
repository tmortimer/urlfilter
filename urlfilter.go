package main

import (
	"log"
	"net/http"
)

// Conifgure and launch URL filtering service which can be used
// to check against known malicious URLs.
func main() {
    log.Fatal(http.ListenAndServe(":8080", nil))
}
