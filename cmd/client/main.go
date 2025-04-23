package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"github.com/quanta251/GoYap/internal"
)

func getUsername() (string, error) {
	fmt.Println("-------------- Please Input Your Name Below --------------")
	reader := bufio.NewReader(os.Stdin)
		
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Could not get username from the user...", err)
		return username, err // Does not matter if we send junk as the username. The error will be caught and handled.
	}

	username = strings.TrimSpace(username)

	fmt.Println("You will be connected shortly.")

	return username, nil
}

func checkInput(input string) bool {
	if input == "exit" || input == "quit" {
		return true
	}

	return false
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

	err = helpers.SendMessage(conn, clientName) 	// The first message will be interpreted as defining who we are...
	response, err := helpers.ReceiveMessage(conn)	// Get the server's response to the username submission

	fmt.Printf("Succesfully connected to server '%s'\n\n", serverAddress)
	fmt.Println(string(response))
	fmt.Println("------------------------------------------------------------")

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
			helpers.SendMessage(conn, input)
			return
		}

		err = helpers.SendMessage(conn, input)
		if err != nil {
			continue
		}
		fmt.Printf("Message sent...\n\n")
	}
}
