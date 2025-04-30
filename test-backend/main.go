package main

import (
	"flag"
	"log"
	"net/http"
)

var hostPort string

func init() {
	flag.StringVar(&hostPort, "host-port", "", "host and port of server")
}

func main() {
	flag.Parse()

	if hostPort == "" {
		log.Fatalf("failed to load port and host: empty")
	}

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(hostPort, nil); err != nil {
		log.Fatalf("failed to listen and serve: %s", err)
	}
}
