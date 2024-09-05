import socket
import logging
import signal
import sys
import os
import threading
from common.socket_tcp import SocketTCP
from common.protocol import Protocol
from common.codes import ECHO_MESSAGE, BET_MESSAGE
from common.codes import END_MESSAGE
from common.utils import Bet
from common.utils import store_bets
from common.utils import load_bets
from common.utils import has_won



class Server:
    def __init__(self, port, listen_backlog):
        self._server_socket = SocketTCP(port, listen_backlog)
        self._lotery_agencies_done = 0
        self._clients = {}
        self._client_threads = []
        self._client_sockets = []
        self._server_lock = threading.Lock()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        signal.signal(signal.SIGINT, self.__handle_shutdown)
        signal.signal(signal.SIGTERM, self.__handle_shutdown)
        try:
            while True:
                self.__reap_dead_threads()
                client_sock = self.__accept_new_connection()
                client_thread = threading.Thread(
                    target=self.__handle_client_connection,
                    args=(client_sock,)
                )
                client_thread.start()
                self._client_sockets.append(client_sock)
                self._client_threads.append(client_thread)
        finally:
            for client in self._client_sockets:
                client.close() ## This unblocks any waiting threads
            self._server_socket.close()
            for thread in self._client_threads:
                thread.join()
            self.__reap_dead_threads()

    def __reap_dead_threads(self):
        """
        Reap dead threads
        """
        for thread in self._client_threads:
            if not thread.is_alive():
                self._client_threads.remove(thread)

    def __handle_shutdown(self, signum, frame):
        """
        Handle server shutdown gracefully
        """
        logging.info(
            "action: shutdown_server | result: in_progress | frame: %s",
            frame.f_code.co_name
        )
        self._server_socket.close()
        sig_name = signal.Signals(signum)
        logging.info("action: shutdown_server | result: success | %s", sig_name.name)
        sys.exit(0)

    def __handle_client_connection(self, client_sock: socket) -> None:
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            msg = Protocol.receive_client_message(client_sock)
            if msg['code'] == ECHO_MESSAGE:
                self.__handle_echo(client_sock, msg)
            if msg['code'] == BET_MESSAGE:
                self.__handle_bet(client_sock, msg)
            if msg['code'] == END_MESSAGE:
                self.__handle_end_message(client_sock, msg)
        except ValueError as e:
            logging.error("action: receive_message | result: fail | error: %s", format(e))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: %s", format(e))


    def __handle_end_message(self, client_sock: socket, msg) -> None:
        """
        Handle end message from the client
        """
        with self._server_lock:
            self._lotery_agencies_done += 1
            self._clients[msg['agency_number']] = client_sock
            if self._lotery_agencies_done == int(os.environ.get('TOTAL_AGENCIES')):
                logging.info("action: sorteo | result: success")
                bets = load_bets()
                winning_count = {}
                winning_dnis = []
                for bet in bets:
                    if has_won(bet):
                        winning_count[bet.agency] = winning_count.get(bet.agency, 0) + 1
                        winning_dnis.append(bet.document)
                for agency, agency_socket in self._clients.items():
                    Protocol.send_winners(
                        agency_socket,
                        winning_count.get(agency, 0),
                        winning_dnis
                    )
                    agency_socket.close()
                self._clients.clear()

    def __handle_bet(self, client_sock: socket, msg: dict) -> None:
        """
        Handle bet request from the client
        """
        bets = []
        for bet in msg['bets']:
            if int(bet['number']) < 0:
                logging.error(
                    "action: apuesta_recibida | result: fail | cantidad: %d", len(bets)
                )
                Protocol.send_bet_response_error(client_sock)
                return
            bets.append(Bet(
                bet['agency'],
                bet['name'],
                bet['lastname'],
                bet['dni'],
                bet['date_of_birth'],
                bet['number'],
            ))
        with self._server_lock:
            store_bets(bets)
            logging.info(
                "action: apuesta_recibida | result: success | cantidad: %d",
                len(bets)
            )
        Protocol.send_bet_response_succesful(client_sock)

    def __handle_echo(self, client_sock: socket, msg: dict) -> None:
        """
        Handle echo request from the client
        """
        msg = msg['message']
        addr = client_sock.getpeername()
        logging.info(
            'action: receive_message | result: success | ip: %s | msg: %s',
            addr[0],
            msg,
        )
        Protocol.send_echo_response(client_sock, msg)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        return self._server_socket.accept_new_connection()
