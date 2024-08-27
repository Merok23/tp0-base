# protocol.py
"""
Module dedicating to encoding and encapsulating messages
"""
import socket
import logging
from common.socket_tcp import SocketTCP
from common.codes import ECHO_MESSAGE

CODE_SIZE = 4
SIZE_SIZE = 4

class Protocol:
    """
    Class that abstracts the protocol
    """

    ###################
    # Receive methods #
    ###################

    @staticmethod
    def receive_client_message(client_sock: socket) -> dict:
        """
        Receive message from the client (and handle short reads)
        """
        code = Protocol.__receive_message_code(client_sock)
        if code == ECHO_MESSAGE:
            logging.info("action: receive_client_message | result: success | code: %d", code)
            return {
                "code": code,
                "message": Protocol.__receive_message(client_sock),
            }
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

    #################
    # Send methods #
    #################

    @staticmethod
    def send_echo_response(client_sock: socket, message: str) -> None:
        """
        Send echo response to the client
        """
        Protocol.__send_message_code(client_sock, ECHO_MESSAGE)
        Protocol.__send_message(client_sock, message)

    @staticmethod
    def __send_message(client_sock: socket, message: str) -> None:
        """
        Send message to the client
        """
        encoded_message = message.encode('utf-8')
        size = len(encoded_message)
        Protocol.__send_message_size(client_sock, size)
        logging.info("action: send_message | result: in_progress | size: %d", size)
        SocketTCP.send_all(client_sock, size, encoded_message)

    @staticmethod
    def __send_message_code(client_sock: socket, code: int) -> None:
        """
        Send message code to the client
        """
        logging.info("action: send_message_code | result: in_progress | code: %d", code)
        SocketTCP.send_all(client_sock, CODE_SIZE, code.to_bytes(CODE_SIZE, byteorder='big'))

    @staticmethod
    def __send_message_size(client_sock: socket, size: int) -> None:
        """
        Send message size to the client
        """
        logging.info("action: send_message_size | result: in_progress | size: %d", size)
        SocketTCP.send_all(client_sock, SIZE_SIZE, size.to_bytes(SIZE_SIZE, byteorder='big'))
