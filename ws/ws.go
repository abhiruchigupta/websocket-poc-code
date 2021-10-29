package ws
//
//import (
//	"encoding/json"
//	"io/ioutil"
//	"log"
//	"net/http"
//
//	"github.com/gorilla/websocket"
//)
//
//const EVENT_BUFFER = 100
//
//type Server struct {
//	eventBus chan WsMessage
//	connReg  map[string][]*websocket.Conn
//}
//
//func NewServer() *Server {
//	return &Server{
//		eventBus: make(chan WsMessage, EVENT_BUFFER),
//		connReg:  make(map[string][]*websocket.Conn),
//	}
//}
//
//func (s *Server) SubscribeToEvents(w http.ResponseWriter, r *http.Request) {
//
//	userID := r.Header.Get("X-Compass-WS-User")
//	if userID == "" {
//		// fallback to query params if header is not provided
//		queryParam := r.URL.Query()["userId"]
//		if len(queryParam) != 1 {
//			w.WriteHeader(500)
//			return
//		}
//		userID = queryParam[0]
//	}
//
//	c, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Print("upgrade:", err)
//		return
//	}
//	if conns, ok := s.connReg[userID]; ok {
//		s.connReg[userID] = append(conns, c)
//	} else {
//		s.connReg[userID] = []*websocket.Conn{c}
//	}
//}
//
//func (s *Server) RunWsManager() {
//
//	for {
//		event := <-s.eventBus
//		conns, ok := s.connReg[event.GetUserID()]
//		if !ok {
//			log.Printf("no connection present for user %s\n", event.GetUserID())
//			continue
//		}
//		message := []byte(event.GetMessage())
//		for _, c := range conns {
//			err := c.WriteMessage(1, message)
//			if err != nil {
//				log.Println("write error:", err)
//			}
//		}
//	}
//}
//
//func (s *Server) PostEvent(w http.ResponseWriter, r *http.Request) {
//	userID := r.Header.Get("X-Compass-WS-User")
//	senderID := r.Header.Get("X-Compass-WS-Sender")
//
//	if userID == "" {
//		w.WriteHeader(500)
//		return
//	}
//
//	if senderID == "" {
//		w.WriteHeader(500)
//		return
//	}
//
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	var info Info
//	err = json.Unmarshal(body, &info)
//	if err != nil {
//		w.WriteHeader(500)
//		return
//	}
//
//	message := InfoWsMessage{
//		userID:  userID,
//		message: info.Message,
//		storeID: 0,
//		senderID: senderID,
//	}
//	log.Printf("userID: %s, message: %s\n", message.GetUserID(), message.GetMessage())
//	s.eventBus <- message
//}
