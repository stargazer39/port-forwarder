package main

import (
	"log"
	"net"
	"time"
)

// Client that sits on the server side and then talks to the proxy server
func ServerClient(proxyServerPort string, forwardPort string) {
	// Create a TCP Connection to the server
retry:
	Info.Printf("Connecting to proxy server %s \n", proxyServerPort)
	conn, err := net.Dial("tcp", proxyServerPort)

	if err != nil {
		Info.Println("Listening on %s failed. Error: %s\n", proxyServerPort, err)
		Info.Println("Retrying...")
		goto retry
	}

	log.Println("Sending message with type 0 connection")
	if err := WriteJSON(conn, IncomingProxyMessage{
		ID:   "0",
		Type: 0,
	}); err != nil {
		Info.Printf("Error: %s\n", err)
		Info.Println("Retrying...")
		goto retry
	}

	for {
		var message MessageToClient

		// Listen for a CONNECTION request from the TCP
		if err := ReadJSON(conn, &message); err != nil {
			Info.Println(err)
			goto retry
		}

		log.Println("Got message " + message.Message)

		if message.Message == "CONNECTION" {
			// Create a connection to the Proxy server
			count := 0
		pServerRetry:
			if count == 5 {
				goto retry
			}

			pServerConn, err := net.Dial("tcp", proxyServerPort)
			count++

			if err != nil {
				Info.Printf("Listening on %s failed. Error: %s\n", proxyServerPort, err)
				time.Sleep(time.Millisecond * 200)

				goto pServerRetry
			}
			count = 0
		fServerRetry:
			if count == 5 {
				goto retry
			}

			count++
			// Dial Server side service
			conn3, err := net.Dial("tcp", forwardPort)

			if err != nil {
				Info.Printf("Listening on %s failed. Error: %s\n", proxyServerPort, err)
				time.Sleep(time.Millisecond * 200)
				goto fServerRetry
			}

			if err := WriteJSON(pServerConn, IncomingProxyConnection{
				ID:   "0",
				Type: 1,
			}); err != nil {
				log.Println(err)
				goto retry
			}

			// Copy it into the Created connection
			go CopyTCP(pServerConn, conn3)
		}
	}
}
