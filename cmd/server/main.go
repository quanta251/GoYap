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
	"github.com/quanta251/GoYap/internal"
)

const bufferSize = 2048
const prefixlength = 4

type Message struct {
    from            string
    payload         []byte
}

var clients = make(map[string]net.Conn) // This is the variable that will store username -> ip address info

type Server struct {
    listenAddr      string
    ln              net.Listener
    quitch          chan struct{}
    msgch           chan Message
    activeConns     int // Track the number of connections here
    shutdownTimeout time.Duration
}

func receiveMessage(conn net.Conn) ([]byte, error) {
	prefix := make([]byte, prefixlength)
	err := helpers.ReadN(conn, prefix)
	if err != nil {
		log.Println("Could not receive prefix from client:", err)
		return nil, err
	}

	var messageLength uint32 = binary.BigEndian.Uint32(prefix)

	payload := make([]byte, messageLength)

	err = helpers.ReadN(conn, payload)
	if err != nil {
		log.Println("Could not receive payload from client:", err)
		return nil, err
	}

	return payload, nil
}

func receiveUsername(addressBook map[string]net.Conn, conn net.Conn) error {
	usernameBytes, err := receiveMessage(conn)
	if err != nil {
		log.Printf("Could not get the username from client '%s': %v\n", conn.RemoteAddr(), err)
		return err
	}

	usernameString := string(usernameBytes)
	addressBook[usernameString] = conn

	return nil
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
		fmt.Println("Waiting for a new connection...")
        conn, err := s.ln.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
            return
        }

        fmt.Println("New connection from:", conn.RemoteAddr())
        s.activeConns++

		// Receive the username here and add it to the map containing other
		// users
		err = receiveUsername(clients, conn)
		if err != nil {
			log.Printf("Could not get the username from client '%s'\n", conn.RemoteAddr())
			conn.Close()			// Close the current connection as it will not be used
			s.activeConns-- 		// Decrement the number of connections
			continue				// Continue to the next connection acceptance
		}

        go s.readLoop(conn)
		fmt.Println("Started a new read loop")
    }
}

func (s *Server) readLoop(conn net.Conn) {
    defer func() {
		conn.Close()
        s.activeConns-- // Decrement the number of activate connections
    }()

    prefix := make([]byte, prefixlength)

    for {
        // Try to get the prefix length
        err := helpers.ReadN(conn, prefix)
        if err != nil {
            fmt.Println("Read error (prefix):", err)
			continue // Continue to the next message without exiting the server
        }

        // Convert the prefix to a length
        var messageLength uint32 = binary.BigEndian.Uint32(prefix)

        // Establish the message buffer which will hold the payload
        messageBuf := make([]byte, messageLength)
        err = helpers.ReadN(conn, messageBuf)
        if err != nil {
            fmt.Println("Read error (body):", err)
			continue
        }

        message := string(messageBuf)

        if message == "quit" || message == "exit" {
            fmt.Printf("Client %s requested to close the connection...\n", conn.RemoteAddr())
            conn.Write([]byte("Goodbye.\n"))
			return
        }

        
        // Otherwise, handle the message normally
        s.msgch <- Message{
            from:       conn.RemoteAddr().String(),
            payload:    messageBuf,
        }

        conn.Write([]byte("Thank you for your message.\n"))

    }
}


func main() {
    server := NewServer(":3000")
    
    go func(){
        for msg := range server.msgch {
            fmt.Printf("Received message from connection(%s):%s\n", msg.from, msg.payload)
        }
    }()

    log.Println(server.Start())
	fmt.Println("The active connections we have are:", clients)
}
