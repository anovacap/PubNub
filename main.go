package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	pn := pubnub.NewPubNub(config)
	fmt.Println("Welcome to Space Race!")
	newLobby("", "", pn) //This creates the lobby for a new game
}
