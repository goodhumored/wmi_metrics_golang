package ws_server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type HandlerFunc func(conn *websocket.Conn, w http.ResponseWriter, r *http.Request)

type WSServer struct {
	upgrader websocket.Upgrader
	path     string
	handlers []HandlerFunc
	router   *mux.Router
}

func New(path string, router *mux.Router) *WSServer {
	return &WSServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		path:   path,
		router: router,
	}
}

func (s *WSServer) Start() {
	s.router.HandleFunc(s.path, s.handleRequest)
	fmt.Printf("Starting WebSocket server on %s\n", s.path)
}

// Adds another handler for the same path
func (s *WSServer) AddHandler(handler HandlerFunc) {
	s.handlers = append(s.handlers, handler)
}

func (s *WSServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup

	wg.Add(len(s.handlers))
	for _, handler := range s.handlers {
		go func() {
			defer wg.Done()
			handler(conn, w, r)
		}()
	}
	wg.Wait()
	err = conn.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
	}
}
