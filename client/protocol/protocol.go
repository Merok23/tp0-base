package protocol

import (
	"fmt"
	"encoding/binary"
	"net"
)


type Bet struct {
	Dni int
	Name string
	Lastname string
	DateOfBirth string
	Number int
}

const (
	CODE_ECHO    = 1
	CODE_BET     = 2
	CODE_SUCCESS = 200
	CODE_ERROR   = 400
)

func htonl(value int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	return bytes
}


func ntohl(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func SendEchoMsg(conn net.Conn,  id int, msgID int) error {
	// Send the code
	codeBytes := htonl(CODE_ECHO)
	conn.Write(codeBytes)
	// Send the size
	msg := fmt.Sprintf("[CLIENT %v] Message N°%v", id, msgID)
	size := len(msg)
	sizeBytes := htonl(size)
	conn.Write(sizeBytes)

	msgBytes := []byte(msg)
	// Send the message
	conn.Write(msgBytes)
	return nil
}

func SendBet(conn net.Conn, bet Bet) error {
	// possible performance improvement (not implemented for time's sake): send all the data with:
	// Size of whole message (4 bytes) + SizeDNI (4 bytes).. etc (fixed size header) (This can be sent with the code)
	// And then we send the data in a single write, the server uses the sizes to read the data (payload)
	dniBytes := htonl(bet.Dni)
	conn.Write(dniBytes)
	numberBytes := htonl(bet.Number)
	conn.Write(numberBytes)
	dateOfBirthBytes := []byte(bet.DateOfBirth)
	conn.Write(dateOfBirthBytes) // always size 10 (YYYY-MM-DD)
	nameBytes := []byte(bet.Name)
	nameBytesSize := htonl(len(nameBytes))
	conn.Write(nameBytesSize)
	conn.Write(nameBytes)
	lastnameBytes := []byte(bet.Lastname)
	lastnameBytesSize := htonl(len(lastnameBytes))
	conn.Write(lastnameBytesSize)
	conn.Write(lastnameBytes)
	return nil
}	

func SendBets(conn net.Conn, bets []Bet) error {
	codeBytes := htonl(CODE_BET)
	conn.Write(codeBytes)
	countBytes := htonl(len(bets))
	conn.Write(countBytes)
	for _, bet := range bets {
		err := SendBet(conn, bet)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReceiveBet(conn net.Conn) (uint32, error) {
	resultCodeBytes := make([]byte, 4)
	conn.Read(resultCodeBytes)
	resultCode := ntohl(resultCodeBytes)
	return resultCode, nil
}

func ReceiveEchoMsg(conn net.Conn, id int) (string, error) {
	// Read the size
	sizeBytes := make([]byte, 4)
	conn.Read(sizeBytes)
	size := ntohl(sizeBytes)
	// Read the message
	msgBytes := make([]byte, size)
	conn.Read(msgBytes)
	msg := string(msgBytes)
	return msg, nil
}
