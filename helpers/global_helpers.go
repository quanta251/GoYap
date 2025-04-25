package helpers

import (
	"net"
	"encoding/binary"
	"log"
)

const prefixLength = 4

// Safely read messages from a connection without the need of a predefined
// buffer length
func readN(conn net.Conn, buf []byte) error {
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
		log.Println("Could not send message (prefix)...")
		return err
	}

	// Send the message
	_, err = conn.Write(payload)
	if err != nil {
		log.Println("Could not send message (body)...")
		return err
	}

	return nil
}

// Listen to a message and decode it.
func ReceiveMessage(conn net.Conn) (string, error) {
	prefix := make([]byte, prefixLength)
	err := readN(conn, prefix)
	if err != nil {
		log.Println("Could not receive prefix from client:", err)
		return "dummy", err
	}

	var messageLength uint32 = binary.BigEndian.Uint32(prefix)

	payload := make([]byte, messageLength)

	err = readN(conn, payload)
	if err != nil {
		log.Println("Could not receive payload from client:", err)
		return "dummy", err
	}

	message := string(payload)

	return message, nil
}
