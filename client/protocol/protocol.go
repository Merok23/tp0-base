package protocol

import (
	"fmt"
	"encoding/binary" // <- PREGUNTAR: Esto se puede usar?
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

func htonl(value int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func SendMsg(conn net.Conn,  id int, msgID int) error {
	// Send the code
	codeBytes := htonl(1)
	conn.Write(codeBytes)
	log.Infof("action: send_code | result: success | client_id: %v | code: 1", id)
	// Send the size
	msg := fmt.Sprintf("[CLIENT %v] Message N°%v", id, msgID)
	size := len(msg)
	sizeBytes := htonl(size)
	conn.Write(sizeBytes)
	log.Infof("action: send_size | result: success | client_id: %v | size: %v", id, size)

	msgBytes := []byte(msg)
	// Send the message
	conn.Write(msgBytes)
	log.Infof("action: send_message | result: success | client_id: %v | message: %s", id, msg)
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
	log.Infof("action: receive_code | result: success | client_id: %v | code: %v", id, code)
	// Read the size
	sizeBytes := make([]byte, 4)
	conn.Read(sizeBytes)
	size := ntohl(sizeBytes)
	log.Infof("action: receive_size | result: success | client_id: %v | size: %v", id, size)
	// Read the message
	msgBytes := make([]byte, size)
	conn.Read(msgBytes)
	msg := string(msgBytes)
	return msg, nil
}
