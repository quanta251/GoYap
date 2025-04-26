package client

import (
	"fmt"
	"github.com/quanta251/GoYap/helpers"
	"net"
	"log"
	"encoding/json"
)

func SomethingToExport(message string) {
	// This is a placeholder for some helper functions that I will add later
}

// Send the command to the server to list users
func RequestUsers(conn net.Conn) error {
	input := "listusers"
	helpers.SendMessage(conn, input)
	fmt.Println("Requesting a list of connected users...")
	response, err := helpers.ReceiveMessage(conn)
	if err != nil {
		log.Println("Could not receive the list of users:", err)
		return err
	}

	// Parse the JSON response
	fmt.Println("Parsing the json")
	var users []string
	err = json.Unmarshal([]byte(response), &users)
	if err != nil {
		log.Println("Error parsing the JSON user list.", err)
		return err
	}

	fmt.Println("Currently Connected Users:")
	for _, user := range users {
		fmt.Println("- " + user)
	}

	return nil
}
