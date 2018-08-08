package server

import (
	"log"
	"reflect"
	"testing"
)

func Test_newChatRoom(t *testing.T) {
	type args struct {
		room         string
		currentUsers map[string]*chatUser
		logger       *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *chatRoom
	}{
		{"HappyPath", args{fakeUser, fakeChatRoomUser(), fakeLogger()}, fakeChatRoom()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newChatRoom(tt.args.room, tt.args.currentUsers, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newChatRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_chatRoom_join(t *testing.T) {
	type fields struct {
		name         string
		users        map[*chatUser]bool
		currentUsers map[string]*chatUser
		logger       *log.Logger
	}
	type args struct {
		u *chatUser
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"HappyPath", fields{fakeRoom, fakeChatRoom().users, fakeChatRoomUser(), fakeLogger()}, args{fakeChatUser(fakeUser)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &chatRoom{
				name:         tt.fields.name,
				users:        tt.fields.users,
				currentUsers: tt.fields.currentUsers,
				logger:       tt.fields.logger,
			}
			cr.join(tt.args.u)
		})
	}
}

func Test_chatRoom_broadcast(t *testing.T) {

	testChatRoom := fakeChatRoom()

	type fields struct {
		name         string
		users        map[*chatUser]bool
		currentUsers map[string]*chatUser
		logger       *log.Logger
	}
	type args struct {
		m *chatMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"HappyPath", fields{fakeRoom, testChatRoom.users, testChatRoom.currentUsers, fakeLogger()}, args{fakeChatMessage()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &chatRoom{
				name:         tt.fields.name,
				users:        tt.fields.users,
				currentUsers: tt.fields.currentUsers,
				logger:       tt.fields.logger,
			}
			cr.broadcast(tt.args.m)
		})
	}
}

func Test_chatRoom_leave(t *testing.T) {

	testChatRoom := fakeChatRoom()

	type fields struct {
		name         string
		users        map[*chatUser]bool
		currentUsers map[string]*chatUser
		logger       *log.Logger
	}
	type args struct {
		u *chatUser
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"HappyPath", fields{fakeRoom, testChatRoom.users, testChatRoom.currentUsers, fakeLogger()}, args{fakeChatUser("troy")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &chatRoom{
				name:         tt.fields.name,
				users:        tt.fields.users,
				currentUsers: tt.fields.currentUsers,
				logger:       tt.fields.logger,
			}
			cr.leave(tt.args.u)
		})
	}
}
