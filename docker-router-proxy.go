package main

import (
	"flag"
	"fmt"
	"net/http"
	"./proxy"
	"./routing"
)

func main() {

	const (
		defaultPort = ":9800"
		defaultPortUsage = "default server port, ':9800', ':8080'..."
		defaultDomain = ".dock"
		defaultDomainUsage = "default docker domain to remove from serialization, eg. '.dock'"
		debugDefault = false
		debugHint = "Print som debug informations"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	debug := flag.Bool("debug", debugDefault, debugHint)
	domain := flag.String("domain", defaultDomain, defaultDomainUsage)
	flag.Parse()

	// debugging?
	if *debug {
		fmt.Println("Domain:", *domain)
		fmt.Println("Server will run on:", *port)
	}

	// create router cache
	cache := new(routing.Cache).Init()

	// create router
	router := &proxy.ProxyRouter{Debug: *debug, Cache: cache, Domain: *domain}

	// start server
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})
	http.ListenAndServe(*port, nil)
}
