package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestHub(t *testing.T) {

	var wg sync.WaitGroup
	// spin up a hub and run it
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	http.HandleFunc("/send", hub.PostEvent)
	go func() {
		err := http.ListenAndServe("localhost:8713", nil)
		if err != nil {
			log.Printf("error starting server: %v\n", err)
		}
	}()

	clientDone := make(chan struct{})
	clientTotal := 2
	// spin up a client and connect
	for i := 0; i < clientTotal; i++ {
		wg.Add(1)
		userID := fmt.Sprintf("user%d", i)
		go RunTestClient(userID, clientDone)
		go SendMessages(userID)
	}

	clientShutdownCount := 0
	for {
		<-clientDone
		clientShutdownCount++
		if clientTotal == clientShutdownCount {
			break
		}
	}
	log.Println("shutting down")
}

func SendMessages(userID string) {
	for {
		msg := fmt.Sprintf("hey %s", userID)
		requestBody, _ := json.Marshal(map[string]string{
			"message": msg,
		})
		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://localhost:8713/send", bytes.NewBuffer(requestBody))
		req.Header.Add("X-Compass-WS-User", userID)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error sending post: %v\n", err)
			continue
		}
		log.Printf("post status: %d", resp.StatusCode)
		time.Sleep(1 * time.Second)
	}
}

func RunTestClient(userID string, clientDone chan struct{}) {
	query := fmt.Sprintf("userId=%s", userID)
	u := url.URL{Scheme: "ws", Host: "localhost:8713", Path: "/subscribe", RawQuery: query}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				done <- struct{}{}
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	pingCounter := 0
	for {
		select {
		case <-done:
			log.Println("done processing")
			clientDone <- struct{}{}
			return
		case <-ticker.C:
			if pingCounter > 3 {
				log.Println("client done firing")
				clientDone <- struct{}{}
				return
			}
			pingCounter++
			log.Println("sending ping")
			c.WriteMessage(websocket.TextMessage, []byte(`{"event":"ping"}`))
		}
	}
}
