// ...existing code...
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	// Parse server port
	var port int
	_, err := fmt.Sscanf(server_port, "%d", &port)
	if err != nil {
		log.Fatalf("Invalid server port %s", server_port)
	}
	if port <= 0 || port > 65535 {
		log.Fatalf("Invalid server port %d", port)
	}
	// Create a TCP socket and listen for incoming connections
	var listener net.Listener
	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error creating TCP listener: %v", err)
	}
	defer listener.Close()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
		}
		// Handle connection in a separate goroutine
		go func(c net.Conn) {
			defer c.Close()
			// Use bufio reader and writer and io.Copy to write to stdout
			reader := bufio.NewReader(c)
			writer := bufio.NewWriter(os.Stdout)
			_, err := io.Copy(writer, reader)
			if err != nil {
				log.Printf("Error copying data to stdout: %v", err)
			}
			// Ensure buffered stdout is flushed
			if err := writer.Flush(); err != nil {
				log.Printf("Error flushing stdout writer: %v", err)
			}
		}(conn)
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
