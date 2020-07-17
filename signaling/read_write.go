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
	user := User{connection: connection}
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
			// set user id
			connection.userId = msg.UserId
			// data in messgae will be room id
			user.roomId = msg.Data.(string)
			// only owner of a room can start a meeting
			user.isOwner = true
			roomManager.register <- user
			// handle one more thing sending the reply back
			// reply should be handled after the registration so handle in room_managers
		case JOIN:
			// set user id
			connection.userId = msg.UserId
			// data in messgae will be room id
			user.roomId = msg.Data.(string)
			// only owner of a room can start a meeting
			user.isOwner = false
			roomManager.register <- user
		case END:
			// handle deregistration
			// only applicable when the requester is the owner of the room
		case LEAVE:
			// handle deregistration
			// remove user from the room and if room becomes empty then delete room
		case MESSAGE:
			if user.roomId != "" {
				broadcastMess := _Message{
					ws:      connection.ws,
					message: msg,
					roomId:  user.roomId,
				}
				roomManager.broadcast <- broadcastMess
			} else {
				log.Printf("Error in broadcast message: %v", err)
				// error reply
				// marshalled, err := json.Marshal(Message{
				// 	Action: READY,
				// 	UserId: room.owner.userId,
				// })
				// if err != nil {
				// 	log.Fatalf("Marshalling Error in Register User: %v", err)
				// }
				// room.owner.send <- marshalled
			}
		}
	}
}

// // write writes a message with the given message type and payload.
// func (c *Connection) write(mt int, payload []byte) error {
// 	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
// 	return c.ws.WriteMessage(mt, payload)
// }

func (connection *Connection) writeMessage() {
	defer func() {
		connection.ws.Close()
	}()
	for {
		// select {
		// case message, ok := <-connection.send:
		// 	if !ok {
		// 		// handle close connection
		// 		return
		// 	}
		// 	// if err := connection.write()
		// }
	}
}
