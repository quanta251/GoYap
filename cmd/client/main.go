package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

// Define the prefix length
const prefixLength = 4

func checkInput(input string) bool {
	if input == "exit" || input == "quit" {
		return true
	}

	return false
}

// Take the message and send it over the connection
func sendMessage(conn net.Conn, message string) error {
	// Cast the message as a byte slice
	payload := []byte(message)

	// Prepare the prefix
	payloadLength := uint32(len(payload))
	prefix := make([]byte, prefixLength)
	binary.BigEndian.PutUint32(prefix, payloadLength)

	// Send the prefix
	_, err := conn.Write(prefix)
	if err != nil {
		fmt.Println("Could not send prefix...")
		return err
	}

	// Send the message
	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println("Could not send message...")
		return err
	}

	return nil
}

func main() {
	serverAddress := "localhost:3000"
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		panic(err)
	}
	defer conn.Close()	


	fmt.Printf("Succesfully connected to server '%s'\n\n", serverAddress)
	fmt.Println("Please input your messages below")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Could not read your message...")
			continue
		}

		// Trim off the white space
		input = strings.TrimSuffix(input, "\n")

		fmt.Printf("The users input is: %b\n", []byte(input))
		fmt.Printf("The byte version of 'quit' is: %b\n", []byte("quit"))
		fmt.Printf("The byte version of 'exit' is: %b\n", []byte("exit"))

		// Check input for quit commands...
		if checkInput(input) {
			fmt.Println("Quitting the client now...")
			return
		}

		err = sendMessage(conn, input)
		if err != nil {
			continue
		}
		fmt.Printf("Message sent...\n\n")
	}
}
