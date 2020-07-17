package signaling

import (
	"encoding/json"
	"log"
)

func (r *RoomManager) handlerChnnels() {
	for {
		select {
		case user := <-r.register:
			// Checks
			// 1. Rooms should exist in database
			// 2. If user is the owner then room should belong to him
			room, ok := r.rooms[user.roomId]
			if !ok { // case room not found
				// create new room
				room = &Room{
					roomId:   user.roomId,
					isLocked: false, // handle later
					users:    []*Connection{user.connection},
				}
				// if this connection is the owner of the room register user as owner
				if user.isOwner {
					room.owner = user.connection
				}
				r.rooms[user.roomId] = room
			} else {
				// if room exists then handle
				// Send READY signal to owner and WAIT signal to other member
				room.users = append(room.users, user.connection)
				if user.isOwner {
					room.owner = user.connection
				}
				// send READY to owner to start offer process
				marshalled, err := json.Marshal(Message{
					Action: READY,
					UserId: room.owner.userId,
				})
				if err != nil {
					log.Fatalf("Marshalling Error in Register User: %v", err)
				}
				room.owner.send <- marshalled
				// send WAIT to other memeber
				replyMess := _Message{
					ws: room.owner.ws,
					message: Message{
						Action: WAIT,
					},
					roomId: user.roomId,
				}
				r.broadcast <- replyMess
			}
		case mess := <-r.broadcast:
			if room, ok := r.rooms[mess.roomId]; ok {
				// get marshal
				marshalled, err := json.Marshal(mess.message)
				if err != nil {
					log.Printf("error in marshalling in broadcast: %v", err)
				} else {
					// loop through all users and broadcast message
					for _, uc := range room.users {
						// don't send message to sender again
						if uc.ws != mess.ws {
							// send to the user channel
							uc.send <- marshalled
						}
					}
				}
			}
		}
	}
}
