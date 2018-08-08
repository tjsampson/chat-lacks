package server

import (
	"context"
	"log"
	"net/http"
)

// "auth" Key is used for Cookie Identifier
// NOT SECURE: DO NOT USE IN PRODUCTION
type authKey int

const (
	cookieKey authKey = iota
)

type adapter func(http.Handler) http.Handler

func adapt(h http.Handler, adapters ...adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func logHandler(logger *log.Logger) adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("incoming request method: %v", r.Method)
			defer logger.Printf("outgoing response")
			h.ServeHTTP(w, r)
		})
	}
}

func authContext(cs *chatServer) adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cs.logger.Printf("checking http auth context")
			cookie, _ := r.Cookie("username")
			if cookie != nil {
				cs.logger.Printf("cookie exists: %v", cookie)
				ctx := context.WithValue(r.Context(), cookieKey, cookie.Value)
				h.ServeHTTP(w, r.WithContext(ctx))
			} else {
				cs.logger.Println("anonymous user")
				w.Write([]byte(`Unknown User. Please send a POST to /users to create your handle`))
			}
		})
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Welcome the Lacks Chat Server!
Try sending a message via a POST request to /messages,
or connecting via websockets by sending a POST request with your desired username to /ws.`))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`STUBBED OUT`))
}

func newMuxRouter(cs *chatServer) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", adapt(http.HandlerFunc(homeHandler), authContext(cs), logHandler(cs.logger)))
	mux.Handle("/messages", adapt(http.HandlerFunc(messageHandler), authContext(cs), logHandler(cs.logger)))
	return mux
}
