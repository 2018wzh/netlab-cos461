###############################################################################
# client-python.py
# Name:
# NetId:
###############################################################################

import sys
import socket

SEND_BUFFER_SIZE = 2048

def client(server_ip, server_port):
    # Open a TCP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    # Connect to server
    sock.connect((server_ip, server_port))
    # Read from stdin and send to server
    while True:
        data = sys.stdin.buffer.read(SEND_BUFFER_SIZE)
        if not data:
            break
        sock.sendall(data)
    # Close the socket
    sock.close()


def main():
    """Parse command-line arguments and call client function """
    if len(sys.argv) != 3:
        sys.exit("Usage: python client-python.py [Server IP] [Server Port] < [message]")
    server_ip = sys.argv[1]
    server_port = int(sys.argv[2])
    client(server_ip, server_port)

if __name__ == "__main__":
    main()
