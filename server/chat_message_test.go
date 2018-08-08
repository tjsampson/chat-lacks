package server

import (
	"testing"
)

func Test_newChatMessage(t *testing.T) {
	type args struct {
		room    string
		user    string
		message string
		mType   messageType
	}
	tests := []struct {
		name string
		args args
		want *chatMessage
	}{
		{"HappyPath", args{fakeRoom, fakeUser, fakeText, create}, fakeChatMessage()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newChatMessage(tt.args.room, tt.args.user, tt.args.message, tt.args.mType)
			if got.Room != tt.want.Room {
				t.Errorf("Invalid Room = %v, want %v", got.Room, tt.want.Room)
			}
			if got.User != tt.want.User {
				t.Errorf("Invalid User = %v, want %v", got.Room, tt.want.Room)
			}
			if got.Text != tt.want.Text {
				t.Errorf("Invalid Text = %v, want %v", got.Room, tt.want.Room)
			}
		})
	}
}
