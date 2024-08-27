# protocol.py
"""
Module dedicating to encoding and encapsulating messages
"""
import socket
import logging
from common.socket_tcp import SocketTCP

CODE_SIZE = 4
SIZE_SIZE = 4

class Protocol:
    """
    Class that abstracts the protocol
    """
    @staticmethod
    def receive_client_message(client_sock: socket) -> str:
        """
        Receive message from the client (and handle short reads)
        """
        code = Protocol.__receive_message_code(client_sock)
        if code == 1:
            logging.info("action: receive_client_message | result: success | code: %d", code)
            return Protocol.__receive_message(client_sock)
        raise ValueError("Invalid code")

    @staticmethod
    def __receive_message(client_sock: socket) -> str:
        """
        Receive message from the client
        """
        size = Protocol.__receive_message_size(client_sock)
        logging.info("action: receive_message | result: in_progress | size: %d", size)
        return SocketTCP.receive_all(client_sock, size).decode('utf-8')

    @staticmethod
    def __receive_message_code(client_sock: socket) -> int:
        """
        Receive message code from the client
        """
        code = client_sock.recv(CODE_SIZE)
        code = int.from_bytes(code, byteorder='big')
        logging.info("action: receive_message_code | result: success | code: %d", code)
        return code

    @staticmethod
    def __receive_message_size(client_sock: socket) -> int:
        """
        Receive message size from the client
        """
        size_bytes = client_sock.recv(SIZE_SIZE)
        size = int.from_bytes(size_bytes, byteorder='big')
        logging.info("action: receive_message_size | result: success | size: %d", size)
        return size
