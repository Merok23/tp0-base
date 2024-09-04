package protocol

import (
	"fmt"
	"encoding/binary"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/socketTCP"
)


type Bet struct {
	Dni int
	Name string
	Lastname string
	DateOfBirth string
	Number int
	AgencyNumber int
}

const (
	CODE_ECHO    = 1
	CODE_BET     = 2
	CODE_END	 = 3
	CODE_SUCCESS = 200
	CODE_ERROR   = 400
	SIZE_INT_32  = 4
	SIZE_DATE	= 10
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
	err := socketTCP.WriteAll(conn, codeBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	// Send the size
	msg := fmt.Sprintf("[CLIENT %v] Message N°%v", id, msgID)
	size := len(msg)
	sizeBytes := htonl(size)
	err = socketTCP.WriteAll(conn, sizeBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	msgBytes := []byte(msg)
	// Send the message
	err = socketTCP.WriteAll(conn, msgBytes, size)
	if err != nil {
		return err
	}
	return nil
}

func SendBet(conn net.Conn, bet Bet) error {
	// possible performance improvement (not implemented for time's sake): send all the data with:
	// Size of whole message (4 bytes) + SizeDNI (4 bytes).. etc (fixed size header) (This can be sent with the code)
	// And then we send the data in a single write, the server uses the sizes to read the data (payload)
	dniBytes := htonl(bet.Dni)
	err := socketTCP.WriteAll(conn, dniBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	numberBytes := htonl(bet.Number)
	err = socketTCP.WriteAll(conn, numberBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	dateOfBirthBytes := []byte(bet.DateOfBirth)
	err = socketTCP.WriteAll(conn, dateOfBirthBytes, SIZE_DATE) // always 10 bytes YYYY-MM-DD
	if err != nil {
		return err
	}
	agencyNumberBytes := htonl(bet.AgencyNumber)
	err = socketTCP.WriteAll(conn, agencyNumberBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	nameBytes := []byte(bet.Name)
	nameBytesSize := htonl(len(nameBytes))
	err = socketTCP.WriteAll(conn, nameBytesSize, SIZE_INT_32)
	if err != nil {
		return err
	}
	err = socketTCP.WriteAll(conn, nameBytes, len(nameBytes))
	if err != nil {
		return err
	}
	lastnameBytes := []byte(bet.Lastname)
	lastnameBytesSize := htonl(len(lastnameBytes))
	err = socketTCP.WriteAll(conn, lastnameBytesSize, SIZE_INT_32)
	if err != nil {
		return err
	}
	err = socketTCP.WriteAll(conn, lastnameBytes, len(lastnameBytes))
	if err != nil {
		return err
	}
	return nil
}	

func SendBets(conn net.Conn, bets []Bet) error {
	codeBytes := htonl(CODE_BET)
	err := socketTCP.WriteAll(conn, codeBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	countBytes := htonl(len(bets))
	err = socketTCP.WriteAll(conn, countBytes, SIZE_INT_32)
	if err != nil {
		return err
	}
	for _, bet := range bets {
		err := SendBet(conn, bet)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendEnd(conn net.Conn, agencyNumber int) (int, []int, int, error) {
	codeBytes := htonl(CODE_END)
	err := socketTCP.WriteAll(conn, codeBytes, SIZE_INT_32)
	if err != nil {
		return 0, nil, 0, err
	}
	agencyNumberBytes := htonl(agencyNumber)
	err = socketTCP.WriteAll(conn, agencyNumberBytes, SIZE_INT_32)
	if err != nil {
		return 0, nil, 0, err
	}
	// Receive the result
	resultCodeBytes := make([]byte, SIZE_INT_32)
	err = socketTCP.ReadAll(conn, resultCodeBytes, SIZE_INT_32)
	resultCode := ntohl(resultCodeBytes)
	if resultCode != CODE_SUCCESS {
		return int(resultCode), nil, 0, fmt.Errorf(
			"Error sending end, error: %v",
			err,
		)
	}
	winnersBytes := make([]byte, SIZE_INT_32)
	err = socketTCP.ReadAll(conn, winnersBytes, SIZE_INT_32)
	if err != nil {
		return 0, nil, 0, err
	}
	winners := ntohl(winnersBytes)
	dnisLenBytes := make([]byte, SIZE_INT_32)
	err = socketTCP.ReadAll(conn, dnisLenBytes, SIZE_INT_32)
	if err != nil {
		return 0, nil, 0, err
	}
	dnisLen := ntohl(dnisLenBytes)
	dniWinners := make([]int, dnisLen)
	for i := 0; i < int(dnisLen); i++ {
		dniBytes := make([]byte, SIZE_INT_32)
		err = socketTCP.ReadAll(conn, dniBytes, SIZE_INT_32)
		if err != nil {
			return 0, nil, 0, err
		}
		dni := ntohl(dniBytes)
		dniWinners[i] = int(dni)
	}
	return int(winners), dniWinners, int(dnisLen), nil
}

func ReceiveBet(conn net.Conn) (uint32, error) {
	resultCodeBytes := make([]byte, SIZE_INT_32)
	err := socketTCP.ReadAll(conn, resultCodeBytes, SIZE_INT_32)
	if err != nil {
		return 0, err
	}
	resultCode := ntohl(resultCodeBytes)
	return resultCode, nil
}

func ReceiveEchoMsg(conn net.Conn, id int) (string, error) {
	// Read the size
	sizeBytes := make([]byte, SIZE_INT_32)
	err := socketTCP.ReadAll(conn, sizeBytes, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	size := ntohl(sizeBytes)
	// Read the message
	msgBytes := make([]byte, size)
	err = socketTCP.ReadAll(conn, msgBytes, int(size))
	if err != nil {
		return "", err
	}
	msg := string(msgBytes)
	return msg, nil
}
