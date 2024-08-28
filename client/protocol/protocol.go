package protocol

import (
	"fmt"
	"encoding/binary" // <- PREGUNTAR: Esto se puede usar?
	"net"
)

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
	conn.Write(codeBytes)
	dniBytes := htonl(dni)
	conn.Write(dniBytes)
	numberBytes := htonl(number)
	conn.Write(numberBytes)
	dateOfBirthBytes := []byte(dateOfBirth)
	conn.Write(dateOfBirthBytes) // always size 10 (YYYY-MM-DD)

	nameBytes := []byte(name)
	nameBytesSize := htonl(len(nameBytes))
	conn.Write(nameBytesSize)
	conn.Write(nameBytes)

	lastnameBytes := []byte(lastname)
	lastnameBytesSize := htonl(len(lastnameBytes))
	conn.Write(lastnameBytesSize)
	conn.Write(lastnameBytes)
	return "", nil
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
