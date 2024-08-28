package common

import (
	"net"
	"time"
	"strconv"
	"os"

	"github.com/op/go-logging"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) StartClientLoop() {
	c.createClientSocket()

	id, err := strconv.Atoi(c.config.ID)
	if err != nil {
		log.Errorf("action: convert_id | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	dni, err := strconv.Atoi(os.Getenv("DNI"))
	number, err := strconv.Atoi(os.Getenv("NUMERO"))
	if err != nil {
		log.Errorf(
			"action: convert | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	name := os.Getenv("NOMBRE")
	lastname := os.Getenv("APELLIDO")
	dateOfBirth := os.Getenv("NACIMIENTO")
	// Send the message
	_, err = protocol.SendBet(c.conn, id, dni, name, lastname, dateOfBirth, number)

	if err != nil {
		log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	response, err := protocol.ReceiveBet(c.conn)
	if err != nil {
		log.Errorf("action: receive_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	if response != protocol.CODE_SUCCESS {
		log.Errorf("action: receive_bet | result: fail | client_id: %v | response: %v",
			c.config.ID,
			response,
		)
		return
	}

	log.Infof(
		"action: apuesta_enviada | result: success | dni: %v | numero: %v",
		dni,
		number,
	)
	// Close the connection
	c.conn.Close()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
}

func (c *Client) StopClientLoop() {
	c.conn.Close()
}
