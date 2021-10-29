package main

import (
	"flag"
	"log"
	"net/http"
	"ws/ws"
)

var addr = flag.String("addr", "localhost:8080", "http service address")


func main() {
	flag.Parse()
	log.SetFlags(0)
	hub := ws.NewHub()
	go hub.Run()

	go hub.HandleMissedEvents()
	hub.InitializeStore()
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	http.HandleFunc("/mls", hub.EventReceiver)

	http.HandleFunc("/send", hub.PostEvent)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	hub.CloseStore()
}
