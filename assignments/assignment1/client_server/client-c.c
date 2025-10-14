/*****************************************************************************
 * client-c.c
 * Name:
 * NetId:
 *****************************************************************************/

#include <errno.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>

#define SEND_BUFFER_SIZE 2048

/* TODO: client()
 * Open socket and send message from stdin.
 * Return 0 on success, non-zero on failure
 */
int client(char *server_ip, char *server_port) {
  // Create socket
  int port = atoi(server_port);
  // Check port number
  if (port <= 0 || port > 65535) {
    fprintf(stderr, "Invalid port number: %s\n", server_port);
    return -1;
  }
  // Create socket
  int sockfd;
  sockfd = socket(AF_INET, SOCK_STREAM, 0);
  if (sockfd < 0) {
    perror("Error opening socket");
    return -1;
  }
  // Get server address
  struct hostent *server;
  server = gethostbyname(server_ip);
  if (server == NULL) {
    fprintf(stderr, "Error, no such host: %s\n", server_ip);
    return -1;
  }
  // Connect to server
  struct sockaddr_in serv_addr;
  memset(&serv_addr, 0, sizeof(serv_addr));
  serv_addr.sin_family = AF_INET;
  memcpy(&serv_addr.sin_addr.s_addr, server->h_addr_list[0],
         server->h_length); // Use first address
  serv_addr.sin_port = htons(port);
  if (connect(sockfd, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
    perror("Error connecting to server");
    return -1;
  }
  // Read message from stdin in chunks of SEND_BUFFER_SIZE and send to server
  char buffer[SEND_BUFFER_SIZE];
  while (1) {
    size_t n = fread(buffer, 1, SEND_BUFFER_SIZE, stdin);
    if (n == 0) {
      if (feof(stdin)) {
        break; // EOF reached
      }
      if (ferror(stdin)) {
        perror("Error reading from stdin");
        close(sockfd);
        return -1;
      }
      // n == 0 but neither feof nor ferror? treat as EOF
      break;
    }
    // Send n bytes to server, loop until all bytes sent
    size_t total_sent = 0;
    while (total_sent < n) {
      ssize_t sent = send(sockfd, buffer + total_sent, n - total_sent, 0);
      if (sent < 0) {
        perror("Error sending to server");
        close(sockfd);
        return -1;
      }
      if (sent == 0) {
        // remote closed connection unexpectedly
        fprintf(stderr, "Peer closed connection while sending\n");
        close(sockfd);
        return -1;
      }
      total_sent += (size_t)sent;
    }
    // printf("Sent %zu bytes to server\n", total_sent);
  }

  // Finished sending data; optionally signal EOF to server and close socket
  if (shutdown(sockfd, SHUT_WR) < 0) {
    perror("Error shutting down socket write side");
    close(sockfd);
    return -1;
  }
  close(sockfd);
  return 0;
}

/*
 * main()
 * Parse command-line arguments and call client function
 */
int main(int argc, char **argv) {
  char *server_ip;
  char *server_port;

  if (argc != 3) {
    fprintf(stderr,
            "Usage: ./client-c [server IP] [server port] < [message]\n");
    exit(EXIT_FAILURE);
  }

  server_ip = argv[1];
  server_port = argv[2];
  return client(server_ip, server_port);
}
