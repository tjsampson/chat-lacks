package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tjsampson/chat-lacks/core"
)

// chatServer data structure
type chatServer struct {
	logger    *log.Logger
	rooms     map[string]*chatRoom
	users     map[string]*chatUser
	userCh    chan *chatUser
	messageCh chan *chatMessage
}

func newChatServer(log *log.Logger) *chatServer {
	return &chatServer{
		logger:    log,
		rooms:     make(map[string]*chatRoom),
		users:     make(map[string]*chatUser),
		userCh:    make(chan *chatUser),
		messageCh: make(chan *chatMessage),
	}
}

func (server *chatServer) newUser(u *chatUser) {
	server.users[u.name] = u
	server.rooms[defaultRoom].join(u)
	go u.conn.read()
}

func (server *chatServer) listUsers(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	var users []string
	if m.Room != "" {
		cr, ok := server.rooms[m.Room]
		if !ok {
			err := user.conn.write(newChatMessage("you", "server", "Room "+m.Room+" doesn't exist.\n", text))
			if err != nil {
				server.logger.Panicf("error writing message: %v", err.Error())
			}
			return
		}
		for u := range cr.users {
			users = append(users, u.name)
		}
	} else {
		for user := range server.users {
			users = append(users, user)
		}
	}
	m.Text = strings.Join(users, ",")
	user.conn.write(m)
}

func (server *chatServer) listRooms(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	var rooms []string
	for room := range server.rooms {
		rooms = append(rooms, room)
	}
	m.Text = strings.Join(rooms, ",")
	user.conn.write(m)
}

func (server *chatServer) joinRoom(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	if room, ok := server.rooms[m.Room]; ok {
		room.join(user)
		return
	}
	m.Text = "Sorry, the room " + m.Room + " doesn't exist.\n"
	m.MessageType = text
	user.conn.write(m)
}

func (server *chatServer) leaveRoom(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	if m.Room == defaultRoom {
		m.Text = "You can't leave the " + defaultRoom + " room).\n"
		m.MessageType = text
		user.conn.write(m)
		return
	}
	ch, ok := server.rooms[m.Room]
	if !ok {
		m.Text = "The " + m.Room + " room does not exist. \n"
		m.MessageType = text
		user.conn.write(m)
		return
	}
	if _, ok = ch.users[user]; !ok {
		m.Text = "You're not a member of the " + m.Room + " room.\n"
		m.MessageType = text
		user.conn.write(m)
		return
	}
	ch.leave(user)
	m.Text = "Left " + m.Room + " room. Returning to the general room.\n"
	user.conn.write(m)
}

func (server *chatServer) createRoom(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	room, ok := server.rooms[m.Room]
	if ok {
		room.join(user)
		return
	}
	newRoom := newChatRoom(m.Room, server.users, server.logger)
	server.rooms[m.Room] = newRoom
	newRoom.join(user)
}

func (server *chatServer) broadcast(m *chatMessage) {
	server.logger.Printf("(%s to %s: %v): %s", m.User, m.Room, time.Now().Format(time.RFC850), m.Text)
	room, ok := server.rooms[m.Room]
	if !ok {
		return
	}
	room.broadcast(m)
}

func (server *chatServer) mute(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	if _, ok := server.users[m.Room]; !ok {
		m.MessageType = text
		m.Text = "The user " + m.Room + " doesn't exist.\n"
		user.conn.write(m)
		return
	}
	user.conn.write(m)
}

func (server *chatServer) unmute(m *chatMessage) {
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	if _, ok := server.users[m.Room]; !ok {
		m.MessageType = text
		m.Text = "The user " + m.Room + " doesn't exist.\n"
		user.conn.write(m)
		return
	}
	user.conn.write(m)
}

func (server *chatServer) dm(m *chatMessage) {
	server.logger.Printf("(%s to %s): %s", m.User, m.Room, m.Text)
	sender, ok := server.users[m.User]
	if !ok {
		return
	}
	recipient, ok := server.users[m.Room]
	if !ok {
		m.MessageType = text
		m.Text = "Sorry, the user " + m.Room + " doesn't exist.\n"
		sender.conn.write(m)
		return
	}
	recipient.conn.write(m)
	sender.conn.write(m)
}

func (server *chatServer) quit(m *chatMessage) {
	server.logger.Printf("(%s to %s): %s", m.User, m.Room, m.Text)
	user, ok := server.users[m.User]
	if !ok {
		return
	}
	user.conn.close()
	delete(server.users, m.User)
	server.rooms[defaultRoom].broadcast(m)
}

func (server *chatServer) run() {
	server.rooms[defaultRoom] = newChatRoom(defaultRoom, server.users, server.logger)
	for {
		select {
		case user := <-server.userCh:
			server.newUser(user)

		case message := <-server.messageCh:
			switch message.MessageType {

			case join:
				server.joinRoom(message)

			case listUsers:
				server.listUsers(message)

			case listRooms:
				server.listRooms(message)

			case create:
				server.createRoom(message)

			case leave:
				server.leaveRoom(message)

			case text:
				server.broadcast(message)

			case mute:
				server.mute(message)

			case unmute:
				server.unmute(message)

			case dm:
				server.dm(message)

			case quit:
				server.quit(message)
			}
		}
	}
}

func (server *chatServer) serveTCP(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	server.logger.Println("tcp server listening on", port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				server.logger.Println(err.Error())
			}
			go func() {
				server.userCh <- createUserConn(conn, server)
			}()
		}
	}()

	server.run()
	return nil
}

// ListenAndServe is the main entry point for the server
// Logger Pointer Param
// AppConfig Pointer Param
func ListenAndServe(l *log.Logger, cfg *core.AppConfig) error {
	server := newChatServer(l)
	errCh := make(chan error, 4)

	// Server up the TCP backend
	go server.serveTCP(":" + cfg.TCPPort)

	// Server up the HTTP Backend
	HTTPPort := cfg.HTTPPort
	l.Printf("http server listening on: %s", HTTPPort)
	go http.ListenAndServe(":"+HTTPPort, newMuxRouter(server))

	// Catch signals for shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errCh:
			if err != nil {
				// panic out to unwind the stack for errors
				l.Panicf("error channel ListenAndServe: %s\n", err.Error())
			}
		case s := <-signalCh:
			// panic out to unwind the stack for caught signals
			l.Panicf("Captured %v. Exiting...\n", s)
		}
	}
}
