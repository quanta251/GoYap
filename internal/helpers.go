package helpers

import (
	"net"
	"fmt"
	"encoding/binary"
	"log"
)


// Safely read messages from a connection without the need of a predefined
// buffer length
func ReadN(conn net.Conn, buf []byte) error {
    total := 0
    for total < len(buf) {
        n, err := conn.Read(buf[total:])
        if err != nil {
            return err
        }
        total += n
    }
    return nil
}

const prefixLength = 4
// Send a message either to server from client or to client from server
func SendMessage(conn net.Conn, message string) error {
	// Cast the message as a byte slice
	payload := []byte(message)

	// Prepare the prefix
	payloadLength := uint32(len(payload))
	prefix := make([]byte, prefixLength)
	binary.BigEndian.PutUint32(prefix, payloadLength)

	// Send the prefix
	_, err := conn.Write(prefix)
	if err != nil {
		fmt.Println("Could not send message (prefix)...")
		return err
	}

	// Send the message
	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println("Could not send message (body)...")
		return err
	}

	return nil
}

func ReceiveMessage(conn net.Conn) ([]byte, error) {
	prefix := make([]byte, prefixLength)
	err := ReadN(conn, prefix)
	if err != nil {
		log.Println("Could not receive prefix from client:", err)
		return nil, err
	}

	var messageLength uint32 = binary.BigEndian.Uint32(prefix)

	payload := make([]byte, messageLength)

	err = ReadN(conn, payload)
	if err != nil {
		log.Println("Could not receive payload from client:", err)
		return nil, err
	}

	return payload, nil
}
