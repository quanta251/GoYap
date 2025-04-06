package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
    "os/signal"
    "syscall"
    "os"
)

const bufferSize = 2048
const prefixlength = 4

type Message struct {
    from            string
    payload         []byte
}

type Server struct {
    listenAddr      string
    ln              net.Listener
    quitch          chan struct{}
    msgch           chan Message
    activeConns     int // Track the number of connections here
    shutdownTimeout time.Duration
}

func NewServer(listenAddr string) *Server {
    return &Server {
        listenAddr:         listenAddr,
        quitch:             make(chan struct{}),
        msgch:              make(chan Message, 10),
        shutdownTimeout:    5 * time.Second,
    }
}

func (s *Server) waitForShutdown() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs // Block here and wait for the signal

    fmt.Println("Received Signal Shutdown")

    close(s.quitch) // Close the blocking channel so the start function can return

    // Wait for all active connections to finish (with timeout)
    timeout := time.After(s.shutdownTimeout)
    select {
    case <-timeout:
        fmt.Println("Shutdown timeout reached!")
    }

    // Stop accepting connections and wait for in-progress messages to finish
    fmt.Println("All connections processed. Server shutting down...")
}

func (s *Server) Start() error {
    ln, err := net.Listen("tcp", s.listenAddr)
    if err != nil {
        fmt.Println("Could not start Listener...")
        return err
    }
    s.ln = ln
    fmt.Println("Server is listening on", s.listenAddr)

    // Start the accept loop
    go s.acceptLoop()

    // Setup graceful shutdown: listen for system signals (e.g., SIGINT)
    go s.waitForShutdown()

    <- s.quitch // Block with this channel

    // Close the listener
    s.ln.Close()
    close(s.msgch)

    return nil
}



func (s *Server) acceptLoop() {
    for {
        conn, err := s.ln.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
            return
        }
        defer conn.Close()

        fmt.Println("New connection from:", conn.RemoteAddr())
        s.activeConns++

        go s.readLoop(conn)
    }
}

func (s *Server) readLoop(conn net.Conn) {
    defer func() {
        s.activeConns-- // Decrement the number of activate connections
    }()

    prefix := make([]byte, prefixlength)

    for {
        // Try to get the prefix length
        err := readN(conn, prefix)
        if err != nil {
            fmt.Println("Read error (prefix):", err)
            return // Breaking the loop here. 
        }

        // Convert the prefix to a length
        var messageLength uint32 = binary.BigEndian.Uint32(prefix)

        // Establish the message buffer which will hold the payload
        messageBuf := make([]byte, messageLength)
        err = readN(conn, messageBuf)
        if err != nil {
            fmt.Println("Read error (body):", err)
            return
        }
        
        s.msgch <- Message{
            from:       conn.RemoteAddr().String(),
            payload:    messageBuf,
        }

        conn.Write([]byte("Thank you for your message.\n"))

        defer func() {
            fmt.Printf("Connection to %s has been closed...\n", conn.RemoteAddr().String())
        }()

        return
    }
}

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

func main() {
    server := NewServer(":3000")
    
    go func(){
        for msg := range server.msgch {
            fmt.Printf("Received message from connection(%s):%s\n", msg.from, msg.payload)
        }
    }()

    log.Fatal(server.Start())
}
