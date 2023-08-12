package main

import (
	"log"
	"os"

	"golang.org/x/exp/slices"
)

var (
	Debug *log.Logger
	Info  *log.Logger
)

func main() {
	Debug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	proxyServerPort := ":8070"
	forwardPort := ":8080"

	proxyPort := ":8070"
	listenPort := ":8071"

	if slices.Contains(os.Args, "server-client") {
		Info.Println("Starting in server-client mode")
		ServerClient(proxyServerPort, forwardPort)
		return
	}

	Info.Println("Starting in proxy-server mode")
	ProxyServer(proxyPort, listenPort)
}
