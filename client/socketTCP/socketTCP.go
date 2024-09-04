package socketTCP

import (
	"net"
)

func WriteAll(conn net.Conn, data []byte, size int) error {
	written := 0
	for written < size {
		n, err := conn.Write(data[written:size])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}

func ReadAll(conn net.Conn, data []byte, size int) error {
	read := 0
	for read < size {
		n, err := conn.Read(data[read:size])
		if err != nil {
			return err
		}
		read += n
	}
	return nil
}
