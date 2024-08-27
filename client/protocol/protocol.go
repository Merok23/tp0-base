package protocol

import (
	"fmt"
	"encoding/binary" // <- PREGUNTAR: Esto se puede usar?
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

func SendMsg(conn net.Conn,  id int, msgID int) error {
	// Send the size (4 bytes)
	msg := fmt.Sprintf("[CLIENT %v] Message N°%v\n", id, msgID)
	size := len(msg)
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))
	conn.Write(sizeBytes)
	log.Infof("action: send_size | result: success | client_id: %v | size: %v", id, size)

	msgBytes := []byte(msg)
	// Send the message
	conn.Write(msgBytes)
	log.Infof("action: send_message | result: success | client_id: %v | message: %s", id, msg)
	return nil
}