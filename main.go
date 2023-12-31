package main

import (
	"log"
	"os"
	"strings"
)

var (
	Debug *log.Logger
	Info  *log.Logger
)

func main() {
	Debug = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	proxyServerPort := os.Getenv("PROXY_SERVER_ADDRESS")
	forwardPort := os.Getenv("PROXY_FORWARD_PORT")

	proxyPort := os.Getenv("PROXY_CLIENT_PORT")
	listenPort := os.Getenv("PROXY_LISTEN_PORT")

	if strings.Contains(strings.Join(os.Args, " "), "server-client") {
		Info.Println("Starting in server-client mode")
		ServerClient(proxyServerPort, forwardPort)
		return
	}

	Info.Println("Starting in proxy-server mode")
	ProxyServer(proxyPort, listenPort)
}
