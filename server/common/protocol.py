# protocol.py
"""
Module dedicating to encoding and encapsulating messages
"""
import socket
import logging
from common.socket_tcp import SocketTCP
from common.codes import ECHO_MESSAGE, BET_MESSAGE
from common.codes import SUCCESS_CODE, ERROR_CODE

CODE_SIZE = 4
SIZE_SIZE = 4
DATE_OF_BIRTH_SIZE = 10

class Protocol:
    """
    Class that abstracts the protocol
    """

    ###################
    # Receive methods #
    ###################

    @staticmethod
    def __receive_bet(client_sock: socket) -> dict:
        """
        Receive bet message from the client
        """
        dni = Protocol.__receive_uint32(client_sock)
        number = Protocol.__receive_uint32(client_sock)
        date_of_birth = SocketTCP.receive_all(client_sock, DATE_OF_BIRTH_SIZE)
        agency = Protocol.__receive_uint32(client_sock)
        name = Protocol.__receive_variable_string(client_sock)
        lastname = Protocol.__receive_variable_string(client_sock)
        return {
            "code": BET_MESSAGE,
            "dni": str(dni),
            "number": str(number),
            "date_of_birth": str(date_of_birth, 'utf-8'),   
            "name": name,   
            "lastname": lastname,
            "agency": agency,
        }

    @staticmethod
    def __receive_bets(client_sock: socket) -> dict:
        """
        Receive bet message from the client
        """
        bets = []
        bets_size = Protocol.__receive_uint32(client_sock)
        for _ in range(bets_size):
            bet = Protocol.__receive_bet(client_sock)
            bets.append(bet)

        return {
            "code": BET_MESSAGE,
            "bets": bets
        }

    @staticmethod
    def receive_client_message(client_sock: socket) -> dict:
        """
        Receive message from the client
        """
        code = Protocol.__receive_uint32(client_sock)
        if code == ECHO_MESSAGE:
            message = Protocol.__receive_variable_string(client_sock)
            return {
                "code": code,
                "message": message
            }
        if code == BET_MESSAGE:
            return Protocol.__receive_bets(client_sock)
        raise ValueError("Invalid code")

    @staticmethod
    def __receive_variable_string(client_sock: socket) -> str:
        """
        Receive message from the client
        """
        size = Protocol.__receive_message_size(client_sock)
        return SocketTCP.receive_all(client_sock, size).decode('utf-8')

    @staticmethod
    def __receive_uint32(client_sock: socket) -> int:
        """
        Receive message code from the client
        """
        code = client_sock.recv(CODE_SIZE)
        code = int.from_bytes(code, byteorder='big')
        return code

    @staticmethod
    def __receive_message_size(client_sock: socket) -> int:
        """
        Receive message size from the client
        """
        size_bytes = client_sock.recv(SIZE_SIZE)
        size = int.from_bytes(size_bytes, byteorder='big')
        return size

    #################
    # Send methods #
    #################

    @staticmethod
    def send_echo_response(client_sock: socket, message: str) -> None:
        """
        Send echo response to the client
        """
        Protocol.__send_message(client_sock, message)

    @staticmethod
    def send_bet_response_succesful(client_sock: socket) -> None:
        """
        Send bet response to the client
        """
        Protocol.__send_message_code(client_sock, SUCCESS_CODE)

    @staticmethod
    def send_bet_response_error(client_sock: socket) -> None:
        """
        Send bet response to the client
        """
        Protocol.__send_message_code(client_sock, ERROR_CODE)

    @staticmethod
    def __send_message(client_sock: socket, message: str) -> None:
        """
        Send message to the client
        """
        encoded_message = message.encode('utf-8')
        size = len(encoded_message)
        Protocol.__send_message_size(client_sock, size)
        SocketTCP.send_all(client_sock, size, encoded_message)

    @staticmethod
    def __send_message_code(client_sock: socket, code: int) -> None:
        """
        Send message code to the client
        """
        SocketTCP.send_all(
            client_sock,
            CODE_SIZE,
            code.to_bytes(CODE_SIZE, byteorder='big', signed=False)
        )

    @staticmethod
    def __send_message_size(client_sock: socket, size: int) -> None:
        """
        Send message size to the client
        """
        SocketTCP.send_all(
            client_sock,
            SIZE_SIZE,
            size.to_bytes(SIZE_SIZE, byteorder='big', signed=False)
        )
