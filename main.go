package main

import (
	"fmt"
	"log"
	"net"
)

const bufferSize = 9

type Message struct {
    from            string
    payload         []byte
}

type Server struct {
    listenAddr      string
    ln              net.Listener
    quitch          chan struct{}
    msgch           chan Message
}

func NewServer(listenAddr string) *Server {
    return &Server {
        listenAddr:     listenAddr,
        quitch:         make(chan struct{}),
        msgch:          make(chan Message, 10),
    }
}

func (s *Server) Start() error {
    ln, err := net.Listen("tcp", s.listenAddr)
    if err != nil {
        fmt.Println("Could not start Listener...")
        return err
    }
    defer ln.Close()
    s.ln = ln

    go s.acceptLoop()

    <- s.quitch // Start blocking with this channel
    close(s.msgch)

    return nil
}

func (s *Server) acceptLoop() {
    for {
        conn, err := s.ln.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
            continue
        }

        fmt.Println("New connection from:", conn.RemoteAddr())

        go s.readLoop(conn)
    }
}

func (s *Server) readLoop(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, bufferSize)

    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("read error:", err)
            continue
        }
        s.msgch <- Message{
            from: conn.RemoteAddr().String(),
            payload: buf[:n],
        }

        conn.Write([]byte("Thank you for your message!\n"))
    }
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
