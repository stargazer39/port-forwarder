package main

func main() {

}

// Client that sits on the server side and then talks to the proxy server
func ServerClient() {
	// Create a TCP Connection to the server
	// Listen for a CONNECTION request from the TCP
	// Create a connection to the Proxy server
	// Dial Server side service
	// Copy it into the Created connection
}

// Server that sits in the proxy server
func ProxyServer() {
	// Listen for ServerClient connection
	// Listen on a port that allow others to connect
	// When connection comes in get ref to it
	// Send a CONNECTION request down the ServerClient
	// Wait for a connection to get created
	// Copy the connection ref to the new connection
}
