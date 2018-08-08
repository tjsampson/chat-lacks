package server

import (
	"fmt"
	"log"
	"time"
)

const (
	fakeRoom = "fakeRoom"
	fakeUser = "fakeUser"
	fakeText = "fakeText"
)

var (
	fakeTimeStamp, _ = time.Parse(time.RFC822, "01 Jan 18 10:00 UTC")
)

type fakeChatConnection struct{}

func (fc *fakeChatConnection) read() error {
	fmt.Println("Fake Read")
	return nil
}

func (fc *fakeChatConnection) write(message *chatMessage) error {
	fmt.Println("Fake write")
	return nil
}

func (fc *fakeChatConnection) close() {
	fmt.Println("Fake close")
}

func (fc *fakeChatConnection) Close() error {
	fmt.Println("Fake close")
	return nil
}

func newFakeConnection() *fakeChatConnection {
	return &fakeChatConnection{}
}

func fakeChatUser(name string) *chatUser {
	var chatUserName string
	if len(name) > 0 {
		chatUserName = name
	} else {
		chatUserName = fakeUser
	}
	return &chatUser{
		name: chatUserName,
		conn: newFakeConnection(),
	}
}

func fakeChatMessage() *chatMessage {
	return &chatMessage{
		Room:        fakeRoom,
		User:        fakeUser,
		Text:        fakeText,
		MessageType: create,
		Timestamp:   fakeTimeStamp,
	}
}

func fakeLogger() *log.Logger {
	return &log.Logger{}
}

func fakeChatRoom() *chatRoom {
	return &chatRoom{
		name:         fakeUser,
		users:        make(map[*chatUser]bool),
		currentUsers: fakeChatRoomUser(),
		logger:       fakeLogger(),
	}
}

func fakeChatRoomUser() map[string]*chatUser {
	var testData map[string]*chatUser
	testData = make(map[string]*chatUser)
	testData[fakeUser] = fakeChatUser(fakeUser)
	return testData
}

func fakeChatRooms() map[string]*chatRoom {
	var testData map[string]*chatRoom
	testData = make(map[string]*chatRoom)
	testData[fakeRoom] = fakeChatRoom()
	return testData
}

func fakeChatServer() *chatServer {
	return &chatServer{
		logger:    fakeLogger(),
		rooms:     fakeChatRooms(),
		users:     fakeChatRoomUser(),
		userCh:    make(chan *chatUser),
		messageCh: make(chan *chatMessage),
	}
}
