package ws

import (
	//"fmt"
	"log"
	"net/http"
	//"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to wait for next ping message from the peer.
	maxPingInterval = 10 * time.Second

	// Frequency of checking whether the connection is alive from the client
	pingCheckPeriod = 20 * time.Second
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// the userID for the client
	userID string
	// The hub providing the messages
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan WsMessage
	// Channel of inbound pings to be responded to
	pings chan struct{}
}

// readMessages reads messages from the the websocket connection.
//
// A goroutine running readMessages is started for each connection. This
// is used to read the ping messages coming from the client to keep the connection alive
func (c *Client) readMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			break
		}
		if string(message) == `{"event":"ping"}` {
			c.pings <- struct{}{}
		}
	}
}

// writeMessages sends messages from the hub to the websocket connection.
//
// A goroutine running writeMessages is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingCheckPeriod)
	pingReceived := time.Now()
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			log.Printf("Sending a msg now")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.closeConnection()
				return
			}

			a := InfoWsMessage{
				Message: message.GetMessage(),
				UserID: message.GetUserID(),
				StoreID: message.GetStoreID(),
				SenderID: message.GetSenderID(),
				Type: message.GetMessageType(),
			}

			err := c.conn.WriteJSON(a)

			if err != nil {
				log.Printf("error writing message: %v", err)
				c.closeConnection()
				return
			}
			c.hub.store.ConfirmSentMessage(message.GetStoreID())
		case <-c.pings:
			pingReceived = time.Now()
			if err := c.conn.WriteControl(websocket.PongMessage, nil, time.Now().Add(writeWait)); err != nil {
				log.Printf("error sending pong: %v", err)
				return
			}
		case <-ticker.C:
			if pingReceived.Before(time.Now().Add(-maxPingInterval)) {
				c.closeConnection()
				return
			}
		}

	}
}

func (c *Client) closeConnection() {
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

var upgrader = websocket.Upgrader{} // use default options

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// get the userID
	userID := r.Header.Get("X-Compass-WS-User")
	if userID == "" {
		// fallback to query params if header is not provided
		queryParam := r.URL.Query()["userId"]
		if len(queryParam) != 1 {
			w.WriteHeader(500)
			return
		}
		userID = queryParam[0]
	}

	log.Println("received subscription request from client")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan WsMessage, 256),
		pings:  make(chan struct{}),
		userID: userID,
	}
	client.hub.register <- client
	hub.users <- userID

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessages()
	go client.readMessages()
}
