package ws_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type WSHandler func(data any) error

type WSClient struct {
	server          string
	connection      *websocket.Conn
	handlers        []WSHandler
	handlersChanged chan struct{}
	closed          chan struct{}
}

func New(server string) *WSClient {
	return &WSClient{server, nil, make([]WSHandler, 0), make(chan struct{}), make(chan struct{})}
}

func (c *WSClient) Connect() error {
	fmt.Printf("Connecting %v\n", c.server)
	conn, _, err := websocket.DefaultDialer.Dial(c.server, http.Header{})
	if err != nil {
		return fmt.Errorf("failed connecting: %w", err)
	}
	c.connection = conn

	// go c.handleMessages()

	return nil
}

// Sends data in JSON format requires json tags
func (c *WSClient) Send(data any) error {
	return c.connection.WriteJSON(data)
}

func (c *WSClient) handleMessages() {
	errors := make(chan error, len(c.handlers))
	go func() {
		for {
			select {
			case error, ok := <-errors:
				if !ok {
					return
				}
				if error != nil {
					fmt.Printf("failed handling message: %v", error)
				}
			case <-c.handlersChanged:
				errors = make(chan error, len(c.handlers))
			case <-c.closed:
				errors = nil
			}
		}
	}()
	for {
		_, body, err := c.connection.ReadMessage()
		if err != nil {
			fmt.Printf("Failed reading message: %v\n", err)
			if err := c.Close(); err != nil {
				panic(fmt.Errorf("failed closing websocket client: %w", err))
			}
			return
		}
		strbody := string(body)
		fmt.Printf("Got new message of type: %v", strbody)

		for _, handler := range c.handlers {
			go func() {
				if error := handler(strbody); error != nil {
					errors <- error
				}
			}()

		}
	}
}

func (c *WSClient) AddTopicHandler(topic string, handler WSHandler) {
	c.handlers = append(c.handlers, handler)
	c.handlersChanged <- struct{}{}
}

func (c *WSClient) Close() error {
	defer func() { c.closed <- struct{}{} }()
	return c.connection.Close()
}
