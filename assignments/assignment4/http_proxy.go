/*****************************************************************************
 * http_proxy.go
 * Names:
 * NetIds:
 *****************************************************************************/

// TODO: implement an HTTP proxy

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func server(listen_port string) {
	var port int
	_, err := fmt.Sscanf(listen_port, "%d", &port)
	if err != nil {
		log.Fatalf("Invalid listen port %s", listen_port)
	}
	if port <= 0 || port > 65535 {
		log.Fatalf("Invalid listen port %d", port)
	}
	var listener net.Listener
	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error creating TCP listener: %v", err)
	}
	defer listener.Close()
	log.Printf("Listening on :%d", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handle(conn)
	}
}

func errhandle(conn net.Conn) {
	defer conn.Close()
	resp := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\n\r\n")
	_, _ = conn.Write([]byte(resp))
}

func handle(conn net.Conn) {
	defer conn.Close()

	br := bufio.NewReader(conn)
	for {
		req, err := http.ReadRequest(br)
		if err != nil {
			if err == io.EOF {
				return
			}
			// If it's a timeout or temporary, just log and return
			log.Printf("Read request error: %v", err)
			return
		}

		// Print the request
		log.Printf("----- Received HTTP request -----\n")
		log.Printf("%s %s %s\n", req.Method, req.URL.String(), req.Proto)
		for k, vals := range req.Header {
			for _, v := range vals {
				log.Printf("%s: %s\n", k, v)
			}
		}

		// Only handle GET requests for now
		if req.Method != "GET" {
			log.Printf("Unsupported method: %s", req.Method)
			errhandle(conn)
			continue
		}

		// Build outgoing request to forward
		outReq, err := http.NewRequest(req.Method, req.URL.String(), nil)
		if err != nil {
			log.Printf("Error creating new request: %v", err)
			errhandle(conn)
			return
		}

		// Copy headers
		for k, vals := range req.Header {
			for _, v := range vals {
				outReq.Header.Add(k, v)
			}
		}
		// Set header connection close
		outReq.Header.Set("Connection", "close")

		// Perform the request
		resp, err := http.DefaultTransport.RoundTrip(outReq)
		if err != nil {
			log.Printf("Error forwarding request: %v", err)
			errhandle(conn)
			return
		}

		// Write response back to the client
		err = resp.Write(conn)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error writing response to client: %v", err)
			return
		}

		// If either side requested connection close, break
		if req.Close || resp.Close {
			return
		}
	}
}

func main() {
	log.SetOutput(os.Stderr)
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./http_proxy [listen port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
