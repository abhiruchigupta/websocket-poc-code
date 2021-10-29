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


type WsMessage interface {
	GetUserID() string
	GetStoreID() int64
	GetMessage() string
	GetSenderID() string
}

type InfoWsMessage struct {
	message string
	userID  string
	storeID int64
	senderID string
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
			//var msg string

			//err := c.ReadJSON(&msg)

			_, message, err := c.ReadMessage()
			//data := []byte(message)

			log.Printf("recv: %s", message)

			//log.Println("msg rcvd: ", message)

			//var input InfoWsMessage
			//
			//err = json.Unmarshal(data, &input)
			//fmt.Println(err)

			//fmt.Printf("%v , %v, %v, %v", input.message, input.senderID, input.storeID, input.userID)

			//log.Println("error: , msg: , final json: input", err, msg, input)



			if err != nil {
				log.Println("read:", err)
				done <- struct{}{}
				return
			}

			//log.Print("recv: %v from sender: %v", input.message, input.senderID)
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
