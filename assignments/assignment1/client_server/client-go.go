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

	// Read into a buffer of size SEND_BUFFER_SIZE and write each chunk fully to the connection.
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, SEND_BUFFER_SIZE)
	for {
		n, rerr := reader.Read(buf)
		if n > 0 {
			// ensure full write of buf[:n]
			written := 0
			for written < n {
				wn, werr := conn.Write(buf[written:n])
				if werr != nil {
					log.Fatalf("Error sending data to server: %v", werr)
				}
				written += wn
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				break
			}
			log.Fatalf("Error reading from stdin: %v", rerr)
		}
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
