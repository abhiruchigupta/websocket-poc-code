package ws

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ListingEvent struct {
	UserId string
	StoreId int64
	Message string
}

func (p ListingEvent) GetUserID() string {
	return p.UserId
}

func (p ListingEvent) GetMessage() string {
	return p.Message
}

func (p ListingEvent) GetStoreID() int64 {
	return p.StoreId
}

type ListingEventMessage struct {
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

var userEmailMap = map[string]string{
	"ls@compass.com":              "538c8ffd85699309abc794e9",
	"grant.harper@compass.com":    "5d52cbb8ad8e4425f8bb471e",
	"alice.yoon@compass.com":      "5c38f2399474a808898e6ba4",
	"alex.zaman@compass.com":      "5ccaf14b5c240259deb5de23",
	"praveen.solanki@compass.com": "5f329df673ba8700018866ea",
}

func (h *Hub) EventReceiver(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var eventMessage ListingEventMessage
	err = json.Unmarshal(body, &eventMessage)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	userID := userEmailMap[eventMessage.Email]
	if userID == "" {
		log.Printf("user email not found: %s, returning 500 error", eventMessage.Email)
		w.WriteHeader(500)
		return
	}
	message, _ := json.Marshal(eventMessage)
	event := ListingEvent{
		UserId:  userID,
		Message: string(message),
	}
	log.Printf("Listing event received User: %s, Message : %s\n",
		event.UserId, event.Message)

	// Before adding it to channel we want to persist this listing event
	event.StoreId, _ = h.store.StoreMessage(userID, event.Message, time.Now())
	h.BroadcastMessage(event)
	w.Write([]byte(`{"status": "ok"}`))
}
