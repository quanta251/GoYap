package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Define the prefix length
const prefixLength = 4

func getUsername() (string, error) {
	fmt.Println("-------------- Please Input Your Name Below --------------")
	reader := bufio.NewReader(os.Stdin)
		
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Could not get username from the user...", err)
		return username, err // Does not matter if we send junk as the username. The error will be caught and handled.
	}

	username = strings.TrimSpace(username)

	fmt.Printf("Welcome, %s. You will be connected shortly...\n", username)

	return username, nil
}

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

func main() {
	clientName, err := getUsername()
	if err != nil {
		log.Fatalln("Could not get username from client...", err)
	}

	serverAddress := "localhost:3000"
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		panic(err)
	}
	defer conn.Close()	
	err = sendMessage(conn, clientName) // The first message will be interpreted as defining who we are...


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

		// Check input for quit commands...
		if checkInput(input) {
			fmt.Println("Quitting the client now...")
			sendMessage(conn, input)
			return
		}

		err = sendMessage(conn, input)
		if err != nil {
			continue
		}
		fmt.Printf("Message sent...\n\n")
	}
}
