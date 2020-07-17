package signaling

import (
	"encoding/json"
	"log"
)

func (r *RoomManager) handlerChnnels() {
	for {
		select {
		case user := <-r.register:
			// if the user is the owner of this room then check database
			// is user has this room or not only proceed then
			room, ok := r.rooms[user.roomId]
			if !ok { // case room not found
				// create new room
				room = &Room{
					roomId:   user.roomId,
					isLocked: false, // handle later
					users:    []*Connection{user.connection},
				}
				// if this connection is the owner of the room regis user as owner
				if user.isOwner {
					room.owner = user.connection
				}
				r.rooms[user.roomId] = room
			} else {
				room.users = append(room.users, user.connection)
				if user.isOwner {
					room.owner = user.connection
				}
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
