package server

import (
	"fmt"
	"strings"
)

const chatHelp = `Welcome to the Lacks Chat server!
Commands:
  /help       see this help message again    (example: /help)
  /listusers  see all connected users        (example: /listusers)
  /listrooms  see all rooms                  (example: /listrooms)
  /newroom    create and join new room       (example: /newroom random)
  /join       join a room                    (example: /join random)
  /leave      leave a room                   (example: /leave random)
  /mute       mute a user                    (example: /mute rob)
  /unmute     unmute a user                  (example: /unmute rob)
  /mutes      see who you've muted           (example: /mutes)
  /quit       disconnect from chat server    (example: /quit)
  /dm         send a message to a user       (example: /dm troy: hello!)
`

type command func(tc *tcpUser, arg string)

// commands are each action a client is able to perform besides just sending
// plaintext.
var commands = map[string]command{
	"/help":      helpCmd,
	"/listusers": listUsersCmd,
	"/listrooms": listRoomsCmd,
	"/newroom":   newRoomCmd,
	"/join":      joinRoomCmd,
	"/leave":     leaveRoomCmd,
	"/mute":      muteCmd,
	"/unmute":    unmuteCmd,
	"/mutes":     mutesCmd,
	"/quit":      quitCmd,
	"/dm":        dmCmd,
}

// Show Help Message Command
func helpCmd(tc *tcpUser, _ string) {
	tc.write(defChatMsg(chatHelp, text))
}

// Default Chat Message Wrapper
// Syntax sugar to remove duplication
func defChatMsg(msg string, msgType messageType) *chatMessage {
	return newChatMessage("you", "server", msg, msgType)
}

// Create New ChatRoom Command
func newRoomCmd(tc *tcpUser, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		tc.write(defChatMsg("What room are you trying to create? Room name is empty.\n", text))
		return
	}
	if arg == tc.room {
		tc.write(defChatMsg("That room has already been created. Use '\\join "+tc.room+"'\n", text))
		return
	}
	tc.send <- newChatMessage(arg, tc.user, tc.user+" created new room "+arg, create)
}

// Join ChatRoom Commmand
func joinRoomCmd(tc *tcpUser, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		tc.write(defChatMsg("What room are you trying to join? Room name is empty.\n", text))
		return
	}
	if arg == tc.room {
		tc.write(defChatMsg("Easy buddy! You've already joined this room\n", text))
		return
	}
	tc.send <- newChatMessage(arg, tc.user, tc.user+" joined room "+arg, join)
}

// Leave ChatRoom Command
func leaveRoomCmd(tc *tcpUser, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		tc.write(defChatMsg("What room are you trying to leave? Room name cannot be blank\n", text))
	}
	tc.send <- newChatMessage(arg, tc.user, tc.user+" joined room "+arg, leave)
}

// Mute a User Command
func muteCmd(tc *tcpUser, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		tc.writeText("Who are you trying to mute? Username cannot be empty\n")
		return
	}
	if arg == tc.user {
		tc.writeText("Think, McFly! You can't mute yourself\n")
		return
	}

	tc.send <- newChatMessage(arg, tc.user, "Muted user "+arg+".\n", mute)
}

func unmuteCmd(tc *tcpUser, arg string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		tc.writeText("Username cannot be blank\n")
		return
	}
	tc.send <- newChatMessage(arg, tc.user, "Unmuted user "+arg+".\n", unmute)
}

func mutesCmd(tc *tcpUser, _ string) {
	if len(tc.muted) < 1 {
		tc.writeText("Hello! You haven't muted anyone.\n")
		return
	}
	var mutes []string
	for mute := range tc.muted {
		mutes = append(mutes, mute)
	}
	muteList := strings.Join(mutes, "\n  - ")
	tc.writeText("You've muted:\n  - " + muteList + "\n")
}

func dmCmd(tc *tcpUser, arg string) {
	i := strings.Index(arg, ":")
	if i == -1 {
		tc.writeText("/dm command not understood, you're missing the colon separator ':' from your command.\n")
		tc.writeText("(example: /dm troy: hello!)\n")
		return
	}
	msgs := strings.SplitAfterN(arg, ":", 2)
	if len(msgs) < 2 {
		tc.writeText("/dm command not understood, commands appears to be malformed. Type '/help' to see how to use each command.\n")
		tc.writeText("(example: /dm troy: hello!)\n")
		return
	}

	user := strings.TrimSpace(strings.TrimRight(msgs[0], ":"))
	if user == "" {
		tc.writeText("/dm command not understood, you're missing a username.\n")
		return
	}

	if user == tc.user {
		tc.writeText("You can't send a dm to yourself.\n")
		return
	}

	msg := strings.TrimSpace(msgs[1])
	if msg == "" {
		tc.writeText("/dm command not understood, it looks like your message is blank.\n")
		return
	}
	tc.send <- (newChatMessage(user, tc.user, msg, dm))
}

func listUsersCmd(tc *tcpUser, arg string) {
	if arg != "" {
		tc.send <- newChatMessage(arg, tc.user, "", listUsers)
		return
	}
	tc.send <- newChatMessage("", tc.user, "", listUsers)
}

func listRoomsCmd(tc *tcpUser, _ string) {
	tc.send <- newChatMessage("", tc.user, "", listRooms)
}

func quitCmd(tc *tcpUser, _ string) {
	tc.send <- newChatMessage(defaultRoom, tc.user, fmt.Sprintf("%s has disconnected", tc.user), quit)
}
