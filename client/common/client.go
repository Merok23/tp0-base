package common

import (
	"net"
	"time"
	"strconv"
	"os"
	"bufio"
	"strings"

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
	id, err := strconv.Atoi(c.config.ID)
	if err != nil {
		log.Errorf("action: convert_id | result: fail | client_id: %v | error: %v",
		c.config.ID,
		err,
	)
	return
}
file := os.Getenv("FILE")
f, err := os.Open(file)
	if err != nil {
		log.Errorf("action: open_file | result: fail | client_id: %v | error: %v",
		c.config.ID,
		err,
		)
		return
	}
	defer f.Close()
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		c.createClientSocket()
		line := scanner.Text()
		data := strings.Split(line, ",")
		name := data[0]
		lastname := data[1]
		dni, err := strconv.Atoi(data[2])
		dateOfBirth := data[3]
		number, err := strconv.Atoi(data[4])
		if err != nil {
			log.Errorf("action: convert_data | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		_, err = protocol.SendBet(c.conn, id, dni, name, lastname, dateOfBirth, number)
		if err != nil {
			log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof(
			"action: apuesta_enviada | result: success | dni: %v | numero: %v",
			dni,
			number,
		)
		c.conn.Close()
	}

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
