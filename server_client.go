package main

import (
	"log"
	"net"
)

// Client that sits on the server side and then talks to the proxy server
func ServerClient(proxyServerPort string, forwardPort string) {
	// Create a TCP Connection to the server
	conn, err := net.Dial("tcp", proxyServerPort)

	if err != nil {
		log.Panicf("Listening on %s failed. Error: %s\n", proxyServerPort, err)
	}

	log.Println("Sending message with type 0 connection")
	if err := WriteJSON(conn, IncomingProxyMessage{
		ID:   "0",
		Type: 0,
	}); err != nil {
		log.Panicln("Error: ", err)
	}

	for {
		var message MessageToClient

		// Listen for a CONNECTION request from the TCP
		if err := ReadJSON(conn, &message); err != nil {
			log.Panicln(err)
		}

		log.Println("Got message " + message.Message)

		if message.Message == "CONNECTION" {
			// Create a connection to the Proxy server
			connn, err := net.Dial("tcp", proxyServerPort)

			if err != nil {
				log.Panicf("Listening on %s failed. Error: %s\n", proxyServerPort, err)
			}

			// Dial Server side service
			conn3, err := net.Dial("tcp", forwardPort)

			if err != nil {
				log.Panicf("Listening on %s failed. Error: %s\n", proxyServerPort, err)
			}

			if err := WriteJSON(connn, IncomingProxyConnection{
				ID:   "0",
				Type: 1,
			}); err != nil {
				log.Println(err)
				continue
			}

			// Copy it into the Created connection
			go HandleTCPConn(connn, conn3)
		}
	}
}
