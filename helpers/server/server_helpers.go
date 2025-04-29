package server

import (
	"encoding/json"
	"net"
	"log"
	"fmt"
	"github.com/quanta251/GoYap/helpers"
)

func ParseMessage(message string) {

}

func SendUsernameList(conn net.Conn, usernameMap map[string]net.Conn) error {
	// Send the list of usernames to the connection
	var usernames []string
	for username := range usernameMap {
		usernames = append(usernames, username)
	}

	response, err := json.Marshal(usernames)
	if err != nil {
		log.Println("Could not encode the usernames to JSON.")
		return err
	}

	fmt.Println("the length of the json message is", len(response))
	err = helpers.SendMessage(conn, string(response))
	if err != nil {
		log.Println("Could not send the stringified JSON usernames list.")
		return err
	}
	fmt.Println("Successfully sent the list of usernames to a user.")
	return nil
}
