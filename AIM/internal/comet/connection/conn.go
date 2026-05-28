package connection

import (
	"AIM/internal/comet/protocol"
	"AIM/internal/storage/redis" // 只保留Redis，去掉friend依赖

	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	conn *websocket.Conn
	uid  string
	send chan []byte
	mu   sync.Mutex
}

var (
	connections = make(map[string]*Connection)
	mu          sync.RWMutex
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "uid required", 400)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := &Connection{
		conn: ws,
		uid:  uid,
		send: make(chan []byte, 256),
	}

	mu.Lock()
	connections[uid] = conn
	mu.Unlock()

	// 用户上线：标记Redis在线
	_ = redis.SetOnline(uid)

	go conn.writeLoop()
	conn.readLoop()
}

func (c *Connection) readLoop() {
	defer func() {
		mu.Lock()
		delete(connections, c.uid)
		mu.Unlock()
		close(c.send)
		_ = c.conn.Close()

		// 用户下线：标记Redis离线
		_ = redis.SetOffline(c.uid)
	}()

	for {
		var msg protocol.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}
	}
}

func (c *Connection) writeLoop() {
	defer func() {
		recover()
		_ = c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

// GetConn 获取用户连接（给推送模块使用）
func GetConn(uid string) *Connection {
	mu.RLock()
	defer mu.RUnlock()
	return connections[uid]
}

func (c *Connection) SendMessage(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- data:
	default:
	}
}
