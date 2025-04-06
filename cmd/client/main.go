package main

import (
    "fmt"
    "bufio"
    "encoding/binary"
    "net"
    "os"
    "strings"
)

const prefixLength = 4

func main() {

    conn, err := net.Dial("tcp", "sabretooth:3000")
    if err != nil {
        panic(err)
    }

    defer conn.Close()


    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("Client -- Enter Message to Server: ")
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Input error:", err)
            return
        }

        input = strings.TrimSpace(input)

        if input == "" {
            continue
        }

        payload := []byte(input)
        length := uint32(len(payload))

        prefix := make([]byte, prefixLength)
        binary.BigEndian.PutUint32(prefix, length)

        _, err = conn.Write(prefix)
        if err != nil {
            fmt.Println("Error sending prefix:", err)
            return
        }
        
        _, err = conn.Write(payload)
        if err != nil {
            fmt.Println("Error sending payload", err)
            return
        }


        response := make([]byte, 1024)
        n, err := conn.Read(response)
        if err != nil {
            fmt.Println("Server Closed Connection.")
            return
        }

        fmt.Printf("Server -- %s\n", response[:n])
        if input == "QUIT" {
            fmt.Println("Exiting Client.")
            return
        }
    }
}
