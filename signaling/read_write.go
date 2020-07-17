package signaling

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// readMessage will constantly read message from the websocket connection
func (connection *Connection) readMessage() {
	defer func() {
		// unregister will come here
		connection.ws.Close()
	}()

	for {
		_, byteMsg, err := connection.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg Message
		err = json.Unmarshal(byteMsg, &msg)
		if err != nil {
			log.Printf("error in unmarshalling in readMessage: %v", err)
		}
		// take suitable actions
		switch msg.Action {
		case START:
			user := User{
				connection: connection,
				// data in messgae will be room id
				roomId: msg.Data.(string),
				// only owner of a room can start a meeting
				isOwner: true,
			}
			roomManager.register <- user
			// handle one more thing sending the reply back
			// reply should be handled after the registration so handle in room_managers
		}
	}
}
