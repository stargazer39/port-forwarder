package main

import (
	"log"
	"os"

	"golang.org/x/exp/slices"
)

func main() {
	proxyServerPort := ":8070"
	forwardPort := ":8080"

	proxyPort := ":8070"
	listenPort := ":8071"

	if slices.Contains(os.Args, "server-client") {
		log.Println("Starting in server-client mode")
		ServerClient(proxyServerPort, forwardPort)
		return
	}

	log.Println("Starting in proxy-server mode")
	ProxyServer(proxyPort, listenPort)
}
