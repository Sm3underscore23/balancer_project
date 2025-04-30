package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var hostPort string

func init() {
	flag.StringVar(&hostPort, "host-port", "", "host and port of server")
}

func main() {
	flag.Parse()

	if hostPort == "" {
		hostPort = os.Getenv("HOST_PORT")
	}

	log.Println(hostPort)

	if hostPort == "" {
		log.Fatalf("failed to get host and port: empty")
	}

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(hostPort, nil); err != nil {
		log.Fatalf("failed to listen and serve: %s", err)
	}
}
