package server

import (
	"log"
)

const (
	defaultRoom = "general"
)

type chatRoom struct {
	name         string
	users        map[*chatUser]bool
	currentUsers map[string]*chatUser
	logger       *log.Logger
}

func newChatRoom(room string, currentUsers map[string]*chatUser, logger *log.Logger) *chatRoom {
	return &chatRoom{
		name:         room,
		users:        make(map[*chatUser]bool),
		currentUsers: currentUsers,
		logger:       logger,
	}
}

func (cr *chatRoom) join(u *chatUser) {
	if _, ok := cr.users[u]; ok {
		u.conn.write(newChatMessage(cr.name, u.name, "Changing to room "+cr.name+"\n", join))
		return
	}
	cr.users[u] = true
	cr.broadcast(newChatMessage(cr.name, u.name, u.name+" has joined "+cr.name+"\n", join))
}

func (cr *chatRoom) broadcast(m *chatMessage) {
	for u := range cr.users {
		if _, ok := cr.currentUsers[u.name]; !ok {
			continue
		}
		err := u.conn.write(m)
		if err != nil {
			cr.logger.Println("Broadcast error: ", err.Error())
		}
	}
}

func (cr *chatRoom) leave(u *chatUser) {
	delete(cr.users, u)
}
