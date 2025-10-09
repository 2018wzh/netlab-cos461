/*****************************************************************************
 * client-go.go
 * Name:
 * NetId:
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const SEND_BUFFER_SIZE = 2048

/* TODO: client()
 * Open socket and send message from stdin.
 */
func client(server_ip string, server_port string) {
	// Parse server address and port
	var serverAddr string
	var port int
	_, err := fmt.Sscanf(server_port, "%d", &port)
	if err != nil {
		log.Fatalf("Invalid server port %s", server_port)
	}
	if port <= 0 || port > 65535 {
		log.Fatalf("Invalid server port %d", port)
	}
	serverAddr = fmt.Sprintf("%s:%d", server_ip, port)
	// Create a TCP socket and connect to the server
	var conn net.Conn
	conn, err = net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// Read message from stdin and send in SEND_BUFFER_SIZE chunks for compatibility
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		message := scanner.Text()
		for i := 0; i < len(message); i += SEND_BUFFER_SIZE {
			end := i + SEND_BUFFER_SIZE
			if end > len(message) {
				end = len(message)
			}
			_, err := io.WriteString(conn, message[i:end])
			if err != nil {
				log.Fatalf("Error sending message to server: %v", err)
			}
		}
	}
	// Send a final newline to indicate end of message
	_, err = io.WriteString(conn, "\n")
	if err != nil {
		log.Fatalf("Error sending final newline to server: %v", err)
	}
}

// Main parses command-line arguments and calls client function
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./client-go [server IP] [server port] < [message file]")
	}
	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
