package protocol

import (
	"fmt"
	"encoding/binary"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/socketTCP"
)

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

func SendBet(
	conn net.Conn,
	id int,
	dni int,
	name string,
	lastname string,
	dateOfBirth string,
	number int,
) (string, error) {
	// possible performance improvement: send all the data with:
	// CODE
	// size of whole message (4 bytes) SeparatorDNI (4 bytes) ... 
	// Send the code
	codeBytes := htonl(CODE_BET)
	err := socketTCP.WriteAll(conn, codeBytes, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	dniBytes := htonl(dni)
	err = socketTCP.WriteAll(conn, dniBytes, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	numberBytes := htonl(number)
	err = socketTCP.WriteAll(conn, numberBytes, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	dateOfBirthBytes := []byte(dateOfBirth)
	err = socketTCP.WriteAll(conn, dateOfBirthBytes, SIZE_DATE)
	if err != nil {
		return "", err
	}
	nameBytes := []byte(name)
	nameBytesSize := htonl(len(nameBytes))
	err = socketTCP.WriteAll(conn, nameBytesSize, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	err = socketTCP.WriteAll(conn, nameBytes, len(nameBytes))
	if err != nil {
		return "", err
	}
	lastnameBytes := []byte(lastname)
	lastnameBytesSize := htonl(len(lastnameBytes))
	err = socketTCP.WriteAll(conn, lastnameBytesSize, SIZE_INT_32)
	if err != nil {
		return "", err
	}
	err = socketTCP.WriteAll(conn, lastnameBytes, len(lastnameBytes))
	if err != nil {
		return "", err
	}
	return "", nil
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

