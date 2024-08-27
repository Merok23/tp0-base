package protocol

import (
	"fmt"
	"encoding/binary" // <- PREGUNTAR: Esto se puede usar?
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

func encodeInt32(value int) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	return bytes
}

func SendMsg(conn net.Conn,  id int, msgID int) error {
	// Send the code
	codeBytes := encodeInt32(1)
	conn.Write(codeBytes)
	log.Infof("action: send_code | result: success | client_id: %v | code: 1", id)
	// Send the size
	msg := fmt.Sprintf("[CLIENT %v] Message N°%v\n", id, msgID)
	size := len(msg)
	sizeBytes := encodeInt32(size)
	conn.Write(sizeBytes)
	log.Infof("action: send_size | result: success | client_id: %v | size: %v", id, size)

	msgBytes := []byte(msg)
	// Send the message
	conn.Write(msgBytes)
	log.Infof("action: send_message | result: success | client_id: %v | message: %s", id, msg)
	return nil
}