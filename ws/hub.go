package ws

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const SQLITE_STORE = "sqlite"

func storeFactory(storeType string) MessageStore {
	if storeType == SQLITE_STORE {
		return &SqliteMessageStore{
			Db: nil,
		}
	}
	return &DummyMessageStore{}
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registry by userID
	clientReg map[string][]*Client

	// Inbound messages from the clients.
	broadcast chan WsMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Event store
	store MessageStore

	// register clients channel for handling missed events.
	users chan string
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan WsMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clientReg:  make(map[string][]*Client),
		store:      storeFactory(SQLITE_STORE),
		users: 		make(chan string, 10),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Println("hub registering client")
			if conns, ok := h.clientReg[client.userID]; ok {
				h.clientReg[client.userID] = append(conns, client)
			} else {
				h.clientReg[client.userID] = []*Client{client}
			}
		case client := <-h.unregister:
			log.Println("hub unregistering client")
			h.removeClient(client)
		case message := <-h.broadcast:
			log.Println("hub broadcasting message")
			clients := h.clientReg[message.GetUserID()]
			for _, client := range clients {
				log.Printf("found client to send msg to")
				select {
				case client.send <- message:
				default:
					close(client.send)
					// do we need to remove the client here if this happens?
				}
			}
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	clients := h.clientReg[client.userID]
	targetIndex := -1
	for i, c := range clients {
		if client == c {
			targetIndex = i
		}
	}
	if targetIndex == -1 {
		log.Println("error attempting to remove client that is no longer in the registry")
		return
	}
	clients[targetIndex] = clients[len(clients)-1]
	clients[len(clients)-1] = nil
	clients = clients[:len(clients)-1]
	h.clientReg[client.userID] = clients
}

type Info struct {
	Message string `json:"message"`
}

type InfoWsMessage struct {
	message string
	userID  string
	storeID int64
	senderID string
}

func (i InfoWsMessage) GetUserID() string {
	return i.userID
}

func (i InfoWsMessage) GetMessage() string {
	return i.message
}

func (i InfoWsMessage) GetStoreID() int64 {
	return i.storeID
}

func (i InfoWsMessage) GetSenderID() string {
	return i.senderID
}

func (h *Hub) PostEvent(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Compass-WS-User")
	senderID := r.Header.Get("X-Compass-WS-Sender")

	if userID == "" {
		w.WriteHeader(500)
		return
	}

	if senderID == "" {
		w.WriteHeader(500)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var info Info
	err = json.Unmarshal(body, &info)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	message := InfoWsMessage{
		userID:  userID,
		message: info.Message,
		storeID: 0,
		senderID: senderID,
	}
	h.BroadcastMessage(message)
}

func (h *Hub) BroadcastMessage(message WsMessage) {
	log.Printf("userID: %s, message: %s, senderID: %s\n", message.GetUserID(), message.GetMessage(), message.GetSenderID())
	h.broadcast <- message
}

func (h *Hub) InitializeStore() {
	h.store.Initialize()
}

func (h *Hub) CloseStore() {
	h.store.Close()
}

func (h *Hub) HandleMissedEvents () {
	for {
		userId := <-h.users
		log.Printf("Checking missed events for user %s", userId)
		dbMessages := h.store.RetrieveUnsentMessages(userId)
		for _, msg := range dbMessages {
			eventMessage := MessageEvent{
				UserId: userId,
				Message: *msg.Message,
				StoreId: *msg.ID,
				SenderId: *msg.SenderID,
			}
			h.BroadcastMessage(eventMessage)
		}
	}
}
