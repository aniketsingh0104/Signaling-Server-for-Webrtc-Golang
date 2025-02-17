package signaling

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the peer.
	readWait = 60 * time.Second
)

// readMessage will constantly read message from the websocket connection
func (connection *Connection) readMessage() {
	defer func() {
		// unregister will come here
		connection.ws.Close()
	}()
	// set maximum time limit for reading a messgage
	connection.ws.SetReadDeadline(time.Now().Add(readWait))
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
			//TODO: Checks
			// 1. Rooms should exist in database
			// 2. If user is the owner then room should belong to him
			// set user id
			connection.userId = msg.UserId
			// data in messgae will be room id
			user.roomId = msg.Data.(string)
			// only owner of a room can start a meeting
			user.isOwner = true
			RManager.register <- user
			// handle one more thing sending the reply back
			// reply should be handled after the registration so handle in room_managers
		case JOIN:
			//TODO: Checks
			// 1. Rooms should exist in database
			// set user id
			connection.userId = msg.UserId
			// data in messgae will be room id
			user.roomId = msg.Data.(string)
			// only owner of a room can start a meeting
			user.isOwner = false
			RManager.register <- user
		case END:
			// handle deregistration
			// only applicable when the requester is the owner of the room
			// check if user is the owner
			if user.isOwner {
				// iterate though all users and remove them one by one
				RManager.unregister <- Unregister{user: user, action: ALL}
			} else {
				log.Println("User is not the owner of the room END not applicable")
				// TODO: send some reply to user
			}
		case LEAVE:
			// just send user to leave
			RManager.unregister <- Unregister{user: user, action: SELF}
		case MESSAGE:
			if user.roomId != "" {
				broadcastMess := _Message{
					ws:      connection.ws,
					message: msg,
					roomId:  user.roomId,
				}
				RManager.broadcast <- broadcastMess
			} else {
				log.Printf("Error in broadcast message: %v", err)
				// TODO: send some reply of error
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (connection *Connection) writeMessage() {
	defer func() {
		connection.ws.Close()
	}()
	for {
		message, ok := <-connection.send
		if !ok {
			connection.write(websocket.CloseMessage, []byte{})
			return
		}
		if err := connection.write(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
