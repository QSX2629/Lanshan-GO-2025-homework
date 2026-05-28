// Package ws provides WebSocket connection management.
package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 65536
)

// Conn wraps a WebSocket connection with metadata.
type Conn struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	ws     *websocket.Conn
	mu     sync.Mutex
	send   chan []byte
	closed bool
}

// NewConn creates a new Conn from an upgraded WebSocket connection.
func NewConn(id, userID uint, raw *websocket.Conn) *Conn {
	return &Conn{
		ID:     id,
		UserID: userID,
		ws:     raw,
		send:   make(chan []byte, 256),
	}
}

// ReadPump reads messages from the WebSocket connection and pushes them to a read channel.
func (c *Conn) ReadPump(readCh chan<- *RawMessage) {
	defer func() {
		c.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		readCh <- &RawMessage{ConnID: c.ID, UserID: c.UserID, Data: data}
	}
}

// WritePump writes messages from the send channel to the WebSocket connection.
func (c *Conn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send enqueues a message to be sent on this connection.
func (c *Conn) Send(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	select {
	case c.send <- data:
	default:
		// Buffer full, drop message.
	}
}

// SendJSON marshals and sends a JSON message.
func (c *Conn) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Send(data)
	return nil
}

// Close shuts down the connection.
func (c *Conn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	if c.ws != nil {
		c.ws.Close()
	}
	close(c.send)
}

// RawMessage is a message received from a connection.
type RawMessage struct {
	ConnID uint   `json:"conn_id"`
	UserID uint   `json:"user_id"`
	Data   []byte `json:"data"`
}
