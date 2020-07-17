package signaling

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomManager = RoomManager{
	rooms:     make(map[string]*Room),
	broadcast: make(chan _Message),
	register:  make(chan User),
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Error WebSocketHandler: ", err)
		return
	}
	// make connection here and start readMessage in a thread
	ws.ReadMessage()
}
