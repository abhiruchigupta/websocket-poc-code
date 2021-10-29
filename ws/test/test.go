package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Game struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Platform    []string `json:"platform"`
	Message string `json:"message"`

}

type InfoWsMessage struct {
	Message string `json:"message"`
	UserID  string `json:"userID"`
	StoreID int64 `json:"storeID"`
	SenderID string `json:"senderID"`
}



func main() {

	myGame := Game{
		Title:       "Fifa 19",
		Description: "Football simulation game, based on UEFA players",
		Platform:    []string{"PS4"},
		Message: "hello",
	}


	a := InfoWsMessage{
		Message: "message.GetMessage()",
		UserID: "message.GetUserID()",
		StoreID: 0,
		SenderID: "message.GetSenderID()",
	}

	out, err2 := json.MarshalIndent(a, "", "    ")
	if err2 != nil {
		log.Fatal("Failed to generate json", err2)
	}
	fmt.Printf("%s\n", string(out))


	// For comparison, the usual way would be: j, err := json.Marshal(myGame)

	// MarshalIndent accepts:
	// 1) the data
	// 2) a prefix to place on all lines but 1st
	// 3) an indent to place before lines based on indent level

	// Print Json with indents, the pretty way:
	prettyJSON, err := json.MarshalIndent(myGame, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
