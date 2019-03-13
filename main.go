package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

func main() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "sub-c-6548f134-1950-11e9-9cda-0ee81137d4bc"
	config.PublishKey = "pub-c-56dc14e3-fc0c-44f6-8fe0-9c0989edbebb"
	pn := pubnub.NewPubNub(config)
	fmt.Println("Welcome to Space Race!")
	newLobby("", "", pn) //This creates the lobby for a new game
}
