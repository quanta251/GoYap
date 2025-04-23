package helpers

import (
	"net"
)

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
