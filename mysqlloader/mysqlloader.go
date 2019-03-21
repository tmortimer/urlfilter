// Utility to load the MySQL back end.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tjarratt/babble"
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/connectors"
	"log"
	"math/rand"
	"os"
)

// Adds URLs one at a time. Using domains-only.txt which has 26k+ domains in it, this is pretty slow.
// Grabbed that domain list from: http://mirror1.malwaredomains.com/files/domains.txt
func main() {
	configPath := flag.String("config", "", "Path to config file.")
	listPath := flag.String("list", "", "Path to config list of domains.")
	pathDepth := flag.Int("mpdepth", 0, "Max depth of path to add to domains.")
	queryDepth := flag.Int("mqdepth", 0, "Max depth of query to add to domains.")

	flag.Parse()

	config, err := config.ParseConfigFile(*configPath)
	if err != nil {
		log.Fatalf("Unable to load config: %s", err)
	}

	conn, err := connectors.NewMySQL(config.MySQL)
	if err != nil {
		log.Fatalf("Unable to connect to MySQL: %s", err)
	}

	list, err := os.Open(*listPath)
	defer list.Close()
	if err != nil {
		log.Fatalf("Unable to open URL list: %s", err)
	}

	babbler := babble.NewBabbler()
	scanner := bufio.NewScanner(list)
	count := 0
    for scanner.Scan() {
    	url := scanner.Text()

    	if *pathDepth != 0 {
    		depth := rand.Intn(*pathDepth)
    		if depth != 0 {
	    		babbler.Count = depth
	    		babbler.Separator = "/"
	    		url += "/"
	    		url += babbler.Babble()
	    	}
    	}

    	if *queryDepth != 0 {
    		babbler.Count = 2
    		depth := rand.Intn(*queryDepth)
    		if depth != 0 {
    			url += "?"
    			for i := 0; i < depth; i++ {
		    		babbler.Separator = "="
		    		url += babbler.Babble()

		    		if depth > i + 1 {
		    			url += "&"
		    		}
		    	}
	    	}
    	}

        conn.AddURL(url)
        fmt.Println(url)
        count++
    }

    fmt.Printf("Added %d URLs to the DB.\n", count)

    if err = scanner.Err(); err != nil {
        log.Fatalf("URL list scanner failed: %s", err)
    }
}