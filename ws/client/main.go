package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type InfoWsMessage struct {
	Message string `json:"message"`
	UserID  string `json:"userID"`
	StoreID int64 `json:"storeID"`
	SenderID string `json:"senderID"`
	Type string `json:"messageType"`
}

func (i InfoWsMessage) GetUserID() string {
	return i.UserID
}

func (i InfoWsMessage) GetMessage() string {
	return i.Message
}

func (i InfoWsMessage) GetStoreID() int64 {
	return i.StoreID
}

func (i InfoWsMessage) GetSenderID() string {
	return i.SenderID
}

func (i InfoWsMessage) GetMessageType() string {
	return i.Type
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var userID = flag.String("userid", "foo", "userID for the subscribing client")


func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	query := fmt.Sprintf("userId=%s", *userID)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/subscribe", RawQuery: query}
	log.Printf("connecting to %s", u.String())

	headers := http.Header{
		"X-Compass-WS-User": []string{*userID},
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		for {
			var msg InfoWsMessage
			err := c.ReadJSON(&msg)

			if err != nil {
				log.Println("read:", err)
				done <- struct{}{}
				return
			}
			log.Printf("recv: %s of type %s from %s", msg.Message, msg.Type, msg.SenderID)

		}
	}()

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		case <-done:
			log.Println("done processing")
			return
		case <-ticker.C:
			c.WriteMessage(websocket.TextMessage, []byte(`{"event":"ping"}`))
		}
	}
}
