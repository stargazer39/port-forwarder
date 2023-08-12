package main

import (
	"log"
	"net"
	"sync"
)

// Server that sits in the proxy server
func ProxyServer(proxyPort string, listenPort string) {
	var refs []*IncomingProxyConnection
	var refLock sync.Mutex

	refChan := make(chan bool)

	// Listen on a port that allow others to connect
	ll, err := net.Listen("tcp", listenPort)

	if err != nil {
		log.Panicf("Listening on %s failed. Error: %s\n", listenPort, err)
	}

	// Listen for ServerClient connection
	l, err := net.Listen("tcp", proxyPort)

	if err != nil {
		log.Panicf("Listening on %s failed. Error: %s\n", proxyPort, err)
	}

	go func() {
		log.Println("Listening on proxy port " + proxyPort)
		for {
			conn, err := l.Accept()

			if err != nil {
				log.Printf("Failed to accept. Error: %s", err)
				continue
			}

			// When connection comes in get ref to it
			var initialMessage IncomingProxyMessage
			var ipc IncomingProxyConnection

			if err := ReadJSON(conn, &initialMessage); err != nil {
				log.Printf("Failed to read. Error: %s", err)
				continue
			}
			log.Println("Read message ", initialMessage)

			ipc.ID = initialMessage.ID
			ipc.Type = initialMessage.Type
			ipc.Connection = conn

			refLock.Lock()
			refs = append(refs, &ipc)
			refLock.Unlock()
			refChan <- true
		}
	}()

	for {
		log.Println("Listening on listening port " + listenPort)
		conn, err := ll.Accept()

		log.Println("Got connection to " + listenPort)

		if err != nil {
			log.Printf("Failed to accept. Error: %s", err)
		}
		errr := false

	out:
		for {
			for _, r := range refs {
				if r.Type == 0 {
					log.Println("Found a type 0 connection")
					// Send a CONNECTION request down the ServerClient
					if err := WriteJSON(r.Connection, MessageToClient{
						Message: "CONNECTION",
					}); err != nil {
						log.Println(err)
						errr = true
						break out
					}

					log.Println("Wrote a message tp  type 0", "CONNECTION")
					break out
				}
			}
		}
		if errr {
			continue
		}

		done := false
	pp:
		for {
			<-refChan
			for _, r := range refs {
				if r.Type == 1 && !r.Locked {
					log.Println("Got a type 1 unlocked connection")
					r.Locked = true
					conn2 := r.Connection
					done = true

					// Copy the connection ref to the new connection
					go HandleTCPConn(conn2, conn)
				}
			}
			if done {
				break pp
			}
			// Wait for a connection to get created
			// time.Sleep(time.Millisecond * 100)
		}

	}
}
