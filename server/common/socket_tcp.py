# socket_tcp.py
"""
Module for socket abstraction
"""
import socket
import logging

class SocketTCP:
    """
    Class that abstracts the TCP socket
    """
    def __init__(self, port: int, listen_backlog: int = 1000):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def accept_new_connection(self) -> socket.socket:
        """
        Accept new connection from the client
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(
            'action: accept_connections | result: success | ip: %s',
            addr[0]
        )
        return c

    def close(self) -> None:
        """
        Close the server socket
        """
        self._server_socket.close()

    def getpeername(self, client_sock: socket) -> tuple:
        """
        Get the client address
        """
        return client_sock.getpeername()

    @staticmethod
    def send(client_sock: socket, encoded_msg: bytes) -> None:
        """
        Send message to the client (and handle short writes)
        """
        length = len(encoded_msg)
        totalsent = 0
        while totalsent < length:
            sent = client_sock.send(encoded_msg[totalsent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            totalsent += sent

    @staticmethod
    def receive_all(client_sock: socket, size: int) -> bytes:
        """
        Receive message from the client (and handle short reads)
        """
        received = 0
        chunks = []
        while received < size:
            chunk = client_sock.recv(size - received)
            if chunk == b'':
                raise RuntimeError("socket connection broken")
            chunks.append(chunk)
            received += len(chunk)
        return b''.join(chunks) # Encoded message

    @staticmethod
    def send_all(client_sock: socket, size: int, message: bytes) -> None:
        """
        Send message to the client (and handle short writes)
        """
        totalsent = 0
        while totalsent < size:
            sent = client_sock.send(message[totalsent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            totalsent += sent
