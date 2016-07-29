package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Allow all requests
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("websocket: upgrade error:", err)
		return
	}
	defer c.Close()

	watch := NewWatch("http://192.168.250.100:4646", c, 1*time.Minute)
	go watch.Start()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("websocket: connection error:", err)
			break
		}
		log.Printf("recv: %s", message)
		break
	}

	watch.Stop()
}
