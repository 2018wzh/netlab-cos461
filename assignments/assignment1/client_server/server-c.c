/*****************************************************************************
 * server-c.c
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

#define QUEUE_LENGTH 10
#define RECV_BUFFER_SIZE 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 * Return 0 on success, non-zero on failure
 */
int server(char *server_port) {
  // Listen to socket on the given port
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
  // Bind socket to port
  struct sockaddr_in serv_addr;
  memset(&serv_addr, 0, sizeof(serv_addr));
  serv_addr.sin_family = AF_INET;
  serv_addr.sin_addr.s_addr = INADDR_ANY;
  serv_addr.sin_port = htons(port);
  if (bind(sockfd, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
    perror("Error binding socket");
    return -1;
  }
  // Listen for incoming connections
  if (listen(sockfd, QUEUE_LENGTH) < 0) {
    perror("Error listening on socket");
    return -1;
  }
  // printf("Server ready, listening on port %d\n", port);
  //  Accept incoming connections
  while (1) {
    // Get client connection
    struct sockaddr_in cli_addr;
    socklen_t clilen = sizeof(cli_addr);
    int clisockfd = accept(sockfd, (struct sockaddr *)&cli_addr, &clilen);
    if (clisockfd < 0) {
      perror("Error accepting connection");
      continue;
    }
    // printf("Client connected\n");
    while (1) {
      // Receive message from client
      char buffer[RECV_BUFFER_SIZE];
      memset(buffer, 0, RECV_BUFFER_SIZE);
      ssize_t n = recv(clisockfd, buffer, RECV_BUFFER_SIZE - 1, 0);
      if (n < 0) {
        perror("Error receiving from client");
        close(clisockfd);
        break;
      } else if (n == 0) {
        // printf("Client disconnected\n");
        close(clisockfd);
        break;
      }
      // Print received message to stdout
      size_t total_written = 0;
      while (total_written < n) {
        ssize_t written =
            fwrite(buffer + total_written, 1, n - total_written, stdout);
        if (written < 0) {
          perror("Error writing to stdout");
          close(clisockfd);
          return -1;
        }
        total_written += written;
      }
    }
    fflush(stdout);
    // Close client socket
    close(clisockfd);
  }
  close(sockfd);
  return 0;
}

/*
 * main():
 * Parse command-line arguments and call server function
 */
int main(int argc, char **argv) {
  char *server_port;

  if (argc != 2) {
    fprintf(stderr, "Usage: ./server-c [server port]\n");
    exit(EXIT_FAILURE);
  }

  server_port = argv[1];
  return server(server_port);
}
