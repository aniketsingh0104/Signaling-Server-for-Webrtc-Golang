package signaling

import (
	"encoding/json"
	"log"
)

func (r *RoomManager) deleteUser(connection *Connection, room *Room) {
	// remove user from the room
	delete(room.users, connection)
	// if user is owner then remove from owner field
	if connection == room.owner {
		room.owner = nil
	}
	// close channel
	close(connection.send)
	// if no more users are left then delete the room
	if len(room.users) == 0 {
		delete(r.rooms, room.roomId)
	}
}

func (r *RoomManager) handleChannels() {
	for {
		select {
		case user := <-r.register:
			room, ok := r.rooms[user.roomId]
			if !ok { // case room not found
				// create new room
				room = &Room{
					roomId:   user.roomId,
					isLocked: false, // handle later
					users:    make(map[*Connection]bool),
				}
				room.users[user.connection] = true
				// if this connection is the owner of the room register user as owner
				if user.isOwner {
					room.owner = user.connection
				}
				r.rooms[user.roomId] = room
			} else {
				// if room exists then handle
				// Send READY signal to owner and WAIT signal to other member
				if _, ok := room.users[user.connection]; !ok {
					room.users[user.connection] = true
				}
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
		case unregis := <-r.unregister:
			user := unregis.user
			room, ok := r.rooms[user.roomId]
			if !ok {
				log.Printf("Room %s Not found for user: %s", user.roomId, user.connection.userId)
			} else {
				if unregis.action == ALL {
					for uc := range room.users {
						r.deleteUser(uc, room)
					}
				} else {
					if _, ok := room.users[user.connection]; ok {
						r.deleteUser(user.connection, room)
					} else {
						log.Printf("User %s not present in room %s", user.connection.userId, user.roomId)
					}
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
					for uc := range room.users {
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
