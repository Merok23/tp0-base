# protocol.py
"""
Module dedicating to encoding and encapsulating messages
"""
import socket
import logging
from common.socket_tcp import SocketTCP

class Protocol:
    """
    Class that abstracts the protocol
    """
    @staticmethod
    def receive_client_message(client_sock: socket) -> str:
        """
        Receive message from the client (and handle short reads)
        """
        size_bytes = client_sock.recv(4) # TODO: Formalize protocol.
        size = int.from_bytes(size_bytes, byteorder='big') # network to host
        logging.info("action: receive_message | result: in_progress | size: %d", size)
        return SocketTCP.receive_all(client_sock, size).decode('utf-8')
