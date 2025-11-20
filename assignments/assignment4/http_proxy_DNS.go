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
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
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

		// DNS prefetching
		if resp.Header.Get("Content-Type") == "text/html" {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				errhandle(conn)
				return
			}
			bodyStr := string(bodyBytes)
			go prefetch(bodyStr)
			// Restore the body for writing back to client
			resp.Body = ioutil.NopCloser(strings.NewReader(bodyStr))
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

// DNS prefetch
func prefetch(body string) {
	// Parse the HTML body
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return
	}

	// Traverse the HTML nodes to find <a> tags with href attributes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
					hostname := extractHostname(attr.Val)
					if hostname != "" {
						go func(host string) {
							addrs, err := net.LookupHost(host)
							if err != nil {
								log.Printf("DNS lookup failed for %s: %v", host, err)
							} else {
								log.Printf("DNS lookup successful for %s: %v", host, addrs)
							}
						}(hostname)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
}

func extractHostname(url string) string {
	// Remove the "http://" or "https://" prefix
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	// Extract the hostname (up to the first "/")
	parts := strings.SplitN(url, "/", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./http_proxy [listen port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
