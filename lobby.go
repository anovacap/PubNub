package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pubnub "github.com/pubnub/go"
)

func userInput(lobby string, username string) (string, string) {
	var (
		newlobby    string
		newUsername string
		err         error
	)
	reader := bufio.NewReader(os.Stdin)
	if lobby == "" {
		fmt.Print("Enter Lobby Name: ")
	} else {
		fmt.Print("Enter Lobby Name (" + lobby + "): ")
	}
	newlobby, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	newlobby = strings.Replace(newlobby, "\n", "", -1) // Convert CRLF to LF
	if newlobby != "" {
		lobby = newlobby //Use last lobby name when player doesn't provide
	}
	if username == "" { //Ask for username
		fmt.Print("Enter Your Name: ")
	} else {
		fmt.Print("Enter Your Name (" + username + "): ")
	}
	newUsername, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	newUsername = strings.Replace(newUsername, "\n", "", -1) // Convert CRLF to LF
	if newUsername != "" {
		username = newUsername
	}
	if (lobby == "") || (username == "") { //The player must have a lobby and a username
		fmt.Println("You Must Provide a Lobby and Name! ")
		userInput(lobby, username) // Start over
	}
	return lobby, username
}

func hereNow(channel string, pn *pubnub.PubNub) int {
	// Return count of occupants for a channel
	res, _, err := pn.HereNow().
		Channels([]string{channel}).
		Execute()
	if err != nil {
		panic(err)
	}
	return res.TotalOccupancy
}

func newLobby(lobby string, username string, pn *pubnub.PubNub) {
	var (
		guestName string
		hostName  string
		isHost    bool
	)
	data := make(map[string]interface{})
	lobby, username = userInput(lobby, username)
	lobbylistener := pubnub.NewListener()
	endLobby := make(chan bool)
	go func() {
		for {
			select {
			case status := <-lobbylistener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					game_occupants := hereNow(lobby, pn)
					lobby_occupants := hereNow(lobby+"_lobby", pn)
					if game_occupants > 0 || lobby_occupants > 2 {
						fmt.Println("Game already in progress! Please try another lobby.")
						fmt.Print("Game: ")
						fmt.Println(game_occupants)
						fmt.Print("Lobby: ")
						fmt.Println(lobby_occupants)
						pn.RemoveListener(lobbylistener)
						pn.Unsubscribe().
							Channels([]string{lobby + "_lobby"}).
							Execute()
						newLobby(lobby, username, pn)
						return
					}
					if lobby_occupants == 0 {
						isHost = true
						fmt.Println("Waiting for guest...")
					} else if lobby_occupants == 1 {
						fmt.Println("Waiting for host...")
						data["guestName"] = username
						guestName = username
						pn.Publish().
							Channel(lobby + "_lobby").
							Message(data).
							Execute()
					}
				}
			case message := <-lobbylistener.Message:
				if msg, ok := message.Message.(map[string]interface{}); ok {
					if !isHost {
						if val, ok := msg["hostName"]; ok { // When the guest receives the host username then the game is ready to start.
							hostName = val.(string)
							endLobby <- true
							return
						}
					} else {
						if val, ok := msg["guestName"]; ok { // Receives the guest username then the host sends the host username and starts a game.
							guestName = val.(string)
							data["hostName"] = username
							hostName = username
							pn.Publish().
								Channel(lobby + "_lobby").
								Message(data).
								Execute()
							endLobby <- true
							return
						}
					}
				}
			}
		}
	}()
	pn.AddListener(lobbylistener)
	pn.Subscribe().
		Channels([]string{lobby + " lobby"}).
		Execute()
	<-endLobby
	pn.RemoveListener(lobbylistener)
	pn.Unsubscribe().
		Channels([]string{lobby + "_lobby"}).
		Execute()
	startGame(isHost, lobby, strings.Title(hostName),
		strings.Title(guestName), pn)
}
