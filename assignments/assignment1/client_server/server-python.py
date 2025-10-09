###############################################################################
# server-python.py
# Name:
# NetId:
###############################################################################

import sys
import socket

RECV_BUFFER_SIZE = 2048
QUEUE_LENGTH = 10

def server(server_port):
    # Open a TCP socket and listen for connections
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.bind(('', server_port))
    # Use a FIFO queue (list) of client sockets.
    sock.listen(QUEUE_LENGTH)
    # Use non-blocking mode to accept new connections while serving clients.
    sock.setblocking(False)
    queue = []
    while True:
        # Drain any pending connections into the queue without blocking
        while True:
            try:
                conn, addr = sock.accept()
            except BlockingIOError:
                break
            # handle this connection later; use blocking mode for client
            conn.setblocking(True)
            queue.append(conn)

        # If no queued clients, block until we get one
        if not queue:
            sock.setblocking(True)
            conn, addr = sock.accept()
            conn.setblocking(True)
            sock.setblocking(False)
            queue.append(conn)

        # Serve the next client in FIFO order
        conn = queue.pop(0)
        while True:
            # Read blocks from client and write to stdout
            data = conn.recv(RECV_BUFFER_SIZE)
            if not data:
                break
            sys.stdout.buffer.write(data)
        sys.stdout.flush()
        conn.close()


def main():
    """Parse command-line argument and call server function """
    if len(sys.argv) != 2:
        sys.exit("Usage: python server-python.py [Server Port]")
    server_port = int(sys.argv[1])
    server(server_port)

if __name__ == "__main__":
    main()
