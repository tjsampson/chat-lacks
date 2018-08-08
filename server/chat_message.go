package server

import "time"

type messageType int

const (
	join = messageType(iota)
	listUsers
	listRooms
	create
	leave
	text
	mute
	unmute
	dm
	quit
)

type chatMessage struct {
	Room        string
	User        string
	Text        string
	Timestamp   time.Time
	MessageType messageType
}

func newChatMessage(room, user, message string, mType messageType) *chatMessage {
	return &chatMessage{
		Room:        room,
		User:        user,
		Text:        message,
		Timestamp:   time.Now(),
		MessageType: mType,
	}
}
