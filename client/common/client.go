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
	Batch         int
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

func createBetFromLine(line string, agencyNumber int) (protocol.Bet, error) {
	data := strings.Split(line, ",")
	name := data[0]
	lastname := data[1]
	dni, err := strconv.Atoi(data[2])
	dateOfBirth := data[3]
	number, err := strconv.Atoi(data[4])
	bet := protocol.Bet{
		Dni: dni,
		Name: name,
		Lastname: lastname,
		DateOfBirth: dateOfBirth,
		Number: number,
		AgencyNumber: agencyNumber,
	}
	return bet, err
}

func checkForErrorsBet(response uint32, err error) bool {
	if err != nil {
		log.Errorf("action: receive_bet | result: fail | error: %v",
			err,
		)
		return true
	}
	if response != protocol.CODE_SUCCESS {
		log.Errorf("action: receive_bet | result: fail | error: %v",
			err,
		)
		return true
	}
	return false
}

func (c *Client) StartClientLoop() {
	file := os.Getenv("FILE")
	agencyNumber := strings.TrimPrefix(file, "/agency-")
	agencyNumber = strings.TrimSuffix(agencyNumber, ".csv")
	number, err := strconv.Atoi(agencyNumber)
	if err != nil {
		log.Errorf("action: convert_agency_number | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		number = 1
	}
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
	count := 0
	bets := []protocol.Bet{}
	for scanner.Scan() {
		line := scanner.Text()
		bet, err := createBetFromLine(line, number)
		if err != nil {
			log.Errorf("action: create_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	bets = append(bets, bet)
	count++
	if count >= c.config.Batch {
			c.createClientSocket()
			err = protocol.SendBets(c.conn, bets)
			if err != nil {
				log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
				return
			}
			response, err := protocol.ReceiveBet(c.conn)
			if checkForErrorsBet(response, err) {
				return
			}
			count = 0
			bets = []protocol.Bet{}
			c.conn.Close()
		}
	}
	if count > 0 {
		c.createClientSocket()
		err = protocol.SendBets(c.conn, bets)
		if err != nil {
			log.Errorf("action: send_bets | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		response, err := protocol.ReceiveBet(c.conn)
		if checkForErrorsBet(response, err) {
			return
		}
		c.conn.Close()
	}
}

func (c *Client) StopClientLoop() {
	c.conn.Close()
}
