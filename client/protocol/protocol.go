package protocol

import (
	"fmt"
	"encoding/binary" // <- PREGUNTAR: Esto se puede usar?
	"net"
)


func htonl(value int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func SendMsg(conn net.Conn,  id int, msgID int) error {
	// Send the code
	codeBytes := htonl(1)
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

func ntohl(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func ReceiveMsg(conn net.Conn, id int) (string, error) {
	// Read the code
	codeBytes := make([]byte, 4)
	conn.Read(codeBytes)
	code := ntohl(codeBytes)
	if code != 1 {
		return "", fmt.Errorf("Invalid code received: %v", code)
	}
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
