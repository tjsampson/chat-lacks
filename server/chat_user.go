package server

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type chatConn interface {
	read() error
	write(message *chatMessage) error
	close()
}

type chatUser struct {
	name string
	conn chatConn
}

type connUser struct {
	room  string
	user  string
	muted map[string]bool
	send  chan<- *chatMessage
}

type tcpUser struct {
	conn net.Conn
	r    *bufio.Reader
	*connUser
}

func createUserConn(conn net.Conn, server *chatServer) *chatUser {
	u := newTCPUser(conn, server)

	err := u.write(newChatMessage(u.room, u.user, chatHelp, text))
	if err != nil {
		server.logger.Panicf("failed to write chat message: %v \n", err.Error())
	}

	return &chatUser{
		name: u.name(),
		conn: u,
	}
}

func attemptConnWrite(logger *log.Logger, conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		logger.Panicf("failed to write connection: %v \n", err.Error())
	}
}

func newTCPUser(conn net.Conn, server *chatServer) *tcpUser {
	var name string
	r := bufio.NewReader(conn)
	attemptConnWrite(server.logger, conn, "Please enter your handle (i.e. username): ")
	for {
		n, err := r.ReadString('\n')
		if err != nil {
			conn.Close()
		}
		n = strings.TrimSpace(n)
		if n == "" {
			attemptConnWrite(server.logger, conn, "No blank names. Try again: ")
			continue
		}
		if _, ok := server.users[n]; !ok {
			name = n
			break
		}
		attemptConnWrite(server.logger, conn, "Seats taken! Try again: ")
		continue
	}

	return &tcpUser{
		r:    bufio.NewReader(conn),
		conn: conn,
		connUser: &connUser{
			room:  defaultRoom,
			muted: make(map[string]bool),
			user:  name,
			send:  server.messageCh,
		},
	}
}

func (tc *tcpUser) read() error {
	for {
		msg, err := tc.r.ReadString('\n')
		if err != nil {
			tc.send <- newChatMessage("everyone", tc.user, tc.user+" has left the chat\n", quit)
			return err
		}
		if ok := tc.handleCommand(msg); ok {
			continue
		}
		tc.send <- newChatMessage(tc.room, tc.user, msg, text)
	}
}

func (tc *tcpUser) write(message *chatMessage) error {
	if _, ok := tc.muted[message.User]; ok {
		return nil
	}
	switch message.MessageType {
	case text:
		return tc.writeText(message.Timestamp.Format("Mon Jan 2, 2006 03:04:05PM") + " - " + "(" + message.User + " to " + message.Room + "): " + message.Text)

	case listUsers, listRooms:
		return tc.writeText(message.Text + "\n")

	case join, create:
		tc.room = message.Room
		return tc.writeText(message.Text)

	case leave:
		tc.room = defaultRoom
		return tc.writeText(message.Text)

	case mute:
		if _, ok := tc.muted[message.Room]; !ok {
			tc.muted[message.Room] = true
			return tc.writeText(message.Text)
		}
		return tc.writeText("User " + message.Room + " is already muted.\n")

	case unmute:
		if _, ok := tc.muted[message.Room]; ok {
			delete(tc.muted, message.Room)
			return tc.writeText(message.Text)
		}
		return tc.writeText("User " + message.Room + " isn't muted.\n")

	case dm:
		tc.writeText("(" + message.User + " to " + message.Room + "): " + message.Text + "\n")
	}
	return nil
}

func (tc *tcpUser) writeText(text string) error {
	_, err := tc.conn.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

func (tc *tcpUser) close() {
	tc.conn.Close()
}

func (tc *tcpUser) name() string {
	return tc.user
}

func (tc *tcpUser) handleCommand(s string) bool {
	if !strings.HasPrefix(s, "/") {
		return false
	}
	cmd := strings.TrimSpace(strings.Split(s, " ")[0])
	cmdFunc, ok := commands[cmd]
	if !ok {
		tc.write(newChatMessage(tc.room, tc.user, "Command "+cmd+" doesn't exist\n", text))
		return true
	}
	cmdArg := strings.TrimSpace(strings.TrimPrefix(s, cmd))
	cmdFunc(tc, cmdArg)
	return true
}
