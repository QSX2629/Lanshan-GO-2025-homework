package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aim/aim/configs"
	gwrouter "github.com/aim/aim/internal/gateway/router"
	"github.com/aim/aim/internal/gateway/session"
	"github.com/aim/aim/internal/gateway/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	cfg := configs.Load()
	log.Printf("[gateway] starting on %s:%d", cfg.Server.Gateway.Host, cfg.Server.Gateway.Port)

	sm := session.NewManager()
	router := gwrouter.NewRouter(sm)

	http.HandleFunc(cfg.Server.Gateway.WSPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[gateway] upgrade error: %v", err)
			return
		}

		var connID uint
		c := ws.NewConn(connID, 0, conn)
		sm.Add(c)

		readCh := make(chan *ws.RawMessage, 128)
		go c.ReadPump(readCh)
		go c.WritePump()

		go func() {
			for msg := range readCh {
				router.Route(msg)
			}
			sm.Remove(c.ID, c.UserID)
		}()
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Gateway.Host, cfg.Server.Gateway.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
